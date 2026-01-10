requireAuth();

let selectedProduct = null;
let currentProducts = [];
let currentPage = 1;
let totalPages = 1;
let totalCount = 0;
const pageSize = 20;

function setProductStatus(text) {
  setText("productStatus", text || "");
}

function normalizeCategoryValues(values) {
  if (Array.isArray(values)) {
    return values.filter(function(value) { return !!value; });
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
    .filter(function(value) { return !!value; });
}

function parseStockValue(value) {
  const stock = Number(value);
  if (!Number.isFinite(stock) || stock < 0) return null;
  return stock;
}

function normalizeBarcode(value) {
  if (value === null || value === undefined) return "";
  return String(value).trim();
}

function normalizeBrand(value) {
  if (value === null || value === undefined) return "";
  return String(value).trim();
}

function normalizeDescription(value) {
  if (value === null || value === undefined) return "";
  return String(value).trim();
}

function buildProductFormData(values) {
  const formData = new FormData();
  formData.set("name", values.name);
  formData.set("price", String(values.price));
  formData.set("brand", values.brand || "");
  formData.set("barcode", values.barcode || "");
  formData.set("description", values.description || "");
  formData.set("stock", String(values.stock));
  values.categories.forEach(function(category) {
    formData.append("category", category);
  });
  if (values.isCampaign !== undefined) {
    formData.set("isCampaign", values.isCampaign ? "true" : "false");
  }
  if (values.isActive !== undefined) {
    formData.set("isActive", values.isActive ? "true" : "false");
  }
  if (values.imageFile) {
    formData.set("image", values.imageFile, values.imageFile.name);
  }
  return formData;
}

// ✅ DÜZELTME 1: targetSelect parametresi eklendi. Sadece istenen kutuyu günceller.
async function populateProductCategorySelects(selectedValues, preloadedCategories, targetSelect) {
  const desiredSelection = normalizeCategoryValues(selectedValues);
  const categoryData = Array.isArray(preloadedCategories) && preloadedCategories.length > 0
    ? preloadedCategories
    : null;

  let categories = categoryData;

  if (!categories) {
    const res = await fetch("/categories");
    if (handleUnauthorized(res)) return;
    const payload = await safeJson(res);
    categories = (payload && payload.data) ? payload.data : (payload || []);
  }

  const activeCategories = (categories || []).filter(function(category) { return category && category.isActive; });
  const activeNames = new Set(activeCategories.map(function(category) { return category.name; }));

  // Eğer hedef belirtildiyse onu, yoksa hepsini seç (eski uyumluluk)
  const selects = targetSelect ? [targetSelect] : document.querySelectorAll(".product-category-select");

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

    activeCategories.forEach(function(category) {
      const opt = document.createElement("option");
      opt.value = category.name;
      opt.textContent = category.name;
      select.appendChild(opt);
    });

    preserved.forEach(function(value) {
      if (!activeNames.has(value)) return;
      const opt = Array.from(select.options).find(function(option) { return option.value === value; });
      if (opt) opt.selected = true;
    });
  });
}

async function loadCategories() {
  const filterSelect = document.getElementById("categoryFilter");
  const preserved = filterSelect ? filterSelect.value : "";

  const res = await fetch("/categories");
  if (handleUnauthorized(res)) return;
  const payload = await safeJson(res);
  const data = (payload && payload.data) ? payload.data : (payload || []);

  // ✅ Sadece "Yeni Ürün Ekle" formundaki select'i doldur
  const addProductSelect = document.getElementById("addProductCategorySelect");
  if (addProductSelect) {
    await populateProductCategorySelects(undefined, data, addProductSelect);
  }

  if (filterSelect) {
    filterSelect.innerHTML = "";
    const def = document.createElement("option");
    def.value = "";
    def.textContent = "Tüm Kategoriler";
    filterSelect.appendChild(def);

    (data || []).forEach(function(category) {
      const opt = document.createElement("option");
      opt.value = category.name;
      opt.textContent = category.name;
      filterSelect.appendChild(opt);
    });

    const exists = (data || []).some(function(category){ return category.name === preserved; });
    filterSelect.value = exists ? preserved : "";
  }
}

