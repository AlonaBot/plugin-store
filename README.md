# AlonaBot 官方插件源

AlonaBot 官方插件源，官方 WebUI 的商店把它列为缺省源。索引只有一个文件 `index.json`，格式见 Core 的 [CONTRACTS §9.2](https://github.com/AlonaBot/AlonaBot/blob/main/spec/CONTRACTS-v0.1.md)：

```json
{"version": 1, "plugins": [{"repo": "owner/name", "role": "plugin", "tags": ["weather"]}]}
```

条目只有三项：`repo`（`owner/name`）、`role`（`plugin`、`adapter`、`agent`、`webui` 之一）和可省略的 `tags`。名称、描述、版本与图标不写进索引，由商店从组件仓库最新 Release 的 `release.json` 与包内 manifest 读取，所以索引不固定发行版本，也不需要随每次发布更新。

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

本仓库只托管索引，不下载或执行任何社区组件；组件由使用方的 Core 按各自仓库的 Release 与签名规则安装。

## 许可

[MIT](LICENSE)。
