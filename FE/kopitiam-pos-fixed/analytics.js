/**
 * ============================================================
 * KopiTiam POS — Analytics Dashboard (analytics.js) [FIXED]
 * Real-time insights dashboard with Chart.js
 * All event listeners attached after DOM ready
 * Data fetched from API with loading states & error handling
 * ============================================================
 */

// ── Configuration ───────────────────────────────────────────
const CONFIG = {
  API_BASE_URL: '',
  TOKEN_KEY: 'kopitiam_token',
  USER_KEY: 'kopitiam_user',
  REFRESH_KEY: 'kopitiam_refresh_token',
};

// ── State ───────────────────────────────────────────────────
let salesChart = null;
let paymentChart = null;
let currentPeriod = '7D';
let isLoading = false;
let dateRange = { start: '', end: '' };
let autoRefreshInterval = null;
let retryCount = 0;
const MAX_RETRIES = 3;
let auditLogLimit = 100;

// ── DOM Helpers ─────────────────────────────────────────────
const $ = (id) => document.getElementById(id);

// ── Config Loader ───────────────────────────────────────────
async function loadConfig() {
  try {
    const res = await fetch('./config.json');
    if (res.ok) {
      const data = await res.json();
      if (data && data.API_BASE_URL !== undefined) {
        CONFIG.API_BASE_URL = data.API_BASE_URL;
      }
    }
  } catch (e) { /* fallback */ }
}

// ── Auth Check ──────────────────────────────────────────────
function checkAuth() {
  const token = localStorage.getItem(CONFIG.TOKEN_KEY);
  const user = JSON.parse(localStorage.getItem(CONFIG.USER_KEY) || '{}');

  if (!token) {
    window.location.href = './index.html';
    return false;
  }

  // Restrict access to Manager role only
  const role = (user.role || '').toLowerCase();
  if (role !== 'manager') {
    if (role === 'cashier') {
      window.location.href = './cashier.html';
    } else if (role === 'kitchen') {
      window.location.href = './kitchen.html';
    } else {
      window.location.href = './index.html';
    }
    return false;
  }

  // Populate user info
  if (user.name) {
    const userNameEl = $('user-name');
    const userRoleEl = $('user-role');
    const userInitialEl = $('user-initial');
    if (userNameEl) userNameEl.textContent = user.name;
    if (userRoleEl) userRoleEl.textContent = user.role || 'Manager';
    if (userInitialEl) userInitialEl.textContent = (user.name || 'M').charAt(0).toUpperCase();
  }

  return token;
}

function logout() {
  localStorage.removeItem(CONFIG.TOKEN_KEY);
  localStorage.removeItem(CONFIG.REFRESH_KEY);
  localStorage.removeItem(CONFIG.USER_KEY);
  window.location.href = './index.html';
}

// ── API Helpers with Retry ──────────────────────────────────
async function fetchAPI(endpoint, method = 'GET', body = null) {
  const token = localStorage.getItem(CONFIG.TOKEN_KEY);
  
  let formattedEndpoint = endpoint;
  if (!endpoint.startsWith('/api/v1/')) {
    formattedEndpoint = '/api/v1' + (endpoint.startsWith('/') ? endpoint : '/' + endpoint);
  }

  const options = {
    method,
    headers: {
      'Authorization': `Bearer ${token}`,
      'Accept': 'application/json',
    }
  };

  if (body) {
    options.headers['Content-Type'] = 'application/json';
    options.body = JSON.stringify(body);
  }

  const res = await fetch(`${CONFIG.API_BASE_URL}${formattedEndpoint}`, options);

  if (res.status === 401) {
    logout();
    throw new Error('Session expired. Please login again.');
  }

  const data = await res.json();
  return data;
}

async function apiGet(endpoint, retries = 0) {
  const token = localStorage.getItem(CONFIG.TOKEN_KEY);

  try {
    const res = await fetch(`${CONFIG.API_BASE_URL}${endpoint}`, {
      headers: {
        'Authorization': `Bearer ${token}`,
        'Accept': 'application/json',
      }
    });

    if (res.status === 401) {
      logout();
      throw new Error('Session expired. Please login again.');
    }

    const data = await res.json();
    if (!res.ok || !data.success) {
      throw new Error(data.message || data.error || `HTTP ${res.status}`);
    }
    retryCount = 0;
    return data;
  } catch (err) {
    if (retries < MAX_RETRIES && (err.name === 'TypeError' || err.message.includes('network'))) {
      const delay = Math.pow(2, retries) * 1000;
      showToast(`Retrying... (${retries + 1}/${MAX_RETRIES})`, 'warning');
      await new Promise(r => setTimeout(r, delay));
      return apiGet(endpoint, retries + 1);
    }
    throw err;
  }
}

// ── Date Helpers ────────────────────────────────────────────
const ID_DAYS = ['Min', 'Sen', 'Sel', 'Rab', 'Kam', 'Jum', 'Sab'];
const ID_MONTHS = ['Jan', 'Feb', 'Mar', 'Apr', 'Mei', 'Jun', 'Jul', 'Agu', 'Sep', 'Okt', 'Nov', 'Des'];

function formatDateID(date) {
  const d = new Date(date);
  return `${ID_DAYS[d.getDay()]}, ${d.getDate()} ${ID_MONTHS[d.getMonth()]} ${d.getFullYear()}`;
}

function formatDateShort(date) {
  const d = new Date(date);
  return `${ID_DAYS[d.getDay()]}, ${d.getDate()} ${ID_MONTHS[d.getMonth()]}`;
}

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
  if (num === undefined || num === null || isNaN(num)) return '0';
  return num.toString().replace(/\B(?=(\d{3})+(?!\d))/g, '.');
}

function getDateRange(period) {
  const end = new Date();
  const start = new Date();
  switch (period) {
    case '7D': start.setDate(end.getDate() - 7); break;
    case '30D': start.setDate(end.getDate() - 30); break;
    case '90D': start.setDate(end.getDate() - 90); break;
  }
  return {
    start: start.toISOString().split('T')[0],
    end: end.toISOString().split('T')[0],
  };
}

// ── Loading State ───────────────────────────────────────────
function setDashboardLoading(loading) {
  isLoading = loading;
  const refreshBtn = $('refresh-btn');
  if (refreshBtn) {
    refreshBtn.disabled = loading;
    refreshBtn.innerHTML = loading
      ? `<span class="kt-spinner mr-2"></span> Refreshing...`
      : `<svg class="w-4 h-4 mr-1.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/></svg> Refresh`;
  }

  // Show shimmer on stat cards
  const statValues = document.querySelectorAll('.stat-value');
  statValues.forEach(el => {
    if (loading) {
      el.dataset.originalText = el.textContent;
      el.innerHTML = '<span class="shimmer" style="display:inline-block;width:80%;height:1.5rem;"></span>';
    } else {
      // Remove shimmer, restore text will be done by update functions
    }
  });

  // Show chart loading overlays
  const salesLoading = $('sales-chart-loading');
  const paymentLoading = $('payment-chart-loading');
  if (salesLoading) salesLoading.classList.toggle('hidden', !loading);
  if (paymentLoading) paymentLoading.classList.toggle('hidden', !loading);
}