function buildProductsUrl(page) {
  const selected = document.getElementById("categoryFilter").value;
  const params = new URLSearchParams({
    page: String(page),
    limit: String(pageSize)
  });

  if (selected) {
    params.set("category", selected);
  }

  return "/admin/api/products?" + params.toString();
}

function renderPagination() {
  const container = document.getElementById("productPagination");
  if (!container) return;

  container.innerHTML = "";

  const label = document.createElement("span");
  label.className = "pagination-label";
  label.textContent = "Page " + currentPage + " / " + totalPages;
  container.appendChild(label);

  const addButton = function(label, page, disabled) {
    const button = document.createElement("button");
    button.type = "button";
    button.textContent = label;
    button.disabled = disabled;
    button.addEventListener("click", function() {
      if (page === currentPage) return;
      loadProducts(page);
    });
    container.appendChild(button);
  };

  // Kısa sayfalama kontrolü (Önceki / Sonraki)
  addButton("Önceki", currentPage - 1, currentPage <= 1);
  addButton("Sonraki", currentPage + 1, currentPage >= totalPages);
}

async function toggleCampaign(checkbox) {
  const checked = checkbox.checked;
  const id = checkbox.dataset.id;

  if (!id) {
    checkbox.checked = !checked;
    alert("Ürün id yok");
    return;
  }

  checkbox.disabled = true;

  try {
    console.log("Toggle campaign payload:", { isCampaign: checked });
    const res = await fetch("/admin/api/products/" + id, {
      method: "PUT",
      headers: authHeaders(),
      body: JSON.stringify({
        isCampaign: checked
      })
    });

    console.log("Toggle campaign response status:", res.status);
    const responseBody = await safeJson(res);
    console.log("Toggle campaign response body:", responseBody);

    if (handleUnauthorized(res)) {
      checkbox.checked = !checked;
      alert("Kampanya güncellenemedi");
      return;
    }

    if (!res.ok) {
      checkbox.checked = !checked;
      console.error("Toggle campaign failed:", responseBody || res.statusText);
      alert("Kampanya güncellenemedi");
      return;
    }

    const updated = currentProducts.find(function(item) { return getId(item) === id; });
    if (updated) {
      updated.isCampaign = checked;
    }
  } catch (err) {
    checkbox.checked = !checked;
    alert("Kampanya güncellenemedi");
  } finally {
    checkbox.disabled = false;
  }
}

async function handleQuickSaveProduct(product, fields) {
  if (!product) return;

  const id = getId(product);
  if (!id) {
    alert("Ürün id yok");
    return;
  }

  const stock = parseStockValue(fields.stockInput.value);
  if (stock === null) {
    alert("Stok 0 veya daha büyük olmalı");
    return;
  }

  const payload = {
    stock: stock,
    brand: normalizeBrand(fields.brandInput.value),
    barcode: normalizeBarcode(fields.barcodeInput.value)
  };

  fields.saveButton.disabled = true;
  const originalText = fields.saveButton.textContent;
  fields.saveButton.textContent = "Kaydediliyor...";

  try {
    console.log("Quick save payload:", payload);
    const res = await fetch("/admin/api/products/" + id, {
      method: "PUT",
      headers: authHeaders(),
      body: JSON.stringify(payload)
    });

    if (handleUnauthorized(res)) return;

    console.log("Quick save response status:", res.status);
    const data = await safeJson(res);
    console.log("Quick save response body:", data);
    if (!res.ok) {
      console.error("Quick save failed:", data || res.statusText);
      alert("Güncelleme başarısız: " + ((data && data.error) ? data.error : res.statusText));
      fields.saveButton.textContent = originalText;
      return;
    }

    // Modeli güncelle
    const index = currentProducts.findIndex(function(item) { return getId(item) === id; });
    if (index >= 0 && data) {
      currentProducts[index] = data;
    }

    // ✅ DÜZELTME 3: Tabloyu yeniden çizmek yerine butonu güncelle.
    // Bu sayede diğer satırlardaki veriler kaybolmaz ve odak bozulmaz.
    fields.saveButton.textContent = "Kaydedildi ✓";
    fields.saveButton.style.backgroundColor = "#166534"; // Yeşil renk
    fields.saveButton.style.color = "white";

    setTimeout(() => {
        fields.saveButton.textContent = "Kaydet";
        fields.saveButton.style.backgroundColor = ""; 
        fields.saveButton.style.color = "";
        fields.saveButton.disabled = false;
    }, 2000);

    setProductStatus("Ürün güncellendi");
  } catch (err) {
    console.error(err);
    alert("Güncelleme başarısız");
    fields.saveButton.textContent = originalText;
    fields.saveButton.disabled = false;
  }
}

