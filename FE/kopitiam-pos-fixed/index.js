/**
 * ============================================================
 * KopiTiam POS — Login Page (index.js) [FIXED]
 * All event listeners attached after DOM ready
 * Functions exposed to window for HTML onclick compatibility
 * ============================================================
 */

// ── Configuration ───────────────────────────────────────────
const CONFIG = {
  API_BASE_URL: '',
  TOKEN_KEY: 'kopitiam_token',
  REFRESH_KEY: 'kopitiam_refresh_token',
  USER_KEY: 'kopitiam_user',
};

// Try to load from config.json file if served via HTTP
async function loadConfig() {
  try {
    const res = await fetch('./config.json');
    if (res.ok) {
      const data = await res.json();
      if (data && data.API_BASE_URL !== undefined) {
        CONFIG.API_BASE_URL = data.API_BASE_URL;
      }
    }
  } catch (e) {
    // Fallback to default
  }
}

// ── State ───────────────────────────────────────────────────
let currentTab = 'manager';
let pinValue = '';

// ── DOM References ──────────────────────────────────────────
const $ = (id) => document.getElementById(id);
const $$ = (sel) => document.querySelectorAll(sel);

// ── Tab Switching ───────────────────────────────────────────
function switchTab(tab) {
  currentTab = tab;

  const managerTab = $('tab-manager');
  const staffTab = $('tab-staff');

  if (tab === 'manager') {
    managerTab.className = 'flex-1 flex items-center justify-center gap-2 py-3 px-4 rounded-xl text-sm font-semibold transition-all duration-200 tab-active';
    staffTab.className = 'flex-1 flex items-center justify-center gap-2 py-3 px-4 rounded-xl text-sm font-semibold transition-all duration-200 tab-inactive';

    $('form-manager').classList.remove('hidden');
    $('form-staff').classList.add('hidden');
    $('manager-email').focus();
  } else {
    staffTab.className = 'flex-1 flex items-center justify-center gap-2 py-3 px-4 rounded-xl text-sm font-semibold transition-all duration-200 tab-active';
    managerTab.className = 'flex-1 flex items-center justify-center gap-2 py-3 px-4 rounded-xl text-sm font-semibold transition-all duration-200 tab-inactive';

    $('form-staff').classList.remove('hidden');
    $('form-staff').classList.add('slide-in');
    $('form-manager').classList.add('hidden');
    clearPin();
  }

  hideError();
  hideSuccess();
}

// ── PIN Handling ────────────────────────────────────────────
function appendPin(digit) {
  if (pinValue.length < 6) {
    pinValue += digit;
    updatePinDisplay();

    if (pinValue.length === 6) {
      setTimeout(() => handleStaffLogin(), 200);
    }
  }
}

function backspacePin() {
  pinValue = pinValue.slice(0, -1);
  updatePinDisplay();
}

function clearPin() {
  pinValue = '';
  updatePinDisplay();
}

function updatePinDisplay() {
  const dots = document.querySelectorAll('.pin-dot');
  dots.forEach((dot, index) => {
    dot.classList.toggle('filled', index < pinValue.length);
  });
}

// ── Message Handling ────────────────────────────────────────
function showError(message) {
  const container = $('error-message');
  const text = $('error-text');
  text.textContent = message;
  container.classList.remove('hidden');
  hideSuccess();

  const card = $('login-card');
  card.classList.remove('shake');
  void card.offsetWidth;
  card.classList.add('shake');
}

function hideError() {
  $('error-message').classList.add('hidden');
}

function showSuccess(message, subtext) {
  const container = $('success-message');
  const text = $('success-text');
  const sub = $('success-subtext');
  text.textContent = message;
  sub.textContent = subtext || '';
  container.classList.remove('hidden');
  hideError();

  const card = $('login-card');
  card.classList.add('success-pulse');
  setTimeout(() => card.classList.remove('success-pulse'), 600);
}

function hideSuccess() {
  $('success-message').classList.add('hidden');
}

// ── Loading State ────────────────────────────────────────────
function setLoading(isLoading, type) {
  const btnText = $(`${type}-btn-text`);
  const btnIcon = $(`${type}-btn-icon`);
  const btn = $(`${type}-submit`);

  if (isLoading) {
    btn.disabled = true;
    btnText.textContent = 'Loading...';
    btnIcon.innerHTML = `<svg class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path></svg>`;
    btn.classList.add('opacity-75', 'cursor-not-allowed');
  } else {
    btn.disabled = false;
    btnText.textContent = type === 'manager' ? 'Masuk sebagai Manager' : 'Masuk sebagai Staff';
    btnIcon.innerHTML = `<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7l5 5m0 0l-5 5m5-5H6"/></svg>`;
    btn.classList.remove('opacity-75', 'cursor-not-allowed');
  }
}

// ── Redirect with Fallback ──────────────────────────────────
function safeRedirect(url, fallbackMessage) {
  if (window.location.protocol === 'file:') {
    showSuccess(
      'Login berhasil!',
      `${fallbackMessage} File ini dibuka dari disk lokal (file://). Serve via Caddy untuk redirect otomatis.`
    );
    return;
  }

  fetch(url, { method: 'HEAD', cache: 'no-cache' })
    .then(response => {
      if (response.ok || response.status === 304) {
        window.location.href = url;
      } else {
        showSuccess(
          'Login berhasil!',
          `${fallbackMessage} Halaman tujuan (${url}) belum tersedia. Status: ${response.status}`
        );
      }
    })
    .catch(() => {
      window.location.href = url;
    });
}