// ── Fetch Dashboard Data ────────────────────────────────────
async function fetchDashboardData() {
  if (isLoading) return;

  setDashboardLoading(true);

  // Use dateRange if set, otherwise use period
  let start, end;
  if (dateRange.start && dateRange.end) {
    start = dateRange.start;
    end = dateRange.end;
  } else {
    const range = getDateRange(currentPeriod);
    start = range.start;
    end = range.end;
  }

  try {
    // Parallel fetch all endpoints with start_date and end_date filtering
    const [salesSummary, bestSellers, dailySales, paymentSummary, auditLogs] = await Promise.all([
      apiGet(`/api/v1/analytics/sales-summary?start_date=${start}&end_date=${end}`).catch(err => {
        console.warn('Sales summary fetch failed:', err.message);
        return { data: {} };
      }),
      apiGet(`/api/v1/analytics/best-sellers?limit=5&start_date=${start}&end_date=${end}`).catch(err => {
        console.warn('Best sellers fetch failed:', err.message);
        return { data: [] };
      }),
      apiGet(`/api/v1/analytics/daily-sales?start_date=${start}&end_date=${end}`).catch(err => {
        console.warn('Daily sales fetch failed:', err.message);
        return { data: [] };
      }),
      apiGet(`/api/v1/analytics/payment-summary?start_date=${start}&end_date=${end}`).catch(err => {
        console.warn('Payment summary fetch failed:', err.message);
        return { data: [] };
      }),
      apiGet(`/api/v1/analytics/audit-logs?limit=${auditLogLimit}`).catch(err => {
        console.warn('Audit logs fetch failed:', err.message);
        return { data: [] };
      }),
    ]);

    updateStatCards(salesSummary.data || {});
    updateSalesChart(dailySales.data || [], start, end);
    updatePaymentChart(paymentSummary.data || []);
    updateBestSellersTable(bestSellers.data || []);
    renderAuditLogs(auditLogs.data || []);

  } catch (err) {
    console.error('Dashboard load error:', err);
    showToast(err.message, 'error');
  } finally {
    setDashboardLoading(false);
  }
}

// ── Update Stat Cards ───────────────────────────────────────
function updateStatCards(data) {
  const gross = data.gross_revenue || 0;
  const net = data.net_revenue || 0;
  const transactions = data.total_transactions || 0;
  const aov = transactions > 0 ? Math.round(net / transactions) : 0;

  const stats = [
    { id: 'stat-gross', value: gross, prefix: 'Rp ', isMoney: true },
    { id: 'stat-net', value: net, prefix: 'Rp ', isMoney: true },
    { id: 'stat-orders', value: transactions, prefix: '', isMoney: false },
    { id: 'stat-aov', value: aov, prefix: 'Rp ', isMoney: true },
  ];

  stats.forEach(stat => {
    const valueEl = $(`${stat.id}-value`);
    if (valueEl) {
      valueEl.textContent = stat.prefix + (stat.isMoney ? formatNumber(stat.value) : formatNumber(stat.value));
      valueEl.classList.remove('shimmer');
    }
  });

  // Update trends using actual calculated values from the backend comparative analysis
  const grossChange = data.gross_revenue_change || 0;
  const netChange = data.net_revenue_change || 0;
  const ordersChange = data.total_transactions_change || 0;
  const aovChange = data.avg_order_value_change || 0;

  updateTrend('stat-gross-trend', grossChange, grossChange >= 0, 'vs last period');
  updateTrend('stat-net-trend', netChange, netChange >= 0, 'vs last period');
  updateTrend('stat-orders-trend', ordersChange, ordersChange >= 0, 'vs last period');
  updateTrend('stat-aov-trend', aovChange, aovChange >= 0, 'vs last period');
}

function updateTrend(elementId, percent, isUp, label) {
  const el = $(elementId);
  if (!el) return;

  const absPercent = Math.abs(percent).toFixed(1);
  const colorClass = isUp ? 'text-green-600' : 'text-red-500';
  const arrowPath = isUp 
    ? 'M5 10l7-7m0 0l7 7m-7-7v18' 
    : 'M19 14l-7 7m0 0l-7-7m7 7V3';

  el.className = `flex items-center gap-1 text-xs font-semibold ${colorClass}`;
  el.innerHTML = `
    <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="${arrowPath}"/>
    </svg>
    ${absPercent}% ${label}
  `;
}

// ── Update Sales Trend Chart ────────────────────────────────
function updateSalesChart(data, startDate, endDate) {
  const ctx = $('sales-chart')?.getContext('2d');
  if (!ctx) return;

  // Generate labels for the date range
  const labels = [];
  const values = [];

  const start = new Date(startDate);
  const end = new Date(endDate);
  const dayDiff = Math.ceil((end - start) / (1000 * 60 * 60 * 24));

  // Build date map from API data
  const dataMap = {};
  if (Array.isArray(data)) {
    data.forEach(item => {
      if (item.sale_date) {
        const dateStr = item.sale_date.split('T')[0];
        dataMap[dateStr] = item.gross_revenue || item.total_revenue || item.net_profit || 0;
      }
    });
  }

  // Generate labels and values for each day in range
  for (let i = 0; i <= dayDiff && i <= 30; i++) {
    const d = new Date(start);
    d.setDate(d.getDate() + i);
    const dateStr = d.toISOString().split('T')[0];
    labels.push(formatDateShort(d));
    values.push(dataMap[dateStr] || 0);
  }

  if (salesChart) salesChart.destroy();

  salesChart = new Chart(ctx, {
    type: 'line',
    data: {
      labels,
      datasets: [{
        label: 'Revenue',
        data: values,
        borderColor: '#16a34a',
        backgroundColor: 'rgba(22, 163, 74, 0.1)',
        borderWidth: 2.5,
        tension: 0.4,
        fill: true,
        pointBackgroundColor: '#16a34a',
        pointBorderColor: '#ffffff',
        pointBorderWidth: 2,
        pointRadius: 4,
        pointHoverRadius: 6,
      }]
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      plugins: {
        legend: { display: false },
        tooltip: {
          backgroundColor: 'rgba(17, 24, 39, 0.9)',
          titleColor: '#f9fafb',
          bodyColor: '#f9fafb',
          padding: 12,
          cornerRadius: 8,
          displayColors: false,
          callbacks: {
            label: (ctx) => `Revenue: ${formatRupiah(ctx.raw)}`,
          }
        }
      },
      scales: {
        x: {
          grid: { display: false },
          ticks: { color: '#9ca3af', font: { size: 11 } }
        },
        y: {
          grid: { color: 'rgba(0,0,0,0.05)' },
          ticks: {
            color: '#9ca3af',
            font: { size: 11 },
            callback: (val) => {
              if (val >= 1000000) return 'Rp ' + (val / 1000000).toFixed(1) + 'M';
              if (val >= 1000) return 'Rp ' + (val / 1000).toFixed(0) + 'K';
              return 'Rp ' + val;
            },
          },
          beginAtZero: true,
        }
      },
      interaction: { intersect: false, mode: 'index' },
    }
  });
}

