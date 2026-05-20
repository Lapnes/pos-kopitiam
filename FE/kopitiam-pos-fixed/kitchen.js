/**
 * ============================================================
 * KopiTiam POS — Kitchen Display System (kitchen.js)
 * Real-time order management for kitchen staff
 * ============================================================
 */

const CONFIG = {
  API_BASE_URL: '',
  TOKEN_KEY: 'kopitiam_token',
  USER_KEY: 'kopitiam_user',
};

let orders = [];
let currentFilter = 'all';
let refreshInterval = null;
let orderTimers = {};

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
    if (elRole) elRole.textContent = user.role || 'Kitchen';
    if (elInitial) elInitial.textContent = (user.name || 'K').charAt(0).toUpperCase();

    // Role-based sidebar tab access control
    const linkCashier = $('link-cashier');
    const linkMenu = $('link-menu');
    const linkStock = $('link-stock');
    
    if (user.role === 'cashier') {
      if (linkCashier) linkCashier.classList.remove('hidden');
      if (linkMenu) linkMenu.classList.add('hidden');
      if (linkStock) linkStock.classList.add('hidden');
    } else if (user.role === 'manager') {
      if (linkCashier) linkCashier.classList.remove('hidden');
      if (linkMenu) linkMenu.classList.remove('hidden');
      if (linkStock) linkStock.classList.remove('hidden');
    } else { // kitchen
      if (linkCashier) linkCashier.classList.add('hidden');
      if (linkMenu) linkMenu.classList.remove('hidden');
      if (linkStock) linkStock.classList.remove('hidden');
    }

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

