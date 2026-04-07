<script>
  import { ArrowRight, Eye, EyeOff, RotateCcw, Settings2 } from 'lucide-svelte';
  import { generate } from './lib/wasm.js';
  import { loadSettings } from './lib/storage.js';
  import SettingsPanel from './lib/SettingsPanel.svelte';

  let state = 'idle'; // 'idle' | 'generated'
  let settingsOpen = false;
  let passphrase = '';
  let password = '';
  let revealed = false;
  let errorMsg = '';
  let settings = loadSettings();

  async function handleGenerate() {
    if (!passphrase.trim()) return;
    errorMsg = '';

    const opts = {
      passphrase,
      salt: settings.salt || null,
      length: settings.length || null,
      version: settings.version || null,
      noUppers: settings.noUppers,
      noNumbers: settings.noNumbers,
      noSpecials: settings.noSpecials,
      customSpecials: settings.customSpecials,
      noAmbiguous: settings.noAmbiguous,
      customAmbiguous: settings.customAmbiguous,
    };

    const result = await generate(opts);

    if (result.error) {
      errorMsg = result.error;
      return;
    }

    password = result.password;
    revealed = false;

    try {
      await navigator.clipboard.writeText(password);
    } catch {
      // Clipboard permission denied — password still shown
    }

    state = 'generated';
    settingsOpen = false;
  }

  function handleKeydown(e) {
    if (e.key === 'Enter') handleGenerate();
  }

  function startOver() {
    password = '';
    passphrase = '';
    errorMsg = '';
    revealed = false;
    state = 'idle';
  }
</script>

<main>
  <h1>passgen</h1>

  <div class="field" class:output-mode={state === 'generated'}>
    {#if state === 'idle'}
      <!-- Passphrase input -->
      <input
        class="field-input"
        type="password"
        placeholder="enter passphrase"
        bind:value={passphrase}
        on:keydown={handleKeydown}
        autocomplete="off"
        autocorrect="off"
        autocapitalize="off"
        spellcheck="false"
      />
      <div class="btn-slot">
        <button class="action-btn btn-go" on:click={handleGenerate} title="Generate">
          <ArrowRight size={18} strokeWidth={2.5} />
        </button>
      </div>
    {:else}
      <!-- Password output -->
      <div class="field-display" class:revealed>
        {#if revealed}
          {password}
        {:else}
          {'●'.repeat(Math.min(password.length, 28))}
        {/if}
      </div>
      <div class="btn-slot">
        <button
          class="action-btn btn-icon"
          on:click={() => (revealed = !revealed)}
          title={revealed ? 'Hide password' : 'Reveal password'}
        >
          {#if revealed}
            <EyeOff size={17} strokeWidth={2} />
          {:else}
            <Eye size={17} strokeWidth={2} />
          {/if}
        </button>
        <button class="action-btn btn-icon" on:click={startOver} title="Start over">
          <RotateCcw size={17} strokeWidth={2} />
        </button>
      </div>
    {/if}
  </div>

  <div class="copied-hint" class:visible={state === 'generated'}>✓ copied to clipboard</div>

  {#if errorMsg}
    <div class="error-hint">{errorMsg}</div>
  {/if}

  <div class="settings-row">
    <button
      class="settings-toggle"
      class:open={settingsOpen}
      disabled={state === 'generated'}
      on:click={() => (settingsOpen = !settingsOpen)}
    >
      <Settings2 size={13} strokeWidth={2} />
      settings
    </button>
  </div>

  {#if settingsOpen && state === 'idle'}
    <SettingsPanel bind:settings />
  {/if}
</main>

<style>
  :global(*, *::before, *::after) { box-sizing: border-box; margin: 0; padding: 0; }

  :global(body) {
    background: #0d0d0d;
    color: #f0f0f0;
    font-family: system-ui, -apple-system, BlinkMacSystemFont, sans-serif;
    min-height: 100vh;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 1.5rem;
  }

  main {
    width: 100%;
    max-width: 480px;
  }

  h1 {
    font-size: clamp(2.2rem, 8vw, 3.2rem);
    font-weight: 900;
    letter-spacing: 0.18em;
    text-transform: uppercase;
    color: #1e1e1e;
    margin-bottom: 2rem;
    line-height: 1;
    user-select: none;
  }

  /* ── Field (identical dimensions in both states) ── */
  .field {
    background: #1a1a1a;
    border: 1px solid #2e2e2e;
    border-radius: 12px;
    display: flex;
    align-items: center;
    overflow: hidden;
    height: 56px;
    transition: border-color 0.3s ease;
  }
  .field.output-mode { border-color: #c8ff0033; }

  .field-input {
    flex: 1;
    background: transparent;
    border: none;
    outline: none;
    padding: 0 1rem;
    font-size: 0.95rem;
    color: #888;
    font-family: inherit;
    height: 100%;
    min-width: 0;
  }
  .field-input::placeholder { color: #333; }

  .field-display {
    flex: 1;
    padding: 0 1rem;
    height: 100%;
    display: flex;
    align-items: center;
    font-family: monospace;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    min-width: 0;
    letter-spacing: 0.2em;
    font-size: 1rem;
    color: #2e2e2e;
    transition: color 0.2s ease, letter-spacing 0.2s ease;
  }
  .field-display.revealed {
    color: #c8ff00;
    letter-spacing: 0.05em;
    font-size: 0.88rem;
  }

  /* ── Buttons ── */
  .btn-slot {
    display: flex;
    gap: 4px;
    padding: 6px 6px 6px 0;
    flex-shrink: 0;
  }

  .action-btn {
    background: #c8ff00;
    border: none;
    cursor: pointer;
    height: 44px;
    border-radius: 8px;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #0d0d0d;
    transition: background 0.15s;
    flex-shrink: 0;
  }
  .action-btn:hover { background: #d4ff33; }
  .btn-go { width: 44px; }
  .btn-icon { width: 40px; }

  /* ── Hints ── */
  .copied-hint {
    font-size: 0.65rem;
    color: #c8ff0055;
    text-align: right;
    margin-top: 0.45rem;
    letter-spacing: 0.1em;
    text-transform: uppercase;
    font-weight: 600;
    visibility: hidden;
  }
  .copied-hint.visible { visibility: visible; }

  .error-hint {
    font-size: 0.72rem;
    color: #ff4444;
    margin-top: 0.4rem;
    text-align: right;
  }

  /* ── Settings toggle ── */
  .settings-row {
    display: flex;
    justify-content: flex-end;
    margin-top: 0.6rem;
  }

  .settings-toggle {
    display: flex;
    align-items: center;
    gap: 0.35rem;
    font-size: 0.72rem;
    font-weight: 600;
    letter-spacing: 0.08em;
    text-transform: uppercase;
    color: #3a3a3a;
    border: none;
    background: none;
    cursor: pointer;
    padding: 0.25rem 0;
    transition: color 0.15s;
  }
  .settings-toggle:hover:not(:disabled) { color: #555; }
  .settings-toggle.open { color: #c8ff00; }
  .settings-toggle:disabled { opacity: 0.2; cursor: default; }
</style>
