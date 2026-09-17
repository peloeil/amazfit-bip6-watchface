const path = require('path');
const { spawnSync } = require('child_process');

const cliRoot = path.dirname(require.resolve('@zeppos/zeus-cli/package.json'));
const result = spawnSync(
  process.execPath,
  [path.join(cliRoot, 'bin/main.js'), ...process.argv.slice(2)],
  {
    cwd: __dirname,
    stdio: 'inherit',
    env: {
      ...process.env,
      NODE_PATH: [
        path.join(cliRoot, 'private-modules'),
        process.env.NODE_PATH
      ].filter(Boolean).join(path.delimiter)
    }
  }
);

if (result.error) {
  console.error(result.error.message);
  process.exit(1);
}

process.exit(result.status === null ? 1 : result.status);