function renderProductList(data) {
  const table = document.getElementById("productList");
  const tbody = table.querySelector("tbody");
  tbody.innerHTML = "";

  if (!Array.isArray(data) || data.length === 0) {
    const emptyRow = document.createElement("tr");
    const emptyCell = document.createElement("td");
    emptyCell.colSpan = 6;
    emptyCell.className = "muted";
    emptyCell.textContent = "Ürün yok";
    emptyRow.appendChild(emptyCell);
    tbody.appendChild(emptyRow);
    return;
  }

  data.forEach(function(product) {
    const categoryLabel = Array.isArray(product.category)
      ? (product.category.length ? product.category.join(", ") : "-")
      : (product.category || "-");

    const row = document.createElement("tr");
    row.className = "product-row";

    const stockValue = Number.isFinite(Number(product.stock)) ? Number(product.stock) : null;
    if (stockValue === 0) {
      row.classList.add("out-of-stock-row");
    }

    const info = document.createElement("td");
    info.className = "stacked-text clickable";
    info.innerHTML =
      "<div><strong>" + (product.name || "-") + "</strong></div>" +
      "<div class='muted'>" +
        (product.price ?? "-") + " • " + categoryLabel + " • " + (product.isActive ? "Aktif" : "Pasif") +
      "</div>";
    info.onclick = function() { selectProduct(product); };

    // --- Input Alanları ---

    const brandCell = document.createElement("td");
    const brandInput = document.createElement("input");
    brandInput.type = "text";
    brandInput.className = "table-input";
    brandInput.placeholder = "Marka";
    brandInput.value = product.brand || "";
    // ✅ DÜZELTME 2: Input değişince hafızadaki (currentProducts) veriyi de güncelle
    brandInput.addEventListener("input", function(e) {
        product.brand = e.target.value; 
    });
    brandInput.addEventListener("click", function(event) { event.stopPropagation(); });
    brandCell.appendChild(brandInput);

    const barcodeCell = document.createElement("td");
    const barcodeInput = document.createElement("input");
    barcodeInput.type = "text";
    barcodeInput.className = "table-input";
    barcodeInput.placeholder = "Barkod";
    barcodeInput.value = product.barcode || "";
    barcodeInput.addEventListener("input", function(e) {
        product.barcode = e.target.value; 
    });
    barcodeInput.addEventListener("click", function(event) { event.stopPropagation(); });
    barcodeCell.appendChild(barcodeInput);

    const stockCell = document.createElement("td");
    const stockInput = document.createElement("input");
    stockInput.type = "number";
    stockInput.min = "0";
    stockInput.step = "1";
    stockInput.className = "table-input";
    stockInput.placeholder = "Stok";
    stockInput.value = stockValue === null ? "" : stockValue;
    stockInput.addEventListener("input", function(e) {
        product.stock = e.target.value; 
    });
    stockInput.addEventListener("click", function(event) { event.stopPropagation(); });
    stockCell.appendChild(stockInput);

    // --- Kampanya Toggle ---

    const campaignCell = document.createElement("td");
    campaignCell.className = "campaign-cell";
    const campaignToggle = document.createElement("input");
    campaignToggle.type = "checkbox";
    campaignToggle.className = "campaign-toggle";
    campaignToggle.checked = !!product.isCampaign;
    campaignToggle.dataset.id = getId(product) || "";
    campaignToggle.addEventListener("click", function(event) { event.stopPropagation(); });
    campaignToggle.addEventListener("change", function(event) {
      event.stopPropagation();
      toggleCampaign(campaignToggle);
    });
    campaignCell.appendChild(campaignToggle);

    // --- Butonlar ---

    const actions = document.createElement("td");
    actions.className = "inline-actions";

    const saveBtn = document.createElement("button");
    saveBtn.type = "button";
    saveBtn.className = "small";
    saveBtn.textContent = "Kaydet";
    saveBtn.onclick = function(ev) {
      ev.stopPropagation();
      handleQuickSaveProduct(product, {
        brandInput: brandInput,
        barcodeInput: barcodeInput,
        stockInput: stockInput,
        saveButton: saveBtn
      });
    };

    const deleteBtn = document.createElement("button");
    deleteBtn.type = "button";
    deleteBtn.className = "danger ghost small";
    deleteBtn.textContent = "Sil";
    deleteBtn.onclick = function(ev) {
      ev.stopPropagation();
      handleDeleteProduct(product);
    };

    actions.appendChild(saveBtn);
    actions.appendChild(deleteBtn);

    row.appendChild(info);
    row.appendChild(brandCell);
    row.appendChild(barcodeCell);
    row.appendChild(stockCell);
    row.appendChild(campaignCell);
    row.appendChild(actions);

    tbody.appendChild(row);
  });
}

