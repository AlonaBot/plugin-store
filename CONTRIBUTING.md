# 贡献索引条目

本仓库只维护 `index.json`：官方 WebUI 的缺省插件源。格式以 Core 的 [CONTRACTS §9.2](https://github.com/AlonaBot/AlonaBot/blob/main/spec/CONTRACTS-v0.1.md) 为准，不在这里定义字段或增加未来字段。

## 申请收录

1. 组件在独立的公开仓库发布过 GitHub Release，带 `release.json` 与包内 manifest；官方组件始终签名。
2. 通过 [Issue 表单](.github/ISSUE_TEMPLATE/submit-plugin.yml)提交，填写仓库、角色、可选标签与 Release 地址。WebUI 商店的发布弹窗会预填这四项，不需要手抄格式。
3. 维护者审核后更新 `index.json` 并合并。条目由维护者维护，申请者不直接改这个仓库。

## 审核要点

- 仓库确实可公开访问，最新 Release 提供组件包与 `release.json`。
- `role` 与组件实际角色一致：`plugin` 唯一可多装，`adapter` 按 `adapter_id` 去重，`agent` 与 `webui` 全局各一个。
- 组件不是索引仓库的镜像或空壳，不把索引当成安装包下载点。
- 标签能帮助检索，不堆砌无关词；同一个仓库不在索引里出现两次（大小写不同也不行）。

## 改索引

只改 `index.json`，保持 `version` 为 1、条目只含 `repo`、`role`、`tags`。提交前在仓库根目录运行：

```sh
node scripts/check-index.mjs
```

不要往索引里写名称、描述、版本、摘要、图标或下载地址：这些从组件仓库最新 Release 的 `release.json` 与包内 manifest 取得，写进索引就是重复元数据，还会在发布后过期。

提交信息格式 `type/scope: 摘要`，例如 `feat/index: 收录 weather 插件`；正文说明收录的仓库、审核结果与实际运行的命令。代码注释、commit message 与文档用简体中文。

## 安全性

本仓库的 CI 只读取 `index.json`，不下载、不构建、不执行任何被索引的组件，也不接触发布私钥。索引条目不代表对组件行为的担保：安装时的签名、摘要与平台校验始终在 Core 侧完成。