// ── Update Payment Methods Chart ──────────────────────────
function updatePaymentChart(data) {
  const ctx = $('payment-chart')?.getContext('2d');
  if (!ctx) return;

  const colors = {
    cash: '#22c55e',
    qris: '#3b82f6',
    debit_card: '#a855f7',
    credit_card: '#f59e0b',
    transfer: '#06b6d4',
    e_wallet: '#ec4899',
  };

  const labels = [];
  const values = [];
  const bgColors = [];

  if (!Array.isArray(data) || data.length === 0) {
    // Empty state
    labels.push('No Data');
    values.push(1);
    bgColors.push('#e5e7eb');
  } else {
    data.forEach(item => {
      const method = item.payment_method || 'unknown';
      labels.push(method.replace(/_/g, ' ').toUpperCase());
      values.push(item.total_amount || 0);
      bgColors.push(colors[method] || '#6b7280');
    });
  }

  if (paymentChart) paymentChart.destroy();

  paymentChart = new Chart(ctx, {
    type: 'doughnut',
    data: {
      labels,
      datasets: [{
        data: values,
        backgroundColor: bgColors,
        borderWidth: 0,
        hoverOffset: 8,
      }]
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      cutout: '65%',
      plugins: {
        legend: {
          position: 'bottom',
          labels: {
            usePointStyle: true,
            pointStyle: 'circle',
            padding: 16,
            color: '#6b7280',
            font: { size: 12 }
          }
        },
        tooltip: {
          backgroundColor: 'rgba(17, 24, 39, 0.9)',
          padding: 12,
          cornerRadius: 8,
          callbacks: {
            label: (ctx) => {
              const total = ctx.dataset.data.reduce((a, b) => a + b, 0);
              if (total === 0) return 'No data';
              const pct = ((ctx.raw / total) * 100).toFixed(1);
              return `${ctx.label}: ${pct}% (${formatRupiah(ctx.raw)})`;
            }
          }
        }
      }
    }
  });
}

// ── Render Audit Logs ───────────────────────────────────────
function renderAuditLogs(logs) {
  const terminalBody = $('terminal-body');
  if (!terminalBody) return;

  if (!Array.isArray(logs) || logs.length === 0) {
    terminalBody.innerHTML = `<div class="text-slate-500">// No audit logs recorded yet.</div>`;
    return;
  }

  terminalBody.innerHTML = logs.map(log => {
    // Format timestamp nicely
    const date = new Date(log.timestamp);
    const timeStr = date.toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit', second: '2-digit' });
    const dateStr = date.toLocaleDateString('id-ID', { year: '2-digit', month: '2-digit', day: '2-digit' });
    
    // Determine action color tag
    let actionColor = 'text-yellow-400';
    if (log.action === 'CREATE') actionColor = 'text-green-400';
    if (log.action === 'DELETE') actionColor = 'text-red-400';
    if (log.action === 'VOID') actionColor = 'text-pink-500';

    // Format new_value if it's a JSON payload
    let detailStr = '';
    try {
      const parsed = JSON.parse(log.new_value);
      // Simplify printout details
      let details = [];
      if (parsed.name) details.push(`name: "${parsed.name}"`);
      if (parsed.price) details.push(`price: ${parsed.price}`);
      if (parsed.total) details.push(`total: ${parsed.total}`);
      if (parsed.table_number) details.push(`table: "${parsed.table_number}"`);
      if (parsed.role) details.push(`role: "${parsed.role}"`);
      if (parsed.notes) details.push(`notes: "${parsed.notes}"`);
      if (parsed.reason) details.push(`reason: "${parsed.reason}"`);
      
      if (details.length > 0) {
        detailStr = `<span class="text-slate-400 font-semibold">[${details.join(', ')}]</span>`;
      } else {
        const sliced = log.new_value.slice(0, 80);
        detailStr = `<span class="text-slate-500">${sliced}${log.new_value.length > 80 ? '...' : ''}</span>`;
      }
    } catch (e) {
      if (log.new_value) {
        const sliced = log.new_value.slice(0, 80);
        detailStr = `<span class="text-slate-500">${sliced}${log.new_value.length > 80 ? '...' : ''}</span>`;
      }
    }

    const shortUserId = log.user_id ? log.user_id.slice(0, 8) : 'system';

    return `
      <div class="leading-relaxed hover:bg-slate-800/30 px-2 py-0.5 rounded transition-colors whitespace-nowrap">
        <span class="text-[#27c93f]">[${dateStr} ${timeStr}]</span>
        <span class="text-blue-400 font-semibold">$ ${log.entity}</span>
        <span class="${actionColor} font-bold">${log.action}</span>
        ${log.entity_id && log.entity_id !== '00000000-0000-0000-0000-000000000000' ? `<span class="text-purple-400 font-semibold">(${log.entity_id.slice(0, 8)})</span>` : ''}
        ${detailStr}
        <span class="text-slate-500 text-[10px] float-right ml-4">by: ${shortUserId}</span>
      </div>
    `;
  }).join('');
}

// ── Update Best Sellers Table ───────────────────────────────
function updateBestSellersTable(data) {
  const tbody = $('best-sellers-body');
  if (!tbody) return;

  if (!Array.isArray(data) || data.length === 0) {
    tbody.innerHTML = `
      <tr>
        <td colspan="6" class="py-8 text-center text-[var(--text-secondary)]">
          <svg class="w-12 h-12 mx-auto mb-3 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"/>
          </svg>
          <p class="text-sm font-medium">No sales data available</p>
          <p class="text-xs mt-1">Try adjusting the date range</p>
        </td>
      </tr>
    `;
    return;
  }

  const maxSales = Math.max(...data.map(d => d.total_sales || 0));
  const rankColors = ['bg-green-500', 'bg-amber-500', 'bg-blue-500', 'bg-gray-400', 'bg-gray-400'];

  tbody.innerHTML = data.map((item, index) => {
    const barWidth = maxSales > 0 ? (item.total_sales / maxSales) * 100 : 0;
    const rankColor = rankColors[index] || 'bg-gray-400';

    return `
      <tr class="border-b border-gray-100 hover:bg-gray-50/50 transition-colors">
        <td class="py-3 px-4">
          <span class="inline-flex items-center justify-center w-7 h-7 rounded-lg ${rankColor} text-white text-xs font-bold">
            #${index + 1}
          </span>
        </td>
        <td class="py-3 px-4">
          <div class="font-semibold text-gray-900 text-sm">${item.menu_name || 'Unknown'}</div>
        </td>
        <td class="py-3 px-4">
          <span class="inline-flex items-center px-2.5 py-1 rounded-full text-xs font-medium bg-gray-100 text-gray-700">
            ${item.category || 'Menu'}
          </span>
        </td>
        <td class="py-3 px-4 text-sm font-semibold text-gray-700">${formatNumber(item.total_quantity || 0)}</td>
        <td class="py-3 px-4 text-sm font-bold text-gray-900">${formatRupiah(item.total_sales)}</td>
        <td class="py-3 px-4">
          <div class="w-24 h-2 bg-gray-100 rounded-full overflow-hidden">
            <div class="h-full ${rankColor} rounded-full transition-all duration-500" style="width:${barWidth}%"></div>
          </div>
        </td>
      </tr>
    `;
  }).join('');
}

