<script>
  import { generate } from './lib/wasm.js';
  import { loadSettings, saveSettings } from './lib/storage.js';

  let settings = loadSettings();
  let result = '';

  async function test() {
    const r = await generate({ passphrase: 'hello', ...settings });
    result = r.password ?? r.error;
    saveSettings(settings);
  }
</script>

<main>
  <h1>passgen</h1>
  <p>salt: "{settings.salt}" | length: {settings.length}</p>
  <button on:click={test}>Test</button>
  {#if result}<p>{result}</p>{/if}
</main>
