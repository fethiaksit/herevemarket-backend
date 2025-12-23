package handlers

import "github.com/gin-gonic/gin"

func Home() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(200, `
<!doctype html>
<html lang="tr">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Hereve Market Admin</title>

<style>
body {
  margin: 0;
  font-family: system-ui, -apple-system, "Segoe UI", sans-serif;
  background: #f1f5f9;
}
header {
  background: #0b2d66;
  color: white;
  padding: 20px;
}
main {
  max-width: 1100px;
  margin: 0 auto;
  padding: 20px;
  display: grid;
  gap: 20px;
}
section {
  background: white;
  border-radius: 12px;
  border: 1px solid #e2e8f0;
  padding: 16px;
}
form {
  display: grid;
  gap: 8px;
}
label { font-weight: 600; }
input, select {
  padding: 8px;
  border-radius: 6px;
  border: 1px solid #cbd5e1;
}
button {
  padding: 10px;
  background: #2563eb;
  color: white;
  border: none;
  border-radius: 6px;
  font-weight: 700;
  cursor: pointer;
}
button.small { padding: 6px 10px; font-weight: 600; }
button.ghost {
  background: white;
  color: #0b2d66;
  border: 1px solid #cbd5e1;
}
button.danger { background: #dc2626; }
button.danger.ghost {
  color: #dc2626;
  border-color: #fecdd3;
  background: #fff1f2;
}
.clickable { cursor: pointer; }
.card {
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  padding: 8px;
  cursor: pointer;
}
.card.product-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}
.order-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}
.inline-actions { display: flex; gap: 8px; }
.stacked-text { display: grid; gap: 4px; }
.list { display: grid; gap: 6px; }
.muted { color: #475569; font-size: 0.9rem; }
.stacked { display: grid; gap: 10px; }
.stacked.labels { gap: 4px; }
.badge {
  display: inline-block;
  padding: 4px 8px;
  border-radius: 999px;
  font-weight: 700;
  font-size: 0.85rem;
  background: #e2e8f0;
  color: #0b2d66;
}
.badge.pending { background: #fef9c3; color: #854d0e; }
.badge.completed { background: #dcfce7; color: #166534; }
.badge.canceled { background: #fee2e2; color: #991b1b; }
.order-items {
  width: 100%;
  border-collapse: collapse;
}
.order-items th, .order-items td {
  padding: 6px;
  border-bottom: 1px solid #e2e8f0;
  text-align: left;
}
.order-items th {
  font-size: 0.9rem;
  color: #475569;
}
.order-items .numeric { text-align: right; }
.address-block {
  display: grid;
  gap: 2px;
}
.address-label {
  font-weight: 600;
}
.order-summary {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 4px;
}
hr { border: 0; border-top: 1px solid #e2e8f0; margin: 12px 0; }
</style>
</head>

<body>
<header>
  <h1>Hereve Market – Admin Panel</h1>
</header>

<main>

<section>
  <h2>Admin Login</h2>
  <form id="loginForm">
    <label>Email</label>
    <input name="email" value="admin@market.com">
    <label>Password</label>
    <input name="password" type="password" value="123456">
    <button type="submit">Login</button>
  </form>
  <div id="loginStatus" class="muted"></div>
</section>

<section>
  <h2>Kategoriler</h2>

  <form id="addCategory">
    <input name="name" placeholder="Kategori adı">
    <button type="submit">Ekle</button>
  </form>

  <div id="categoryList" class="list"></div>

  <form id="editCategory" class="stacked" style="display:none; margin-top:10px;">
    <hr>
    <div class="muted">Seçilen kategori: <strong id="catName"></strong> <span id="catId" class="muted"></span></div>
    <label>Ad</label>
    <input name="name" placeholder="Yeni ad">
    <label><input type="checkbox" name="isActive"> Aktif</label>
    <button type="submit">Güncelle</button>
    <button type="button" id="deleteCategory" class="danger">Pasifleştir</button>
  </form>
</section>

<section>
  <h2>Ürünler</h2>

  <label>Kategori Filtresi</label>
  <select id="categoryFilter">
    <option value="">Tüm Kategoriler</option>
  </select>

  <form id="addProduct">
    <input name="name" placeholder="Ürün adı">
    <input name="price" placeholder="Fiyat (örn: 24.90)">
    <select name="category" id="productCategorySelect" multiple>
      <option value="">Kategori Seç</option>
    </select>
    <input name="imageUrl" placeholder="Görsel URL">
    <button type="submit">Ekle</button>
  </form>

  <div id="productList" class="list"></div>
  <div id="productStatus" class="muted"></div>

  <form id="editProduct" class="stacked" style="display:none; margin-top:10px;">
    <hr>
    <div class="muted">Seçilen ürün: <strong id="prodName"></strong> <span id="prodId" class="muted"></span></div>
    <label>Ad</label>
    <input name="name" placeholder="Ürün adı">
    <label>Fiyat</label>
    <input name="price" placeholder="Fiyat">
    <label>Kategori</label>
    <select name="category" id="productCategorySelect" multiple>
      <option value="">Kategori Seç</option>
    </select>
    <label>Görsel URL</label>
    <input name="imageUrl" placeholder="Görsel URL">
    <label><input type="checkbox" name="isActive"> Aktif</label>
    <button type="submit">Güncelle</button>
    <button type="button" id="deleteProduct" class="danger">Pasifleştir</button>
  </form>
</section>

<section>
  <h2>Siparişler</h2>
  <div id="ordersStatus" class="muted"></div>
  <div id="ordersList" class="list"></div>
</section>

</main>

<script>
let token = "";
let selectedCategory = null;
let selectedProduct = null;
let currentProducts = [];

/* id helper: mongo bazen _id, bazen id döner */
function getId(obj) {
  return (obj && (obj._id || obj.id)) ? (obj._id || obj.id) : null;
}

function hasToken() {
  return typeof token === "string" && token.length > 10;
}

function authHeaders() {
  return {
    "Content-Type": "application/json",
    "Authorization": "Bearer " + token
  };
}

async function safeJson(res) {
  try { return await res.json(); } catch { return null; }
}

function setText(id, text) {
  const el = document.getElementById(id);
  if (el) el.innerText = text || "";
}

function setProductStatus(text) {
  setText("productStatus", text || "");
}

function normalizeCategoryValues(values) {
  if (Array.isArray(values)) {
    return values.filter(function(v) { return !!v; });
  }

  if (typeof values === "string" && values) {
    return [values];
  }

  return [];
}

function getSelectedCategories(select) {
  if (!select) return [];

  return Array.from(select.selectedOptions || [])
    .map(function(opt) { return opt.value; })
    .filter(function(v) { return !!v; });
}

async function populateProductCategorySelects(selectedValues, preloadedCategories) {
  const desiredSelection = normalizeCategoryValues(selectedValues);
  const categoryData = Array.isArray(preloadedCategories) && preloadedCategories.length > 0
    ? preloadedCategories
    : null;

  let categories = categoryData;

  if (!categories) {
    const res = await fetch("/categories");
    const payload = await safeJson(res);
    categories = (payload && payload.data) ? payload.data : (payload || []);
  }

  const activeCategories = (categories || []).filter(function(c) { return c && c.isActive; });
  const activeNames = new Set(activeCategories.map(function(c) { return c.name; }));
  const selects = document.querySelectorAll("#productCategorySelect");

  selects.forEach(function(select) {
    const preserved = desiredSelection.length > 0 ? desiredSelection : getSelectedCategories(select);
    select.innerHTML = "";
    select.multiple = true;

    const def = document.createElement("option");
    def.value = "";
    def.textContent = "Kategori Seç";
    def.disabled = true;
    if (preserved.length === 0) def.selected = true;
    select.appendChild(def);

    activeCategories.forEach(function(c) {
      const opt = document.createElement("option");
      opt.value = c.name;
      opt.textContent = c.name;
      select.appendChild(opt);
    });

    preserved.forEach(function(val) {
      if (!activeNames.has(val)) return;
      const opt = Array.from(select.options).find(function(o) { return o.value === val; });
      if (opt) opt.selected = true;
    });
  });
}

/* LOGIN */
document.getElementById("loginForm").onsubmit = async function(e) {
  e.preventDefault();
  setText("loginStatus", "Giriş yapılıyor...");

  const f = new FormData(e.target);

  const res = await fetch("/admin/login", {
    method: "POST",
    headers: {"Content-Type":"application/json"},
    body: JSON.stringify({
      email: f.get("email"),
      password: f.get("password")
    })
  });

  const j = await safeJson(res);
  if (!res.ok || !j || !j.token) {
    setText("loginStatus", "Hata: " + JSON.stringify(j));
    return;
  }

  token = j.token;
  setText("loginStatus", "Giriş başarılı ✅");
  await loadCategories();
  await loadProducts();
};

/* CATEGORIES */
async function loadCategories() {
  // 1) dropdown için public categories
  const filterSelect = document.getElementById("categoryFilter");
  const preserved = filterSelect ? filterSelect.value : "";

  const catRes = await fetch("/categories");
  const catPayload = await safeJson(catRes);
  const catData = (catPayload && catPayload.data) ? catPayload.data : (catPayload || []);

  await populateProductCategorySelects(undefined, catData);

  if (filterSelect) {
    filterSelect.innerHTML = "";
    const def = document.createElement("option");
    def.value = "";
    def.textContent = "Tüm Kategoriler";
    filterSelect.appendChild(def);

    (catData || []).forEach(function(c) {
      const opt = document.createElement("option");
      opt.value = c.name;
      opt.textContent = c.name;
      filterSelect.appendChild(opt);
    });

    // önceki seçimi koru
    const exists = (catData || []).some(function(c){ return c.name === preserved; });
    filterSelect.value = exists ? preserved : "";
  }

  // 2) liste: admin varsa admin categories (aktif/pasif), yoksa public
  let listUrl = "/categories";
  let listInit = undefined;
  if (hasToken()) {
    listUrl = "/admin/categories";
    listInit = { headers: authHeaders() };
  }

  const res = await fetch(listUrl, listInit);
  const payload = await safeJson(res);
  const data = (payload && payload.data) ? payload.data : (payload || []);

  const el = document.getElementById("categoryList");
  el.innerHTML = "";

  if (!Array.isArray(data) || data.length === 0) {
    el.innerHTML = "<div class='muted'>Kategori yok</div>";
    return;
  }

  data.forEach(function(c) {
    const card = document.createElement("div");
    card.className = "card clickable";
    card.innerHTML = "<div><strong>" + (c.name || "-") + "</strong></div>" +
      "<div class='muted'>" + (c.isActive ? "Aktif" : "Pasif") + "</div>";
    card.onclick = function() { selectCategory(c); };
    el.appendChild(card);
  });
}

document.getElementById("categoryFilter").onchange = function() {
  loadProducts();
};

/* PRODUCTS */
function renderProductList(data) {
  const el = document.getElementById("productList");
  el.innerHTML = "";

  if (!Array.isArray(data) || data.length === 0) {
    el.innerHTML = "<div class='muted'>Ürün yok</div>";
    return;
  }

  data.forEach(function(p) {
    const categoryLabel = Array.isArray(p.category)
      ? (p.category.length ? p.category.join(", ") : "-")
      : (p.category || "-");

    const row = document.createElement("div");
    row.className = "card product-row";

    const info = document.createElement("div");
    info.className = "stacked-text clickable";
    info.innerHTML =
      "<div><strong>" + (p.name || "-") + "</strong></div>" +
      "<div class='muted'>" +
        (p.price ?? "-") + " • " + categoryLabel + " • " + (p.isActive ? "Aktif" : "Pasif") +
      "</div>";
    info.onclick = function() { selectProduct(p); };

    row.appendChild(info);

    if (hasToken()) {
      const actions = document.createElement("div");
      actions.className = "inline-actions";

      const deleteBtn = document.createElement("button");
      deleteBtn.type = "button";
      deleteBtn.className = "danger ghost small";
      deleteBtn.textContent = "Sil";
      deleteBtn.onclick = function(ev) {
        ev.stopPropagation();
        handleDeleteProduct(p);
      };

      actions.appendChild(deleteBtn);
      row.appendChild(actions);
    }

    el.appendChild(row);
  });
}

async function loadProducts() {
  const selected = document.getElementById("categoryFilter").value;
  let url = hasToken() ? "/admin/products" : "/products";

  if (selected) {
    url += "?" + new URLSearchParams({ category: selected }).toString();
  }

  setProductStatus("Ürünler yükleniyor...");

  const res = await fetch(url, hasToken() ? { headers: authHeaders() } : undefined);
  const payload = await safeJson(res);
  if (!res.ok) {
    setProductStatus("Hata: ürünler getirilemedi");
    return;
  }

  const data = (payload && payload.data) ? payload.data : (payload || []);
  currentProducts = Array.isArray(data) ? data : [];

  renderProductList(currentProducts);
  setProductStatus("");
}

/* SELECT */
function selectCategory(c) {
  selectedCategory = c;
  const id = getId(c);

  document.getElementById("editCategory").style.display = "grid";
  document.getElementById("catName").innerText = c.name || "-";
  document.getElementById("catId").innerText = id ? ("(id: " + id + ")") : "(id yok)";

  const f = document.getElementById("editCategory");
  f.elements.name.value = c.name || "";
  f.elements.isActive.checked = !!c.isActive;
}

async function selectProduct(p) {
  selectedProduct = p;
  const id = getId(p);

  const categories = normalizeCategoryValues(p.category);

  document.getElementById("editProduct").style.display = "grid";
  document.getElementById("prodName").innerText = p.name || "-";
  document.getElementById("prodId").innerText = id ? ("(id: " + id + ")") : "(id yok)";

  await populateProductCategorySelects(categories);

  const f = document.getElementById("editProduct");
  f.elements.name.value = p.name || "";
  f.elements.price.value = (p.price ?? "");
  f.elements.imageUrl.value = p.imageUrl || "";
  f.elements.isActive.checked = !!p.isActive;
}

async function handleDeleteProduct(product) {
  if (!hasToken()) { alert("Önce admin login ol"); return; }
  if (!product) return;

  const id = getId(product);
  if (!id) { alert("Ürün id yok"); return; }

  const confirmed = confirm("Bu ürünü silmek istediğinize emin misiniz?");
  if (!confirmed) return;

  const res = await fetch("/admin/products/" + id, {
    method: "DELETE",
    headers: authHeaders()
  });
  const payload = await safeJson(res);

  if (!res.ok) {
    alert("Silme başarısız: " + ((payload && payload.error) ? payload.error : res.statusText));
    return;
  }

  currentProducts = currentProducts.filter(function(item) { return getId(item) !== id; });
  if (selectedProduct && getId(selectedProduct) === id) {
    selectedProduct = null;
    document.getElementById("editProduct").style.display = "none";
  }

  renderProductList(currentProducts);
  setProductStatus("Ürün silindi");
  alert("Ürün silindi");
}

/* CATEGORY CRUD (admin required) */
document.getElementById("addCategory").onsubmit = async function(e) {
  e.preventDefault();
  if (!hasToken()) { alert("Önce admin login ol"); return; }

  const f = new FormData(e.target);
  await fetch("/admin/categories", {
    method: "POST",
    headers: authHeaders(),
    body: JSON.stringify({ name: f.get("name"), isActive: true })
  });

  e.target.reset();
  loadCategories();
};

document.getElementById("editCategory").onsubmit = async function(e) {
  e.preventDefault();
  if (!hasToken()) { alert("Önce admin login ol"); return; }
  if (!selectedCategory) return;

  const id = getId(selectedCategory);
  if (!id) { alert("Kategori id yok"); return; }

  const f = new FormData(e.target);

  await fetch("/admin/categories/" + id, {
    method: "PUT",
    headers: authHeaders(),
    body: JSON.stringify({
      name: f.get("name"),
      isActive: f.get("isActive") === "on"
    })
  });

  loadCategories();
};

document.getElementById("deleteCategory").onclick = async function() {
  if (!hasToken()) { alert("Önce admin login ol"); return; }
  if (!selectedCategory) return;

  const id = getId(selectedCategory);
  if (!id) { alert("Kategori id yok"); return; }

  await fetch("/admin/categories/" + id, {
    method: "DELETE",
    headers: authHeaders()
  });

  selectedCategory = null;
  document.getElementById("editCategory").style.display = "none";
  loadCategories();
};

// --- ORDERS UI START ---
const ORDER_API_URL = "http://52.57.82.30/orders";

function formatDateTime(value) {
  const d = value ? new Date(value) : null;
  if (!d || isNaN(d.getTime())) return "-";
  return d.toLocaleString("tr-TR");
}

function formatCurrency(value) {
  if (typeof value !== "number") return "-";
  return value.toLocaleString("tr-TR", { style: "currency", currency: "TRY", minimumFractionDigits: 2 });
}

function statusBadgeClass(status) {
  var normalized = (status || "").toLowerCase();
  if (normalized === "completed") return "badge completed";
  if (normalized === "canceled" || normalized === "cancelled") return "badge canceled";
  return "badge pending";
}

function renderOrderItems(items) {
  const table = document.createElement("table");
  table.className = "order-items";

  const thead = document.createElement("thead");
  thead.innerHTML = "<tr><th>Ürün</th><th class='numeric'>Adet</th><th class='numeric'>Birim</th><th class='numeric'>Satır</th></tr>";
  table.appendChild(thead);

  const tbody = document.createElement("tbody");
  (items || []).forEach(function(item) {
    const tr = document.createElement("tr");
    const lineTotal = (item && typeof item.price === "number" && typeof item.quantity === "number")
      ? item.price * item.quantity
      : null;

    tr.innerHTML =
      "<td>" + (item && item.name ? item.name : "-") + "</td>" +
      "<td class='numeric'>" + (item && item.quantity != null ? item.quantity : "-") + "</td>" +
      "<td class='numeric'>" + formatCurrency(item && typeof item.price === "number" ? item.price : null) + "</td>" +
      "<td class='numeric'>" + formatCurrency(lineTotal) + "</td>";
    tbody.appendChild(tr);
  });

  if (!items || items.length === 0) {
    const tr = document.createElement("tr");
    tr.innerHTML = "<td colspan='4' class='muted'>Ürün bulunamadı</td>";
    tbody.appendChild(tr);
  }

  table.appendChild(tbody);
  return table;
}

function renderOrders(data) {
  const el = document.getElementById("ordersList");
  el.innerHTML = "";

  if (!Array.isArray(data) || data.length === 0) {
    el.innerHTML = "<div class='muted'>Sipariş yok</div>";
    return;
  }

  data.forEach(function(order) {
    const card = document.createElement("div");
    card.className = "card";

    const header = document.createElement("div");
    header.className = "order-row clickable";
    header.innerHTML =
      "<div class='stacked-text'>" +
        "<div><strong>Sipariş #" + (order && getId(order) ? getId(order) : "-") + "</strong></div>" +
        "<div class='muted'>Oluşturma: " + formatDateTime(order && order.createdAt) + "</div>" +
      "</div>" +
      "<div class='order-summary'>" +
        "<div><span class='" + statusBadgeClass(order && order.status) + "'>" + (order && order.status ? order.status : "Bilinmiyor") + "</span></div>" +
        "<div class='muted'>" + (order && order.paymentMethod ? order.paymentMethod : "-") + "</div>" +
        "<div><strong>" + formatCurrency(order && typeof order.totalPrice === "number" ? order.totalPrice : null) + "</strong></div>" +
      "</div>";

    const details = document.createElement("div");
    details.className = "stacked";
    details.style.display = "none";

    const customer = order && order.customer ? order.customer : {};
    const customerBox = document.createElement("div");
    customerBox.className = "stacked labels";
    customerBox.innerHTML =
      "<div class='address-label'>Başlık:</div>" +
      "<div>" + (customer.title || "-") + "</div>" +
      "<div class='address-label'>Adres:</div>" +
      "<div class='address-block'>" + (customer.detail || "-").split("\n").join("<br>") + "</div>" +
      (customer.note ? ("<div class='address-label'>Not:</div><div>" + customer.note + "</div>") : "");

    const itemsBox = document.createElement("div");
    itemsBox.appendChild(renderOrderItems(order && order.items ? order.items : []));

    const totalLine = document.createElement("div");
    totalLine.style.textAlign = "right";
    totalLine.innerHTML = "<strong>Toplam: " + formatCurrency(order && typeof order.totalPrice === "number" ? order.totalPrice : null) + "</strong>";

    details.appendChild(customerBox);
    details.appendChild(itemsBox);
    details.appendChild(totalLine);

    header.onclick = function() {
      details.style.display = details.style.display === "none" ? "grid" : "none";
    };

    card.appendChild(header);
    card.appendChild(details);
    el.appendChild(card);
  });
}

async function loadOrders() {
  setText("ordersStatus", "Siparişler yükleniyor...");
  const res = await fetch(ORDER_API_URL);
  const payload = await safeJson(res);
  if (!res.ok) {
    setText("ordersStatus", "Hata: siparişler getirilemedi");
    return;
  }
  const data = (payload && payload.data) ? payload.data : (payload || []);
  renderOrders(data);
  setText("ordersStatus", "");
}
// --- ORDERS UI END ---

/* PRODUCT CRUD (admin required) */
document.getElementById("addProduct").onsubmit = async function(e) {
  e.preventDefault();
  if (!hasToken()) { alert("Önce admin login ol"); return; }

  const f = new FormData(e.target);
  const price = parseFloat(f.get("price"));
  if (Number.isNaN(price)) { alert("Fiyat sayı olmalı (örn 24.90)"); return; }

  const categories = getSelectedCategories(e.target.querySelector('select[name="category"]'));
  if (categories.length === 0) { alert("En az bir kategori seç"); return; }

  await fetch("/admin/products", {
    method: "POST",
    headers: authHeaders(),
    body: JSON.stringify({
      name: f.get("name"),
      price: price,
      category: categories,
      imageUrl: f.get("imageUrl"),
      isActive: true
    })
  });

  e.target.reset();
  loadProducts();
};

document.getElementById("editProduct").onsubmit = async function(e) {
  e.preventDefault();
  if (!hasToken()) { alert("Önce admin login ol"); return; }
  if (!selectedProduct) return;

  const id = getId(selectedProduct);
  if (!id) { alert("Ürün id yok"); return; }

  const f = new FormData(e.target);
  const price = parseFloat(f.get("price"));
  if (Number.isNaN(price)) { alert("Fiyat sayı olmalı"); return; }

  const categories = getSelectedCategories(e.target.querySelector('select[name="category"]'));
  if (categories.length === 0) { alert("En az bir kategori seç"); return; }

  await fetch("/admin/products/" + id, {
    method: "PUT",
    headers: authHeaders(),
    body: JSON.stringify({
      name: f.get("name"),
      price: price,
      category: categories,
      imageUrl: f.get("imageUrl"),
      isActive: f.get("isActive") === "on"
    })
  });

  loadProducts();
};

document.getElementById("deleteProduct").onclick = async function() {
  if (!hasToken()) { alert("Önce admin login ol"); return; }
  if (!selectedProduct) return;

  await handleDeleteProduct(selectedProduct);
};

/* initial load (public) */
loadCategories();
loadProducts();
loadOrders();
</script>

</body>
</html>
`)
	}
}