// ── Period Toggle ───────────────────────────────────────────
function setPeriod(period) {
  currentPeriod = period;
  dateRange = { start: '', end: '' }; // Reset custom range

  // Update button styles
  ['7D', '30D', '90D'].forEach(p => {
    const btn = $(`btn-${p.toLowerCase()}`);
    if (btn) {
      btn.classList.toggle('active', p === period);
    }
  });

  // Update date range label
  const range = getDateRange(period);
  const label = $(`date-range-label`);
  if (label) {
    label.textContent = `${formatDateShort(range.start)} - ${formatDateShort(range.end)}`;
  }

  fetchDashboardData();
}

// ── Date Range Picker ─────────────────────────────────────
function initDateRangePicker() {
  const pickerEl = $('date-range-picker');
  if (!pickerEl || typeof Litepicker === 'undefined') return;

  const today = new Date();
  const lastWeek = new Date();
  lastWeek.setDate(today.getDate() - 7);

  const picker = new Litepicker({
    element: pickerEl,
    singleMode: false,
    startDate: lastWeek,
    endDate: today,
    format: 'YYYY-MM-DD',
    numberOfMonths: 2,
    numberOfColumns: 2,
    tooltipText: {
      one: 'day',
      other: 'days',
    },
    setup: (picker) => {
      picker.on('selected', (date1, date2) => {
        const start = date1.format('YYYY-MM-DD');
        const end = date2.format('YYYY-MM-DD');
        dateRange = { start, end };
        currentPeriod = 'CUSTOM';

        // Reset period buttons
        ['7D', '30D', '90D'].forEach(p => {
          const btn = $(`btn-${p.toLowerCase()}`);
          if (btn) btn.classList.remove('active');
        });

        // Update label
        const label = $('date-range-label');
        if (label) {
          label.textContent = `${formatDateShort(start)} - ${formatDateShort(end)}`;
        }

        fetchDashboardData();
      });
    },
  });
}

// ── Toast Notifications ─────────────────────────────────────
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
  }, 4000);
}

// ── Sidebar Toggle (Mobile) ─────────────────────────────────
function toggleSidebar() {
  const sidebar = $('sidebar');
  const overlay = $('sidebar-overlay');
  const isOpen = !sidebar.classList.contains('-translate-x-full');

  if (isOpen) {
    sidebar.classList.add('-translate-x-full');
    overlay.classList.add('hidden');
  } else {
    sidebar.classList.remove('-translate-x-full');
    overlay.classList.remove('hidden');
  }
}

// ── Dark Mode Toggle ────────────────────────────────────────
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

// ── Setup Event Listeners ───────────────────────────────────
function setupEventListeners() {
  // Sidebar toggle
  const toggleBtn = $('btn-toggle-sidebar');
  if (toggleBtn) toggleBtn.addEventListener('click', toggleSidebar);

  const overlay = $('sidebar-overlay');
  if (overlay) overlay.addEventListener('click', toggleSidebar);

  // Dark mode
  const darkModeBtn = $('btn-dark-mode');
  if (darkModeBtn) darkModeBtn.addEventListener('click', toggleDarkMode);

  // Logout
  const logoutBtn = $('btn-logout');
  if (logoutBtn) logoutBtn.addEventListener('click', logout);

  // Period buttons
  const btn7d = $('btn-7d');
  const btn30d = $('btn-30d');
  const btn90d = $('btn-90d');
  if (btn7d) btn7d.addEventListener('click', () => setPeriod('7D'));
  if (btn30d) btn30d.addEventListener('click', () => setPeriod('30D'));
  if (btn90d) btn90d.addEventListener('click', () => setPeriod('90D'));

  // Refresh button
  const refreshBtn = $('refresh-btn');
  if (refreshBtn) refreshBtn.addEventListener('click', fetchDashboardData);

  // Refresh Audit Log button and Limit Slider
  const refreshAuditBtn = $('btn-refresh-audit');
  const auditLimitSlider = $('audit-limit-slider');
  const auditLimitVal = $('audit-limit-val');

  if (auditLimitSlider && auditLimitVal) {
    auditLimitSlider.addEventListener('input', (e) => {
      auditLogLimit = parseInt(e.target.value);
      auditLimitVal.textContent = auditLogLimit;
    });
    auditLimitSlider.addEventListener('change', () => {
      fetchDashboardData();
    });
  }

  if (refreshAuditBtn) {
    refreshAuditBtn.addEventListener('click', async (e) => {
      e.preventDefault();
      try {
        refreshAuditBtn.innerHTML = '🔄 Loading...';
        const res = await apiGet(`/api/v1/analytics/audit-logs?limit=${auditLogLimit}`);
        renderAuditLogs(res.data || []);
      } catch (err) {
        console.error(err);
      } finally {
        refreshAuditBtn.innerHTML = '🔄 Refresh';
      }
    });
  }

  // View all button
  const viewAllBtn = $('btn-view-all');
  if (viewAllBtn) {
    viewAllBtn.addEventListener('click', () => {
      switchTab('sales');
    });
  }

  // Sidebar links
  document.querySelectorAll('.sidebar-link[data-page]').forEach(link => {
    link.addEventListener('click', (e) => {
      e.preventDefault();
      const page = link.getAttribute('data-page');
      switchTab(page, link);
    });
  });
}

