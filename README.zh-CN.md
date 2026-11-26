# AlonaBot 官方组件索引

[English](README.md)

官方 WebUI 使用的缺省组件源。[index.json](index.json) 列出组件仓库，Core 从对应仓库安装已发布的平台归档。

## 申请收录

在组件仓库发布 GitHub Release 后，填写[收录表单](https://github.com/AlonaBot/plugin-store/issues/new?template=submit-plugin.yml)。WebUI 商店可预填表单，维护者审核后更新索引。

## 文档

- [索引格式与校验](docs/index-format.zh-CN.md)
- [贡献指南](CONTRIBUTING.zh-CN.md)
- [构建与发布组件](https://alonabot.github.io/docs/advanced/publish)
- [使用组件商店](https://alonabot.github.io/docs/usage/webui)

修改索引后，使用 Node.js 校验：

```sh
node scripts/check-index.mjs
```
