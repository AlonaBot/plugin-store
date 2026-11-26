# AlonaBot component index

[简体中文](README.zh-CN.md)

The default component source used by the official AlonaBot WebUI. [index.json](index.json) lists component repositories; Core installs their published platform archives.

## Submit a component

Publish a GitHub Release in your component repository, then complete the [submission form](https://github.com/AlonaBot/plugin-store/issues/new?template=submit-plugin.yml). The WebUI store can prefill the form. Maintainers review submissions and update the index.

## Documentation

- [Index format and validation](docs/index-format.md)
- [Contributing](CONTRIBUTING.md)
- [Build and publish a component](https://alonabot.github.io/en/docs/advanced/publish)
- [Use the component store](https://alonabot.github.io/en/docs/usage/webui)

Validate an index change with Node.js:

```sh
node scripts/check-index.mjs
```
