'use strict';
const assert = require('node:assert/strict');
const fs = require('node:fs');
const vm = require('node:vm');
const { test } = require('node:test');

test('unchanged state skips rendering, changes and connection recovery render', async () => {
  const elements = new Map();
  let writes = 0;
  let state = { phase: 'ready', selected: false, message: 'Ready' };
  let connected = true;
  const context = vm.createContext({
    URL,
    location: { href: 'http://127.0.0.1:12345/pc/' },
    document: {
      getElementById(id) {
        if (!elements.has(id)) elements.set(id, new Proxy({ addEventListener() {} }, {
          set(target, key, value) { writes++; target[key] = value; return true; }
        }));
        return elements.get(id);
      }
    },
    fetch: async () => ({ ok: connected, text: async () => JSON.stringify(state) }),
    setTimeout() {}
  });
  vm.runInContext(fs.readFileSync(`${__dirname}/web/app.js`, 'utf8'), context);
  await new Promise(setImmediate); // Finish the initial poll.
  assert.equal(elements.get('phase').textContent, '接続の準備完了');
  writes = 0;
  await vm.runInContext('poll()', context);
  assert.equal(writes, 0);

  state = { ...state, lastRequest: 'GET → HTTP 200' };
  await vm.runInContext('poll()', context);
  assert.equal(elements.get('request').textContent, 'ファイル要求：GET → HTTP 200');
  assert.ok(writes > 0);

  connected = false;
  await vm.runInContext('poll()', context);
  assert.equal(elements.get('phase').textContent, '接続が終了、または中断しました');
  connected = true;
  await vm.runInContext('poll()', context);
  assert.equal(elements.get('phase').textContent, '接続の準備完了');
});
