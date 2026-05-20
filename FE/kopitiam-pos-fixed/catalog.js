/**
 * ============================================================
 * KopiTiam POS — Digital Self-Order Catalog (catalog.js)
 * Customer-facing tablet interface with Midtrans Snap sandbox
 * ============================================================
 */

const CONFIG = {
  // Use relative path since API is served on the same host under Basic Auth Group
  API_BASE_URL: '/catalog/api'
};

let menuItems = [];
let cart = [];
let selectedCategory = 'all';
let selectedPaymentMethod = 'online';
let categories = [];

const $ = (id) => document.getElementById(id);

// Format pricing helper
function formatRupiah(num) {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0,
    maximumFractionDigits: 0
  }).format(num);
}

// Fetch menus and initialize catalog
async function initCatalog() {
  showLoading("Memuat daftar menu...");
  try {
    const res = await fetch(`${CONFIG.API_BASE_URL}/menus`);
    const data = await res.json();
    if (!res.ok) throw new Error(data.message || "Gagal memuat menu");

    menuItems = data.data || [];
    
    // Extract unique categories from active menus
    const categorySet = new Set();
    menuItems.forEach(item => {
      if (item.Categories) {
        item.Categories.forEach(cat => categorySet.add(JSON.stringify(cat)));
      }
    });

    categories = Array.from(categorySet).map(s => JSON.parse(s));
    
    renderCategories();
    renderMenuItems();
    updateCartUI();
  } catch (err) {
    alert("Katalog error: " + err.message);
  } finally {
    hideLoading();
  }
}

// Render horizontal categories buttons
function renderCategories() {
  const container = $('categories-container');
  container.innerHTML = `
    <button class="px-4 py-2 rounded-xl text-sm font-medium transition cursor-pointer ${selectedCategory === 'all' ? 'bg-kopitiam-500 text-white shadow-sm' : 'bg-white text-gray-600 border border-gray-100 hover:bg-gray-50'}" id="cat-all">
        Semua Menu
    </button>
  `;

  // Attach event listener for "Semua Menu"
  $('cat-all').addEventListener('click', () => {
    selectedCategory = 'all';
    renderCategories();
    renderMenuItems();
  });

  categories.forEach(cat => {
    const isSelected = selectedCategory === cat.id;
    const btn = document.createElement('button');
    btn.className = `px-4 py-2 rounded-xl text-sm font-medium transition cursor-pointer ${isSelected ? 'bg-kopitiam-500 text-white shadow-sm' : 'bg-white text-gray-600 border border-gray-100 hover:bg-gray-50'}`;
    btn.textContent = cat.name;
    btn.addEventListener('click', () => {
      selectedCategory = cat.id;
      renderCategories();
      renderMenuItems();
    });
    container.appendChild(btn);
  });
}

// Render main menu grid cards list
function renderMenuItems() {
  const container = $('menu-items-grid');
  container.innerHTML = '';

  // Filter items
  const filtered = selectedCategory === 'all' 
    ? menuItems 
    : menuItems.filter(item => item.Categories && item.Categories.some(cat => cat.id === selectedCategory));

  if (filtered.length === 0) {
    container.innerHTML = `
      <div class="col-span-full py-16 text-center text-gray-400 text-sm">
        Tidak ada menu yang tersedia untuk kategori ini.
      </div>
    `;
    return;
  }

  filtered.forEach(item => {
    const card = document.createElement('div');
    card.className = "menu-card rounded-2xl overflow-hidden cursor-pointer flex flex-col h-full bg-white relative";
    
    // Best seller badge
    const badgeHTML = item.is_best_seller 
      ? `<span class="absolute top-3 left-3 bestseller-badge px-3 py-1 rounded-full text-[10px] font-bold tracking-wide uppercase flex items-center gap-1 z-10">🔥 Best Seller</span>`
      : '';

    // Image fallback
    const imageURL = item.image_url || 'https://images.unsplash.com/photo-1509042239860-f550ce710b93?w=500&auto=format&fit=crop&q=60';

    card.innerHTML = `
      ${badgeHTML}
      <div class="h-44 w-full bg-gray-100 relative overflow-hidden">
        <img src="${imageURL}" alt="${item.name}" class="w-full h-full object-cover transition duration-300 hover:scale-105">
      </div>
      <div class="p-4 flex flex-col flex-grow">
        <h3 class="font-bold text-gray-800 text-base font-outfit truncate">${item.name}</h3>
        <p class="text-xs text-gray-400 mt-1 line-clamp-2 h-8">${item.description || 'Tidak ada deskripsi'}</p>
        <div class="mt-4 flex items-center justify-between">
          <span class="font-bold text-kopitiam-600 text-base">${formatRupiah(item.price)}</span>
          <span class="text-[10px] px-2 py-1 bg-gray-50 border border-gray-100 rounded text-gray-500 uppercase font-bold tracking-wider">${item.station}</span>
        </div>
      </div>
    `;

    card.addEventListener('click', () => addToCart(item));
    container.appendChild(card);
  });
}

