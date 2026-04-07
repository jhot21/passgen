// Lazy WASM initializer. The Go runtime (window.Go) is loaded by wasm_exec.js
// via the <script> tag in index.html before this module runs.

let initialized = false;
let initPromise = null;

async function init() {
  if (initialized) return;
  const go = new window.Go();
  const result = await WebAssembly.instantiateStreaming(
    fetch('/passgen.wasm'),
    go.importObject,
  );
  go.run(result.instance);
  initialized = true;
}

/**
 * Generate a password using the passgen WASM library.
 * @param {Object} opts
 * @param {string} opts.passphrase - required
 * @param {string} [opts.salt]
 * @param {number} [opts.length]
 * @param {number} [opts.version]
 * @param {boolean} [opts.noUppers]
 * @param {boolean} [opts.noNumbers]
 * @param {boolean} [opts.noSpecials]
 * @param {string|null} [opts.customSpecials] - null = use library default
 * @param {boolean} [opts.noAmbiguous]
 * @param {string|null} [opts.customAmbiguous] - null = use library default
 * @returns {Promise<{password: string}|{error: string}>}
 */
export async function generate(opts) {
  if (!initPromise) {
    initPromise = init();
  }
  await initPromise;

  // Strip null values — the WASM binding treats missing/undefined as "not set"
  const cleaned = Object.fromEntries(
    Object.entries(opts).filter(([, v]) => v !== null && v !== undefined && v !== ''),
  );

  return window.passgenGenerate(cleaned);
}