// ── Dynamic Tab Switching & Loaders ──────────────────────────
async function switchTab(page, activeLink) {
  // Update sidebar active states
  if (!activeLink) {
    activeLink = document.querySelector(`.sidebar-link[data-page="${page}"]`);
  }
  document.querySelectorAll('.sidebar-link[data-page]').forEach(link => {
    link.classList.remove('active');
    const badge = link.querySelector('span.rounded-full');
    if (badge) badge.remove();
  });
  if (activeLink) {
    activeLink.classList.add('active');
    // Add small green active dot
    const dot = document.createElement('span');
    dot.className = 'ml-auto w-2 h-2 rounded-full bg-green-400';
    activeLink.appendChild(dot);
  }

  // Hide all sections
  const sections = ['dashboard', 'sales', 'inventory', 'orders', 'shifts', 'employees', 'menus'];
  sections.forEach(s => {
    const el = $(`${s}-section`);
    if (el) el.classList.add('hidden');
  });

  // Show target section
  const targetSection = $(`${page}-section`);
  if (targetSection) {
    targetSection.classList.remove('hidden');
  }

  // Load target data
  try {
    if (page === 'dashboard') {
      await fetchDashboardData();
    } else if (page === 'sales') {
      await loadSalesSection();
    } else if (page === 'inventory') {
      await loadInventorySection();
    } else if (page === 'orders') {
      await loadOrdersSection();
    } else if (page === 'shifts') {
      await loadShiftsSection();
    } else if (page === 'employees') {
      await loadEmployeesSection();
    } else if (page === 'menus') {
      await loadMenusSection();
    }
  } catch (error) {
    console.error(`Failed to load page: ${page}`, error);
    showToast(`Error loading ${page} data`, 'error');
  }
}

// 1. Sales Section Loader
async function loadSalesSection() {
  const perfResponse = await fetchAPI('/analytics/cashier-performance');
  const voidResponse = await fetchAPI('/analytics/void-return-logs');

  // Render Cashier Performance
  const perfBody = $('cashier-performance-body');
  if (perfBody) {
    perfBody.innerHTML = '';
    const perfData = perfResponse.data || [];
    if (perfData.length === 0) {
      perfBody.innerHTML = `<tr><td colspan="4" class="py-4 text-center text-xs text-[var(--text-muted)]">No performance data found</td></tr>`;
    } else {
      perfData.forEach(row => {
        perfBody.innerHTML += `
          <tr class="border-b border-[var(--border-color)] text-sm">
            <td class="py-3 px-3 font-semibold text-[var(--text-primary)]">${row.cashier_name || 'Unknown'}</td>
            <td class="py-3 px-3 text-[var(--text-secondary)]">${row.total_orders || 0} orders</td>
            <td class="py-3 px-3 font-bold text-[var(--text-primary)]">${formatRupiah(row.total_sales || 0)}</td>
            <td class="py-3 px-3 text-[var(--text-secondary)]">${Number(row.void_rate_pct || 0).toFixed(1)}%</td>
          </tr>
        `;
      });
    }
  }

  // Render Void & Return Logs
  const voidBody = $('void-return-body');
  if (voidBody) {
    voidBody.innerHTML = '';
    const logs = voidResponse.data || [];
    if (logs.length === 0) {
      voidBody.innerHTML = `<tr><td colspan="4" class="py-4 text-center text-xs text-[var(--text-muted)]">No void/return audits recorded</td></tr>`;
    } else {
      logs.forEach(row => {
        const typeBadge = row.event_type === 'return' 
          ? `<span class="px-2 py-0.5 text-xs font-bold rounded-full bg-amber-100 text-amber-800">RETURN</span>`
          : `<span class="px-2 py-0.5 text-xs font-bold rounded-full bg-red-100 text-red-800">VOID</span>`;
        voidBody.innerHTML += `
          <tr class="border-b border-[var(--border-color)] text-sm">
            <td class="py-3 px-3 text-[var(--text-secondary)]">${row.event_date ? row.event_date.split('T')[0] : 'N/A'}</td>
            <td class="py-3 px-3">${typeBadge}</td>
            <td class="py-3 px-3 font-bold text-red-500">${formatRupiah(row.amount || 0)}</td>
            <td class="py-3 px-3 text-[var(--text-secondary)] italic">"${row.reason || 'No reason provided'}"</td>
          </tr>
        `;
      });
    }
  }
}

// 2. Inventory Section Loader
let rawMaterialsCache = [];
async function loadInventorySection() {
  const response = await fetchAPI('/inventory/raw-materials');
  const materials = response.data || [];
  rawMaterialsCache = materials;

  const invBody = $('inventory-body');
  if (invBody) {
    invBody.innerHTML = '';
    if (materials.length === 0) {
      invBody.innerHTML = `<tr><td colspan="7" class="py-4 text-center text-xs text-[var(--text-muted)]">No raw materials seeded</td></tr>`;
    } else {
      materials.forEach(rm => {
        const isLow = rm.current_stock <= rm.min_stock_level;
        const statusBadge = isLow
          ? `<span class="px-2.5 py-0.5 text-xs font-bold rounded-full bg-red-100 text-red-800">⚠️ LOW STOCK</span>`
          : `<span class="px-2.5 py-0.5 text-xs font-bold rounded-full bg-green-100 text-green-800">✓ NORMAL</span>`;
        
        invBody.innerHTML += `
          <tr class="border-b border-[var(--border-color)] text-sm">
            <td class="py-3.5 px-4 font-mono text-[var(--text-muted)] text-xs">${rm.id.substring(0, 8)}...</td>
            <td class="py-3.5 px-4 font-semibold text-[var(--text-primary)]">${rm.name}</td>
            <td class="py-3.5 px-4 font-bold ${isLow ? 'text-red-500' : 'text-[var(--text-primary)]'}">${formatNumber(rm.current_stock)}</td>
            <td class="py-3.5 px-4 text-[var(--text-secondary)]">${rm.unit}</td>
            <td class="py-3.5 px-4 font-medium text-[var(--text-secondary)]">${formatRupiah(rm.cost_per_unit)}</td>
            <td class="py-3.5 px-4">${statusBadge}</td>
            <td class="py-3.5 px-4 text-right">
              <div class="flex items-center justify-end gap-2">
                <button onclick="editRawMaterial('${rm.id}')" class="p-1 text-blue-600 hover:text-blue-800 hover:bg-blue-50 rounded transition-colors" title="Edit">
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/></svg>
                </button>
                <button onclick="deleteRawMaterial('${rm.id}')" class="p-1 text-red-600 hover:text-red-800 hover:bg-red-50 rounded transition-colors" title="Delete">
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/></svg>
                </button>
              </div>
            </td>
          </tr>
        `;
      });
    }
  }

  // Update Stock Adjustment Dropdown
  const select = $('adjust-material-id');
  if (select) {
    select.innerHTML = materials.map(m => `<option value="${m.id}">${m.name} (${m.current_stock} ${m.unit})</option>`).join('');
  }
}

function showRawMaterialModal() {
  const modal = $('material-modal');
  const title = $('material-modal-title');
  const form = $('material-form');
  
  if (modal) {
    title.textContent = '📦 Add Raw Material';
    form.reset();
    $('material-id').value = '';
    modal.classList.remove('hidden');
  }
}

