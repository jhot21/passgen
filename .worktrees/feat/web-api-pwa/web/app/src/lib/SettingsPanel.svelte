<script>
  import { Eye, EyeOff } from 'lucide-svelte';
  import { saveSettings } from './storage.js';

  export let settings;

  let saltRevealed = false;

  function save() {
    saveSettings(settings);
  }

  function resetSpecials() {
    settings.customSpecials = null;
    save();
  }

  function resetAmbiguous() {
    settings.customAmbiguous = null;
    save();
  }
</script>

<div class="settings-panel">

  <!-- Salt -->
  <div class="setting-row">
    <span class="setting-label">Salt</span>
    <div class="salt-wrap">
      <input
        type={saltRevealed ? 'text' : 'password'}
        class="setting-input setting-input-wide"
        placeholder="personal salt"
        bind:value={settings.salt}
        on:input={save}
      />
      <button
        class="action-btn btn-icon-sm"
        on:click={() => (saltRevealed = !saltRevealed)}
        title={saltRevealed ? 'Hide salt' : 'Reveal salt'}
      >
        {#if saltRevealed}
          <EyeOff size={14} />
        {:else}
          <Eye size={14} />
        {/if}
      </button>
    </div>
  </div>

  <div class="divider" />

  <!-- Length -->
  <div class="setting-row">
    <span class="setting-label">Length</span>
    <input
      type="number"
      class="setting-input setting-num"
      min="1"
      max="512"
      bind:value={settings.length}
      on:input={save}
    />
  </div>

  <!-- Version -->
  <div class="setting-row">
    <span class="setting-label">Version</span>
    <select class="setting-select" bind:value={settings.version} on:change={save}>
      <option value={1}>1</option>
      <option value={2}>2</option>
      <option value={3}>3</option>
    </select>
  </div>

  <div class="divider" />

  <!-- No Uppercase -->
  <div class="setting-row">
    <span class="setting-label">No uppercase</span>
    <button
      class="toggle"
      class:on={settings.noUppers}
      on:click={() => { settings.noUppers = !settings.noUppers; save(); }}
      role="switch"
      aria-checked={settings.noUppers}
    />
  </div>

  <!-- No Numbers -->
  <div class="setting-row">
    <span class="setting-label">No numbers</span>
    <button
      class="toggle"
      class:on={settings.noNumbers}
      on:click={() => { settings.noNumbers = !settings.noNumbers; save(); }}
      role="switch"
      aria-checked={settings.noNumbers}
    />
  </div>

  <!-- No Specials -->
  <div class="setting-row">
    <span class="setting-label">No specials</span>
    <button
      class="toggle"
      class:on={settings.noSpecials}
      on:click={() => { settings.noSpecials = !settings.noSpecials; save(); }}
      role="switch"
      aria-checked={settings.noSpecials}
    />
  </div>

  <!-- Custom Specials -->
  <div class="setting-row">
    <span class="setting-label">Custom specials</span>
    <div class="input-reset-wrap">
      <input
        type="text"
        class="setting-input"
        placeholder="!@#$%^&*"
        value={settings.customSpecials ?? ''}
        on:input={(e) => { settings.customSpecials = e.target.value || null; save(); }}
      />
      <button class="reset-link" on:click={resetSpecials}>reset</button>
    </div>
  </div>

  <div class="divider" />

  <!-- No Ambiguous -->
  <div class="setting-row">
    <span class="setting-label">No ambiguous</span>
    <button
      class="toggle"
      class:on={settings.noAmbiguous}
      on:click={() => { settings.noAmbiguous = !settings.noAmbiguous; save(); }}
      role="switch"
      aria-checked={settings.noAmbiguous}
    />
  </div>

  <!-- Custom Ambiguous -->
  <div class="setting-row">
    <span class="setting-label">Custom ambiguous</span>
    <div class="input-reset-wrap">
      <input
        type="text"
        class="setting-input"
        placeholder="Il1O05S"
        value={settings.customAmbiguous ?? ''}
        on:input={(e) => { settings.customAmbiguous = e.target.value || null; save(); }}
      />
      <button class="reset-link" on:click={resetAmbiguous}>reset</button>
    </div>
  </div>

</div>

<style>
  .settings-panel {
    margin-top: 0.75rem;
    background: #111;
    border: 1px solid #1e1e1e;
    border-radius: 12px;
    padding: 1.25rem 1.25rem 1rem;
    display: flex;
    flex-direction: column;
    gap: 0.9rem;
  }

  .setting-row {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 1rem;
  }

  .setting-label {
    font-size: 0.78rem;
    color: #555;
    font-weight: 500;
    letter-spacing: 0.04em;
    white-space: nowrap;
  }

  .setting-input {
    background: #1a1a1a;
    border: 1px solid #2a2a2a;
    border-radius: 6px;
    padding: 0.35rem 0.6rem;
    font-size: 0.78rem;
    color: #888;
    outline: none;
    font-family: monospace;
    letter-spacing: 0.05em;
    width: 100px;
  }
  .setting-input:focus { border-color: #3a3a3a; color: #aaa; }
  .setting-input-wide { width: 130px; }
  .setting-num { width: 60px; text-align: center; font-family: inherit; }

  .setting-select {
    background: #1a1a1a;
    border: 1px solid #2a2a2a;
    border-radius: 6px;
    padding: 0.35rem 0.5rem;
    font-size: 0.78rem;
    color: #888;
    outline: none;
    width: 60px;
    text-align: center;
    cursor: pointer;
    appearance: none;
  }

  .salt-wrap {
    display: flex;
    align-items: center;
    gap: 4px;
  }

  .action-btn {
    background: #c8ff00;
    border: none;
    cursor: pointer;
    border-radius: 6px;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #0d0d0d;
    transition: background 0.15s;
    flex-shrink: 0;
  }
  .action-btn:hover { background: #d4ff33; }
  .btn-icon-sm { width: 30px; height: 30px; }

  .toggle {
    width: 34px;
    height: 18px;
    background: #222;
    border-radius: 9px;
    position: relative;
    cursor: pointer;
    border: 1px solid #2e2e2e;
    flex-shrink: 0;
    padding: 0;
    transition: background 0.15s, border-color 0.15s;
  }
  .toggle::after {
    content: '';
    position: absolute;
    width: 12px;
    height: 12px;
    background: #444;
    border-radius: 50%;
    top: 2px;
    left: 2px;
    transition: transform 0.15s, background 0.15s;
  }
  .toggle.on { background: #c8ff00; border-color: #c8ff00; }
  .toggle.on::after { transform: translateX(16px); background: #0d0d0d; }

  .input-reset-wrap {
    display: flex;
    align-items: center;
    gap: 0.4rem;
  }

  .reset-link {
    font-size: 0.65rem;
    color: #333;
    cursor: pointer;
    text-decoration: underline;
    text-underline-offset: 2px;
    background: none;
    border: none;
    padding: 0;
    white-space: nowrap;
  }
  .reset-link:hover { color: #555; }

  .divider {
    height: 1px;
    background: #1e1e1e;
  }
</style>
