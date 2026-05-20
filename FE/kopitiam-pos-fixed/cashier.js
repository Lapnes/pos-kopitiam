/**
 * ============================================================
 * KopiTiam POS — Cashier Terminal (cashier.js)
 * Full-featured POS interface with menu, cart, and payment
 * ============================================================
 */

const CONFIG = {
  API_BASE_URL: '',
  TOKEN_KEY: 'kopitiam_token',
  USER_KEY: 'kopitiam_user',
};

let menuItems = [];
let currentOrder = [];
let selectedCategory = 'all';
let selectedPayment = 'cash';
let isProcessing = false;

const $ = (id) => document.getElementById(id);

// ── Auth ────────────────────────────────────────────────────
function checkAuth() {
  const token = localStorage.getItem(CONFIG.TOKEN_KEY);
  const user = JSON.parse(localStorage.getItem(CONFIG.USER_KEY) || '{}');
  if (!token) {
    window.location.href = './index.html';
    return false;
  }
  if (user.name) {
    const elName = $('user-name');
    const elRole = $('user-role');
    const elInitial = $('user-initial');
    if (elName) elName.textContent = user.name;
    if (elRole) elRole.textContent = user.role || 'Cashier';
    if (elInitial) elInitial.textContent = (user.name || 'C').charAt(0).toUpperCase();
    
    // Hide Analytics link if not manager
    if (user.role !== 'manager') {
      const analyticsLinks = document.querySelectorAll('a[href*="analytics.html"]');
      analyticsLinks.forEach(link => link.style.display = 'none');
    }
  }
  return token;
}

function logout() {
  localStorage.removeItem(CONFIG.TOKEN_KEY);
  localStorage.removeItem(CONFIG.USER_KEY);
  window.location.href = './index.html';
}

// ── API ─────────────────────────────────────────────────────
async function apiGet(endpoint) {
  const token = localStorage.getItem(CONFIG.TOKEN_KEY);
  const res = await fetch(`${CONFIG.API_BASE_URL}${endpoint}`, {
    headers: { 'Authorization': `Bearer ${token}`, 'Accept': 'application/json' }
  });
  if (res.status === 401) { logout(); throw new Error('Session expired'); }
  const data = await res.json();
  if (!res.ok || !data.success) throw new Error(data.message || 'Request failed');
  return data;
}

async function apiPost(endpoint, body) {
  const token = localStorage.getItem(CONFIG.TOKEN_KEY);
  const res = await fetch(`${CONFIG.API_BASE_URL}${endpoint}`, {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json',
      'Accept': 'application/json'
    },
    body: JSON.stringify(body)
  });
  if (res.status === 401) { logout(); throw new Error('Session expired'); }
  const data = await res.json();
  if (!res.ok || !data.success) throw new Error(data.message || 'Request failed');
  return data;
}

// ── Formatters ────────────────────────────────────────────────
function formatRupiah(num) {
  if (num === undefined || num === null || num === 'undefined' || num === 'null' || isNaN(num) || (typeof num === 'string' && (num.toLowerCase().includes('undefined') || num.toLowerCase().includes('null')))) {
    return 'Rp 0';
  }
  const parsed = parseFloat(num);
  if (isNaN(parsed)) return 'Rp 0';
  if (parsed > 0 && parsed < 10) {
    return 'Rp ' + parsed.toFixed(2);
  }
  return 'Rp ' + Math.round(parsed).toString().replace(/\B(?=(\d{3})+(?!\d))/g, '.');
}

function formatNumber(num) {
  if (!num || isNaN(num)) return '0';
  return num.toString().replace(/\B(?=(\d{3})+(?!\d))/g, '.');
}

let bestSellers = [];

