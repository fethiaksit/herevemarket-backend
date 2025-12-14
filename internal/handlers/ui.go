package handlers

import "github.com/gin-gonic/gin"

// Home renders a lightweight HTML landing page to help humans explore the API without other tooling.
func Home() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(200, `<!doctype html>
<html lang="en">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>Hereve Market API</title>
  <style>
    :root { color-scheme: light dark; }
    body { font-family: system-ui, -apple-system, "Segoe UI", sans-serif; margin: 24px; line-height: 1.5; }
    header { margin-bottom: 24px; }
    h1 { margin: 0 0 8px 0; }
    section { margin: 16px 0; padding: 16px; border: 1px solid #d1d5db; border-radius: 10px; }
    code { background: #11182710; padding: 2px 6px; border-radius: 6px; }
    .endpoint { display: flex; align-items: center; gap: 8px; margin: 6px 0; }
    .method { font-weight: 700; padding: 2px 8px; border-radius: 999px; font-size: 0.9rem; }
    .GET { background: #d1fae5; color: #065f46; }
    .POST { background: #e0f2fe; color: #075985; }
    .PUT { background: #fef3c7; color: #92400e; }
    .DELETE { background: #fee2e2; color: #991b1b; }
    details summary { cursor: pointer; }
    form { display: grid; gap: 8px; margin-top: 8px; }
    input, textarea { padding: 8px 10px; border-radius: 8px; border: 1px solid #d1d5db; font-size: 1rem; }
    button { padding: 10px 14px; border: none; border-radius: 8px; font-weight: 600; background: #2563eb; color: white; cursor: pointer; }
    footer { margin-top: 24px; color: #6b7280; font-size: 0.95rem; }
  </style>
</head>
<body>
  <header>
    <h1>Hereve Market API</h1>
    <p>Hızlıca başlamak için aşağıdaki uç noktaları ve örnek istekleri kullanabilirsiniz.</p>
  </header>

  <section>
    <h2>Hızlı Bakış</h2>
    <div class="endpoint"><span class="method GET">GET</span><code>/products</code></div>
    <div class="endpoint"><span class="method GET">GET</span><code>/categories</code></div>
    <div class="endpoint"><span class="method POST">POST</span><code>/admin/login</code></div>
    <div class="endpoint"><span class="method GET">GET</span><code>/admin/products</code> <small>(JWT gerekli)</small></div>
    <div class="endpoint"><span class="method POST">POST</span><code>/admin/categories</code> <small>(JWT gerekli)</small></div>
  </section>

  <section>
    <h2>JWT Alma</h2>
    <details open>
      <summary>Örnek admin girişi</summary>
      <p><code>POST /admin/login</code> gövdesi:</p>
      <pre>{
  "username": "demo",
  "password": "demo"
}</pre>
      <p>Yanıttaki <code>token</code> değerini aşağıdaki isteklerde kullanın.</p>
    </details>
  </section>

  <section>
    <h2>Canlı İstekler</h2>
    <details open>
      <summary>Kategorileri getir (public)</summary>
      <p><code>GET /categories</code></p>
    </details>

    <details>
      <summary>Kategori oluştur (admin)</summary>
      <form id="create-category">
        <label>JWT Token<input name="token" placeholder="Bearer ..." required /></label>
        <label>Adı<input name="name" placeholder="Örn: Elektronik" required /></label>
        <label>Aktif mi?<input name="isActive" value="true" /></label>
        <button type="submit">Kategori oluştur</button>
      </form>
      <pre id="create-category-result"></pre>
    </details>
  </section>

  <footer>
    <p>Postman / curl yerine hızlı denemeler için bu sayfayı kullanabilirsiniz. Alanlar isteğe göre düzenlenebilir.</p>
  </footer>

  <script>
    const form = document.getElementById('create-category');
    const output = document.getElementById('create-category-result');
    form?.addEventListener('submit', async (e) => {
      e.preventDefault();
      const data = new FormData(form);
      const body = {
        name: data.get('name'),
        isActive: data.get('isActive') === 'true'
      };
      try {
        const res = await fetch('/admin/categories', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': data.get('token')
          },
          body: JSON.stringify(body)
        });
        const json = await res.json();
        output.textContent = JSON.stringify(json, null, 2);
      } catch (err) {
        output.textContent = err?.message || 'İstek gönderilemedi';
      }
    });
  </script>
</body>
</html>`)
	}
}
