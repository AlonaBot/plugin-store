# AlonaBot 官方插件源

AlonaBot 官方插件源，官方 WebUI 的商店把它列为缺省源。索引只有一个文件 `index.json`，格式见 Core 的 [CONTRACTS §9.2](https://github.com/AlonaBot/AlonaBot/blob/main/spec/CONTRACTS-v0.1.md)：

```json
{"version": 1, "plugins": [{"repo": "owner/name", "role": "plugin", "tags": ["weather"]}]}
```

条目只有三项：`repo`（`owner/name`）、`role`（`plugin`、`adapter`、`agent`、`webui` 之一）和可省略的 `tags`。名称、描述、版本与图标不写进索引，所以索引不固定发行版本，也不需要随每次发布更新。

## 商店快照

商店不逐条请求 GitHub API，而是读取 CI 生成的静态快照 `store.json`。[Store snapshot](.github/workflows/store.yml) 在 `main` 更新、每 6 小时和手动触发时运行 `cmd/storegen`：读取每个条目的最新 Release，下载一个组件归档并按 Release 提供的 sha256 校验，取出包内 `plugin.yaml` 的名称、描述、作者与图标，再记下星数和各平台归档的地址、大小与 sha256。结果强制推送到 `index` 分支，可从以下任一地址读取：

| 来源 | 地址 |
| --- | --- |
| GitHub | `https://raw.githubusercontent.com/AlonaBot/plugin-store/index/store.json` |
| jsDelivr | `https://cdn.jsdelivr.net/gh/AlonaBot/plugin-store@index/store.json` |

单个条目失败时写进快照的 `errors`，其余条目照常发布；全部失败时任务失败，`index` 分支保留上一次的快照。PR 上只生成不发布，用来提前发现新条目的 Release 问题。本地生成：

```sh
GITHUB_TOKEN=$(gh auth token) go run ./cmd/storegen -out dist/store.json
```

## 收录的组件

收录了什么只看 [index.json](index.json)，不在这里另列一份，免得与索引脱节。

## 校验索引

改动 `index.json` 后在仓库根目录运行，只依赖 Node 标准库：

```sh
node scripts/check-index.mjs
```

脚本检查格式版本、字段集合、`owner/name` 形式、角色取值、标签类型与重复仓库（按大小写不敏感比较），报错时给出 JSON 路径。CI 用同一脚本校验，提交索引前本地跑一次能省一个来回。

## 收录新组件

收录申请走 [Issue 表单](.github/ISSUE_TEMPLATE/submit-plugin.yml)，WebUI 的发布弹窗会把它预填好。维护者审核后更新 `index.json`，申请者不需要自己改这个仓库。审核要点见 [CONTRIBUTING.md](CONTRIBUTING.md)。

本仓库只托管索引与快照。生成快照时下载归档只读取声明与图标，不执行任何社区组件；组件由使用方的 Core 按各自仓库的 Release 安装，并按 sha256 校验。
