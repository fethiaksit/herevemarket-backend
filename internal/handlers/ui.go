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
  <input name="price" placeholder="Fiyat">
  <input name="category" placeholder="Kategori">
  <input name="imageUrl" placeholder="Görsel URL">
  <button>Ekle</button>
</form>

  <div id="productList" class="list"></div>

<form id="editProduct" style="display:none; margin-top:10px;">
  <div class="muted">Seçilen ürün: <strong id="prodName"></strong></div>
  <input name="name" placeholder="Ad">
  <input name="price" placeholder="Fiyat">
  <input name="category" placeholder="Kategori">
  <input name="imageUrl" placeholder="Görsel URL">
  <label><input type="checkbox" name="isActive"> Aktif</label>
  <button>Güncelle</button>
  <button type="button" id="deleteProduct" class="danger">Pasifleştir</button>
</form>
</section>

</main>

<script>
let token = "";
let selectedCategory = null;
let selectedProduct = null;

function authHeaders() {
  return {
    "Content-Type": "application/json",
    "Authorization": "Bearer " + token
  };
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
  const res = await fetch(token ? "/admin/categories" : "/categories",
    token ? { headers: authHeaders() } : undefined
  );
  const payload = await res.json();
  const data = payload.data || payload;

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

function selectCategory(c) {
  selectedCategory = c;
  document.getElementById("editCategory").style.display = "grid";
  document.getElementById("catName").innerText = c.name;
  const f = document.getElementById("editCategory");
  f.elements.name.value = c.name;
  f.elements.isActive.checked = !!c.isActive;
}

document.getElementById("addCategory").onsubmit = async e => {
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

document.getElementById("addProduct").onsubmit = async e => {
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
  const f = new FormData(e.target);
  await fetch("/admin/products/" + selectedProduct._id, {
    method: "PUT",
    headers: authHeaders(),
    body: JSON.stringify({
      name: f.get("name"),
      price: parseFloat(f.get("price")),
      category: f.get("category"),
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
