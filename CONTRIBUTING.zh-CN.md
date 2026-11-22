# 参与组件索引维护

[English](CONTRIBUTING.md)

范围确认、依赖审批、提交语言、PR 格式与审查流程遵循 [Core 贡献指南](https://github.com/AlonaBot/AlonaBot/blob/main/CONTRIBUTING.zh-CN.md)。本页补充索引仓库的具体步骤。

## 申请收录

1. 在公开可访问的组件仓库发布 GitHub Release，提供含 `plugin.yaml` 声明及必要资源的平台归档。
2. 填写[收录表单](.github/ISSUE_TEMPLATE/submit-plugin.yml)，提供仓库、角色、可选标签与 Release 地址。WebUI 发布弹窗可预填这些字段。
3. 维护者审核后更新 `index.json`。

审核时确认仓库和 Release 可访问、包内声明与仓库及角色一致、条目有实际用途且没有重复。标签用于检索，审核者在收录前检查组件源码与文档。

## 修改索引

条目仅使用 `repo`、`role` 与可选的 `tags`，遵循[索引格式](docs/index-format.zh-CN.md)。保留与本次修改无关的现有条目。

运行：

```sh
node scripts/check-index.mjs
git diff --check
```

commit message 使用简体中文。PR 标题遵循 Core 的 `type/scope: 摘要`，例如 `docs/index: 更新收录指南`。按 PR 模板填写原因、实际验证与兼容性或风险；维护者在当前提交检查通过后审查并 squash 合并。

CI 使用 Node.js 校验本地索引并检查仓库改动，读取仓库数据，不下载或执行被收录的组件。用户安装包时，由 Core 执行组件安装检查。

安全漏洞通过 [Core 私密报告入口](https://github.com/AlonaBot/AlonaBot/security/advisories/new)提交。公开报告先移除凭据与真实用户数据。
