// 校验 index.json 是否符合 Core 的插件源格式（CONTRACTS §9.2）。
// 用法：node scripts/check-index.mjs [文件路径]，缺省校验仓库根目录的 index.json。
import { readFileSync } from 'node:fs';
import { fileURLToPath } from 'node:url';

const indexPath = process.argv[2] ?? fileURLToPath(new URL('../index.json', import.meta.url));

/** 抛出带 JSON 路径的错误，与 Core 的 pluginindex 报错风格一致。 */
function fail(path, message) {
  throw new Error(`index.json:${path}: ${message}`);
}

let document;
try {
  document = JSON.parse(readFileSync(indexPath, 'utf8'));
} catch (err) {
  fail('$', err.message);
}

// 1. 顶层只有 version 与 plugins；顶层数组或 null 在下一行按对象拒绝。
if (typeof document !== 'object' || Array.isArray(document)) fail('$', 'expected an object');
for (const key of Object.keys(document)) {
  if (key !== 'version' && key !== 'plugins') fail('$', `unexpected field ${JSON.stringify(key)}`);
}
if (document.version !== 1) fail('$.version', `expected 1, got ${JSON.stringify(document.version)}`);
if (!Array.isArray(document.plugins)) fail('$.plugins', 'expected an array');

/** 角色取自 CONTRACTS §9.1 的四种组件类型。 */
const roles = new Set(['plugin', 'adapter', 'agent', 'webui']);

/** owner/name 形式，与 Core 的 repositoryName 正则逐字相同。 */
const repositoryName = /^[A-Za-z0-9_.-]+\/[A-Za-z0-9_.-]+$/;

const seen = new Map();

document.plugins.forEach((entry, i) => {
  const path = `$.plugins[${i}]`;

  // 2. 每条只允许 repo、role、tags；名称、版本与描述来自组件仓库的 Release。
  if (typeof entry !== 'object' || entry === null || Array.isArray(entry)) fail(path, 'expected an object');
  for (const key of Object.keys(entry)) {
    if (key !== 'repo' && key !== 'role' && key !== 'tags') fail(path, `unexpected field ${JSON.stringify(key)}`);
  }

  if (typeof entry.repo !== 'string' || !repositoryName.test(entry.repo)) fail(`${path}.repo`, 'expected owner/name');
  if (entry.repo.split('/').some(part => part === '.' || part === '..')) fail(`${path}.repo`, 'expected owner/name');

  if (!roles.has(entry.role)) fail(`${path}.role`, `unknown role ${JSON.stringify(entry.role)}`);

  if (entry.tags !== undefined) {
    if (!Array.isArray(entry.tags)) fail(`${path}.tags`, 'expected an array');
    entry.tags.forEach((tag, j) => {
      if (typeof tag !== 'string') fail(`${path}.tags[${j}]`, 'expected a string');
    });
  }

  // 3. 同一仓库只出现一次；GitHub 名称不区分大小写，报错带上首次出现的路径。
  const key = entry.repo.toLowerCase();
  if (seen.has(key)) fail(`${path}.repo`, `duplicate repository ${entry.repo}, first listed at ${seen.get(key)}`);
  seen.set(key, `${path}.repo`);
});

console.log(`index.json ok: version ${document.version}, ${document.plugins.length} entries`);