async function apiPatch(endpoint, body) {
  const token = localStorage.getItem(CONFIG.TOKEN_KEY);
  const res = await fetch(`${CONFIG.API_BASE_URL}${endpoint}`, {
    method: 'PATCH',
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

// ── Formatters ──────────────────────────────────────────────
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

function formatElapsed(seconds) {
  const mins = Math.floor(seconds / 60);
  const secs = seconds % 60;
  return `${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`;
}

// ── Orders ──────────────────────────────────────────────────
async function loadOrders() {
  try {
    const data = await apiGet('/api/v1/orders?status=pending,confirmed,paid,cooking,ready&limit=50');
    orders = (data.data || []).map(o => {
      let kdsStatus = o.status;
      if (o.status === 'pending' || o.status === 'confirmed' || o.status === 'paid') kdsStatus = 'pending';
      else if (o.status === 'cooking') kdsStatus = 'preparing';
      else if (o.status === 'ready') kdsStatus = 'ready';
      else if (o.status === 'served') kdsStatus = 'completed';
      return { ...o, status: kdsStatus };
    });
    renderOrders();
    updateCounts();
  } catch (err) {
    showToast('Failed to load orders: ' + err.message, 'error');
    // Fallback demo data
    orders = [
      {
        id: '1001',
        order_number: 'ORD-001',
        status: 'pending',
        items: [
          { menu_item_name: 'Espresso', quantity: 2, notes: 'Extra hot' },
          { menu_item_name: 'Croissant', quantity: 1, notes: '' }
        ],
        total_amount: 51000,
        created_at: new Date(Date.now() - 120000).toISOString(),
        order_type: 'dine_in',
        table_number: 'A3'
      },
      {
        id: '1002',
        order_number: 'ORD-002',
        status: 'preparing',
        items: [
          { menu_item_name: 'Cappuccino', quantity: 1, notes: 'Oat milk' },
          { menu_item_name: 'Latte', quantity: 2, notes: '' },
          { menu_item_name: 'Cheesecake', quantity: 1, notes: '' }
        ],
        total_amount: 108000,
        created_at: new Date(Date.now() - 300000).toISOString(),
        order_type: 'takeaway',
        table_number: null
      },
      {
        id: '1003',
        order_number: 'ORD-003',
        status: 'ready',
        items: [
          { menu_item_name: 'Green Tea', quantity: 1, notes: 'Less sugar' }
        ],
        total_amount: 20000,
        created_at: new Date(Date.now() - 600000).toISOString(),
        order_type: 'dine_in',
        table_number: 'B1'
      }
    ];
    renderOrders();
    updateCounts();
  }
}

function updateCounts() {
  const counts = { pending: 0, preparing: 0, ready: 0, completed: 0 };
  orders.forEach(o => { if (counts[o.status] !== undefined) counts[o.status]++; });

  $('count-pending').textContent = counts.pending;
  $('count-preparing').textContent = counts.preparing;
  $('count-ready').textContent = counts.ready;
  $('count-completed').textContent = counts.completed;
}

function renderOrders() {
  const grid = $('orders-grid');

  const filtered = currentFilter === 'all' 
    ? orders 
    : orders.filter(o => o.status === currentFilter);

  if (filtered.length === 0) {
    grid.innerHTML = `
      <div class="col-span-full py-12 text-center text-[var(--text-secondary)]">
        <svg class="w-12 h-12 mx-auto mb-3 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"/>
        </svg>
        <p class="text-sm font-medium">No ${currentFilter} orders</p>
      </div>
    `;
    return;
  }

  const user = JSON.parse(localStorage.getItem(CONFIG.USER_KEY) || '{}');
  const isCashier = user.role === 'cashier';
  const isManager = user.role === 'manager';
  const canCancel = isCashier || isManager;
  const canProcess = user.role === 'kitchen' || isManager;

  grid.innerHTML = filtered.map(order => {
    const elapsed = Math.floor((Date.now() - new Date(order.created_at)) / 1000);
    const isUrgent = elapsed > 600 && order.status !== 'ready' && order.status !== 'completed';
    const statusLabels = { pending: 'Terima Pesanan', preparing: 'Proses Pesanan', ready: 'Pesanan Selesai', completed: 'Selesai' };
    const nextActions = {
      pending: { label: 'Terima Pesanan', action: 'preparing', class: 'kds-btn-primary' },
      preparing: { label: 'Proses Pesanan', action: 'ready', class: 'kds-btn-success' },
      ready: { label: 'Pesanan Selesai', action: 'completed', class: 'kds-btn-secondary' }
    };
    const next = nextActions[order.status];
    const showNext = canProcess && next;
    const showCancel = canCancel && order.status !== 'completed' && order.status !== 'cancelled';

    return `
      <div class="order-card status-${order.status} kt-fade-in">
        <div class="flex items-start justify-between mb-3">
          <div>
            <div class="flex items-center gap-2 mb-1">
              <h3 class="font-bold text-lg text-[var(--text-primary)]">${order.order_number || `#${order.id}`}</h3>
              <span class="status-badge ${order.status}">${statusLabels[order.status]}</span>
            </div>
            <div class="flex items-center gap-3 text-xs text-[var(--text-secondary)]">
              <span class="flex items-center gap-1">
                <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>
                <span class="timer ${isUrgent ? 'urgent' : ''}" data-order-id="${order.id}" data-created="${order.created_at}">${formatElapsed(elapsed)}</span>
              </span>
              ${order.table_number ? `<span class="flex items-center gap-1"><svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"/></svg> Table ${order.table_number}</span>` : '<span class="flex items-center gap-1"><svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"/></svg> Takeaway</span>'}
            </div>
          </div>
          <span class="text-sm font-bold text-kopitiam-600">${formatRupiah(order.total || order.total_amount || 0)}</span>
        </div>

        <div class="space-y-2 mb-4">
          ${(order.items || []).map(item => `
            <div class="flex items-start justify-between text-sm">
              <div class="flex items-start gap-2">
                <span class="font-bold text-kopitiam-600 min-w-[1.5rem]">${item.quantity}x</span>
                <div>
                  <span class="font-medium text-[var(--text-primary)]">${item.menu_item_name || item.name || 'Item'}</span>
                  <p class="text-xs ${item.notes ? 'text-amber-600 font-medium' : 'text-gray-400 italic'} mt-0.5">
                    📝 ${item.notes || 'tidak ada catatan'}
                  </p>
                </div>
              </div>
            </div>
          `).join('')}
        </div>

        <div class="flex gap-2 pt-3 border-t border-[var(--border-color)]">
          ${showNext ? `
            <button onclick="updateOrderStatus('${order.id}', '${next.action}')" class="${next.class} kds-btn flex-1 justify-center">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/></svg>
              ${next.label}
            </button>
          ` : ''}
          ${showCancel ? `
            <button onclick="cancelOrder('${order.id}')" class="kds-btn kds-btn-secondary flex-1 justify-center gap-1.5" title="Cancel Order">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/></svg>
              Cancel Order
            </button>
          ` : ''}
        </div>
      </div>
    `;
  }).join('');
}

// ── Order Actions ───────────────────────────────────────────
async function updateOrderStatus(orderId, newStatus) {
  try {
    const backendStatus = newStatus === 'preparing' ? 'cooking' : (newStatus === 'completed' ? 'served' : newStatus);
    await apiPatch(`/api/v1/orders/${orderId}/status`, { status: backendStatus });
    showToast(`Order marked as ${newStatus}`, 'success');
    loadOrders();
  } catch (err) {
    showToast('Failed to update: ' + err.message, 'error');
  }
}

async function cancelOrder(orderId) {
  if (!confirm('Cancel this order?')) return;
  try {
    await apiPatch(`/api/v1/orders/${orderId}/status`, { status: 'cancelled' });
    showToast('Order cancelled', 'warning');
    loadOrders();
  } catch (err) {
    showToast('Failed to cancel: ' + err.message, 'error');
  }
}

// ── Timer Update ────────────────────────────────────────────
function startTimerUpdates() {
  setInterval(() => {
    document.querySelectorAll('.timer[data-created]').forEach(el => {
      const created = new Date(el.dataset.created);
      const elapsed = Math.floor((Date.now() - created) / 1000);
      el.textContent = formatElapsed(elapsed);

      const orderCard = el.closest('.order-card');
      if (elapsed > 600 && orderCard && !orderCard.classList.contains('status-ready') && !orderCard.classList.contains('status-completed')) {
        el.classList.add('urgent');
      }
    });
  }, 1000);
}

// ── Toast ───────────────────────────────────────────────────
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
  document.querySelectorAll('.category-btn').forEach(btn => {
    btn.addEventListener('click', () => {
      currentFilter = btn.dataset.filter;
      document.querySelectorAll('.category-btn').forEach(b => b.classList.remove('active'));
      btn.classList.add('active');
      renderOrders();
    });
  });

  $('btn-dark-mode')?.addEventListener('click', toggleDarkMode);
  $('btn-logout')?.addEventListener('click', logout);
}

// ── Init ────────────────────────────────────────────────────
document.addEventListener('DOMContentLoaded', () => {
  initDarkMode();
  setupEventListeners();
  if (!checkAuth()) return;
  loadOrders();
  startTimerUpdates();

  // Auto refresh every 30 seconds
  refreshInterval = setInterval(loadOrders, 30000);
});

// ── Stock View & Material Inventory Tracking ────────────────
let rawMaterials = [];
let allMenus = [];
let selectedKitchenTab = 'orders';

function switchKitchenTab(tab) {
  selectedKitchenTab = tab;
  document.querySelectorAll('.sidebar-link').forEach(link => link.classList.remove('active'));
  
  if (tab === 'orders') {
    $('link-kds')?.classList.add('active');
    $('orders-view').classList.remove('hidden');
    $('stock-view').classList.add('hidden');
    $('menu-view').classList.add('hidden');
    loadOrders();
    // Resume polling KDS orders
    if (!refreshInterval) {
      refreshInterval = setInterval(loadOrders, 30000);
    }
  } else if (tab === 'stock') {
    $('link-stock')?.classList.add('active');
    $('orders-view').classList.add('hidden');
    $('stock-view').classList.remove('hidden');
    $('menu-view').classList.add('hidden');
    // Pause KDS order auto-refresh to save bandwidth
    if (refreshInterval) {
      clearInterval(refreshInterval);
      refreshInterval = null;
    }
    loadStockData();
  } else if (tab === 'menu') {
    $('link-menu')?.classList.add('active');
    $('orders-view').classList.add('hidden');
    $('stock-view').classList.add('hidden');
    $('menu-view').classList.remove('hidden');
    if (refreshInterval) {
      clearInterval(refreshInterval);
      refreshInterval = null;
    }
    loadMenuData();
  }
}

async function loadStockData() {
  const tbody = $('stock-table-body');
  if (!tbody) return;
  
  tbody.innerHTML = `
    <tr>
      <td colspan="7" class="py-8 text-center text-xs text-[var(--text-secondary)]">
        <div class="kt-spinner mx-auto mb-2"></div>
        Fetching daily stock estimation data...
      </td>
    </tr>
  `;
  
  try {
    const res = await apiGet('/api/v1/inventory/stock-estimation');
    rawMaterials = res.data || [];
    renderStockTable(rawMaterials);
  } catch (err) {
    console.error('Failed to load raw material stock estimation:', err);
    tbody.innerHTML = `
      <tr>
        <td colspan="7" class="py-8 text-center text-xs text-red-500 font-semibold">
          ⚠️ Failed to load raw materials estimation. Click Refresh to try again.
        </td>
      </tr>
    `;
    showToast('Failed to load raw material stock estimation', 'error');
  }
}

function renderStockTable(items) {
  const tbody = $('stock-table-body');
  if (!tbody) return;
  
  if (items.length === 0) {
    tbody.innerHTML = `
      <tr>
        <td colspan="7" class="py-8 text-center text-xs text-[var(--text-secondary)]">
          No matching raw materials found.
        </td>
      </tr>
    `;
    return;
  }
  
  tbody.innerHTML = items.map(item => {
    const isInsufficient = item.sufficiency === 'insufficient';
    const statusText = isInsufficient ? '⚠️ Kurang' : '✅ Cukup / Aman';
    const statusClass = isInsufficient 
      ? 'bg-red-500/10 text-red-500 font-bold border border-red-500/20' 
      : 'bg-green-500/10 text-green-500 border border-green-500/20';
    
    return `
      <tr class="hover:bg-gray-50/50 transition-colors text-sm text-[var(--text-primary)]">
        <td class="py-3.5 px-4 font-semibold">${item.material_name}</td>
        <td class="py-3.5 px-4 font-bold text-[var(--text-primary)]">${formatNumber(item.current_stock)}</td>
        <td class="py-3.5 px-4 text-[var(--text-secondary)]">${formatNumber(item.min_stock_level)}</td>
        <td class="py-3.5 px-4 font-semibold text-blue-600">${formatNumber(item.daily_required)}</td>
        <td class="py-3.5 px-4">
          <span class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-semibold ${statusClass}">
            ${statusText}
          </span>
        </td>
        <td class="py-3.5 px-4 font-bold ${isInsufficient ? 'text-red-500' : 'text-[var(--text-secondary)]'}">
          ${item.shortage > 0 ? formatNumber(item.shortage) : '-'}
        </td>
        <td class="py-3.5 px-4 text-xs font-semibold text-[var(--text-secondary)]">${item.unit}</td>
      </tr>
    `;
  }).join('');
}

function searchMaterials() {
  const query = $('stock-search').value.toLowerCase().trim();
  if (!query) {
    renderStockTable(rawMaterials);
    return;
  }
  const filtered = rawMaterials.filter(item => 
    item.material_name.toLowerCase().includes(query) ||
    item.unit.toLowerCase().includes(query)
  );
  renderStockTable(filtered);
}

// ── Menu Availability ──────────────────────────────────────────
async function loadMenuData() {
  const tbody = $('menu-table-body');
  if (!tbody) return;
  
  tbody.innerHTML = `
    <tr>
      <td colspan="7" class="py-8 text-center text-xs text-[var(--text-secondary)]">
        <div class="kt-spinner mx-auto mb-2"></div>
        Fetching all menus...
      </td>
    </tr>
  `;
  
  try {
    const res = await apiGet('/api/v1/menus/all');
    allMenus = res.data || [];
    renderMenuTable(allMenus);
  } catch (err) {
    console.error('Failed to load menus:', err);
    tbody.innerHTML = `
      <tr>
        <td colspan="7" class="py-8 text-center text-xs text-red-500 font-semibold">
          ⚠️ Failed to load menus. Click Refresh to try again.
        </td>
      </tr>
    `;
    showToast('Failed to load menus', 'error');
  }
}

function renderMenuTable(items) {
  const tbody = $('menu-table-body');
  if (!tbody) return;
  
  if (items.length === 0) {
    tbody.innerHTML = `
      <tr>
        <td colspan="7" class="py-8 text-center text-xs text-[var(--text-secondary)]">
          No matching menu items found.
        </td>
      </tr>
    `;
    return;
  }
  
  tbody.innerHTML = items.map(item => {
    const statusText = item.is_active ? '✅ Available' : '❌ Sold Out';
    const statusClass = item.is_active 
      ? 'bg-green-500/10 text-green-600 border border-green-500/20' 
      : 'bg-red-500/10 text-red-600 border border-red-500/20';
    
    const cats = (item.Categories || []).map(c => c.name).join(', ') || '-';
    
    const hasRecipe = item.is_recipe_based ? 'Yes' : 'No';
    const recipeClass = item.is_recipe_based 
      ? 'bg-blue-500/10 text-blue-600 border border-blue-500/20' 
      : 'bg-gray-500/10 text-gray-500 border border-gray-500/20';
    
    const actionText = item.is_active ? 'Set Unavailable' : 'Set Available';
    const actionClass = item.is_active 
      ? 'bg-red-500 hover:bg-red-600 text-white font-semibold text-xs rounded-lg px-3 py-1.5 transition-all shadow-sm cursor-pointer'
      : 'bg-green-500 hover:bg-green-600 text-white font-semibold text-xs rounded-lg px-3 py-1.5 transition-all shadow-sm cursor-pointer';

    return `
      <tr class="hover:bg-gray-50/50 transition-colors text-sm text-[var(--text-primary)]">
        <td class="py-3.5 px-4 font-semibold">${item.name}</td>
        <td class="py-3.5 px-4 text-[var(--text-secondary)]">${cats}</td>
        <td class="py-3.5 px-4 font-semibold text-kopitiam-600">${formatRupiah(item.price)}</td>
        <td class="py-3.5 px-4 text-[var(--text-secondary)] font-bold">${formatNumber(item.daily_stock || 0)}</td>
        <td class="py-3.5 px-4">
          <span class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-semibold ${recipeClass}">
            ${hasRecipe}
          </span>
        </td>
        <td class="py-3.5 px-4">
          <span class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-semibold ${statusClass}">
            ${statusText}
          </span>
        </td>
        <td class="py-3.5 px-4 text-right">
          <button onclick="toggleMenuAvailability('${item.id}', ${item.is_active})" class="${actionClass}">
            ${actionText}
          </button>
        </td>
      </tr>
    `;
  }).join('');
}

async function toggleMenuAvailability(id, currentActive) {
  try {
    const token = localStorage.getItem(CONFIG.TOKEN_KEY);
    const targetUrl = `${CONFIG.API_BASE_URL}/api/v1/menus/${id}`;
    
    const currentItem = allMenus.find(m => m.id === id);
    if (!currentItem) return;

    // Build categories payload to ensure correct GORM association saving
    const catsPayload = (currentItem.Categories || []).map(c => ({ id: c.id, name: c.name }));

    const body = {
      ...currentItem,
      is_active: !currentActive,
      categories: catsPayload
    };

    const res = await fetch(targetUrl, {
      method: 'PUT',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
        'Accept': 'application/json'
      },
      body: JSON.stringify(body)
    });

    if (res.status === 401) { logout(); return; }
    const data = await res.json();
    if (!res.ok || !data.success) throw new Error(data.message || 'Failed to update menu');

    showToast(`Menu status updated successfully!`, 'success');
    loadMenuData();
  } catch (err) {
    console.error('Error toggling menu availability:', err);
    showToast('Failed to toggle menu: ' + err.message, 'error');
  }
}

function searchMenus() {
  const query = $('menu-search').value.toLowerCase().trim();
  if (!query) {
    renderMenuTable(allMenus);
    return;
  }
  const filtered = allMenus.filter(item => 
    item.name.toLowerCase().includes(query) ||
    (item.Categories || []).some(c => c.name.toLowerCase().includes(query))
  );
  renderMenuTable(filtered);
}

function formatNumber(num) {
  if (!num && num !== 0) return '0';
  return Math.round(num).toString().replace(/\B(?=(\d{3})+(?!\d))/g, '.');
}

window.updateOrderStatus = updateOrderStatus;
window.cancelOrder = cancelOrder;
window.logout = logout;
window.toggleDarkMode = toggleDarkMode;
window.switchKitchenTab = switchKitchenTab;
window.loadStockData = loadStockData;
window.searchMaterials = searchMaterials;
window.loadMenuData = loadMenuData;
window.toggleMenuAvailability = toggleMenuAvailability;
window.searchMenus = searchMenus;