function closeRawMaterialModal() {
  const modal = $('material-modal');
  if (modal) modal.classList.add('hidden');
}

function editRawMaterial(id) {
  const rm = rawMaterialsCache.find(x => x.id === id);
  if (!rm) return;

  const modal = $('material-modal');
  const title = $('material-modal-title');
  
  if (modal) {
    title.textContent = '✏️ Edit Raw Material';
    $('material-id').value = rm.id;
    $('material-name').value = rm.name;
    $('material-unit').value = rm.unit;
    $('material-stock').value = rm.current_stock;
    $('material-min-stock').value = rm.min_stock_level;
    $('material-cost').value = rm.cost_per_unit;
    modal.classList.remove('hidden');
  }
}

async function deleteRawMaterial(id) {
  if (!confirm('Are you sure you want to delete this raw material? Recipes referencing this material will also be impacted.')) return;
  try {
    await fetchAPI(`/inventory/raw-materials/${id}`, 'DELETE');
    showToast('Raw material deleted successfully', 'success');
    await loadInventorySection();
  } catch (err) {
    showToast('Failed to delete raw material: ' + err.message, 'error');
  }
}

async function handleRawMaterialSubmit(e) {
  e.preventDefault();
  const id = $('material-id').value;
  const payload = {
    name: $('material-name').value,
    unit: $('material-unit').value,
    current_stock: parseFloat($('material-stock').value),
    min_stock_level: parseFloat($('material-min-stock').value),
    cost_per_unit: parseFloat($('material-cost').value)
  };

  try {
    if (id) {
      await fetchAPI(`/inventory/raw-materials/${id}`, 'PUT', payload);
      showToast('Raw material updated successfully', 'success');
    } else {
      await fetchAPI('/inventory/raw-materials', 'POST', payload);
      showToast('Raw material added successfully', 'success');
    }
    closeRawMaterialModal();
    await loadInventorySection();
  } catch (err) {
    showToast('Failed to save raw material: ' + err.message, 'error');
  }
}

function showAdjustStockModal() {
  const modal = $('adjust-stock-modal');
  if (modal) modal.classList.remove('hidden');
}

function closeAdjustStockModal() {
  const modal = $('adjust-stock-modal');
  if (modal) modal.classList.add('hidden');
  $('adjust-stock-form')?.reset();
}

async function handleAdjustStock(e) {
  e.preventDefault();
  const materialId = $('adjust-material-id').value;
  const adjustType = $('adjust-type').value;
  const qty = parseFloat($('adjust-qty').value);
  const notes = $('adjust-notes').value;

  if (isNaN(qty) || qty <= 0) {
    showToast('Please enter a valid quantity', 'error');
    return;
  }

  // Find current stock
  const currentItem = rawMaterialsCache.find(m => m.id === materialId);
  if (!currentItem) return;

  let finalQty = currentItem.current_stock;
  if (adjustType === 'received') {
    finalQty += qty;
  } else if (adjustType === 'damaged') {
    finalQty = Math.max(0, finalQty - qty);
  } else if (adjustType === 'audit') {
    finalQty = qty;
  }

  try {
    const response = await fetchAPI('/stock/adjust', 'POST', {
      raw_material_id: materialId,
      quantity_after: finalQty,
      reason: `${adjustType.toUpperCase()} - ${notes}`
    });

    if (response.success) {
      showToast('Stock adjusted successfully!', 'success');
      closeAdjustStockModal();
      await loadInventorySection();
    } else {
      showToast(response.message || 'Failed to adjust stock', 'error');
    }
  } catch (err) {
    showToast('Failed to adjust stock due to network error', 'error');
  }
}

// 3. Orders Section Loader
let ordersCache = [];
async function loadOrdersSection() {
  const response = await fetchAPI('/orders');
  ordersCache = response.data || [];
  renderOrdersTable(ordersCache);
}

function renderOrdersTable(orders) {
  const body = $('orders-body');
  if (!body) return;

  body.innerHTML = '';
  if (orders.length === 0) {
    body.innerHTML = `<tr><td colspan="6" class="py-4 text-center text-xs text-[var(--text-muted)]">No orders found</td></tr>`;
    return;
  }

  orders.forEach(o => {
    let statusClass = 'bg-gray-100 text-gray-800';
    if (o.status === 'confirmed') statusClass = 'bg-blue-100 text-blue-800';
    if (o.status === 'paid') statusClass = 'bg-green-100 text-green-800';
    if (o.status === 'served') statusClass = 'bg-purple-100 text-purple-800';
    if (o.status === 'cancelled') statusClass = 'bg-red-100 text-red-800';

    const dateStr = o.created_at ? o.created_at.split('T')[0] : 'N/A';
    
    // Actions: Cancel is only available for active unpaid/unserved orders
    const canCancel = o.status !== 'paid' && o.status !== 'served' && o.status !== 'cancelled';
    const actionButton = canCancel 
      ? `<button onclick="cancelOrder('${o.id}')" class="kt-btn text-xs bg-red-500 text-white py-1 px-2.5 rounded-md hover:bg-red-600 transition">Cancel</button>`
      : `<span class="text-xs text-[var(--text-muted)] italic">No actions available</span>`;

    body.innerHTML += `
      <tr class="border-b border-[var(--border-color)] text-sm">
        <td class="py-3.5 px-4 font-semibold text-[var(--text-primary)]">#${o.order_number}</td>
        <td class="py-3.5 px-4 text-xs font-bold text-[var(--text-secondary)] uppercase">${o.order_type.replace('_', ' ')}</td>
        <td class="py-3.5 px-4 font-bold text-[var(--text-primary)]">${formatRupiah(o.total)}</td>
        <td class="py-3.5 px-4"><span class="px-2 py-0.5 text-xs font-bold rounded-full ${statusClass}">${o.status.toUpperCase()}</span></td>
        <td class="py-3.5 px-4 text-[var(--text-secondary)] text-xs">${dateStr}</td>
        <td class="py-3.5 px-4">${actionButton}</td>
      </tr>
    `;
  });
}

function searchOrders() {
  const query = $('order-search').value.toLowerCase().trim();
  if (!query) {
    renderOrdersTable(ordersCache);
    return;
  }
  const filtered = ordersCache.filter(o => o.order_number.toLowerCase().includes(query));
  renderOrdersTable(filtered);
}

async function cancelOrder(orderId) {
  if (!confirm('Are you sure you want to cancel and void this order?')) return;
  try {
    const response = await fetchAPI(`/orders/${orderId}/status`, 'PUT', { status: 'cancelled' });
    if (response.success) {
      showToast('Order cancelled successfully', 'success');
      await loadOrdersSection();
    } else {
      showToast(response.message || 'Failed to cancel order', 'error');
    }
  } catch (err) {
    showToast('Failed to cancel order due to server error', 'error');
  }
}

