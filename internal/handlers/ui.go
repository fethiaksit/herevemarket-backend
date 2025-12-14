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
.clickable {
  cursor: pointer;
}
.muted {
  color: #475569;
}
.stacked {
  display: grid;
  gap: 10px;
}
.card {
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  padding: 8px;
}
.list {
  display: grid;
  gap: 6px;
}
pre {
  background: #020617;
  color: #e5e7eb;
  padding: 10px;
  border-radius: 6px;
}
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
  <button>Login</button>
</form>
<div id="loginStatus"></div>
</section>

<section>
<h2>Kategoriler</h2>
<form id="addCategory">
  <input name="name" placeholder="Kategori adı">
  <button>Ekle</button>
</form>
<div id="categoryList" class="list"></div>
<form id="editCategory" class="stacked" style="display: none; margin-top: 10px;">
  <div class="muted">Seçilen kategori: <strong id="categorySelection"></strong></div>
  <label>Ad</label>
  <input name="name" placeholder="Yeni ad">
  <label><input type="checkbox" name="isActive"> Aktif</label>
  <div class="list">
    <button type="submit">Güncelle</button>
    <button type="button" id="deleteCategory" class="danger">Sil</button>
  </div>
</form>
</section>

<section>
<h2>Ürünler</h2>
<form id="addProduct">
  <input name="name" placeholder="Ürün adı">
  <input name="price" placeholder="Fiyat">
  <input name="category" placeholder="Kategori">
  <input name="imageUrl" placeholder="Görsel URL">
  <button>Ekle</button>
</form>
<div id="productList" class="list"></div>
<form id="editProduct" class="stacked" style="display: none; margin-top: 10px;">
  <div class="muted">Seçilen ürün: <strong id="productSelection"></strong></div>
  <label>Ad</label>
  <input name="name" placeholder="Ürün adı">
  <label>Fiyat</label>
  <input name="price" placeholder="Fiyat">
  <label>Kategori</label>
  <input name="category" placeholder="Kategori">
  <label>Görsel URL</label>
  <input name="imageUrl" placeholder="Görsel URL">
  <label><input type="checkbox" name="isActive"> Aktif</label>
  <div class="list">
    <button type="submit">Güncelle</button>
    <button type="button" id="deleteProduct" class="danger">Sil</button>
  </div>
</form>
</section>

</main>

<script>
let token = "";
let selectedCategory = null;
let selectedProduct = null;

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
  token = j.token;
  document.getElementById("loginStatus").innerText = "Giriş OK";
  loadCategories();
  loadProducts();
};

function authHeaders() {
  return {
    "Content-Type": "application/json",
    "Authorization": "Bearer " + token
  };
}

async function loadCategories() {
  const res = await fetch(token ? "/admin/categories" : "/categories", token ? { headers: authHeaders() } : undefined);
  const data = await res.json();
  const el = document.getElementById("categoryList");
  el.innerHTML = "";
  data.forEach(function(c) {
    const card = document.createElement("div");
    card.className = "card clickable";
    card.innerHTML = "<div>" + c.name + "</div><div class='muted'>" + (c.isActive ? "Aktif" : "Pasif") + "</div>";
    card.onclick = function() { selectCategory(c); };
    el.appendChild(card);
  });
}

async function loadProducts() {
  const res = await fetch(token ? "/admin/products" : "/products", token ? { headers: authHeaders() } : undefined);
  const payload = await res.json();
  const data = payload.data || payload; // supports both array and paginated responses
  const el = document.getElementById("productList");
  el.innerHTML = "";
  data.forEach(function(p) {
    const card = document.createElement("div");
    card.className = "card clickable";
    card.innerHTML = "<div>" + p.name + "</div><div class='muted'>" + p.price + " • " + p.category + " • " + (p.isActive ? "Aktif" : "Pasif") + "</div>";
    card.onclick = function() { selectProduct(p); };
    el.appendChild(card);
  });
}

function selectCategory(category) {
  selectedCategory = category;
  document.getElementById("categorySelection").innerText = category.name;
  const form = document.getElementById("editCategory");
  form.style.display = "grid";
  form.elements.name.value = category.name;
  form.elements.isActive.checked = !!category.isActive;
}

function selectProduct(product) {
  selectedProduct = product;
  document.getElementById("productSelection").innerText = product.name;
  const form = document.getElementById("editProduct");
  form.style.display = "grid";
  form.elements.name.value = product.name;
  form.elements.price.value = product.price;
  form.elements.category.value = product.category;
  form.elements.imageUrl.value = product.imageUrl;
  form.elements.isActive.checked = !!product.isActive;
}

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

document.getElementById("editCategory").onsubmit = async function(e) {
  e.preventDefault();
  if (!selectedCategory) return;
  const f = new FormData(e.target);
  await fetch(`/admin/categories/${selectedCategory.id}`, {
    method: "PUT",
    headers: authHeaders(),
    body: JSON.stringify({
      name: f.get("name"),
      isActive: f.get("isActive") === "on",
    })
  });
  selectedCategory = null;
  document.getElementById("editCategory").style.display = "none";
  loadCategories();
};

document.getElementById("deleteCategory").onclick = async function() {
  if (!selectedCategory) return;
  await fetch(`/admin/categories/${selectedCategory.id}`, {
    method: "DELETE",
    headers: authHeaders()
  });
  selectedCategory = null;
  document.getElementById("editCategory").style.display = "none";
  loadCategories();
};

document.getElementById("addProduct").onsubmit = async function(e) {
  e.preventDefault();
  const f = new FormData(e.target);
  await fetch("/admin/products", {
    method: "POST",
    headers: authHeaders(),
    body: JSON.stringify({
      name: f.get("name"),
      price: parseFloat(f.get("price")),
      category: f.get("category"),
      imageUrl: f.get("imageUrl"),
      isActive: true
    })
  });
  loadProducts();
};

document.getElementById("editProduct").onsubmit = async function(e) {
  e.preventDefault();
  if (!selectedProduct) return;
  const f = new FormData(e.target);
  await fetch(`/admin/products/${selectedProduct.id}`, {
    method: "PUT",
    headers: authHeaders(),
    body: JSON.stringify({
      name: f.get("name"),
      price: parseFloat(f.get("price")),
      category: f.get("category"),
      imageUrl: f.get("imageUrl"),
      isActive: f.get("isActive") === "on",
    })
  });
  selectedProduct = null;
  document.getElementById("editProduct").style.display = "none";
  loadProducts();
};

document.getElementById("deleteProduct").onclick = async function() {
  if (!selectedProduct) return;
  await fetch(`/admin/products/${selectedProduct.id}`, {
    method: "DELETE",
    headers: authHeaders()
  });
  selectedProduct = null;
  document.getElementById("editProduct").style.display = "none";
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
