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
label {
  font-weight: 600;
}
input {
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
button.danger {
  background: #dc2626;
}
.card {
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  padding: 8px;
  cursor: pointer;
}
.list {
  display: grid;
  gap: 6px;
}
.muted {
  color: #475569;
  font-size: 0.9rem;
}
</style>
</head>

<body>
<header>
  <h1>Hereve Market – Admin Panel</h1>
  <div id="loginStatus" class="muted">Giriş yapılmadı.</div>
</header>

<main>

<section>
<h2>Admin Login</h2>
<form id="loginForm">
  <label>Email</label>
  <input name="email" value="admin@market.com">
  <label>Password</label>
  <input name="password" type="password" value="123456">
  <button>Login</button>
</form>
<div id="loginStatus" class="muted"></div>
</section>

<section>
<h2>Kategoriler</h2>

<form id="addCategory">
  <input name="name" placeholder="Kategori adı">
  <button>Ekle</button>
</form>

  <div id="categoryList" class="list"></div>

<form id="editCategory" style="display:none; margin-top:10px;">
  <div class="muted">Seçilen kategori: <strong id="catName"></strong></div>
  <input name="name" placeholder="Yeni ad">
  <label><input type="checkbox" name="isActive"> Aktif</label>
  <button>Güncelle</button>
  <button type="button" id="deleteCategory" class="danger">Pasifleştir</button>
</form>
</section>

<section>
<h2>Ürünler</h2>

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
async function safeJson(res) {
  try { return await res.json(); } catch { return null; }
}

function setText(id, text) {
  const el = document.getElementById(id);
  if (el) el.innerText = text || "";
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
  const f = new FormData(e.target);

  const res = await fetch("/admin/login", {
    method: "POST",
    headers: {"Content-Type": "application/json"},
    body: JSON.stringify({
      email: f.get("email"),
      password: f.get("password")
    })
  });

  const j = await res.json();
  if (!res.ok) {
    document.getElementById("loginStatus").innerText = "Hata";
    return;
  }

  token = j.token;
  document.getElementById("loginStatus").innerText = "Giriş başarılı";
  loadCategories();
  loadProducts();
};

/* LIST LOADERS */
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
    setMsg("categoryMsg", "");
    return;
  }

  data.forEach(c => {
    const d = document.createElement("div");
    d.className = "card";
    d.innerHTML = c.name + " (" + (c.isActive ? "Aktif" : "Pasif") + ")";
    d.onclick = () => selectCategory(c);
    el.appendChild(d);
  });
}

document.getElementById("categoryFilter").onchange = function() {
  loadProducts();
};

/* PRODUCTS */
async function loadProducts() {
  const selected = document.getElementById("categoryFilter").value;
  let url = hasToken() ? "/admin/products" : "/products";

  if (selected) {
    url += "?" + new URLSearchParams({ category: selected }).toString();
  }

  const res = await fetch(url, hasToken() ? { headers: authHeaders() } : undefined);
  const payload = await safeJson(res);
  const data = (payload && payload.data) ? payload.data : (payload || []);

  const el = document.getElementById("productList");
  el.innerHTML = "";

  if (!Array.isArray(data) || data.length === 0) {
    el.innerHTML = "<div class='muted'>Ürün yok</div>";
    return;
  }

  data.forEach(function(p) {
    const card = document.createElement("div");
    card.className = "card clickable";
    const categoryLabel = Array.isArray(p.category)
      ? (p.category.length ? p.category.join(", ") : "-")
      : (p.category || "-");
    card.innerHTML =
      "<div><strong>" + (p.name || "-") + "</strong></div>" +
      "<div class='muted'>" +
        (p.price ?? "-") + " • " + categoryLabel + " • " + (p.isActive ? "Aktif" : "Pasif") +
      "</div>";
    card.onclick = function() { selectProduct(p); };
    el.appendChild(card);
  });
}

/* SELECT */
function selectCategory(c) {
  selectedCategory = c;
  document.getElementById("editCategory").style.display = "grid";
  document.getElementById("catName").innerText = c.name;
  const f = document.getElementById("editCategory");
  f.elements.name.value = c.name;
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

/* CATEGORY CRUD (admin required) */
document.getElementById("addCategory").onsubmit = async function(e) {
  e.preventDefault();
  const f = new FormData(e.target);
  await fetch("/admin/categories", {
    method: "POST",
    headers: authHeaders(),
    body: JSON.stringify({ name: f.get("name"), isActive: true })
  });
  loadCategories();
};

document.getElementById("editCategory").onsubmit = async e => {
  e.preventDefault();
  const f = new FormData(e.target);
  await fetch("/admin/categories/" + selectedCategory._id, {
    method: "PUT",
    headers: authHeaders(),
    body: JSON.stringify({
      name: f.get("name"),
      isActive: f.get("isActive") === "on"
    })
  });
  loadCategories();
};

document.getElementById("deleteCategory").onclick = async () => {
  await fetch("/admin/categories/" + selectedCategory._id, {
    method: "DELETE",
    headers: authHeaders()
  });
  loadCategories();
};

/* PRODUCTS */
async function loadProducts() {
  const res = await fetch(token ? "/admin/products" : "/products",
    token ? { headers: authHeaders() } : undefined
  );
  const payload = await res.json();
  const data = payload.data || payload;

  const el = document.getElementById("productList");
  el.innerHTML = "";

  data.forEach(p => {
    const d = document.createElement("div");
    d.className = "card";
    d.innerHTML = p.name + " - " + p.price + " (" + (p.isActive ? "Aktif" : "Pasif") + ")";
    d.onclick = () => selectProduct(p);
    el.appendChild(d);
  });
}

function selectProduct(p) {
  selectedProduct = p;
  document.getElementById("editProduct").style.display = "grid";
  document.getElementById("prodName").innerText = p.name;
  const f = document.getElementById("editProduct");
  f.elements.name.value = p.name;
  f.elements.price.value = p.price;
  f.elements.category.value = p.category;
  f.elements.imageUrl.value = p.imageUrl;
  f.elements.isActive.checked = !!p.isActive;
}

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
  loadProducts();
};

document.getElementById("editProduct").onsubmit = async function(e) {
  e.preventDefault();
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

document.getElementById("deleteProduct").onclick = async () => {
  await fetch("/admin/products/" + selectedProduct._id, {
    method: "DELETE",
    headers: authHeaders()
  });
  loadProducts();
};

loadCategories();
loadProducts();
</script>

</body>
</html>
`)
	}
}