// ── Menu ──────────────────────────────────────────────────────
async function loadMenu() {
  try {
    try {
      const bsRes = await apiGet('/api/v1/best-sellers?limit=10');
      bestSellers = bsRes.data || [];
    } catch (e) {
      console.warn('Failed to load best sellers:', e);
      bestSellers = [];
    }

    const data = await apiGet('/api/v1/menus?limit=100');
    menuItems = data.data || [];
    renderMenu();
  } catch (err) {
    showToast('Failed to load menu: ' + err.message, 'error');
    // Fallback demo data
    menuItems = [
      { id: '1', name: 'Espresso', price: 18000, category: 'coffee', image_url: '' },
      { id: '2', name: 'Cappuccino', price: 25000, category: 'coffee', image_url: '' },
      { id: '3', name: 'Latte', price: 28000, category: 'coffee', image_url: '' },
      { id: '4', name: 'Americano', price: 22000, category: 'coffee', image_url: '' },
      { id: '5', name: 'Green Tea', price: 20000, category: 'tea', image_url: '' },
      { id: '6', name: 'Thai Tea', price: 24000, category: 'tea', image_url: '' },
      { id: '7', name: 'Croissant', price: 15000, category: 'food', image_url: '' },
      { id: '8', name: 'Sandwich', price: 35000, category: 'food', image_url: '' },
      { id: '9', name: 'Cheesecake', price: 30000, category: 'dessert', image_url: '' },
      { id: '10', name: 'Tiramisu', price: 32000, category: 'dessert', image_url: '' },
    ];
    renderMenu();
  }
}

function renderMenu() {
  const grid = $('menu-grid');
  const search = ($('menu-search')?.value || '').toLowerCase();

  const filtered = menuItems.filter(item => {
    const matchesCategory = selectedCategory === 'all' || 
      item.category === selectedCategory || 
      (item.categories && item.categories.some(c => c.name === selectedCategory));
    const matchesSearch = !search || (item.name || '').toLowerCase().includes(search);
    return matchesCategory && matchesSearch;
  });

  if (filtered.length === 0) {
    grid.innerHTML = `
      <div class="col-span-full py-12 text-center text-[var(--text-secondary)]">
        <svg class="w-12 h-12 mx-auto mb-3 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
        </svg>
        <p class="text-sm font-medium">No items found</p>
      </div>
    `;
    return;
  }

  grid.innerHTML = filtered.map(item => {
    const itemCat = item.category || (item.categories && item.categories.length > 0 ? item.categories[0].name : 'Menu');
    const isBestSeller = bestSellers.some(bs => String(bs.menu_id) === String(item.id));
    
    return `
      <div class="menu-card p-3 relative cursor-pointer" data-id="${item.id}" onclick="addToOrder('${item.id}')">
        ${isBestSeller ? `
          <div class="absolute top-2 right-2 bg-gradient-to-r from-amber-500 to-orange-500 text-white font-black text-[9px] uppercase tracking-wider px-2 py-0.5 rounded-full shadow-md z-10 flex items-center gap-0.5">
            🔥 Best Seller
          </div>
        ` : ''}
        <div class="w-full h-24 rounded-lg bg-gradient-to-br from-gray-100 to-gray-200 flex items-center justify-center mb-3 text-3xl">
          ${getCategoryEmoji(itemCat)}
        </div>
        <h4 class="font-semibold text-sm text-[var(--text-primary)] truncate">${item.name}</h4>
        <p class="text-xs text-[var(--text-secondary)] mt-0.5">${itemCat}</p>
        <p class="text-sm font-bold text-kopitiam-600 mt-2">${formatRupiah(item.price)}</p>
      </div>
    `;
  }).join('');
}

function getCategoryEmoji(cat) {
  const map = { coffee: '☕', tea: '🍵', food: '🥐', dessert: '🍰' };
  return map[cat] || '🍽️';
}

// ── Order Management ────────────────────────────────────────
function addToOrder(itemId) {
  const item = menuItems.find(i => String(i.id) === String(itemId));
  if (!item) return;

  const existing = currentOrder.find(o => String(o.id) === String(itemId));
  if (existing) {
    existing.quantity += 1;
  } else {
    currentOrder.push({ ...item, quantity: 1, notes: '' });
  }
  renderOrder();
  showToast(`${item.name} added`, 'success');
}

