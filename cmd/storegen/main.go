// Command storegen 把 index.json 展开为商店直接读取的静态快照 store.json。
//
// 每个条目读取最新 Release、下载一个组件归档并按 sha256 校验，只取包内
// plugin.yaml 与图标，不执行归档内容。单个条目失败记录在 errors 中，
// 不影响其余条目；全部失败时以非零状态退出，保留上一次发布的快照。
package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/AlonaBot/AlonaBot/contracts/manifest"
	"github.com/AlonaBot/AlonaBot/contracts/pluginindex"
)

// 远端内容不可信；与 WebUI 原先在线抓取时的上限一致。
const (
	requestTimeout       = 60 * time.Second
	metadataBytes        = 1 << 20
	archiveBytes         = 64 << 20
	archiveExpandedBytes = 256 << 20
	archiveFileBytes     = 128 << 20
	archiveEntries       = 2048
	manifestBytes        = 256 << 10
	iconBytes            = 1 << 20
)

// store 是 store.json 的顶层结构，字段名与 WebUI 商店接口一致。
type store struct {
	Version     int          `json:"version"`
	GeneratedAt string       `json:"generated_at"`
	Items       []entry      `json:"items"`
	Errors      []entryError `json:"errors"`
}

type entry struct {
	Repo        string   `json:"repo"`
	Role        string   `json:"role"`
	Tags        []string `json:"tags"`
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Version     string   `json:"version"`
	Author      string   `json:"author,omitempty"`
	Stars       *int     `json:"stars,omitempty"`
	IconURL     string   `json:"icon_url,omitempty"`
	ReleaseURL  string   `json:"release_url"`
	PublishedAt string   `json:"published_at"`
	Assets      []asset  `json:"assets"`
}

// asset 是一个平台归档；静态 WebUI 的 OS 与 Arch 为 any。
type asset struct {
	Name   string `json:"name"`
	OS     string `json:"os"`
	Arch   string `json:"arch"`
	URL    string `json:"url"`
	Size   int64  `json:"size"`
	SHA256 string `json:"sha256"`
}

type entryError struct {
	Repo  string `json:"repo"`
	Error string `json:"error"`
}

type release struct {
	Tag         string `json:"tag_name"`
	URL         string `json:"html_url"`
	PublishedAt string `json:"published_at"`
	Assets      []struct {
		Name   string `json:"name"`
		URL    string `json:"browser_download_url"`
		Size   int64  `json:"size"`
		Digest string `json:"digest"`
	} `json:"assets"`
}

type client struct {
	http  *http.Client
	token string
}

func main() {
	indexPath := flag.String("index", "index.json", "source index path")
	out := flag.String("out", "dist/store.json", "output path")
	flag.Parse()
	if err := run(context.Background(), *indexPath, *out); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context, indexPath, out string) error {
	raw, err := os.ReadFile(indexPath)
	if err != nil {
		return err
	}
	index, err := pluginindex.Parse(raw)
	if err != nil {
		return err
	}
	c := &client{http: &http.Client{Timeout: requestTimeout}, token: os.Getenv("GITHUB_TOKEN")}
	result := store{Version: 1, GeneratedAt: time.Now().UTC().Format(time.RFC3339), Items: []entry{}, Errors: []entryError{}}
	for _, plugin := range index.Plugins {
		value, err := c.entry(ctx, plugin)
		if err != nil {
			log.Printf("%s: %v", plugin.Repo, err)
			result.Errors = append(result.Errors, entryError{Repo: plugin.Repo, Error: err.Error()})
			continue
		}
		result.Items = append(result.Items, value)
	}
	if len(result.Items) == 0 && len(index.Plugins) > 0 {
		return errors.New("every index entry failed; keeping the previous snapshot")
	}
	encoded, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(out, append(encoded, '\n'), 0o644); err != nil {
		return err
	}
	log.Printf("wrote %s: %d items, %d errors", out, len(result.Items), len(result.Errors))
	return nil
}

