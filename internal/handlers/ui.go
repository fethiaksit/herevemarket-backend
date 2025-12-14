package handlers

import "github.com/gin-gonic/gin"

// Home renders a lightweight HTML landing page to explore the API.
func Home() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(200, `<!doctype html>
<html lang="tr">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>Hereve Market API</title>
  <style>
    :root { color-scheme: light; }
    body {
      font-family: system-ui, -apple-system, "Segoe UI", sans-serif;
      margin: 0;
      background: #f7f9fc;
      color: #0f172a;
    }
    header {
      padding: 24px;
      background: #0b2d66;
      color: #fff;
    }
    h1 { margin: 0 0 6px 0; }
    main { padding: 24px; max-width: 1000px; margin: 0 auto; display: grid; gap: 20px; }
    section {
      background: #fff;
      border-radius: 12px;
      border: 1px solid #e2e8f0;
      padding: 16px;
    }
    details summary { cursor: pointer; font-weight: 700; }
    label { display: grid; gap: 6px; font-weight: 600; }
    input {
      padding: 8px 10px;
      border-radius: 8px;
      border: 1px solid #cbd5e1;
    }
    button {
      padding: 10px 14px;
      border: none;
      border-radius: 8px;
      font-weight: 700;
      background: #2563eb;
      color: white;
      cursor: pointer;
    }
    pre {
      background: #020617;
      color: #e5e7eb;
      padding: 12px;
      border-radius: 8px;
      overflow: auto;
    }
    .list { display: grid; gap: 8px; }
    .card {
      border: 1px solid #e2e8f0;
      border-radius: 8px;
      padding: 10px;
    }
  </style>
</head>
<body>
<header>
  <h1>Hereve Market API</h1>
  <p>JWT girerek ürün ve kategori CRUD işlemlerini test edebilirsiniz.</p>
</header>

<main>
  <section>
    <label>Admin JWT Token
      <input id="token" placeholder="Bearer <token>" />
    </label>
  </section>

  <section>
    <details open>
      <summary>Kategoriler (GET /categories)</summary>
      <div id="category-list" class="list">Yükleniyor...</div>
    </details>
  </section>

  <section>
    <details open>
      <summary>Ürünler (GET /products)</summary>
      <div id="product-list" class="list">Yükleniyor...</div>
    </details>
  </section>
</main>

<script>
const tokenInput = document.getElementById('token');
const categoryList = document.getElementById('category-list');
const productList = document.getElementById('product-list');

const render = (items, el) => {
  if (!Array.isArray(items) || items.length === 0) {
    el.textContent = 'Kayıt yok';
    return;
  }
  el.innerHTML = items.map(i =>
    '<div class="card"><strong>' + (i.name || '-') +
    '</strong><div>Aktif: ' + i.isActive + '</div></div>'
  ).join('');
};

fetch('/categories').then(r => r.json()).then(d => render(d, categoryList));
fetch('/products').then(r => r.json()).then(d => render(d, productList));
</script>
</body>
</html>`)
	}
}