function updateQuantity(itemId, delta) {
  const item = currentOrder.find(o => String(o.id) === String(itemId));
  if (!item) return;

  item.quantity += delta;
  if (item.quantity <= 0) {
    currentOrder = currentOrder.filter(o => String(o.id) !== String(itemId));
  }
  renderOrder();
}

function updateNotes(itemId, value) {
  const item = currentOrder.find(o => String(o.id) === String(itemId));
  if (item) {
    item.notes = value;
  }
}

function clearOrder() {
  currentOrder = [];
  renderOrder();
}

function renderOrder() {
  const container = $('order-items');

  if (currentOrder.length === 0) {
    container.innerHTML = `
      <div id="empty-order" class="text-center py-8 text-[var(--text-secondary)]">
        <svg class="w-12 h-12 mx-auto mb-3 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 11V7a4 4 0 00-8 0v4M5 9h14l1 12H4L5 9z"/>
        </svg>
        <p class="text-sm font-medium">No items added yet</p>
        <p class="text-xs mt-1">Click menu items to add</p>
      </div>
    `;
  } else {
    container.innerHTML = currentOrder.map(item => `
      <div class="order-item p-3 mb-2 bg-[var(--bg-card)] border border-[var(--border-color)] rounded-xl space-y-2">
        <div class="flex items-start justify-between">
          <div>
            <h4 class="font-semibold text-sm text-[var(--text-primary)]">${item.name}</h4>
            <p class="text-xs text-[var(--text-secondary)]">${formatRupiah(item.price)} each</p>
          </div>
          <span class="text-sm font-bold text-kopitiam-600">${formatRupiah(item.price * item.quantity)}</span>
        </div>
        <div class="flex items-center justify-between gap-3">
          <div class="flex items-center gap-2">
            <button onclick="updateQuantity('${item.id}', -1)" class="quantity-btn minus">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 12H4"/></svg>
            </button>
            <span class="text-sm font-bold w-6 text-center">${item.quantity}</span>
            <button onclick="updateQuantity('${item.id}', 1)" class="quantity-btn plus">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/></svg>
            </button>
          </div>
          <input type="text" placeholder="Catatan (e.g. Less sugar)" class="w-full max-w-[160px] px-2 py-1 text-xs bg-[var(--bg-primary)] border border-[var(--border-color)] rounded-lg text-[var(--text-primary)] focus:border-kopitiam-500 outline-none transition-all" value="${item.notes || ''}" oninput="updateNotes('${item.id}', this.value)">
        </div>
      </div>
    `).join('');
  }

  updateTotals();
}

let taxConfigs = [];

async function loadTaxes() {
  try {
    const res = await apiGet('/api/v1/taxes');
    taxConfigs = res.data || [];
  } catch (err) {
    console.warn('Failed to load dynamic taxes, using 10% tax fallback:', err);
    taxConfigs = [
      { name: 'Tax', tax_type: 'tax', percentage: 0.1 }
    ];
  }
}

async function loadCategories() {
  try {
    const res = await apiGet('/api/v1/categories');
    const categories = res.data || [];
    const list = $('category-list');
    if (list) {
      list.innerHTML = `<button class="category-btn active" data-category="all">All</button>`;
      categories.forEach(cat => {
        const emoji = getCategoryEmoji(cat.name);
        list.innerHTML += `<button class="category-btn" data-category="${cat.name}">${emoji} ${cat.name}</button>`;
      });
      
      // Re-attach listeners to dynamically created buttons
      document.querySelectorAll('.category-btn').forEach(btn => {
        btn.addEventListener('click', () => {
          document.querySelectorAll('.category-btn').forEach(b => b.classList.remove('active'));
          btn.classList.add('active');
          selectedCategory = btn.getAttribute('data-category');
          renderMenu();
        });
      });
    }
  } catch (err) {
    console.error('Failed to load categories dynamically', err);
  }
}

