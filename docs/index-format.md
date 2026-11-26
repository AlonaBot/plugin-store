# Index format and validation

[简体中文](index-format.zh-CN.md)

Add a component repository to `index.json` using this structure:

```json
{
  "version": 1,
  "plugins": [
    {"repo": "owner/name", "role": "plugin", "tags": ["weather"]}
  ]
}
```

| Field | Meaning |
| --- | --- |
| `version` | Index format version, currently `1`. |
| `plugins` | Component repository entries. |
| `repo` | GitHub `owner/name`, matching the component's source repository. |
| `role` | `plugin`, `adapter`, `agent` or `webui`, matching the package declaration. |
| `tags` | Optional array of strings for store search. |

The WebUI reads display information from the component's Release and platform archive declaration. Keep names, versions, descriptions, icons and download URLs in the component repository. The index follows repositories without pinning their release versions.

Core defines the format in [`contracts/pluginindex`](https://github.com/AlonaBot/AlonaBot/tree/main/contracts/pluginindex). Its [deployment guide](https://alonabot.github.io/en/docs/usage/webui) explains how users install from the store.

## Validate an edit

```sh
node scripts/check-index.mjs
```

The checker uses Node.js standard libraries. It checks the format version, allowed fields, repository names, roles, tag types and duplicate repositories, treating GitHub repository names as case-insensitive. Errors include the affected JSON path. An alternate index file can be passed as the first argument.

The same checker runs in CI. It reads the index file and never executes a component.
