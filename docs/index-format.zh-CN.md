# 索引格式与校验

[English](index-format.md)

按以下结构将组件仓库加入 `index.json`：

```json
{
  "version": 1,
  "plugins": [
    {"repo": "owner/name", "role": "plugin", "tags": ["weather"]}
  ]
}
```

| 字段 | 含义 |
| --- | --- |
| `version` | 索引格式版本，当前为 `1`。 |
| `plugins` | 组件仓库条目。 |
| `repo` | GitHub `owner/name`，与组件的来源仓库一致。 |
| `role` | `plugin`、`adapter`、`agent` 或 `webui`，与包内声明一致。 |
| `tags` | 可选的字符串数组，用于商店检索。 |

WebUI 从组件的 Release 与平台归档声明读取展示信息。名称、版本、描述、图标与下载地址由组件仓库维护；索引跟随仓库，不固定发行版本。

Core 在 [`contracts/pluginindex`](https://github.com/AlonaBot/AlonaBot/tree/main/contracts/pluginindex) 定义格式。[部署指南](https://github.com/AlonaBot/AlonaBot/blob/main/docs/deployment.zh-CN.md)说明用户如何从商店安装。

## 校验改动

```sh
node scripts/check-index.mjs
```

校验器只使用 Node.js 标准库，检查格式版本、允许的字段、仓库名称、角色、标签类型与重复仓库。GitHub 仓库名按大小写不敏感比较，错误包含对应的 JSON 路径。第一个参数可指定其他索引文件。

CI 运行同一校验器，只读取索引文件，不执行组件。