function getCategoryEmoji(cat) {
  const c = String(cat).toLowerCase();
  if (c.includes('coffee') || c.includes('kopi')) return '☕';
  if (c.includes('tea') || c.includes('teh')) return '🍵';
  if (c.includes('berat') || c.includes('makanan') || c.includes('food')) return '🍛';
  if (c.includes('ringan')) return '🥐';
  if (c.includes('dessert') || c.includes('cake') || c.includes('sweet')) return '🍰';
  if (c.includes('minuman') || c.includes('drink')) return '🥤';
  return '🍽️';
}

function calculateTotals(subtotal) {
  let taxRate = 0.11; // fallback PPN
  let serviceRate = 0.05; // fallback Service
  
  if (taxConfigs && taxConfigs.length > 0) {
    taxRate = 0;
    serviceRate = 0;
    taxConfigs.forEach(t => {
      if (t.tax_type === 'tax' || t.name === 'PPN') {
        taxRate += t.percentage;
      } else if (t.tax_type === 'service' || t.name === 'Service Charge') {
        serviceRate += t.percentage;
      }
    });
  }
  
  const taxAmount = Math.round(subtotal * taxRate);
  const serviceCharge = Math.round(subtotal * serviceRate);
  const total = subtotal + taxAmount + serviceCharge;
  
  return { taxAmount, serviceCharge, total, taxRate, serviceRate };
}

function updateTotals() {
  const subtotal = currentOrder.reduce((sum, item) => sum + (item.price * item.quantity), 0);
  const { taxAmount, serviceCharge, total, taxRate, serviceRate } = calculateTotals(subtotal);

  $('subtotal').textContent = formatRupiah(subtotal);
  
  const taxLabel = $('tax-label');
  if (taxLabel) {
    taxLabel.textContent = `Tax (${Math.round(taxRate * 100)}%) + Service (${Math.round(serviceRate * 100)}%)`;
  }
  $('tax').textContent = formatRupiah(taxAmount + serviceCharge);
  $('total').textContent = formatRupiah(total);
  const chargeTextEl = $('charge-text');
  if (chargeTextEl) {
    chargeTextEl.textContent = `Charge ${formatRupiah(total)}`;
  }
}

// ── Payment ─────────────────────────────────────────────────
function setPayment(method) {
  selectedPayment = method;
  document.querySelectorAll('.payment-method-btn').forEach(btn => {
    btn.classList.toggle('active', btn.dataset.payment === method);
  });
}

async function processPayment() {
  if (currentOrder.length === 0) {
    showToast('Add items to order first', 'warning');
    return;
  }
  if (isProcessing) return;

  isProcessing = true;
  const btn = $('btn-charge');
  btn.disabled = true;
  btn.innerHTML = '<span class="kt-spinner mr-2"></span> Processing...';

  try {
    const subtotal = currentOrder.reduce((sum, item) => sum + (item.price * item.quantity), 0);
    const { taxAmount, serviceCharge, total } = calculateTotals(subtotal);

    const orderData = {
      items: currentOrder.map(item => ({
        menu_id: item.id,
        menu_name: item.name,
        quantity: item.quantity,
        price: item.price,
        notes: item.notes || ''
      })),
      order_type: 'dine_in',
      tax_amount: taxAmount + serviceCharge,
      total_amount: total,
    };

    const res = await apiPost('/api/v1/orders', orderData);
    const orderId = res.data?.id;

    if (!orderId) {
      throw new Error('Order creation did not return an ID');
    }

    const backendPayment = selectedPayment === 'card' ? 'debit_card' : selectedPayment;
    const paymentData = {
      change_amount: 0,
      splits: [
        {
          payment_method: backendPayment,
          amount: total
        }
      ]
    };

    const payRes = await apiPost(`/api/v1/payments?order_id=${orderId}`, paymentData);
    const payment = payRes.data;

    if (payment && payment.snap_token) {
      if (typeof window.snap === 'undefined') {
        showToast('Midtrans Snap SDK is not loaded. Check internet connection and refresh.', 'error');
        isProcessing = false;
        btn.disabled = false;
        btn.innerHTML = `<span id="charge-text">Charge ${$('total').textContent}</span>`;
        return;
      }
      window.snap.pay(payment.snap_token, {
        onSuccess: function(result) {
          showReceipt(orderId, selectedPayment, total);
          showToast('Payment successful!', 'success');
        },
        onPending: function(result) {
          showReceipt(orderId, selectedPayment, total);
          showToast('Payment pending!', 'info');
        },
        onError: function(result) {
          showToast('Payment failed: ' + (result.status_message || ''), 'error');
          isProcessing = false;
          btn.disabled = false;
          btn.innerHTML = `<span id="charge-text">Charge ${$('total').textContent}</span>`;
        },
        onClose: function() {
          showToast('Payment window closed before completion', 'warning');
          isProcessing = false;
          btn.disabled = false;
          btn.innerHTML = `<span id="charge-text">Charge ${$('total').textContent}</span>`;
        }
      });
    } else {
      showReceipt(orderId, selectedPayment, total);
      showToast('Payment successful!', 'success');
    }
  } catch (err) {
    showToast('Payment failed: ' + err.message, 'error');
    isProcessing = false;
    btn.disabled = false;
    btn.innerHTML = `<span id="charge-text">Charge ${$('total').textContent}</span>`;
  }
}