// 4. Shifts Section Loader
async function loadShiftsSection() {
  const response = await fetchAPI('/analytics/shift-reconciliation');
  const shifts = response.data || [];

  const body = $('shifts-body');
  if (body) {
    body.innerHTML = '';
    if (shifts.length === 0) {
      body.innerHTML = `<tr><td colspan="7" class="py-4 text-center text-xs text-[var(--text-muted)]">No shifts history found</td></tr>`;
    } else {
      shifts.forEach(s => {
        const disc = s.cash_difference || 0;
        const diffClass = disc === 0 
          ? 'text-green-500 font-bold' 
          : disc > 0 ? 'text-green-600 font-bold' : 'text-red-500 font-bold';

        const statusClass = s.status === 'open' 
          ? 'bg-blue-100 text-blue-800' 
          : 'bg-gray-100 text-gray-800';

        body.innerHTML += `
          <tr class="border-b border-[var(--border-color)] text-sm">
            <td class="py-3.5 px-4 font-mono text-[var(--text-muted)] text-xs">${s.shift_id.substring(0, 8)}...</td>
            <td class="py-3.5 px-4 font-semibold text-[var(--text-primary)]">${s.cashier_name || 'Cashier'}</td>
            <td class="py-3.5 px-4 text-[var(--text-secondary)]">${formatRupiah(s.opening_cash)}</td>
            <td class="py-3.5 px-4 text-[var(--text-secondary)]">${formatRupiah(s.expected_closing_cash)}</td>
            <td class="py-3.5 px-4 font-bold text-[var(--text-primary)]">${formatRupiah(s.actual_closing_cash)}</td>
            <td class="py-3.5 px-4 ${diffClass}">${disc >= 0 ? '+' : ''}${formatRupiah(disc)}</td>
            <td class="py-3.5 px-4"><span class="px-2 py-0.5 text-xs font-bold rounded-full ${statusClass}">${s.status.toUpperCase()}</span></td>
          </tr>
        `;
      });
    }
  }
}

// 5. Employees Section Loader & CRUD
let employeesCache = [];
async function loadEmployeesSection() {
  const response = await fetchAPI('/employees');
  employeesCache = response.data || [];

  const body = $('employees-body');
  if (body) {
    body.innerHTML = '';
    if (employeesCache.length === 0) {
      body.innerHTML = `<tr><td colspan="5" class="py-4 text-center text-xs text-[var(--text-muted)]">No employees registered</td></tr>`;
    } else {
      employeesCache.forEach(emp => {
        const statusBadge = emp.is_active
          ? `<span class="px-2 py-0.5 text-xs font-bold rounded-full bg-green-100 text-green-800">ACTIVE</span>`
          : `<span class="px-2 py-0.5 text-xs font-bold rounded-full bg-gray-100 text-gray-800">INACTIVE</span>`;
        
        body.innerHTML += `
          <tr class="border-b border-[var(--border-color)] text-sm">
            <td class="py-3.5 px-4 font-semibold text-[var(--text-primary)]">${emp.name}</td>
            <td class="py-3.5 px-4 text-[var(--text-secondary)]">${emp.email || 'No email'}</td>
            <td class="py-3.5 px-4 font-bold text-xs uppercase text-[var(--text-secondary)]">${emp.role}</td>
            <td class="py-3.5 px-4">${statusBadge}</td>
            <td class="py-3.5 px-4 flex items-center gap-2">
              <button onclick="editEmployee('${emp.id}')" class="kt-btn text-xs bg-green-600 text-white py-1 px-2.5 rounded-md hover:bg-green-700 transition">Edit</button>
              <button onclick="deleteEmployee('${emp.id}')" class="kt-btn text-xs bg-red-500 text-white py-1 px-2.5 rounded-md hover:bg-red-600 transition">Delete</button>
            </td>
          </tr>
        `;
      });
    }
  }
}

function showEmployeeModal(emp = null) {
  const title = $('employee-modal-title');
  const modal = $('employee-modal');
  if (!modal) return;

  if (emp) {
    title.textContent = '👤 Edit Employee';
    $('employee-id').value = emp.id;
    $('employee-name').value = emp.name;
    $('employee-email').value = emp.email || '';
    $('employee-email').disabled = true; // Email unique constraint key
    $('employee-password').value = '';
    $('employee-pin').value = '';
    $('employee-pin').required = false; // PIN is hashed, skip if unchanged
    $('employee-role').value = emp.role;
    $('employee-active').checked = emp.is_active;
  } else {
    title.textContent = '👤 Add Employee';
    $('employee-id').value = '';
    $('employee-name').value = '';
    $('employee-email').value = '';
    $('employee-email').disabled = false;
    $('employee-password').value = '';
    $('employee-password').required = true;
    $('employee-pin').value = '';
    $('employee-pin').required = true;
    $('employee-role').value = 'cashier';
    $('employee-active').checked = true;
  }
  modal.classList.remove('hidden');
}

function closeEmployeeModal() {
  $('employee-modal')?.classList.add('hidden');
  $('employee-form')?.reset();
}

async function handleEmployeeSubmit(e) {
  e.preventDefault();
  const id = $('employee-id').value;
  const name = $('employee-name').value;
  const email = $('employee-email').value;
  const password = $('employee-password').value;
  const pin = $('employee-pin').value;
  const role = $('employee-role').value;
  const is_active = $('employee-active').checked;

  const user = JSON.parse(localStorage.getItem(CONFIG.USER_KEY) || '{}');
  const branch_id = user.branch_id;

  if (!branch_id) {
    showToast('Branch context missing', 'error');
    return;
  }

  const payload = {
    name,
    email,
    role,
    is_active
  };

  if (password) payload.password = password;
  if (pin) payload.pin_code = pin;

  try {
    let response;
    if (id) {
      response = await fetchAPI(`/employees/${id}`, 'PUT', payload);
    } else {
      payload.branch_id = branch_id;
      response = await fetchAPI('/employees', 'POST', payload);
    }

    if (response.success) {
      showToast(id ? 'Employee updated successfully' : 'Employee registered successfully', 'success');
      closeEmployeeModal();
      await loadEmployeesSection();
    } else {
      showToast(response.message || response.error || 'Operation failed', 'error');
    }
  } catch (err) {
    showToast('Failed to save employee due to server error', 'error');
  }
}

function editEmployee(id) {
  const emp = employeesCache.find(e => e.id === id);
  if (emp) showEmployeeModal(emp);
}

async function deleteEmployee(id) {
  if (!confirm('Are you sure you want to completely delete this employee?')) return;
  try {
    const response = await fetchAPI(`/employees/${id}`, 'DELETE');
    if (response.success) {
      showToast('Employee deleted successfully', 'success');
      await loadEmployeesSection();
    } else {
      showToast(response.message || 'Failed to delete employee', 'error');
    }
  } catch (err) {
    showToast('Network error deleting employee', 'error');
  }
}

