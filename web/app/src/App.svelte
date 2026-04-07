<script>
  import { generate } from './lib/wasm.js';

  let result = '';
  let error = '';

  async function test() {
    const r = await generate({ passphrase: 'hello', salt: 'world', length: 20 });
    if (r.error) {
      error = r.error;
    } else {
      result = r.password;
    }
  }
</script>

<main>
  <h1>passgen</h1>
  <button on:click={test}>Test WASM</button>
  {#if result}<p>Password: {result}</p>{/if}
  {#if error}<p>Error: {error}</p>{/if}
</main>
