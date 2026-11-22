# 组件索引维护

本仓库维护官方 WebUI 使用的 `index.json` 及其校验器、收录表单和公开文档。

- 贡献流程以 [Core 贡献指南](https://github.com/AlonaBot/AlonaBot/blob/main/CONTRIBUTING.zh-CN.md)为共同依据；本仓差异见 [CONTRIBUTING.zh-CN.md](CONTRIBUTING.zh-CN.md)。
- 格式见[索引说明](docs/index-format.zh-CN.md)，Core 的 `contracts/pluginindex` 定义字段语义。
- 修改前读目标、调用者及当前 diff，保留其他人的改动。Git 写操作使用独立 worktree，提交、推送、合并按用户实际授权执行。
- 维护说明、注释与 commit message 使用简体中文；用户文档维护中英两份。日志与错误使用英文。
- 新依赖先说明用途并取得批准。保持校验器只读取索引，不下载或执行组件。
- 验证运行 `node scripts/check-index.mjs`、`git diff --check`，并核对双语链接、表单字段与 PR 模板对应的检查规则。
