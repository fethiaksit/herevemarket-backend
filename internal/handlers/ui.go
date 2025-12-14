package handlers

import "github.com/gin-gonic/gin"

// Home renders a lightweight HTML landing page to help humans explore the API without other tooling.
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
      font-family: 'Inter', system-ui, -apple-system, "Segoe UI", sans-serif;
      margin: 0;
      background: linear-gradient(180deg, #f7f9fc 0%, #ffffff 50%, #f1f5f9 100%);
      color: #0f172a;
    }
    header {
      padding: 28px;
      background: #0b2d66;
      color: #f8f9fa;
      border-bottom: 6px solid #f5c600;
    }
    h1 { margin: 0 0 6px 0; font-size: 1.8rem; }
    p { margin: 0; }
    main { padding: 24px; display: grid; gap: 20px; max-width: 1000px; margin: 0 auto; }
    section {
      border: 1px solid #e2e8f0;
      border-radius: 14px;
      background: white;
      box-shadow: 0 10px 30px rgba(12, 18, 38, 0.08);
      overflow: hidden;
    }
    section > header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 16px 20px;
      background: linear-gradient(90deg, #0d47a1, #0b2d66);
      color: #fef9c3;
      border-bottom: 4px solid #f5c600;
    }
    section > header h2 { margin: 0; font-size: 1.2rem; }
    .content { padding: 18px 20px 22px; display: grid; gap: 14px; }
    details { border: 1px solid #e2e8f0; border-radius: 12px; overflow: hidden; }
    details summary {
      background: #0d47a1;
      color: #fef9c3;
      padding: 12px 14px;
      cursor: pointer;
      font-weight: 700;
      display: flex;
      align-items: center;
      gap: 10px;
    }
    details summary::marker { color: #fef9c3; }
    .panel { padding: 14px; background: #fff; }
    label { display: grid; gap: 6px; font-weight: 600; color: #0b2d66; }
    input, textarea {
      padding: 10px 12px;
      border-radius: 10px;
      border: 1px solid #cbd5e1;
      font-size: 1rem;
    }
    input:focus, textarea:focus { outline: 2px solid #f5c600; border-color: #0d47a1; }
    button {
      padding: 12px 16px;
      border: none;
      border-radius: 10px;
      font-weight: 700;
      background: linear-gradient(90deg, #f5c600, #fbbf24);
      color: #0b2d66;
      cursor: pointer;
      transition: transform 120ms ease, box-shadow 120ms ease;
      box-shadow: 0 6px 16px rgba(10, 20, 50, 0.2);
    }
    button:hover { transform: translateY(-1px); box-shadow: 0 8px 20px rgba(10, 20, 50, 0.28); }
    .grid { display: grid; gap: 12px; }
    .grid.two { grid-template-columns: repeat(auto-fit, minmax(240px, 1fr)); }
    pre {
      background: #0b2d66;
      color: #fef9c3;
      padding: 12px;
      border-radius: 10px;
      overflow: auto;
      margin: 0;
      font-size: 0.95rem;
    }
    .chip { display: inline-flex; align-items: center; gap: 6px; padding: 6px 10px; background: #0d47a1; color: #fef9c3; border-radius: 999px; font-weight: 700; }
    .list { display: grid; gap: 10px; }
    .card {
      border: 1px solid #e2e8f0;
      border-radius: 12px;
      padding: 12px;
      display: grid;
      gap: 4px;
      background: linear-gradient(180deg, #fef9c3 0%, #ffffff 42%);
    }
    .card h4 { margin: 0; color: #0b2d66; }
    small { color: #475569; }
    .muted { color: #475569; font-size: 0.95rem; }
  </style>
</head>
<body>
  <header>
    <h1>Hereve Market API</h1>
    <p>Giriş yapıp token'ı otomatik alarak ürün / kategori CRUD işlemlerini doğrudan veritabanına yazın.</p>
  </header>

  <main>
    <section>
      <header>
        <h2>Genel Bilgiler</h2>
        <span class="chip">API Köki /</span>
      </header>
      <div class="content grid">
        <div class="grid">
          <strong>Önemli uç noktalar</strong>
          <div class="muted">/products, /categories, /admin/login, /admin/products, /admin/categories</div>
          <pre>{
  "email": "demo",
  "password": "demo"
}</pre>
          <small>Yukarıdaki gövde ile <code>/admin/login</code> isteğini UI otomatik çağırır; token elle girilmez.</small>
        </div>
        <div class="grid">
          <form id="login-form" class="grid">
            <label>Email
              <input name="email" placeholder="admin@example.com" value="demo" required />
            </label>
            <label>Şifre
              <input name="password" type="password" value="demo" required />
            </label>
            <button type="submit">Giriş yap ve token al</button>
          </form>
          <small id="auth-status" class="muted">Token henüz alınmadı.</small>
        </div>
      </div>
    </section>

    <section>
      <header>
        <h2>Kategoriler</h2>
        <span class="chip">CRUD</span>
      </header>
      <div class="content">
        <details open>
          <summary>Kategorileri görüntüle</summary>
          <div class="panel">
            <div id="category-list" class="list muted">Yükleniyor...</div>
          </div>
        </details>
        <details>
          <summary>Yeni kategori ekle</summary>
          <div class="panel">
            <form id="create-category" class="grid two">
              <label>Adı<input name="name" placeholder="Örn: Elektronik" required /></label>
              <label>Aktif mi?<input name="isActive" value="true" /></label>
              <button type="submit">Kaydet</button>
            </form>
            <pre id="create-category-result" class="muted"></pre>
          </div>
        </details>
        <details>
          <summary>Kategori güncelle</summary>
          <div class="panel">
            <form id="update-category" class="grid two">
              <label>Kategori ID<input name="id" placeholder="ObjectID" required /></label>
              <label>Yeni ad<input name="name" placeholder="Opsiyonel" /></label>
              <label>Aktif mi?<input name="isActive" placeholder="true/false" /></label>
              <button type="submit">Güncelle</button>
            </form>
            <pre id="update-category-result" class="muted"></pre>
          </div>
        </details>
        <details>
          <summary>Kategori sil (pasif yap)</summary>
          <div class="panel">
            <form id="delete-category" class="grid two">
              <label>Kategori ID<input name="id" placeholder="ObjectID" required /></label>
              <button type="submit">Sil</button>
            </form>
            <pre id="delete-category-result" class="muted"></pre>
          </div>
        </details>
      </div>
    </section>

    <section>
      <header>
        <h2>Ürünler</h2>
        <span class="chip">CRUD</span>
      </header>
      <div class="content">
        <details open>
          <summary>Ürünleri görüntüle</summary>
          <div class="panel">
            <div id="product-list" class="list muted">Yükleniyor...</div>
          </div>
        </details>
        <details>
          <summary>Yeni ürün ekle</summary>
          <div class="panel">
            <form id="create-product" class="grid two">
              <label>Adı<input name="name" placeholder="Örn: Akıllı Telefon" required /></label>
              <label>Fiyat<input name="price" type="number" step="0.01" placeholder="5999" required /></label>
              <label>Kategori<input name="category" placeholder="Kategori adı" required /></label>
              <label>Görsel URL<input name="imageUrl" placeholder="https://..." required /></label>
              <label>Aktif mi?<input name="isActive" value="true" /></label>
              <button type="submit">Kaydet</button>
            </form>
            <pre id="create-product-result" class="muted"></pre>
          </div>
        </details>
        <details>
          <summary>Ürün güncelle</summary>
          <div class="panel">
            <form id="update-product" class="grid two">
              <label>Ürün ID<input name="id" placeholder="ObjectID" required /></label>
              <label>Ad<input name="name" placeholder="Opsiyonel" /></label>
              <label>Fiyat<input name="price" type="number" step="0.01" placeholder="Opsiyonel" /></label>
              <label>Kategori<input name="category" placeholder="Opsiyonel" /></label>
              <label>Görsel URL<input name="imageUrl" placeholder="Opsiyonel" /></label>
              <label>Aktif mi?<input name="isActive" placeholder="true/false" /></label>
              <button type="submit">Güncelle</button>
            </form>
            <pre id="update-product-result" class="muted"></pre>
          </div>
        </details>
        <details>
          <summary>Ürün sil (pasif yap)</summary>
          <div class="panel">
            <form id="delete-product" class="grid two">
              <label>Ürün ID<input name="id" placeholder="ObjectID" required /></label>
              <button type="submit">Sil</button>
            </form>
            <pre id="delete-product-result" class="muted"></pre>
          </div>
        </details>
      </div>
    </section>
  </main>

  <script>
    const loginForm = document.getElementById('login-form');
    const authStatus = document.getElementById('auth-status');
    let authToken = '';

    const categoryList = document.getElementById('category-list');
    const productList = document.getElementById('product-list');
    const categoryForm = document.getElementById('create-category');
    const productForm = document.getElementById('create-product');
    const updateCategoryForm = document.getElementById('update-category');
    const deleteCategoryForm = document.getElementById('delete-category');
    const updateProductForm = document.getElementById('update-product');
    const deleteProductForm = document.getElementById('delete-product');
    const categoryOutput = document.getElementById('create-category-result');
    const productOutput = document.getElementById('create-product-result');
    const updateCategoryOutput = document.getElementById('update-category-result');
    const deleteCategoryOutput = document.getElementById('delete-category-result');
    const updateProductOutput = document.getElementById('update-product-result');
    const deleteProductOutput = document.getElementById('delete-product-result');

    const renderList = (items, container, type) => {
      if (!Array.isArray(items) || items.length === 0) {
        container.textContent = type + ' bulunamadı.';
        return;
      }
      container.innerHTML = items.map((item) => {
        const priceBlock = item.price ? '<div><strong>Fiyat:</strong> ' + item.price + '</div>' : '';
        const categoryBlock = item.category ? '<div><strong>Kategori:</strong> ' + item.category + '</div>' : '';
        const idValue = item._id || item.id || '-';
        return [
          '<div class="card">',
          '<h4>' + (item.name || 'İsimsiz') + '</h4>',
          priceBlock,
          categoryBlock,
          '<small>ID: ' + idValue + '</small>',
          '<small>Aktif: ' + item.isActive + '</small>',
          '</div>',
        ].join('');
      }).join('');
    };

    const fetchCategories = async () => {
      categoryList.textContent = 'Yükleniyor...';
      try {
        const res = await fetch('/categories');
        const data = await res.json();
        renderList(data, categoryList, 'Kategori');
      } catch (err) {
        categoryList.textContent = 'Kategoriler alınamadı';
      }
    };

    const fetchProducts = async () => {
      productList.textContent = 'Yükleniyor...';
      try {
        const res = await fetch('/products');
        const data = await res.json();
        renderList(data, productList, 'Ürün');
      } catch (err) {
        productList.textContent = 'Ürünler alınamadı';
      }
    };

    const requestJSON = async (method, url, body) => {
      if (!authToken) throw { error: 'Lütfen önce giriş yapın.' };
      const res = await fetch(url, {
        method,
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'Bearer ' + authToken,
        },
        body: body ? JSON.stringify(body) : undefined,
      });
      const json = await res.json().catch(() => ({}));
      if (!res.ok) throw json;
      return json;
    };

    loginForm?.addEventListener('submit', async (e) => {
      e.preventDefault();
      const data = new FormData(loginForm);
      const payload = {
        email: data.get('email'),
        password: data.get('password'),
      };
      authStatus.textContent = 'Token alınıyor...';
      try {
        const res = await fetch('/admin/login', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(payload),
        });
        const json = await res.json();
        if (!res.ok) throw json;
        authToken = json.token;
        authStatus.textContent = 'Giriş başarılı, token eklendi.';
      } catch (err) {
        authStatus.textContent = 'Giriş başarısız: ' + JSON.stringify(err);
      }
    });

    categoryForm?.addEventListener('submit', async (e) => {
      e.preventDefault();
      const data = new FormData(categoryForm);
      const payload = {
        name: data.get('name'),
        isActive: data.get('isActive') === 'true',
      };
      try {
        const result = await requestJSON('POST', '/admin/categories', payload);
        categoryOutput.textContent = JSON.stringify(result, null, 2);
        fetchCategories();
      } catch (err) {
        categoryOutput.textContent = JSON.stringify(err, null, 2);
      }
    });

    updateCategoryForm?.addEventListener('submit', async (e) => {
      e.preventDefault();
      const data = new FormData(updateCategoryForm);
      const body = {};
      const name = data.get('name');
      const isActiveRaw = data.get('isActive');
      if (name) body.name = name;
      if (isActiveRaw) body.isActive = isActiveRaw === 'true';
      try {
        const result = await requestJSON('PUT', '/admin/categories/' + data.get('id'), body);
        updateCategoryOutput.textContent = JSON.stringify(result, null, 2);
        fetchCategories();
      } catch (err) {
        updateCategoryOutput.textContent = JSON.stringify(err, null, 2);
      }
    });

    deleteCategoryForm?.addEventListener('submit', async (e) => {
      e.preventDefault();
      const data = new FormData(deleteCategoryForm);
      try {
        const result = await requestJSON('DELETE', '/admin/categories/' + data.get('id'));
        deleteCategoryOutput.textContent = JSON.stringify(result || { ok: true }, null, 2);
        fetchCategories();
      } catch (err) {
        deleteCategoryOutput.textContent = JSON.stringify(err, null, 2);
      }
    });

    productForm?.addEventListener('submit', async (e) => {
      e.preventDefault();
      const data = new FormData(productForm);
      const payload = {
        name: data.get('name'),
        price: parseFloat(data.get('price')),
        category: data.get('category'),
        imageUrl: data.get('imageUrl'),
        isActive: data.get('isActive') === 'true',
      };
      try {
        const result = await requestJSON('POST', '/admin/products', payload);
        productOutput.textContent = JSON.stringify(result, null, 2);
        fetchProducts();
      } catch (err) {
        productOutput.textContent = JSON.stringify(err, null, 2);
      }
    });

    updateProductForm?.addEventListener('submit', async (e) => {
      e.preventDefault();
      const data = new FormData(updateProductForm);
      const body = {};
      const name = data.get('name');
      const price = data.get('price');
      const category = data.get('category');
      const imageUrl = data.get('imageUrl');
      const isActiveRaw = data.get('isActive');
      if (name) body.name = name;
      if (price) body.price = parseFloat(price);
      if (category) body.category = category;
      if (imageUrl) body.imageUrl = imageUrl;
      if (isActiveRaw) body.isActive = isActiveRaw === 'true';
      try {
        const result = await requestJSON('PUT', '/admin/products/' + data.get('id'), body);
        updateProductOutput.textContent = JSON.stringify(result, null, 2);
        fetchProducts();
      } catch (err) {
        updateProductOutput.textContent = JSON.stringify(err, null, 2);
      }
    });

    deleteProductForm?.addEventListener('submit', async (e) => {
      e.preventDefault();
      const data = new FormData(deleteProductForm);
      try {
        const result = await requestJSON('DELETE', '/admin/products/' + data.get('id'));
        deleteProductOutput.textContent = JSON.stringify(result || { ok: true }, null, 2);
        fetchProducts();
      } catch (err) {
        deleteProductOutput.textContent = JSON.stringify(err, null, 2);
      }
    });

    fetchCategories();
    fetchProducts();
  </script>
</body>
</html>`)
	}
}
