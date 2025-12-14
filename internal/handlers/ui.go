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
</section>

</main>

<script>
let token = "";

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
};

function authHeaders() {
  return {
    "Content-Type": "application/json",
    "Authorization": "Bearer " + token
  };
}

async function loadCategories() {
  const res = await fetch("/categories");
  const data = await res.json();
  const el = document.getElementById("categoryList");
  el.innerHTML = "";
  data.forEach(function(c) {
    el.innerHTML += "<div class='card'>" + c.name + "</div>";
  });
}

async function loadProducts() {
  const res = await fetch("/products");
  const data = await res.json();
  const el = document.getElementById("productList");
  el.innerHTML = "";
  data.forEach(function(p) {
    el.innerHTML += "<div class='card'>" + p.name + " - " + p.price + "</div>";
  });
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

loadCategories();
loadProducts();
</script>

</body>
</html>
`)
	}
}