func (c *client) get(ctx context.Context, address, accept string, limit int64) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, address, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", accept)
	if c.token != "" && strings.HasPrefix(address, "https://api.github.com/") {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	response, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	raw, readErr := io.ReadAll(io.LimitReader(response.Body, limit+1))
	if err := errors.Join(readErr, response.Body.Close()); err != nil {
		return nil, err
	}
	if int64(len(raw)) > limit {
		return nil, fmt.Errorf("GET %s: response exceeds %d bytes", address, limit)
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET %s: HTTP %d: %s", address, response.StatusCode, raw)
	}
	return raw, nil
}

func (c *client) entry(ctx context.Context, plugin pluginindex.Plugin) (entry, error) {
	value := entry{Repo: plugin.Repo, Role: plugin.Role, Tags: plugin.Tags, Assets: []asset{}}
	if value.Tags == nil {
		value.Tags = []string{}
	}
	// 1. 读取最新 Release，校验发布页地址与 SemVer 标签。
	raw, err := c.get(ctx, "https://api.github.com/repos/"+plugin.Repo+"/releases/latest", "application/vnd.github+json", metadataBytes)
	if err != nil {
		return value, err
	}
	var publication release
	if err := json.Unmarshal(raw, &publication); err != nil {
		return value, fmt.Errorf("parse release: %w", err)
	}
	if !strings.EqualFold(publication.URL, "https://github.com/"+plugin.Repo+"/releases/tag/"+publication.Tag) {
		return value, fmt.Errorf("invalid release URL %q", publication.URL)
	}
	version := strings.TrimPrefix(publication.Tag, "v")
	if !manifest.ValidVersion(version) {
		return value, fmt.Errorf("invalid SemVer tag %q", publication.Tag)
	}
	value.Version, value.ReleaseURL, value.PublishedAt = version, publication.URL, publication.PublishedAt

	// 2. 收集组件归档：名称必须是 <id>-<version>-<platform>.tar.gz，地址指向本仓库本标签。
	var first *asset
	for _, a := range publication.Assets {
		platform, ok := archivePlatform(a.Name, version)
		if !ok {
			continue
		}
		expected := "https://github.com/" + plugin.Repo + "/releases/download/" + publication.Tag + "/" + url.PathEscape(a.Name)
		if !strings.EqualFold(a.URL, expected) {
			return value, fmt.Errorf("asset %s URL %q does not match the release", a.Name, a.URL)
		}
		digest, ok := strings.CutPrefix(a.Digest, "sha256:")
		if !ok || len(digest) != sha256.Size*2 {
			return value, fmt.Errorf("asset %s has no sha256 digest", a.Name)
		}
		value.Assets = append(value.Assets, asset{Name: a.Name, OS: platform.OS, Arch: platform.Arch, URL: a.URL, Size: a.Size, SHA256: digest})
		if first == nil || platform.OS == "any" {
			first = &value.Assets[len(value.Assets)-1]
		}
	}
	if first == nil {
		return value, errors.New("release has no component archive")
	}

	// 3. 下载一个归档，按 sha256 与大小校验后读取声明与图标。
	archive, err := c.get(ctx, first.URL, "application/octet-stream", archiveBytes)
	if err != nil {
		return value, err
	}
	sum := sha256.Sum256(archive)
	if int64(len(archive)) != first.Size || hex.EncodeToString(sum[:]) != first.SHA256 {
		return value, fmt.Errorf("asset %s does not match its size or sha256", first.Name)
	}
	declaration, icon, err := readManifest(archive)
	if err != nil {
		return value, fmt.Errorf("%s: %w", first.Name, err)
	}
	m := declaration.Manifest
	switch {
	case m.Role != plugin.Role:
		return value, fmt.Errorf("package role %q does not match the index entry %q", m.Role, plugin.Role)
	case m.Version != version:
		return value, errors.New("package version does not match the release tag")
	case !strings.EqualFold(m.Source, "github.com/"+plugin.Repo):
		return value, fmt.Errorf("package source %q does not match repository %s", m.Source, plugin.Repo)
	}
	for _, a := range value.Assets {
		if !strings.HasPrefix(a.Name, m.ID+"-"+m.Version+"-") {
			return value, fmt.Errorf("archive %s does not match package id %s", a.Name, m.ID)
		}
	}
	// 安装来源用包内声明的规范大小写：Core 会逐字比较。
	value.Repo = strings.TrimPrefix(m.Source, "github.com/")
	value.ID, value.Name, value.Description, value.Author = m.ID, m.Name, m.Description, m.Author
	if icon != nil {
		mime := map[string]string{".svg": "image/svg+xml", ".png": "image/png", ".webp": "image/webp"}[path.Ext(m.Icon)]
		value.IconURL = "data:" + mime + ";base64," + base64.StdEncoding.EncodeToString(icon)
	}

	// 4. 星数只用于排序展示。
	raw, err = c.get(ctx, "https://api.github.com/repos/"+plugin.Repo, "application/vnd.github+json", metadataBytes)
	if err != nil {
		return value, err
	}
	var repository struct {
		Stars *int `json:"stargazers_count"`
	}
	if err := json.Unmarshal(raw, &repository); err != nil {
		return value, err
	}
	value.Stars = repository.Stars
	return value, nil
}

