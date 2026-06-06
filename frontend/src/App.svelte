<script>
  import { fade, fly, scale } from 'svelte/transition';
  import { flip } from 'svelte/animate';
  import { cubicOut } from 'svelte/easing';

  let mode = $state('time');
  let values = $state([1000000, 1000000000]);
  let inputValue = $state('');

  let apiResults = $state(null);
  let fetchError = $state(null);
  let lastReqId = 0;
  let loading = $state(false);

  const PRESETS = {
    time: [
      { label: 'A day',       values: [1, 60, 3600, 86400] },
      { label: '1M vs 1B s',  values: [1000000, 1000000000] },
      { label: 'Geological',  values: [31536000, 3153600000, 4.35e17] },
    ],
    distance: [
      { label: 'Everyday',    values: [1, 100, 1000, 1000000] },
      { label: 'Planetary',   values: [12742000000, 384400000000, 1.496e14] },
      { label: 'Cosmic',      values: [1.496e14, 9.461e18, 9.5e21] },
    ],
  };

  const UNIT_LABELS = {
    time: {
      second: 'second', minute: 'minute', hour: 'hour', day: 'day',
      week: 'week', month: 'month', year: 'year', decade: 'decade',
      century: 'century', millennium: 'millennium',
      geological_era: 'geological era', age_of_earth: 'age of Earth',
      age_of_universe: 'age of universe',
    },
    distance: {
      millimeter: 'millimeter', centimeter: 'centimeter', meter: 'meter',
      football_field: 'football field', kilometer: 'kilometer',
      ten_kilometer_city: '10 km city',
      hundred_kilometer_country: '100 km country',
      half_the_globe: 'half the globe', spain_to_argentina: 'Spain to Argentina',
      london_to_sydney: 'London to Sydney', earth_to_moon: 'Earth to Moon',
      earth_to_sun: 'Earth to Sun', earth_to_mars: 'Earth to Mars',
      width_solar_system: 'solar system width',
    },
  };

  $effect(() => {
    const m = mode;
    const vs = values;
    if (vs.length === 0) {
      apiResults = null;
      loading = false;
      return;
    }
    const id = ++lastReqId;
    loading = true;
    const timer = setTimeout(async () => {
      try {
        const res = await fetch(`/absolute/convert/${m}`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ values: vs }),
        });
        if (!res.ok) throw new Error(`HTTP ${res.status}`);
        const data = await res.json();
        if (id === lastReqId) {
          apiResults = data.results;
          fetchError = null;
        }
      } catch (e) {
        if (id === lastReqId) fetchError = String(e);
      } finally {
        if (id === lastReqId) loading = false;
      }
    }, 200);
    return () => clearTimeout(timer);
  });

  const logWidths = $derived.by(() => {
    if (values.length === 0) return [];
    const logs = values.map((v) => (v > 0 ? Math.log10(v) : 0));
    const min = Math.min(...logs);
    const max = Math.max(...logs);
    const range = max - min || 1;
    return logs.map((l) => ((l - min) / range) * 100);
  });

  function addValue() {
    const text = inputValue.trim();
    if (!text) return;
    const nums = text
      .split(/[\s,;]+/)
      .map(parseFloat)
      .filter((n) => Number.isFinite(n) && n > 0);
    if (nums.length > 0) {
      values = [...values, ...nums];
      inputValue = '';
    }
  }

  function removeValue(i) {
    values = values.filter((_, idx) => idx !== i);
  }

  function loadPreset(p) {
    values = [...p.values];
  }

  function onKeydown(e) {
    if (e.key === 'Enter' || e.key === ',') {
      e.preventDefault();
      addValue();
    } else if (e.key === 'Backspace' && inputValue === '' && values.length > 0) {
      values = values.slice(0, -1);
    }
  }

  function onPaste(e) {
    e.preventDefault();
    const text = e.clipboardData.getData('text');
    inputValue = text;
    addValue();
  }

  function setMode(m) {
    if (m === mode) return;
    mode = m;
    apiResults = null;
  }

  function formatScaled(v) {
    if (!Number.isFinite(v)) return String(v);
    const abs = Math.abs(v);
    if (abs === 0) return '0';
    if (abs >= 1e6 || abs < 1e-3) return v.toExponential(2);
    if (abs >= 100) return v.toFixed(0);
    if (abs >= 1) return v.toFixed(2);
    return v.toFixed(3);
  }

  function formatInput(v) {
    if (!Number.isFinite(v)) return String(v);
    if (Number.isInteger(v) && Math.abs(v) < 1e16) return v.toString();
    return v.toExponential(2);
  }

  function unitLabel(key) {
    return UNIT_LABELS[mode]?.[key] ?? key.replace(/_/g, ' ');
  }

  function pluralize(n, word) {
    return n === 1 ? word : word + 's';
  }
