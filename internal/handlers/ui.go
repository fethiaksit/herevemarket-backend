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
.clickable { cursor: pointer; }
.card {
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  padding: 8px;
  cursor: pointer;
}
.list { display: grid; gap: 6px; }
.muted { color: #475569; font-size: 0.9rem; }
.stacked { display: grid; gap: 10px; }
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
    <input name="category" placeholder="Kategori">
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
    <input name="category" placeholder="Kategori">
    <label>Görsel URL</label>
    <input name="imageUrl" placeholder="Görsel URL">
    <label><input type="checkbox" name="isActive"> Aktif</label>
    <button type="submit">Güncelle</button>
    <button type="button" id="deleteProduct" class="danger">Pasifleştir</button>
  </form>
</section>

</main>

<script>
let token = "";
let selectedCategory = null;
let selectedProduct = null;

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
    card.innerHTML =
      "<div><strong>" + (p.name || "-") + "</strong></div>" +
      "<div class='muted'>" +
        (p.price ?? "-") + " • " + (p.category || "-") + " • " + (p.isActive ? "Aktif" : "Pasif") +
      "</div>";
    card.onclick = function() { selectProduct(p); };
    el.appendChild(card);
  });
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

function selectProduct(p) {
  selectedProduct = p;
  const id = getId(p);

  document.getElementById("editProduct").style.display = "grid";
  document.getElementById("prodName").innerText = p.name || "-";
  document.getElementById("prodId").innerText = id ? ("(id: " + id + ")") : "(id yok)";

  const f = document.getElementById("editProduct");
  f.elements.name.value = p.name || "";
  f.elements.price.value = (p.price ?? "");
  f.elements.category.value = p.category || "";
  f.elements.imageUrl.value = p.imageUrl || "";
  f.elements.isActive.checked = !!p.isActive;
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

  await fetch("/admin/products", {
    method: "POST",
    headers: authHeaders(),
    body: JSON.stringify({
      name: f.get("name"),
      price: price,
      category: f.get("category"),
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

  await fetch("/admin/products/" + id, {
    method: "PUT",
    headers: authHeaders(),
    body: JSON.stringify({
      name: f.get("name"),
      price: price,
      category: f.get("category"),
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

  await fetch("/admin/products/" + id, {
    method: "DELETE",
    headers: authHeaders()
  });

  selectedProduct = null;
  document.getElementById("editProduct").style.display = "none";
  loadProducts();
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