// Add item to Cart
function addToCart(menu) {
  const existing = cart.find(item => item.id === menu.id);
  if (existing) {
    existing.quantity++;
  } else {
    cart.push({
      id: menu.id,
      name: menu.name,
      price: menu.price,
      quantity: 1,
      notes: ''
    });
  }
  updateCartUI();
}

// Update Cart Quantity count
function updateQuantity(id, change) {
  const item = cart.find(x => x.id === id);
  if (!item) return;

  item.quantity += change;
  if (item.quantity <= 0) {
    cart = cart.filter(x => x.id !== id);
  }
  updateCartUI();
}

// Update Cart Item Notes text
function updateNotes(id, val) {
  const item = cart.find(x => x.id === id);
  if (item) item.notes = val;
}

// Remove single Item from Cart
function removeItem(id) {
  cart = cart.filter(x => x.id !== id);
  updateCartUI();
}

// Clear all Cart items
function clearCart() {
  cart = [];
  updateCartUI();
}

// Calculate totals and render Sidebar UI
function updateCartUI() {
  const list = $('cart-list');
  const badge = $('cart-badge');
  badge.textContent = cart.reduce((acc, curr) => acc + curr.quantity, 0);

  if (cart.length === 0) {
    list.innerHTML = `
      <div class="flex flex-col items-center justify-center h-full text-gray-400 text-center p-6">
        <span class="text-4xl mb-2">🛒</span>
        <p class="text-sm font-medium">Keranjang masih kosong</p>
        <p class="text-xs mt-1 text-gray-400">Pilih menu di sebelah kiri untuk ditambahkan</p>
      </div>
    `;
    $('subtotal-val').textContent = 'Rp 0';
    $('tax-val').textContent = 'Rp 0';
    $('service-val').textContent = 'Rp 0';
    $('total-val').textContent = 'Rp 0';
    $('checkout-btn').disabled = true;
    $('checkout-btn').classList.add('opacity-50');
    return;
  }

  $('checkout-btn').disabled = false;
  $('checkout-btn').classList.remove('opacity-50');
  list.innerHTML = '';

  let subtotal = 0;
  cart.forEach(item => {
    subtotal += item.price * item.quantity;
    const row = document.createElement('div');
    row.className = "p-3 bg-white border border-gray-100 rounded-xl space-y-2 relative shadow-sm";
    row.innerHTML = `
      <div class="flex items-start justify-between">
        <div class="pr-6">
          <h4 class="font-bold text-sm text-gray-800">${item.name}</h4>
          <span class="text-xs text-kopitiam-500 font-semibold">${formatRupiah(item.price)}</span>
        </div>
        <button class="absolute top-2 right-2 text-gray-300 hover:text-red-500 text-xs transition" onclick="removeItem('${item.id}')">✕</button>
      </div>
      
      <!-- Notes input field -->
      <input type="text" placeholder="Tambah catatan (contoh: es sedikit)" class="w-full text-xs bg-gray-50 border border-gray-100 rounded-lg p-1.5 focus:outline-none focus:border-kopitiam-400 transition" value="${item.notes}" oninput="updateNotes('${item.id}', this.value)">

      <div class="flex items-center justify-between pt-1">
        <span class="text-xs font-bold text-gray-700">${formatRupiah(item.price * item.quantity)}</span>
        <div class="flex items-center gap-2 border border-gray-100 rounded-lg px-1.5 py-0.5 bg-white">
          <button class="text-gray-400 hover:text-kopitiam-600 font-bold px-1.5 text-sm" onclick="updateQuantity('${item.id}', -1)">-</button>
          <span class="text-xs font-bold text-gray-800 w-4 text-center">${item.quantity}</span>
          <button class="text-gray-400 hover:text-kopitiam-600 font-bold px-1.5 text-sm" onclick="updateQuantity('${item.id}', 1)">+</button>
        </div>
      </div>
    `;
    list.appendChild(row);
  });

  const tax = Math.round(subtotal * 0.1);
  const service = Math.round(subtotal * 0.05); // standard service charge 5%
  const total = subtotal + tax + service;

  $('subtotal-val').textContent = formatRupiah(subtotal);
  $('tax-val').textContent = formatRupiah(tax);
  $('service-val').textContent = formatRupiah(service);
  $('total-val').textContent = formatRupiah(total);
}

