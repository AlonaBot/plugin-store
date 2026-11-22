# Contributing to the component index

[简体中文](CONTRIBUTING.zh-CN.md)

Follow the [Core contribution guide](https://github.com/AlonaBot/AlonaBot/blob/main/CONTRIBUTING.md) for scope, dependency approval, commit language, PR format and review. This guide adds the index-specific steps.

## Request a listing

1. Publish a GitHub Release in a publicly accessible component repository. Provide platform archives containing the component's `plugin.yaml` declaration and required resources.
2. Complete the [submission form](.github/ISSUE_TEMPLATE/submit-plugin.yml) with the repository, role, optional tags and Release URL. The WebUI publish dialog can prefill these fields.
3. A maintainer reviews the submission and updates `index.json`.

Reviewers check that the repository and Release are accessible, the archive declaration matches the repository and role, and the entry is useful and unique. Use relevant search tags. Review the component's source and documentation before listing it.

## Change the index

Keep entries to `repo`, `role` and optional `tags`, following the [index format](docs/index-format.md). Preserve existing entries unless the change specifically concerns them.

Run:

```sh
node scripts/check-index.mjs
git diff --check
```

Use a Simplified Chinese commit message. PR titles follow Core's `type/scope: 摘要` convention, for example `docs/index: 更新收录指南`. Complete the PR template with the reason, actual verification and compatibility or risk information. The maintainer reviews and squash-merges the PR after its current checks pass.

CI validates the local index with Node.js and checks repository changes. It reads repository data without downloading or executing listed components. Core performs component installation checks when a user installs a package.

Report vulnerabilities through [Core's private security reporting channel](https://github.com/AlonaBot/AlonaBot/security/advisories/new). Remove credentials and real user data from public reports.
