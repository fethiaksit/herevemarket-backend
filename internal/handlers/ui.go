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
button.danger { background: #dc2626; }
button.ghost { background: #fff; color: #0f172a; border: 1px solid #cbd5e1; }
button.ghost.danger { color: #b91c1c; border-color: #fecaca; }
.clickable { cursor: pointer; }
.card {
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  padding: 8px;
  cursor: pointer;
  display: grid;
  gap: 6px;
}
.list { display: grid; gap: 6px; }
.muted { color: #475569; font-size: 0.9rem; }
.stacked { display: grid; gap: 10px; }
hr { border: 0; border-top: 1px solid #e2e8f0; margin: 12px 0; }
.inline { display: flex; gap: 8px; align-items: center; }
.toolbar { display: flex; justify-content: space-between; align-items: center; gap: 12px; margin: 8px 0; flex-wrap: wrap; }
.badge {
  display: inline-block;
  padding: 4px 8px;
  border-radius: 9999px;
  font-size: 12px;
  font-weight: 700;
  color: #0f172a;
  background: #e2e8f0;
}
.badge.success { background: #dcfce7; color: #166534; }
.badge.inactive { background: #fee2e2; color: #b91c1c; }
.badge.info { background: #e0f2fe; color: #075985; }
.badge.deleted { background: #f8fafc; color: #0f172a; border: 1px dashed #cbd5e1; }
.modal-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(15, 23, 42, 0.45);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 50;
}
.modal {
  background: white;
  border-radius: 10px;
  padding: 20px;
  width: 320px;
  box-shadow: 0 15px 40px rgba(0,0,0,0.16);
}
.modal-actions { display: flex; justify-content: flex-end; gap: 10px; margin-top: 16px; }
.toast {
  position: fixed;
  bottom: 20px;
  right: 20px;
  background: #0f172a;
  color: white;
  padding: 12px 16px;
  border-radius: 10px;
  box-shadow: 0 10px 30px rgba(0,0,0,0.2);
  z-index: 60;
  min-width: 220px;
}
.toast.error { background: #b91c1c; }
.product-actions { display: flex; gap: 8px; align-items: center; }
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

  <div class="toolbar">
    <div style="display:flex; gap:8px; align-items:center; flex-wrap:wrap;">
      <label style="margin:0;">Kategori Filtresi</label>
      <select id="categoryFilter">
        <option value="">Tüm Kategoriler</option>
      </select>
      <label class="inline" style="margin-left:6px;"><input type="checkbox" id="includeDeleted"> Silinenleri göster</label>
    </div>
    <div class="inline">
      <button type="button" id="bulkDelete" class="danger ghost" disabled>Seçili ürünleri sil</button>
    </div>
  </div>

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

  <form id="editProduct" class="stacked" style="display:none; margin-top:10px;">
    <hr>
    <div class="muted" id="selectedProductInfo">Seçilen ürün: <strong id="prodName"></strong> <span id="prodId" class="muted"></span></div>
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
    <button type="button" id="deleteProduct" class="danger">Sil</button>
  </form>
</section>

</main>

<div id="toastContainer"></div>
<div id="confirmModal" class="modal-backdrop" style="display:none;">
  <div class="modal">
    <div id="confirmMessage" style="font-size:15px; line-height:1.5;"></div>
    <div class="modal-actions">
      <button type="button" id="confirmCancel" class="ghost">Vazgeç</button>
      <button type="button" id="confirmOk" class="danger">Evet, sil</button>
    </div>
  </div>
</div>

<script>
let token = "";
let selectedCategory = null;
let selectedProduct = null;
const selectedProductIds = new Set();

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

function showToast(message, isError) {
  const container = document.getElementById("toastContainer");
  if (!container) return;

  const toast = document.createElement("div");
  toast.className = "toast" + (isError ? " error" : "");
  toast.textContent = message || "";
  container.appendChild(toast);

  setTimeout(() => {
    toast.style.opacity = "0";
    toast.style.transition = "opacity 0.4s ease";
  }, 2200);

  setTimeout(() => {
    toast.remove();
  }, 2700);
}

function showConfirm(message) {
  return new Promise((resolve) => {
    const modal = document.getElementById("confirmModal");
    const messageEl = document.getElementById("confirmMessage");
    const okBtn = document.getElementById("confirmOk");
    const cancelBtn = document.getElementById("confirmCancel");

    if (!modal || !messageEl || !okBtn || !cancelBtn) {
      resolve(window.confirm(message));
      return;
    }

    messageEl.textContent = message;
    modal.style.display = "flex";

    const cleanup = () => {
      modal.style.display = "none";
      okBtn.onclick = null;
      cancelBtn.onclick = null;
    };

    okBtn.onclick = () => { cleanup(); resolve(true); };
    cancelBtn.onclick = () => { cleanup(); resolve(false); };
  });
}

function syncBulkDeleteButton() {
  const btn = document.getElementById("bulkDelete");
  if (!btn) return;

  const count = selectedProductIds.size;
  btn.disabled = !hasToken() || count === 0;
  btn.textContent = count > 0
    ? `Seçili ürünleri sil (${count})`
    : "Seçili ürünleri sil";
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

document.getElementById("includeDeleted").onchange = function() {
  loadProducts();
};

/* PRODUCTS */
async function loadProducts() {
  const includeDeletedCheckbox = document.getElementById("includeDeleted");
  if (includeDeletedCheckbox) {
    includeDeletedCheckbox.disabled = !hasToken();
  }

  const selected = document.getElementById("categoryFilter").value;
  const includeDeleted = includeDeletedCheckbox && includeDeletedCheckbox.checked;
  let url = hasToken() ? "/admin/products" : "/products";

  const params = new URLSearchParams();
  if (selected) params.set("category", selected);
  if (includeDeleted && hasToken()) params.set("includeDeleted", "true");

  const query = params.toString();
  if (query) url += "?" + query;

  const res = await fetch(url, hasToken() ? { headers: authHeaders() } : undefined);
  const payload = await safeJson(res);
  const data = (payload && payload.data) ? payload.data : (payload || []);

  const availableIds = new Set();
  if (Array.isArray(data)) {
    data.forEach(function(p) {
      const id = getId(p);
      if (id) availableIds.add(id);
    });
  }

  Array.from(selectedProductIds).forEach(function(id) {
    if (!availableIds.has(id)) selectedProductIds.delete(id);
  });

  const el = document.getElementById("productList");
  el.innerHTML = "";

  if (!Array.isArray(data) || data.length === 0) {
    el.innerHTML = "<div class='muted'>Ürün yok</div>";
    syncBulkDeleteButton();
    return;
  }

  data.forEach(function(p) {
    const card = document.createElement("div");
    card.className = "card clickable";
    const categoryLabel = Array.isArray(p.category)
      ? (p.category.length ? p.category.join(", ") : "-")
      : (p.category || "-");

    const id = getId(p);

    const title = document.createElement("div");
    title.style.display = "flex";
    title.style.justifyContent = "space-between";
    title.style.alignItems = "center";
    const nameEl = document.createElement("strong");
    nameEl.textContent = p.name || "-";
    title.appendChild(nameEl);

    if (hasToken()) {
      const actions = document.createElement("div");
      actions.className = "product-actions";

      const checkbox = document.createElement("input");
      checkbox.type = "checkbox";
      checkbox.checked = selectedProductIds.has(id);
      checkbox.onclick = function(e) {
        e.stopPropagation();
        if (!id) return;
        if (checkbox.checked) selectedProductIds.add(id); else selectedProductIds.delete(id);
        syncBulkDeleteButton();
      };
      actions.appendChild(checkbox);

      const deleteBtn = document.createElement("button");
      deleteBtn.type = "button";
      deleteBtn.className = "ghost danger";
      deleteBtn.textContent = "Sil";
      deleteBtn.onclick = function(e) {
        e.stopPropagation();
        if (id) handleDeleteProduct(id);
      };
      actions.appendChild(deleteBtn);

      title.appendChild(actions);
    }

    card.appendChild(title);

    const sub = document.createElement("div");
    sub.className = "muted";
    sub.textContent = (p.price ?? "-") + " • " + categoryLabel;
    card.appendChild(sub);

    const badgeRow = document.createElement("div");
    badgeRow.className = "inline";

    const status = document.createElement("span");
    if (p.isDeleted) {
      status.className = "badge deleted";
      status.textContent = "Silinmiş";
    } else if (p.isActive) {
      status.className = "badge success";
      status.textContent = "Aktif";
    } else {
      status.className = "badge inactive";
      status.textContent = "Pasif";
    }
    badgeRow.appendChild(status);

    if (p.isCampaign) {
      const badge = document.createElement("span");
      badge.className = "badge info";
      badge.textContent = "Kampanyalı";
      badgeRow.appendChild(badge);
    }

    card.appendChild(badgeRow);

    card.onclick = function() { selectProduct(p); };
    el.appendChild(card);
  });

  syncBulkDeleteButton();
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
  f.elements.isActive.checked = !!p.isActive && !p.isDeleted;
  f.elements.isActive.disabled = !!p.isDeleted;

  const info = document.getElementById("selectedProductInfo");
  if (info) {
    let text = "Seçilen ürün: " + (p.name || "-");
    if (id) text += " (id: " + id + ")";
    if (p.isDeleted) text += " • Silinmiş";
    info.innerText = text;
  }
}

async function handleDeleteProduct(id) {
  if (!hasToken()) { alert("Önce admin login ol"); return; }

  const confirmed = await showConfirm("Bu ürünü silmek istediğine emin misin? Bu işlem geri alınamaz.");
  if (!confirmed) return;

  const res = await fetch("/admin/products/" + id, {
    method: "DELETE",
    headers: authHeaders()
  });
  const payload = await safeJson(res);

  if (!res.ok) {
    const errMsg = (payload && payload.error) ? payload.error : "Silme başarısız";
    showToast(errMsg, true);
    return;
  }

  if (selectedProduct && getId(selectedProduct) === id) {
    selectedProduct = null;
    document.getElementById("editProduct").style.display = "none";
  }
  selectedProductIds.delete(id);
  await loadProducts();
  showToast((payload && payload.message) ? payload.message : "Ürün silindi");
}

async function handleBulkDelete() {
  if (!hasToken()) { alert("Önce admin login ol"); return; }
  const ids = Array.from(selectedProductIds);
  if (ids.length === 0) return;

  const confirmed = await showConfirm("Seçili ürünleri silmek istediğine emin misin? Bu işlem geri alınamaz.");
  if (!confirmed) return;

  const res = await fetch("/admin/products/bulk-delete", {
    method: "POST",
    headers: authHeaders(),
    body: JSON.stringify({ ids: ids })
  });
  const payload = await safeJson(res);

  if (!res.ok) {
    const errMsg = (payload && payload.error) ? payload.error : "Silme başarısız";
    showToast(errMsg, true);
    return;
  }

  selectedProductIds.clear();
  selectedProduct = null;
  document.getElementById("editProduct").style.display = "none";
  await loadProducts();
  syncBulkDeleteButton();
  showToast((payload && payload.message) ? payload.message : "Ürünler silindi");
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

  const id = getId(selectedProduct);
  if (!id) { alert("Ürün id yok"); return; }

  await handleDeleteProduct(id);
};

document.getElementById("bulkDelete").onclick = function() {
  handleBulkDelete();
};

/* initial load (public) */
loadCategories();
loadProducts();
</script>

</body>
</html>
`)
	}
}