</script>

<div class="container">
  <header>
    <h1>Make Sense <span class="of">of</span> Magnitudes</h1>
    <p class="hint">
      Type a few numbers. See what they mean as <strong>time</strong> or <strong>distance</strong>.
    </p>
  </header>

  <section class="panel">
    <div class="mode">
      <button class:active={mode === 'time'} onclick={() => setMode('time')}>Time</button>
      <button class:active={mode === 'distance'} onclick={() => setMode('distance')}>Distance</button>
      <span class="mode-hint">
        base unit: <code>{mode === 'time' ? 'second' : 'millimeter'}</code>
      </span>
    </div>

    <div class="presets">
      <span class="label">Presets</span>
      {#each PRESETS[mode] as p (p.label)}
        <button class="preset" onclick={() => loadPreset(p)}>{p.label}</button>
      {/each}
    </div>
  </section>

  <section class="panel input-panel">
    <div class="label">Your numbers</div>
    <div class="chips">
      {#each values as v, i (i + ':' + v)}
        <span
          class="chip"
          animate:flip={{ duration: 250 }}
          in:scale={{ duration: 200, start: 0.7, easing: cubicOut }}
          out:fade={{ duration: 150 }}
        >
          {formatInput(v)}
          <button class="x" onclick={() => removeValue(i)} aria-label="Remove {formatInput(v)}">×</button>
        </span>
      {/each}
      <input
        type="text"
        inputmode="decimal"
        bind:value={inputValue}
        onkeydown={onKeydown}
        onpaste={onPaste}
        placeholder={values.length === 0 ? 'Type a number, press Enter…' : '+ add another'}
      />
    </div>
  </section>

  <section class="panel results">
    <div class="results-head">
      <span class="label">Comparison</span>
      {#if loading}<span class="loading">updating…</span>{/if}
    </div>

    {#if fetchError}
      <p class="error">API error: {fetchError}<br /><small>Is the Go server running on :8080?</small></p>
    {:else if values.length === 0}
      <p class="empty">Add some numbers above to see their comparison.</p>
    {:else if apiResults}
      <div class="row head">
        <span>Input</span>
        <span>Equals</span>
        <span class="bar-col">Relative size <small>(log)</small></span>
      </div>
      {#each values as v, i (i + ':' + v)}
        {@const r = apiResults[i]}
        {@const label = unitLabel(r.UnitName)}
        {@const s = formatScaled(r.ScaledValue)}
        <div class="row" animate:flip={{ duration: 250 }} in:fly={{ y: 6, duration: 200 }}>
          <span class="num">{formatInput(v)}</span>
          <span class="equals">
            <span class="scaled">{s}</span>
            <span class="unit">{pluralize(parseFloat(s), label)}</span>
          </span>
          <span class="bar-col">
            <span class="bar-track">
              <span class="bar" style="width: {Math.max(2, logWidths[i])}%"></span>
            </span>
          </span>
        </div>
      {/each}
    {/if}
  </section>
</div>

<style>
  .container {
    max-width: 980px;
    margin: 0 auto;
    padding: 32px 24px 64px;
    display: flex;
    flex-direction: column;
    gap: 16px;
  }

  header h1 {
    margin: 0;
    font-size: 30px;
    letter-spacing: -0.02em;
  }
  header h1 .of {
    color: var(--muted);
    font-weight: 400;
  }
  .hint {
    margin: 6px 0 0;
    color: var(--muted);
    font-size: 14px;
  }

  .panel {
    background: var(--panel);
    border-radius: 12px;
    padding: 16px 18px;
    display: flex;
    flex-direction: column;
    gap: 12px;
  }

  .label {
    color: var(--muted);
    text-transform: uppercase;
    letter-spacing: 0.06em;
    font-size: 11px;
  }

  .mode {
    display: flex;
    align-items: center;
    gap: 8px;
    flex-wrap: wrap;
  }
  .mode button {
    background: transparent;
    color: var(--fg);
    border: 1px solid #2a2f42;
    border-radius: 8px;
    padding: 7px 16px;
    font-size: 13px;
    font-weight: 600;
    cursor: pointer;
    transition: background 0.15s, border-color 0.15s;
  }
  .mode button:hover { border-color: #3a405a; }
  .mode button.active {
    background: var(--accent);
    color: #0b0d12;
    border-color: var(--accent);
  }
  .mode-hint {
    margin-left: auto;
    color: var(--muted);
    font-size: 12px;
  }
  .mode-hint code {
    color: var(--accent);
    font-family: ui-monospace, "SF Mono", Menlo, monospace;
  }

  .presets {
    display: flex;
    align-items: center;
    gap: 8px;
    flex-wrap: wrap;
  }
  .preset {
    background: #1a1f2e;
    color: var(--fg);
    border: 1px solid #262c3e;
    border-radius: 999px;
    padding: 5px 12px;
    font-size: 12px;
    cursor: pointer;
    transition: background 0.15s, border-color 0.15s;
  }
  .preset:hover { background: #232a3d; border-color: #3a405a; }

  .chips {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    align-items: center;
    min-height: 38px;
  }
  .chip {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    background: #1a1f2e;
    color: var(--fg);
    border: 1px solid #2a2f42;
    border-radius: 999px;
    padding: 5px 6px 5px 14px;
    font-family: ui-monospace, "SF Mono", Menlo, monospace;
    font-size: 13px;
  }
  .chip .x {
    background: transparent;
    border: none;
    color: var(--muted);
    cursor: pointer;
    width: 22px;
    height: 22px;
    border-radius: 50%;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    font-size: 16px;
    line-height: 1;
  }
  .chip .x:hover { background: #2a2f42; color: var(--fg); }

  .chips input {
    flex: 1;
    min-width: 160px;
    background: transparent;
    border: none;
    color: var(--fg);
    font-size: 14px;
    font-family: ui-monospace, "SF Mono", Menlo, monospace;
    padding: 6px 4px;
    outline: none;
  }

  .results-head {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
  .loading {
    color: var(--muted);
    font-size: 12px;
    font-style: italic;
  }

  .row {
    display: grid;
    grid-template-columns: 140px 1fr 1.2fr;
    gap: 16px;
    align-items: center;
    padding: 8px 0;
    border-top: 1px solid #1f2433;
  }
  .row.head {
    color: var(--muted);
    font-size: 11px;
    text-transform: uppercase;
    letter-spacing: 0.06em;
    border-top: none;
    padding-bottom: 4px;
  }
  .row.head small { text-transform: none; letter-spacing: 0; color: var(--muted); }

  .num {
    font-family: ui-monospace, "SF Mono", Menlo, monospace;
    color: var(--accent);
    font-size: 14px;
  }
  .equals { display: flex; align-items: baseline; gap: 8px; }
  .scaled { font-weight: 700; font-size: 15px; }
  .unit { color: var(--muted); font-size: 14px; }

  .bar-track {
    display: block;
    width: 100%;
    height: 8px;
    background: #1a1f2e;
    border-radius: 4px;
    overflow: hidden;
  }
  .bar {
    display: block;
    height: 100%;
    background: linear-gradient(90deg, var(--accent), #b07cff);
    border-radius: 4px;
    transition: width 0.25s cubic-bezier(0.2, 0.8, 0.2, 1);
  }

  .empty, .error {
    color: var(--muted);
    font-size: 13px;
    margin: 4px 0 0;
  }
  .error { color: #ff7eb6; }
  .error small { color: var(--muted); }
</style>