function showReceipt(orderId, selectedPayment, total) {
  $('receipt-order-id').textContent = `#${orderId.substring(0, 8)}`;
  $('receipt-payment').textContent = selectedPayment.toUpperCase();
  $('receipt-total').textContent = formatRupiah(total);
  $('payment-modal').classList.remove('hidden');
  $('payment-modal').classList.add('flex');
}

function closePaymentModal() {
  $('payment-modal').classList.add('hidden');
  $('payment-modal').classList.remove('flex');
  currentOrder = [];
  renderOrder();
  isProcessing = false;
  const btn = $('btn-charge');
  btn.disabled = false;
  btn.innerHTML = '<span id="charge-text">Charge Rp 0</span>';
}

// ── Toast ─────────────────────────────────────────────────────
function showToast(message, type = 'info') {
  const container = $('toast-container');
  if (!container) return;

  let msgStr = message;
  if (typeof message === 'object' && message !== null) {
    msgStr = message.message || message.error || JSON.stringify(message);
  }

  const toast = document.createElement('div');
  const colors = {
    error: 'bg-red-50 border-red-200 text-red-700',
    success: 'bg-green-50 border-green-200 text-green-700',
    warning: 'bg-amber-50 border-amber-200 text-amber-700',
    info: 'bg-blue-50 border-blue-200 text-blue-700',
  };
  const icons = {
    error: 'M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z',
    success: 'M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z',
    warning: 'M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z',
    info: 'M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z',
  };

  toast.className = `flex items-center gap-2 px-4 py-3 rounded-xl border shadow-lg ${colors[type]} kt-slide-in`;
  toast.innerHTML = `
    <svg class="w-5 h-5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="${icons[type]}"/>
    </svg>
    <span class="text-sm font-medium">${msgStr}</span>
  `;

  container.insertAdjacentElement('beforeend', toast);
  setTimeout(() => {
    toast.style.opacity = '0';
    toast.style.transform = 'translateY(-8px)';
    toast.style.transition = 'all 0.3s ease';
    setTimeout(() => toast.remove(), 300);
  }, 3000);
}

// ── Dark Mode ───────────────────────────────────────────────
function toggleDarkMode() {
  const html = document.documentElement;
  const isDark = html.getAttribute('data-theme') === 'dark';
  html.setAttribute('data-theme', isDark ? 'light' : 'dark');
  localStorage.setItem('kopitiam_theme', isDark ? 'light' : 'dark');
}

function initDarkMode() {
  const saved = localStorage.getItem('kopitiam_theme');
  const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
  if (saved === 'dark' || (!saved && prefersDark)) {
    document.documentElement.setAttribute('data-theme', 'dark');
  }
}