// Payment method click handlers toggler (simplified for single Online/Midtrans choice)
const payOnlineBtn = $('pay-online');
if (payOnlineBtn) {
  payOnlineBtn.addEventListener('click', () => {
    selectedPaymentMethod = 'online';
  });
}

// Checkout action handler (Midtrans Snap Trigger)
$('checkout-btn').addEventListener('click', async () => {
  if (cart.length === 0) return;

  const tableNumberInput = $('table-number-input');
  const tableNumber = tableNumberInput ? tableNumberInput.value.trim() : '';
  if (!tableNumber) {
    alert("Harap masukkan nomor meja Anda sebelum melakukan checkout!");
    if (tableNumberInput) tableNumberInput.focus();
    return;
  }

  showLoading("Menghubungkan ke Midtrans Snap...");

  try {
    const payload = {
      items: cart.map(item => ({
        menu_id: item.id,
        quantity: item.quantity,
        notes: item.notes || "tidak ada catatan"
      })),
      payment_method: selectedPaymentMethod,
      table_number: tableNumber
    };

    const res = await fetch(`${CONFIG.API_BASE_URL}/checkout`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(payload)
    });

    const data = await res.json();
    if (!res.ok) throw new Error(data.message || "Gagal memproses checkout");

    const snapToken = data.data.snap_token;
    if (!snapToken) throw new Error("Token pembayaran online kosong");

    if (typeof window.snap === 'undefined') {
      throw new Error("Midtrans Snap SDK belum termuat. Hubungkan internet dan refresh halaman.");
    }

    // Launch Midtrans Snap Modal
    window.snap.pay(snapToken, {
      onSuccess: function(result) {
        hideLoading();
        showSuccessOverlay();
        clearCart();
        if (tableNumberInput) tableNumberInput.value = '';
      },
      onPending: function(result) {
        hideLoading();
        alert("Pembayaran tertunda. Harap selesaikan transaksi Anda.");
      },
      onError: function(result) {
        hideLoading();
        alert("Pembayaran gagal! Silakan coba lagi.");
      },
      onClose: function() {
        hideLoading();
        console.log('Customer closed payment modal without completing payment.');
      }
    });

  } catch (err) {
    hideLoading();
    alert("Checkout Error: " + err.message);
  }
});

// Loading screen triggers
function showLoading(msg) {
  $('loading-message').textContent = msg;
  $('loading-overlay').classList.remove('hidden');
}

function hideLoading() {
  $('loading-overlay').classList.add('hidden');
}

// Success screen triggers
function showSuccessOverlay() {
  const overlay = $('success-overlay');
  overlay.classList.remove('hidden');
  let countdown = 4;
  $('countdown-val').textContent = countdown;
  
  const timer = setInterval(() => {
    countdown--;
    $('countdown-val').textContent = countdown;
    if (countdown <= 0) {
      clearInterval(timer);
      overlay.classList.add('hidden');
    }
  }, 1000);
}

// Reset cart button click listener
$('clear-cart-btn').addEventListener('click', clearCart);

// Start
document.addEventListener('DOMContentLoaded', initCatalog);