async function loadProducts(page) {
  const targetPage = page || 1;
  const url = buildProductsUrl(targetPage);

  setProductStatus("Ürünler yükleniyor...");

  const res = await fetch(url, { headers: authHeaders() });
  if (handleUnauthorized(res)) return;
  const payload = await safeJson(res);
  if (!res.ok) {
    setProductStatus("Hata: ürünler getirilemedi");
    return;
  }

  const data = (payload && payload.data) ? payload.data : [];
  currentProducts = Array.isArray(data) ? data : [];
  const pagination = (payload && payload.pagination) ? payload.pagination : {};
  currentPage = (pagination && Number.isFinite(pagination.page)) ? pagination.page : targetPage;
  totalPages = (pagination && Number.isFinite(pagination.totalPages)) ? pagination.totalPages : 1;
  totalCount = (pagination && Number.isFinite(pagination.total)) ? pagination.total : currentProducts.length;

  renderProductList(currentProducts);
  renderPagination();
  setProductStatus("");
}

async function selectProduct(product) {
  selectedProduct = product;
  const id = getId(product);

  const categories = normalizeCategoryValues(product.category);

  document.getElementById("editProduct").style.display = "grid";
  document.getElementById("prodName").innerText = product.name || "-";
  document.getElementById("prodId").innerText = id ? ("(id: " + id + ")") : "(id yok)";

  // ✅ Sadece Düzenleme Formunun Select'ini güncelle
  const editSelect = document.getElementById("editProductCategorySelect");
  if (editSelect) {
    await populateProductCategorySelects(categories, undefined, editSelect);
  }

  const form = document.getElementById("editProduct");
  form.elements.name.value = product.name || "";
  form.elements.price.value = (product.price ?? "");
  form.elements.brand.value = product.brand || "";
  form.elements.barcode.value = product.barcode || "";
  form.elements.stock.value = (product.stock ?? "");
  form.elements.imageUrl.value = product.imageUrl || "";
  form.elements.description.value = product.description || "";
  form.elements.isCampaign.checked = !!product.isCampaign;
  form.elements.isActive.checked = !!product.isActive;
}

async function handleDeleteProduct(product) {
  if (!product) return;

  const id = getId(product);
  if (!id) {
    alert("Ürün id yok");
    return;
  }

  const confirmed = confirm("Bu ürünü silmek istediğinize emin misiniz?");
  if (!confirmed) return;

  const res = await fetch("/admin/api/products/" + id, {
    method: "DELETE",
    headers: authHeaders()
  });
  if (handleUnauthorized(res)) return;
  const payload = await safeJson(res);

  if (!res.ok) {
    alert("Silme başarısız: " + ((payload && payload.error) ? payload.error : res.statusText));
    return;
  }

  if (selectedProduct && getId(selectedProduct) === id) {
    selectedProduct = null;
    document.getElementById("editProduct").style.display = "none";
  }

  // Soft delete sonrası listeyi sayfalama ile yenile
  await loadProducts(currentPage);
  setProductStatus("Ürün pasifleştirildi");
}