// ── Initialization ──────────────────────────────────────────
document.addEventListener('DOMContentLoaded', async () => {
  await loadConfig();
  initDarkMode();
  setupEventListeners();

  if (!checkAuth()) return;

  // Initialize date range picker
  initDateRangePicker();

  // Set initial period and fetch
  setPeriod('7D');

  // Auto refresh every 5 minutes
  autoRefreshInterval = setInterval(fetchDashboardData, 300000);
});

// Cleanup on page unload
window.addEventListener('beforeunload', () => {
  if (autoRefreshInterval) {
    clearInterval(autoRefreshInterval);
  }
});

// ── Menu Management (CRUD) ──────────────────────────────────
let allMenus = [];

async function loadMenusSection() {
  try {
    const res = await fetchAPI('/menus/all');
    allMenus = res.data || [];
    renderMenusTable(allMenus);
  } catch (err) {
    console.error('Failed to fetch menus for CRUD:', err);
    showToast('Failed to load menus', 'error');
  }
}

function renderMenusTable(items) {
  const tbody = $('menus-body');
  if (!tbody) return;

  if (!items || items.length === 0) {
    tbody.innerHTML = `
      <tr>
        <td colspan="7" class="py-8 text-center text-[var(--text-secondary)] text-sm">
          No menu items found. Click 'Add Menu Item' to create one.
        </td>
      </tr>
    `;
    return;
  }

  tbody.innerHTML = items.map(item => `
    <tr class="border-b border-[var(--border-color)] hover:bg-[var(--bg-card)]/50 transition-colors">
      <td class="py-3 px-4 text-sm font-semibold text-[var(--text-primary)]">${item.name}</td>
      <td class="py-3 px-4 text-sm font-bold text-kopitiam-600">${formatRupiah(item.price)}</td>
      <td class="py-3 px-4 text-sm text-[var(--text-secondary)]">${item.daily_stock}</td>
      <td class="py-3 px-4 text-sm text-[var(--text-secondary)] capitalize">${item.station || 'coffee'}</td>
      <td class="py-3 px-4 text-sm">
        <span class="px-2 py-0.5 text-xs font-semibold rounded-full ${item.is_recipe_based ? 'bg-indigo-500/10 text-indigo-400' : 'bg-gray-500/10 text-gray-400'}">
          ${item.is_recipe_based ? 'Yes' : 'No'}
        </span>
      </td>
      <td class="py-3 px-4 text-sm">
        <span class="px-2 py-0.5 text-xs font-semibold rounded-full ${item.is_active ? 'bg-green-500/10 text-green-400' : 'bg-red-500/10 text-red-400'}">
          ${item.is_active ? 'Active' : 'Inactive'}
        </span>
      </td>
      <td class="py-3 px-4 text-sm">
        <div class="flex items-center gap-2">
          <button onclick="editMenu('${item.id}')" class="text-blue-400 hover:text-blue-300 transition-colors p-1" title="Edit Item">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/></svg>
          </button>
          <button onclick="deleteMenu('${item.id}')" class="text-red-400 hover:text-red-300 transition-colors p-1" title="Delete Item">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/></svg>
          </button>
        </div>
      </td>
    </tr>
  `).join('');
}

function showMenuModal(id = null) {
  const modal = $('menu-modal');
  const title = $('menu-modal-title');
  const form = $('menu-form');
  
  if (!modal || !title || !form) return;

  form.reset();
  $('menu-id').value = '';
  title.textContent = '🍽️ Add Menu Item';

  if (id) {
    const item = allMenus.find(m => String(m.id) === String(id));
    if (item) {
      title.textContent = '✏️ Edit Menu Item';
      $('menu-id').value = item.id;
      $('menu-name').value = item.name;
      $('menu-price').value = item.price;
      $('menu-stock').value = item.daily_stock;
      $('menu-station').value = item.station || 'bar';
      $('menu-recipe-based').checked = item.is_recipe_based;
      $('menu-active').checked = item.is_active;
    }
  }

  modal.classList.remove('hidden');
}

function closeMenuModal() {
  const modal = $('menu-modal');
  if (modal) modal.classList.add('hidden');
}

async function handleMenuSubmit(e) {
  e.preventDefault();

  const id = $('menu-id').value;
  const name = $('menu-name').value.trim();
  const price = parseFloat($('menu-price').value);
  const daily_stock = parseInt($('menu-stock').value) || 0;
  const station = $('menu-station').value;
  const is_recipe_based = $('menu-recipe-based').checked;
  const is_active = $('menu-active').checked;

  const payload = {
    name,
    price,
    daily_stock,
    station,
    is_recipe_based,
    is_active
  };

  try {
    let res;
    if (id) {
      // Edit
      res = await fetchAPI(`/menus/${id}`, 'PUT', payload);
    } else {
      // Create
      res = await fetchAPI('/menus', 'POST', payload);
    }

    if (res.success) {
      showToast(id ? 'Menu updated successfully' : 'Menu created successfully', 'success');
      closeMenuModal();
      await loadMenusSection();
    } else {
      showToast(res.message || 'Failed to save menu item', 'error');
    }
  } catch (err) {
    console.error('Failed to submit menu item:', err);
    showToast(err.message || 'Server error saving menu item', 'error');
  }
}

async function deleteMenu(id) {
  if (!confirm('Are you sure you want to delete this menu item?')) return;

  try {
    const res = await fetchAPI(`/menus/${id}`, 'DELETE');
    if (res.success) {
      showToast('Menu item deleted successfully', 'success');
      await loadMenusSection();
    } else {
      showToast(res.message || 'Failed to delete menu item', 'error');
    }
  } catch (err) {
    console.error('Failed to delete menu item:', err);
    showToast('Server error deleting menu item', 'error');
  }
}

// Expose all functions to global scope for HTML inline calls
window.logout = logout;
window.setPeriod = setPeriod;
window.toggleSidebar = toggleSidebar;
window.toggleDarkMode = toggleDarkMode;
window.fetchDashboardData = fetchDashboardData;
window.switchTab = switchTab;
window.showAdjustStockModal = showAdjustStockModal;
window.closeAdjustStockModal = closeAdjustStockModal;
window.handleAdjustStock = handleAdjustStock;
window.searchOrders = searchOrders;
window.cancelOrder = cancelOrder;
window.showEmployeeModal = showEmployeeModal;
window.closeEmployeeModal = closeEmployeeModal;
window.handleEmployeeSubmit = handleEmployeeSubmit;
window.editEmployee = editEmployee;
window.deleteEmployee = deleteEmployee;
window.showMenuModal = showMenuModal;
window.closeMenuModal = closeMenuModal;
window.handleMenuSubmit = handleMenuSubmit;
window.editMenu = showMenuModal;
window.deleteMenu = deleteMenu;