// ── Event Listeners ─────────────────────────────────────────
function setupEventListeners() {
  // Category filters
  document.querySelectorAll('.category-btn').forEach(btn => {
    btn.addEventListener('click', () => {
      selectedCategory = btn.dataset.category;
      document.querySelectorAll('.category-btn').forEach(b => b.classList.remove('active'));
      btn.classList.add('active');
      renderMenu();
    });
  });

  // Search
  const searchInput = $('menu-search');
  if (searchInput) {
    searchInput.addEventListener('input', renderMenu);
  }

  // Payment methods
  document.querySelectorAll('.payment-method-btn').forEach(btn => {
    btn.addEventListener('click', () => setPayment(btn.dataset.payment));
  });

  // Action buttons
  $('btn-charge')?.addEventListener('click', processPayment);
  $('btn-clear')?.addEventListener('click', clearOrder);
  $('btn-new-order')?.addEventListener('click', closePaymentModal);
  $('btn-dark-mode')?.addEventListener('click', toggleDarkMode);
  $('btn-logout')?.addEventListener('click', logout);

    // Close modal on backdrop click
    $('payment-modal')?.addEventListener('click', (e) => {
      if (e.target === $('payment-modal')) closePaymentModal();
    });
    $('modal-best-sellers')?.addEventListener('click', (e) => {
      if (e.target === $('modal-best-sellers')) closeBestSellersModal();
    });
  }

  // ── Best Sellers Modal ─────────────────────────────────────────
  function openBestSellersModal() {
    const modal = $('modal-best-sellers');
    const container = $('best-sellers-list-container');
    if (!modal || !container) return;

    modal.classList.remove('hidden');
    modal.classList.add('flex');

    if (bestSellers.length === 0) {
      container.innerHTML = `
        <div class="text-center py-8 text-xs text-[var(--text-secondary)]">
          No best-selling items calculated yet.
        </div>
      `;
      return;
    }

    container.innerHTML = bestSellers.map((item, index) => {
      const rankColors = [
        'bg-amber-500 text-white shadow-amber-500/20',
        'bg-slate-400 text-white shadow-slate-400/20',
        'bg-amber-700 text-white shadow-amber-700/20'
      ];
      const rankBadgeClass = index < 3 
        ? `${rankColors[index]} w-7 h-7 rounded-full flex items-center justify-center font-bold text-xs shadow-md`
        : 'bg-[var(--bg-secondary)] border border-[var(--border-color)] text-[var(--text-secondary)] w-7 h-7 rounded-full flex items-center justify-center font-bold text-xs';
      
      return `
        <div class="flex items-center justify-between p-3 bg-[var(--bg-secondary)] border border-[var(--border-color)] rounded-xl hover:shadow-sm transition-all">
          <div class="flex items-center gap-3">
            <div class="${rankBadgeClass}">
              ${index + 1}
            </div>
            <div>
              <h4 class="font-bold text-sm text-[var(--text-primary)]">${item.menu_name}</h4>
              <p class="text-[10px] text-[var(--text-secondary)] mt-0.5">${formatNumber(item.total_sold)} items sold</p>
            </div>
          </div>
          <span class="text-xs font-bold text-kopitiam-600">${formatRupiah(item.revenue)}</span>
        </div>
      `;
    }).join('');
  }

  function closeBestSellersModal() {
    const modal = $('modal-best-sellers');
    if (modal) {
      modal.classList.add('hidden');
      modal.classList.remove('flex');
    }
  }

  // ── Init ────────────────────────────────────────────────────
  document.addEventListener('DOMContentLoaded', async () => {
    initDarkMode();
    setupEventListeners();
    if (!checkAuth()) return;
    await loadTaxes();
    await loadCategories();
    loadMenu();
  });

  // Expose to global
  window.addToOrder = addToOrder;
  window.updateQuantity = updateQuantity;
  window.clearOrder = clearOrder;
  window.setPayment = setPayment;
  window.processPayment = processPayment;
  window.closePaymentModal = closePaymentModal;
  window.logout = logout;
  window.toggleDarkMode = toggleDarkMode;
  window.openBestSellersModal = openBestSellersModal;
  window.closeBestSellersModal = closeBestSellersModal;
  window.updateNotes = updateNotes;