// archivePlatform 从 <id>-<version>-<os>-<arch>.tar.gz 或 <id>-<version>-any.tar.gz 取平台。
func archivePlatform(name, version string) (manifest.Platform, bool) {
	stem, ok := strings.CutSuffix(name, ".tar.gz")
	if !ok {
		return manifest.Platform{}, false
	}
	if strings.HasSuffix(stem, "-"+version+"-any") {
		return manifest.Platform{OS: "any", Arch: "any"}, true
	}
	for _, p := range manifest.Platforms {
		if strings.HasSuffix(stem, "-"+version+"-"+p.OS+"-"+p.Arch) {
			return p, true
		}
	}
	return manifest.Platform{}, false
}

func readManifest(archive []byte) (*manifest.Result, []byte, error) {
	declaration, err := archiveFile(archive, "plugin.yaml", manifestBytes)
	if err != nil {
		return nil, nil, err
	}
	if declaration == nil {
		return nil, nil, errors.New("archive has no plugin.yaml")
	}
	parsed, err := manifest.Parse(declaration)
	if err != nil {
		return nil, nil, err
	}
	if parsed.Manifest.Icon == "" {
		return parsed, nil, nil
	}
	icon, err := archiveFile(archive, parsed.Manifest.Icon, iconBytes)
	return parsed, icon, err
}

// archiveFile 在受限展开下读取一个常规文件；不存在时返回 nil。
func archiveFile(archive []byte, name string, limit int64) ([]byte, error) {
	compressed, err := gzip.NewReader(bytes.NewReader(archive))
	if err != nil {
		return nil, err
	}
	defer compressed.Close()
	expanded := &io.LimitedReader{R: compressed, N: archiveExpandedBytes + 1}
	reader := tar.NewReader(expanded)
	var content []byte
	for count := 0; ; count++ {
		header, err := reader.Next()
		if expanded.N == 0 {
			return nil, fmt.Errorf("archive expanded content exceeds %d bytes", archiveExpandedBytes)
		}
		if errors.Is(err, io.EOF) {
			return content, nil
		}
		if err != nil {
			return nil, err
		}
		if count >= archiveEntries {
			return nil, fmt.Errorf("archive exceeds %d entries", archiveEntries)
		}
		if header.Size > archiveFileBytes {
			return nil, fmt.Errorf("archive file %s exceeds %d bytes", header.Name, archiveFileBytes)
		}
		if header.Typeflag != tar.TypeReg || header.Name != name {
			continue
		}
		content, err = io.ReadAll(io.LimitReader(reader, limit+1))
		if err != nil {
			return nil, err
		}
		if int64(len(content)) > limit {
			return nil, fmt.Errorf("archive file %s exceeds %d bytes", name, limit)
		}
	}
}
