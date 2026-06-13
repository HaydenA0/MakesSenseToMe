<script>
  import { flip } from 'svelte/animate';
  import { cubicOut } from 'svelte/easing';
  import { fade, scale } from 'svelte/transition';

  let mode = $state('time');
  let scales = $state([380_000_000_000, 1_200_000_000, 1_200_000, 5_000, 190, 112]);
  let selectedIndex = $state(2);
  let newValue = $state('');
  let addingCard = $state(false);

  let trackEl = $state(null);
  let dragIndex = $state(null);
  let dragOverIndex = $state(null);
  let pointerStart = { x: 0, y: 0 };
  let didDrag = $state(false);

  let apiResults = $state(null);
  let loading = $state(false);
  let lastReqId = 0;
  let fetchError = $state(null);
  let copied = $state(false);

  let editingIndex = $state(null);
  let editValue = $state('');
  let parseError = $state(null);

  let contextMenu = $state(null);

  const DEFAULTS = {
    time: [380_000_000_000, 1_200_000_000, 1_200_000, 5_000, 190, 112],
    distance: [1, 1_000, 1_000_000, 100_000_000, 149_600_000_000_000, 9.461e18],
  };
  const BASE_UNIT = {
    time: 'second',
    distance: 'millimeter',
  };
  const UNIT_LABELS = {
    time: {
      second: 'second', minute: 'minute', hour: 'hour', day: 'day',
      week: 'week', month: 'month', year: 'year', decade: 'decade',
      century: 'century', millennium: 'millennium',
      myriad: 'myriad', ice_age_cycle: 'ice age cycle',
      million_years: 'million years', billion_years: 'billion years',
      age_of_earth: 'age of Earth', age_of_universe: 'age of universe',
      geological_era: 'geological era',
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

  const SUFFIX_MAP = {
    k: 1e3, thousand: 1e3, thousands: 1e3,
    m: 1e6, million: 1e6, millions: 1e6, mn: 1e6,
    b: 1e9, billion: 1e9, billions: 1e9, bn: 1e9,
    t: 1e12, trillion: 1e12, trillions: 1e12, tn: 1e12,
  };

  $effect(() => {
    const m = mode;
    const vs = scales;
    if (vs.length === 0) {
      apiResults = null;
      loading = false;
      return;
    }
    const id = ++lastReqId;
    loading = true;
    const timer = setTimeout(async () => {
      try {
        const res = await fetch(`/convert/${m}`, {
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

  const baseValue = $derived(scales[0] ?? 0);
  const selectedValue = $derived(selectedIndex !== null ? scales[selectedIndex] : null);
  const selectedResult = $derived(selectedIndex !== null && apiResults ? apiResults[selectedIndex] : null);
  const multiplier = $derived(
    selectedValue !== null && baseValue > 0 ? selectedValue / baseValue : 0
  );

  const indicatorLeft = $derived.by(() => {
    if (selectedIndex === null || scales.length === 0) return 0;
    return ((selectedIndex + 0.5) / scales.length) * 100;
  });

  function trimNum(v) {
    return v.toFixed(1).replace(/\.0$/, '');
  }

  function formatShort(v) {
    if (!Number.isFinite(v)) return String(v);
    const abs = Math.abs(v);
    if (abs === 0) return '0';
    if (abs >= 1e12) return trimNum(v / 1e12) + ' T';
    if (abs >= 1e9) return trimNum(v / 1e9) + ' B';
    if (abs >= 1e6) return trimNum(v / 1e6) + ' M';
    if (abs >= 1e3) return trimNum(v / 1e3) + ' K';
    if (abs >= 1) return v.toFixed(0);
    if (abs >= 0.01) return v.toFixed(2);
    if (abs >= 0.001) return trimNum(v * 1000) + ' m';
    return v.toExponential(1);
  }

  function formatNumber(v) {
    if (!Number.isFinite(v)) return String(v);
    const abs = Math.abs(v);
    if (abs === 0) return '0';
    if (abs >= 1e6 || abs < 1e-3) return v.toExponential(2);
    if (abs >= 100) return v.toLocaleString('en-US', { maximumFractionDigits: 0 });
    if (abs >= 1) return v.toLocaleString('en-US', { maximumFractionDigits: 2 });
    return v.toFixed(3);
  }

  function formatRaw(v) {
    if (!Number.isFinite(v)) return String(v);
    if (Math.abs(v) >= 1e16 || (v !== 0 && Math.abs(v) < 1e-3)) return v.toExponential(2);
    return v.toLocaleString('en-US', { maximumFractionDigits: 6 });
  }

  function formatLong(v) {
    if (!Number.isFinite(v)) return String(v);
    const abs = Math.abs(v);
    if (abs >= 1e12) return trimNum(v / 1e12) + ' Trillion';
    if (abs >= 1e9) return trimNum(v / 1e9) + ' Billion';
    if (abs >= 1e6) return trimNum(v / 1e6) + ' Million';
    if (abs >= 1e3) return trimNum(v / 1e3) + ' thousand';
    if (abs >= 1) return v.toLocaleString();
    return v.toString();
  }

  function parseNaturalLanguage(input) {
    const s = String(input ?? '').trim().toLowerCase().replace(/,/g, '');
    if (!s) return null;
    const match = s.match(/^(-?\d+(?:\.\d+)?(?:e[+-]?\d+)?)\s*([a-z]+)?$/);
    if (!match) return null;
    const num = parseFloat(match[1]);
    if (!Number.isFinite(num) || num <= 0) return null;
    const suffix = match[2] || '';
    if (suffix && !(suffix in SUFFIX_MAP)) return null;
    return num * (SUFFIX_MAP[suffix] || 1);
  }

  function unitLabel(key) {
    return UNIT_LABELS[mode]?.[key] ?? key.replace(/_/g, ' ');
  }

  function pluralize(n, word) {
    const v = parseFloat(n);
    return Number.isFinite(v) && v === 1 ? word : word + 's';
  }

  function onCardPointerDown(e, i) {
    if (e.button !== 0) return;
    pointerStart = { x: e.clientX, y: e.clientY };
    didDrag = false;
    dragIndex = i;
  }

  function onPointerMove(e) {
    if (dragIndex === null) return;
    const dx = e.clientX - pointerStart.x;
    if (!didDrag && Math.abs(dx) > 5) {
      didDrag = true;
    }
    if (didDrag) {
      const el = document.elementFromPoint(e.clientX, e.clientY);
      const card = el?.closest('.card');
      if (card) {
        const idx = parseInt(card.dataset.index, 10);
        dragOverIndex = !isNaN(idx) && idx !== dragIndex ? idx : null;
      } else {
        dragOverIndex = null;
      }
    }
  }

  function onPointerUp() {
    if (dragIndex === null) return;
    if (didDrag && dragOverIndex !== null) {
      const from = dragIndex;
      const to = dragOverIndex;
      const next = [...scales];
      const [moved] = next.splice(from, 1);
      next.splice(to, 0, moved);
      scales = next;
      if (selectedIndex !== null) {
        if (selectedIndex === from) selectedIndex = to;
        else if (from < selectedIndex && to >= selectedIndex) selectedIndex--;
        else if (from > selectedIndex && to <= selectedIndex) selectedIndex++;
      }
    } else if (!didDrag) {
      selectedIndex = selectedIndex === dragIndex ? null : dragIndex;
    }
    dragIndex = null;
    dragOverIndex = null;
    didDrag = false;
  }

  function onCardDblClick(e, i) {
    e.preventDefault();
    e.stopPropagation();
    if (didDrag) return;
    openEdit(i);
  }

  function onCardContextMenu(e, i) {
    e.preventDefault();
    e.stopPropagation();
    if (didDrag) return;
    selectedIndex = i;
    const menuW = 200;
    const menuH = 48;
    contextMenu = {
      x: Math.min(e.clientX, window.innerWidth - menuW - 8),
      y: Math.min(e.clientY, window.innerHeight - menuH - 8),
      index: i,
    };
  }

  function closeContextMenu() {
    contextMenu = null;
  }

  function removeCard(i) {
    if (scales.length <= 1) return;
    const next = scales.filter((_, idx) => idx !== i);
    scales = next;
    if (selectedIndex === i) selectedIndex = null;
    else if (selectedIndex !== null && i < selectedIndex) selectedIndex--;
    closeContextMenu();
  }

  function openEdit(i) {
    editingIndex = i;
    editValue = formatShort(scales[i]);
    parseError = null;
    queueMicrotask(() => {
      const el = document.querySelector('.modal input');
      if (el) {
        el.focus();
        el.select();
      }
    });
  }

  function confirmEdit() {
    if (editingIndex === null) return;
    const v = parseNaturalLanguage(editValue);
    if (v === null) {
      parseError = 'Could not parse. Try "3 B", "3 billion", "3.2 billion", or "3001".';
      return;
    }
    const next = [...scales];
    next[editingIndex] = v;
    scales = next;
    editingIndex = null;
    editValue = '';
    parseError = null;
  }

  function cancelEdit() {
    editingIndex = null;
    editValue = '';
    parseError = null;
  }

  function onEditKey(e) {
    if (e.key === 'Enter') {
      e.preventDefault();
      confirmEdit();
    } else if (e.key === 'Escape') {
      e.preventDefault();
      cancelEdit();
    }
  }

  function onAdd() {
    addingCard = true;
    newValue = '';
    queueMicrotask(() => {
      const el = document.querySelector('.card.adding input');
      el?.focus();
    });
  }

  function onAddSubmit() {
    const v = parseNaturalLanguage(newValue);
    if (v === null || v <= 0) {
      addingCard = false;
      newValue = '';
      return;
    }
    scales = [...scales, v];
    selectedIndex = scales.length - 1;
    addingCard = false;
    newValue = '';
  }

  function onAddKey(e) {
    if (e.key === 'Enter') {
      e.preventDefault();
      onAddSubmit();
    } else if (e.key === 'Escape') {
      e.preventDefault();
      addingCard = false;
      newValue = '';
    }
  }

  function switchMode(m) {
    if (m === mode) return;
    mode = m;
    scales = [...DEFAULTS[m]];
    selectedIndex = null;
    apiResults = null;
    closeContextMenu();
    cancelEdit();
    if (addingCard) { addingCard = false; newValue = ''; }
  }

  async function copyResult() {
    if (selectedIndex === null || !selectedResult) return;
    const v = selectedValue;
    const r = selectedResult;
    const baseU = unitLabel(BASE_UNIT[mode]);
    const convU = unitLabel(r.UnitName);
    const text =
      `${formatLong(v)} ${pluralize(v, baseU)} is equivalent to ` +
      `${formatNumber(r.ScaledValue)} ${pluralize(r.ScaledValue, convU)}, ` +
      `which is ${formatNumber(multiplier)}× the base of ` +
      `${formatLong(baseValue)} ${baseU}.`;
    try {
      await navigator.clipboard.writeText(text);
      copied = true;
      setTimeout(() => { copied = false; }, 1500);
    } catch {
      /* ignore */
    }
  }

  function onWindowKeydown(e) {
    if (e.key !== 'Escape') return;
    if (contextMenu) closeContextMenu();
    else if (editingIndex !== null) cancelEdit();
    else if (addingCard) { addingCard = false; newValue = ''; }
  }

  function onWindowPointerDown(e) {
    if (contextMenu && !e.target.closest('.context-menu')) {
      closeContextMenu();
    }
  }
</script>

<svelte:window
  onpointermove={onPointerMove}
  onpointerup={onPointerUp}
  onpointerdown={onWindowPointerDown}
  onkeydown={onWindowKeydown}
/>

<div class="container">
  <header>
    <h1>Make sense 2 me</h1>
    <div class="mode-toggle" role="tablist" aria-label="Mode">
      <button
        class:active={mode === 'time'}
        onclick={() => switchMode('time')}
        role="tab"
        aria-selected={mode === 'time'}
      >Time</button>
      <button
        class:active={mode === 'distance'}
        onclick={() => switchMode('distance')}
        role="tab"
        aria-selected={mode === 'distance'}
      >Distance</button>
    </div>
  </header>

  <div class="cards-row">
    <div class="cards-track" bind:this={trackEl}>
      {#each scales as v, i (i + ':' + v)}
        <div
          class="card"
          class:selected={selectedIndex === i}
          class:dragging={dragIndex === i && didDrag}
          class:drop-target={dragOverIndex === i && dragIndex !== i}
          class:base={i === 0}
          data-index={i}
          animate:flip={{ duration: 250 }}
          onpointerdown={(e) => onCardPointerDown(e, i)}
          ondblclick={(e) => onCardDblClick(e, i)}
          oncontextmenu={(e) => onCardContextMenu(e, i)}
          role="button"
          tabindex="0"
          aria-pressed={selectedIndex === i}
        >
          {#if i === 0}
            <span class="base-badge">BASE</span>
          {/if}
          <span class="value">{formatShort(v)}</span>
          <span class="unit">{BASE_UNIT[mode]}</span>
        </div>
      {/each}
      {#if addingCard}
        <div class="card adding" in:scale={{ duration: 150, start: 0.7, easing: cubicOut }}>
          <input
            type="text"
            inputmode="decimal"
            bind:value={newValue}
            onkeydown={onAddKey}
            onblur={onAddSubmit}
            placeholder="…"
          />
          <span class="unit">{BASE_UNIT[mode]}</span>
        </div>
      {/if}
    </div>

    <button class="add-card" onclick={onAdd} aria-label="Add scale" title="Add scale">
      <svg viewBox="0 0 24 24" width="20" height="20" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="12" y1="5" x2="12" y2="19"></line><line x1="5" y1="12" x2="19" y2="12"></line></svg>
    </button>
  </div>

  <div class="indicator-row">
    <div class="indicator-track" aria-hidden="true">
      <div
        class="indicator-marker"
        style="left: {indicatorLeft}%; opacity: {selectedIndex !== null ? 1 : 0}"
      ></div>
    </div>
    <div class="indicator-spacer"></div>
  </div>

  {#if selectedIndex !== null && selectedResult}
    <section class="detail-panel" in:fade={{ duration: 180 }}>
      <div class="detail-left">
        <div class="icon-circle">
          {#if mode === 'time'}
            <svg viewBox="0 0 24 24" width="30" height="30" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
              <circle cx="12" cy="12" r="9"></circle>
              <polyline points="12 7 12 12 16 14"></polyline>
            </svg>
          {:else}
            <svg viewBox="0 0 24 24" width="30" height="30" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
              <path d="M3 17l6-6 4 4 8-8"></path>
              <polyline points="14 7 21 7 21 14"></polyline>
            </svg>
          {/if}
        </div>
        <div class="detail-text">
          <p class="lead">
            {formatLong(selectedValue)} {pluralize(selectedValue, unitLabel(BASE_UNIT[mode]))} is equivalent to
          </p>
          <p class="big">
            {formatNumber(selectedResult.ScaledValue)}
            <span class="big-unit">{pluralize(selectedResult.ScaledValue, unitLabel(selectedResult.UnitName))}</span>
          </p>
          <p class="muted">
            which is <strong class="hl">{formatNumber(multiplier)}×</strong>
            {multiplier >= 1 ? 'larger' : 'smaller'} than
            {formatLong(baseValue)} {pluralize(baseValue, unitLabel(BASE_UNIT[mode]))}.
          </p>
          <p class="muted-small">
            <strong>{formatShort(baseValue)}</strong> being the base : <code>{BASE_UNIT[mode]}</code>
          </p>
        </div>
      </div>

      <div class="detail-right">
        <div class="info-row">
          <span class="info-label">
            <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><path d="M3 12a9 9 0 1 0 3-6.7"></path><polyline points="3 4 3 10 9 10"></polyline></svg>
            Original value
          </span>
          <span class="info-value">
            {formatRaw(selectedValue)}
            <span class="info-unit">{unitLabel(BASE_UNIT[mode])}</span>
          </span>
        </div>
        <div class="info-row">
          <span class="info-label">
            <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><path d="M12 3v18"></path><path d="M5 10l7-7 7 7"></path></svg>
            Converted value
          </span>
          <span class="info-value">
            {formatNumber(selectedResult.ScaledValue)}
            <span class="info-unit">{unitLabel(selectedResult.UnitName)}</span>
          </span>
        </div>
        <div class="info-row">
          <span class="info-label">
            <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><polyline points="3 17 9 11 13 15 21 7"></polyline><polyline points="14 7 21 7 21 14"></polyline></svg>
            Multiplier
          </span>
          <span class="info-value">{formatNumber(multiplier)}×</span>
        </div>
        <div class="info-row">
          <span class="info-label">
            <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><line x1="4" y1="9" x2="20" y2="9"></line><line x1="4" y1="15" x2="20" y2="15"></line><line x1="10" y1="3" x2="8" y2="21"></line><line x1="16" y1="3" x2="14" y2="21"></line></svg>
            Base unit
          </span>
          <span class="info-value">
            {formatShort(baseValue)}
            <span class="info-unit">({unitLabel(BASE_UNIT[mode])})</span>
          </span>
        </div>
      </div>
    </section>
  {:else if fetchError}
    <section class="detail-panel error">
      <p>API error: {fetchError}</p>
      <p class="muted-small">Is the Go server running on :8080?</p>
    </section>
  {/if}

  <footer>
    <div class="tip">
      <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" class="tip-icon">
        <path d="M9 18h6"></path>
        <path d="M10 22h4"></path>
        <path d="M12 2a7 7 0 0 0-4 12.7c.6.5 1 1.2 1 2v.3h6v-.3c0-.8.4-1.5 1-2A7 7 0 0 0 12 2z"></path>
      </svg>
      <span>
        Drag a card to reorder. Click to select, double-click to edit, right-click to remove.
      </span>
    </div>
    <button class="copy-btn" onclick={copyResult} disabled={selectedIndex === null || !selectedResult}>
      <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
        <rect x="9" y="9" width="13" height="13" rx="2"></rect>
        <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path>
      </svg>
      {copied ? 'Copied!' : 'Copy result'}
    </button>
  </footer>
</div>

{#if contextMenu}
  <div
    class="context-menu"
    style="left: {contextMenu.x}px; top: {contextMenu.y}px"
    role="menu"
  >
    <button role="menuitem" onclick={() => removeCard(contextMenu.index)}>
      <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <polyline points="3 6 5 6 21 6"></polyline>
        <path d="M19 6l-1 14a2 2 0 0 1-2 2H8a2 2 0 0 1-2-2L5 6"></path>
        <path d="M10 11v6"></path>
        <path d="M14 11v6"></path>
        <path d="M9 6V4a2 2 0 0 1 2-2h2a2 2 0 0 1 2 2v2"></path>
      </svg>
      Remove card
    </button>
  </div>
{/if}

{#if editingIndex !== null}
  <div
    class="modal-backdrop"
    onclick={cancelEdit}
    onkeydown={onEditKey}
    oncontextmenu={(e) => e.preventDefault()}
    role="dialog"
    aria-modal="true"
    tabindex="-1"
  >
    <!-- svelte-ignore a11y_no_noninteractive_element_interactions -->
    <div
      class="modal"
      role="document"
      onclick={(e) => e.stopPropagation()}
      onkeydown={(e) => e.stopPropagation()}
      oncontextmenu={(e) => e.stopPropagation()}
    >
      <h3>Edit scale</h3>
      <p class="hint">
        Examples: <code>3 B</code> · <code>3 billion</code> · <code>3.2 billion</code> · <code>3001</code>
      </p>
      <input
        type="text"
        inputmode="decimal"
        bind:value={editValue}
        onkeydown={onEditKey}
        placeholder="Enter a value…"
        spellcheck="false"
        autocomplete="off"
      />
      {#if parseError}
        <p class="error">{parseError}</p>
      {/if}
      <div class="modal-actions">
        <button class="secondary" onclick={cancelEdit}>Cancel</button>
        <button class="primary" onclick={confirmEdit}>Save</button>
      </div>
    </div>
  </div>
{/if}

<style>
  .container {
    max-width: 1120px;
    width: 100%;
    margin: 0 auto;
    padding: 28px 24px 56px;
    display: flex;
    flex-direction: column;
    gap: 14px;
    user-select: none;
  }

  header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 4px 6px;
  }
  header h1 {
    margin: 0;
    font-size: 30px;
    font-weight: 700;
    letter-spacing: -0.02em;
  }

  .mode-toggle {
    display: flex;
    background: var(--panel);
    border: 1px solid var(--grid);
    border-radius: 10px;
    padding: 4px;
    gap: 2px;
  }
  .mode-toggle button {
    background: transparent;
    color: var(--muted);
    border: none;
    border-radius: 7px;
    padding: 7px 20px;
    font-size: 13px;
    font-weight: 600;
    cursor: pointer;
    transition: background 0.15s, color 0.15s;
  }
  .mode-toggle button:hover { color: var(--fg); }
  .mode-toggle button.active {
    background: rgba(124, 109, 255, 0.18);
    color: var(--accent-soft);
  }

  .cards-row {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 0 4px;
  }

  .cards-track {
    display: flex;
    gap: 12px;
    flex: 1;
    overflow-x: auto;
    overflow-y: visible;
    scroll-behavior: smooth;
    padding: 6px 4px 14px;
    scrollbar-width: none;
  }
  .cards-track::-webkit-scrollbar { display: none; }

  .card {
    flex: 0 0 150px;
    height: 80px;
    background: var(--panel);
    border: 1px solid var(--grid);
    border-radius: 12px;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    cursor: grab;
    position: relative;
    transition: border-color 0.15s, background 0.15s, transform 0.12s, box-shadow 0.15s;
    touch-action: none;
  }
  .card:hover {
    border-color: #2f364d;
    background: var(--panel-2);
  }
  .card.selected {
    border-color: var(--accent);
    background: linear-gradient(180deg, rgba(124, 109, 255, 0.10), rgba(124, 109, 255, 0.03));
    box-shadow: 0 0 0 1px var(--accent), 0 4px 18px -8px var(--accent-glow);
  }
  .card.dragging {
    opacity: 0.45;
    cursor: grabbing;
    transform: scale(0.97);
  }
  .card.drop-target {
    border-color: var(--accent);
    border-style: dashed;
    background: rgba(124, 109, 255, 0.06);
  }
  .card.base .value { color: var(--accent-soft); }
  .base-badge {
    position: absolute;
    top: 6px;
    left: 8px;
    font-size: 9px;
    font-weight: 700;
    letter-spacing: 0.08em;
    color: var(--accent-soft);
    background: rgba(124, 109, 255, 0.14);
    padding: 2px 6px;
    border-radius: 4px;
  }
  .card .value {
    font-size: 24px;
    font-weight: 700;
    letter-spacing: -0.01em;
    line-height: 1.1;
  }
  .card .unit {
    color: var(--muted);
    font-size: 12px;
    margin-top: 4px;
  }

  .card.adding {
    border-color: var(--accent);
    border-style: dashed;
    background: rgba(124, 109, 255, 0.06);
    cursor: default;
  }
  .card.adding input {
    font-size: 22px;
    font-weight: 700;
    background: transparent;
    border: none;
    color: var(--fg);
    text-align: center;
    width: 100%;
    outline: none;
    padding: 0 8px;
  }
  .card.adding input::placeholder { color: var(--muted-2); }

  .add-card {
    flex: 0 0 44px;
    width: 44px;
    height: 80px;
    background: var(--panel);
    border: 1px dashed #2a3148;
    border-radius: 12px;
    color: var(--muted);
    cursor: pointer;
    transition: all 0.15s;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 0;
  }
  .add-card:hover {
    border-color: var(--accent);
    color: var(--accent-soft);
    background: rgba(124, 109, 255, 0.05);
  }

  .indicator-row {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 0 4px;
    margin-top: -8px;
  }
  .indicator-track {
    flex: 1;
    position: relative;
    height: 8px;
    display: flex;
    align-items: center;
  }
  .indicator-track::before {
    content: '';
    position: absolute;
    left: 0;
    right: 0;
    top: 50%;
    transform: translateY(-50%);
    height: 2px;
    background: var(--grid);
    border-radius: 1px;
  }
  .indicator-marker {
    position: absolute;
    top: 50%;
    transform: translate(-50%, -50%);
    width: 56px;
    height: 4px;
    background: var(--accent);
    border-radius: 2px;
    box-shadow: 0 0 12px var(--accent-glow);
    transition: opacity 0.2s, left 0.28s cubic-bezier(0.2, 0.8, 0.2, 1);
  }
  .indicator-spacer {
    flex: 0 0 44px;
    width: 44px;
  }

  .detail-panel {
    display: flex;
    background: var(--panel);
    border: 1px solid var(--grid);
    border-radius: 16px;
    padding: 28px 32px;
    gap: 36px;
    min-height: 220px;
  }
  .detail-panel.error {
    color: #ff8ab3;
    flex-direction: column;
    align-items: flex-start;
  }
  .detail-left {
    display: flex;
    align-items: center;
    gap: 24px;
    flex: 1.5;
  }
  .icon-circle {
    flex: 0 0 72px;
    width: 72px;
    height: 72px;
    background: linear-gradient(135deg, #8a7dff, #5a4cd6);
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #fff;
    box-shadow: 0 0 24px var(--accent-glow);
  }
  .detail-text {
    display: flex;
    flex-direction: column;
    gap: 4px;
    min-width: 0;
  }
  .lead {
    margin: 0;
    color: var(--muted);
    font-size: 14px;
  }
  .big {
    margin: 0;
    color: var(--accent-soft);
    font-size: 32px;
    font-weight: 700;
    letter-spacing: -0.02em;
    line-height: 1.15;
  }
  .big-unit {
    color: var(--fg);
    font-weight: 600;
    font-size: 22px;
    margin-left: 4px;
  }
  .muted {
    margin: 4px 0 0;
    color: var(--muted);
    font-size: 13px;
  }
  .muted .hl {
    color: var(--fg);
    font-weight: 700;
  }
  .muted-small {
    margin: 6px 0 0;
    color: var(--muted);
    font-size: 12px;
  }
  .muted-small strong { color: var(--fg); }
  .muted-small code {
    color: var(--accent-soft);
    font-family: ui-monospace, "SF Mono", Menlo, monospace;
  }
  .detail-right {
    flex: 1;
    display: flex;
    flex-direction: column;
    border-left: 1px solid var(--grid);
    padding-left: 32px;
  }
  .info-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 12px 0;
    border-bottom: 1px solid var(--grid);
    font-size: 13px;
    gap: 16px;
  }
  .info-row:last-child { border-bottom: none; }
  .info-label {
    display: inline-flex;
    align-items: center;
    gap: 10px;
    color: var(--muted);
    flex-shrink: 0;
  }
  .info-label svg { color: var(--accent-soft); }
  .info-value {
    color: var(--fg);
    font-weight: 600;
    text-align: right;
  }
  .info-unit {
    color: var(--muted);
    font-weight: 400;
    margin-left: 4px;
  }

  footer {
    display: flex;
    justify-content: space-between;
    align-items: center;
    background: var(--panel);
    border: 1px solid var(--grid);
    border-radius: 12px;
    padding: 14px 20px;
    margin-top: 4px;
    gap: 16px;
  }
  .tip {
    display: flex;
    align-items: center;
    gap: 10px;
    color: var(--muted);
    font-size: 13px;
  }
  .tip-icon { color: var(--accent-soft); flex-shrink: 0; }
  .copy-btn {
    display: inline-flex;
    align-items: center;
    gap: 8px;
    background: transparent;
    color: var(--accent-soft);
    border: 1px solid rgba(124, 109, 255, 0.4);
    border-radius: 8px;
    padding: 7px 14px;
    font-size: 13px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.15s;
    flex-shrink: 0;
  }
  .copy-btn:hover {
    background: rgba(124, 109, 255, 0.1);
    border-color: var(--accent);
    color: #fff;
  }
  .copy-btn:disabled {
    opacity: 0.4;
    cursor: not-allowed;
  }

  .context-menu {
    position: fixed;
    z-index: 100;
    background: var(--panel-2);
    border: 1px solid var(--grid);
    border-radius: 8px;
    padding: 4px;
    min-width: 180px;
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.5);
  }
  .context-menu button {
    display: flex;
    align-items: center;
    gap: 10px;
    width: 100%;
    text-align: left;
    background: transparent;
    border: none;
    color: var(--fg);
    padding: 8px 12px;
    border-radius: 6px;
    font-size: 13px;
    cursor: pointer;
    font-family: inherit;
  }
  .context-menu button:hover {
    background: rgba(124, 109, 255, 0.15);
    color: var(--accent-soft);
  }
  .context-menu button svg {
    color: var(--muted);
    flex-shrink: 0;
  }
  .context-menu button:hover svg {
    color: #ff8ab3;
  }

  .modal-backdrop {
    position: fixed;
    inset: 0;
    z-index: 200;
    background: rgba(0, 0, 0, 0.55);
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 16px;
  }
  .modal {
    background: var(--panel);
    border: 1px solid var(--grid);
    border-radius: 12px;
    padding: 22px 24px;
    min-width: 360px;
    max-width: 480px;
    width: 100%;
    box-shadow: 0 16px 48px rgba(0, 0, 0, 0.5);
  }
  .modal h3 {
    margin: 0 0 6px;
    font-size: 16px;
    font-weight: 600;
  }
  .modal .hint {
    margin: 0 0 14px;
    color: var(--muted);
    font-size: 12px;
  }
  .modal .hint code {
    background: var(--panel-2);
    border: 1px solid var(--grid);
    border-radius: 4px;
    padding: 1px 6px;
    font-family: ui-monospace, "SF Mono", Menlo, monospace;
    font-size: 11px;
    color: var(--accent-soft);
  }
  .modal input {
    width: 100%;
    background: var(--bg);
    border: 1px solid var(--grid);
    border-radius: 8px;
    color: var(--fg);
    padding: 10px 12px;
    font-size: 15px;
    outline: none;
    font-family: ui-monospace, "SF Mono", Menlo, monospace;
    transition: border-color 0.15s;
  }
  .modal input:focus { border-color: var(--accent); }
  .modal .error {
    color: #ff8ab3;
    font-size: 12px;
    margin: 8px 0 0;
  }
  .modal-actions {
    display: flex;
    justify-content: flex-end;
    gap: 8px;
    margin-top: 18px;
  }
  .modal-actions button {
    padding: 8px 16px;
    border-radius: 8px;
    font-size: 13px;
    font-weight: 600;
    cursor: pointer;
    border: 1px solid var(--grid);
    font-family: inherit;
  }
  .modal-actions .secondary {
    background: transparent;
    color: var(--muted);
  }
  .modal-actions .secondary:hover {
    color: var(--fg);
    border-color: #2a3148;
  }
  .modal-actions .primary {
    background: var(--accent);
    color: #fff;
    border-color: var(--accent);
  }
  .modal-actions .primary:hover {
    background: #8a7dff;
    border-color: #8a7dff;
  }

  @media (max-width: 720px) {
    header h1 { font-size: 22px; }
    .mode-toggle button { padding: 6px 12px; font-size: 12px; }
    .card { flex: 0 0 130px; }
    .card .value { font-size: 20px; }
    .detail-panel { flex-direction: column; padding: 20px; gap: 20px; }
    .detail-right { border-left: none; border-top: 1px solid var(--grid); padding-left: 0; padding-top: 12px; }
    .big { font-size: 26px; }
    .icon-circle { flex: 0 0 56px; width: 56px; height: 56px; }
    .modal { min-width: 0; }
  }
</style>