document.getElementById("categoryFilter").addEventListener("change", function() {
  currentPage = 1;
  loadProducts(currentPage);
});

document.getElementById("addProduct").addEventListener("submit", async function(event) {
  event.preventDefault();

  const form = new FormData(event.target);
  const price = parseFloat(form.get("price"));
  if (Number.isNaN(price)) {
    alert("Fiyat sayı olmalı (örn 24.90)");
    return;
  }

  const stock = parseStockValue(form.get("stock"));
  if (stock === null) {
    alert("Stok 0 veya daha büyük olmalı");
    return;
  }

  const categories = getSelectedCategories(event.target.querySelector('select[name="category"]'));
  if (categories.length === 0) {
    alert("En az bir kategori seç");
    return;
  }

  const barcode = normalizeBarcode(form.get("barcode"));
  const imageInput = event.target.querySelector('input[name="image"]');
  const imageFile = imageInput && imageInput.files ? imageInput.files[0] : null;
  if (!imageFile) {
    alert("Görsel seçmelisiniz");
    return;
  }

  const createPayload = buildProductFormData({
    name: form.get("name"),
    price: price,
    brand: normalizeBrand(form.get("brand")),
    barcode: barcode,
    description: normalizeDescription(form.get("description")),
    stock: stock,
    categories: categories,
    imageFile: imageFile,
    isCampaign: form.get("isCampaign") === "on",
    isActive: true
  });

  console.log("Create product payload:", createPayload);

  const res = await fetch("/admin/api/products", {
    method: "POST",
    headers: { "Authorization": "Bearer " + getToken() },
    body: createPayload
  });
  console.log("Create product response status:", res.status);
  const createBody = await safeJson(res);
  console.log("Create product response body:", createBody);
  if (!res.ok) {
    console.error("Create product failed:", createBody || res.statusText);
  }
  if (handleUnauthorized(res)) return;

  event.target.reset();
  loadProducts(currentPage);
});

document.getElementById("editProduct").addEventListener("submit", async function(event) {
  event.preventDefault();
  if (!selectedProduct) return;

  const id = getId(selectedProduct);
  if (!id) {
    alert("Ürün id yok");
    return;
  }

  const form = new FormData(event.target);
  const price = parseFloat(form.get("price"));
  if (Number.isNaN(price)) {
    alert("Fiyat sayı olmalı");
    return;
  }

  const stock = parseStockValue(form.get("stock"));
  if (stock === null) {
    alert("Stok 0 veya daha büyük olmalı");
    return;
  }

  const categories = getSelectedCategories(event.target.querySelector('select[name="category"]'));
  if (categories.length === 0) {
    alert("En az bir kategori seç");
    return;
  }

  const barcode = normalizeBarcode(form.get("barcode"));
  const imageInput = event.target.querySelector('input[name="image"]');
  const imageFile = imageInput && imageInput.files ? imageInput.files[0] : null;

  const updatePayload = buildProductFormData({
    name: form.get("name"),
    price: price,
    brand: normalizeBrand(form.get("brand")),
    barcode: barcode,
    description: normalizeDescription(form.get("description")),
    stock: stock,
    categories: categories,
    imageFile: imageFile,
    isCampaign: form.get("isCampaign") === "on",
    isActive: form.get("isActive") === "on"
  });

  console.log("Update product payload:", updatePayload);

  const res = await fetch("/admin/api/products/" + id, {
    method: "PUT",
    headers: { "Authorization": "Bearer " + getToken() },
    body: updatePayload
  });
  console.log("Update product response status:", res.status);
  const updateBody = await safeJson(res);
  console.log("Update product response body:", updateBody);
  if (!res.ok) {
    console.error("Update product failed:", updateBody || res.statusText);
  }
  if (handleUnauthorized(res)) return;

  loadProducts(currentPage);
});

document.getElementById("deleteProduct").addEventListener("click", async function() {
  if (!selectedProduct) return;

  await handleDeleteProduct(selectedProduct);
});

loadCategories();
loadProducts(currentPage);