// ── API Helpers ─────────────────────────────────────────────
async function apiPost(endpoint, body) {
  const res = await fetch(`${CONFIG.API_BASE_URL}${endpoint}`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Accept': 'application/json'
    },
    body: JSON.stringify(body)
  });
  const data = await res.json();
  if (!res.ok || !data.success) {
    throw new Error(data.message || data.error || 'Request failed');
  }
  return data;
}

// ── Manager Login ────────────────────────────────────────────
async function handleManagerLogin(event) {
  event.preventDefault();
  hideError();
  hideSuccess();

  const email = $('manager-email').value.trim();
  const password = $('manager-password').value;

  if (!email || !password) {
    showError('Email dan password wajib diisi.');
    return;
  }

  setLoading(true, 'manager');

  try {
    const data = await apiPost('/api/v1/auth/login', { email, password });

    if (data.data?.access_token) {
      localStorage.setItem(CONFIG.TOKEN_KEY, data.data.access_token);
      localStorage.setItem(CONFIG.REFRESH_KEY, data.data.refresh_token || '');
      localStorage.setItem(CONFIG.USER_KEY, JSON.stringify(data.data.user || {}));
    }

    safeRedirect('./analytics.html', 'Redirect ke halaman Analytics...');

  } catch (error) {
    console.error('Login error:', error);
    showError(error.message || 'Terjadi kesalahan jaringan. Pastikan backend berjalan.');
    setLoading(false, 'manager');
  }
}

// ── Staff Login (PIN) ───────────────────────────────────────
async function handleStaffLogin() {
  hideError();
  hideSuccess();

  if (pinValue.length !== 6) {
    showError('PIN harus 6 digit.');
    return;
  }

  setLoading(true, 'staff');

  try {
    const data = await apiPost('/api/v1/auth/login-pin', { pin: pinValue });

    if (data.data?.access_token) {
      localStorage.setItem(CONFIG.TOKEN_KEY, data.data.access_token);
      localStorage.setItem(CONFIG.REFRESH_KEY, data.data.refresh_token || '');
      localStorage.setItem(CONFIG.USER_KEY, JSON.stringify(data.data.user || {}));
    }

    const user = data.data?.user || {};
    const role = user.role?.toLowerCase() || '';

    if (role === 'cashier') {
      safeRedirect('./cashier.html', 'Redirect ke halaman Cashier...');
    } else if (role === 'kitchen') {
      safeRedirect('./kitchen.html', 'Redirect ke halaman Kitchen...');
    } else {
      safeRedirect('./analytics.html', 'Redirect ke halaman Analytics...');
    }

  } catch (error) {
    console.error('PIN login error:', error);
    showError(error.message || 'PIN tidak valid atau server tidak merespons.');
    setLoading(false, 'staff');
    clearPin();
  }
}

// ── Auto-redirect if already logged in ───────────────────────
function checkExistingSession() {
  const token = localStorage.getItem(CONFIG.TOKEN_KEY);
  const user = JSON.parse(localStorage.getItem(CONFIG.USER_KEY) || '{}');

  if (token && user.role) {
    const role = user.role.toLowerCase();
    if (role === 'manager') {
      safeRedirect('./analytics.html', 'Sudah login sebagai Manager.');
    } else if (role === 'cashier') {
      safeRedirect('./cashier.html', 'Sudah login sebagai Cashier.');
    } else if (role === 'kitchen') {
      safeRedirect('./kitchen.html', 'Sudah login sebagai Kitchen.');
    }
  }
}

// ── Event Delegation Setup ─────────────────────────────────
function setupEventListeners() {
  // Tab switching
  $('tab-manager').addEventListener('click', () => switchTab('manager'));
  $('tab-staff').addEventListener('click', () => switchTab('staff'));

  // Manager form
  $('manager-form').addEventListener('submit', handleManagerLogin);

  // Numpad buttons - event delegation
  const numpadContainer = document.querySelector('#form-staff .grid');
  if (numpadContainer) {
    numpadContainer.addEventListener('click', (e) => {
      const btn = e.target.closest('button[data-pin]');
      if (!btn) return;

      const action = btn.getAttribute('data-pin');
      if (action === 'clear') {
        clearPin();
      } else if (action === 'backspace') {
        backspacePin();
      } else if (/^[0-9]$/.test(action)) {
        appendPin(action);
      }
    });
  }

  // Staff submit button
  $('staff-submit').addEventListener('click', handleStaffLogin);

  // Keyboard support for PIN
  document.addEventListener('keydown', (e) => {
    if (currentTab !== 'staff') return;

    if (e.key >= '0' && e.key <= '9') {
      appendPin(e.key);
    } else if (e.key === 'Backspace') {
      backspacePin();
    } else if (e.key === 'Escape') {
      clearPin();
    } else if (e.key === 'Enter' && pinValue.length === 6) {
      handleStaffLogin();
    }
  });
}

// ── Initialization ────────────────────────────────────────────
document.addEventListener('DOMContentLoaded', async () => {
  await loadConfig();
  setupEventListeners();
  checkExistingSession();
  $('manager-email').focus();
});

// Expose functions to global scope for HTML onclick handlers (fallback)
window.switchTab = switchTab;
window.appendPin = appendPin;
window.backspacePin = backspacePin;
window.clearPin = clearPin;
window.handleManagerLogin = handleManagerLogin;
window.handleStaffLogin = handleStaffLogin;
