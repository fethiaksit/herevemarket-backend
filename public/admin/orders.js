requireAuth();

const ORDERS_API_URL = "/orders";
const DELETE_ORDER_API_URL = "/admin/api/orders";

function formatDateTime(value) {
  const date = value ? new Date(value) : null;
  if (!date || Number.isNaN(date.getTime())) return "-";
  return date.toLocaleString("tr-TR");
}

function formatCurrency(value) {
  if (typeof value !== "number") return "-";
  return value.toLocaleString("tr-TR", {
    style: "currency",
    currency: "TRY",
    minimumFractionDigits: 2,
  });
}

function statusBadgeClass(status) {
  const normalized = (status || "").toLowerCase();
  if (normalized === "completed") return "badge completed";
  if (normalized === "canceled" || normalized === "cancelled") return "badge canceled";
  return "badge pending";
}

function tableBody() {
  return document.getElementById("ordersTableBody");
}

function clearTable() {
  const tbody = tableBody();
  if (tbody) tbody.innerHTML = "";
}

function addEmptyRow(message) {
  const tbody = tableBody();
  if (!tbody) return;
  const row = document.createElement("tr");
  const cell = document.createElement("td");
  cell.colSpan = 8;
  cell.className = "muted";
  cell.textContent = message;
  row.appendChild(cell);
  tbody.appendChild(row);
}

function renderOrders(orders) {
  clearTable();
  if (!Array.isArray(orders) || orders.length === 0) {
    addEmptyRow("Sipariş yok");
    return;
  }

  const tbody = tableBody();
  if (!tbody) return;

  orders.forEach((order) => {
    const row = document.createElement("tr");
    const orderId = getId(order) || "-";
    const customerTitle = order && order.customer && order.customer.title ? order.customer.title : "-";
    const itemCount = Array.isArray(order && order.items) ? order.items.length : 0;

    const cells = [
      orderId,
      formatDateTime(order && order.createdAt),
      customerTitle,
      order && order.paymentMethod ? order.paymentMethod : "-",
      itemCount,
      formatCurrency(order && order.totalPrice),
    ];

    cells.forEach((value, index) => {
      const cell = document.createElement("td");
      if (index >= 4 && index <= 5) {
        cell.className = "numeric";
      }
      cell.textContent = value;
      row.appendChild(cell);
    });

    const statusCell = document.createElement("td");
    const badge = document.createElement("span");
    badge.className = statusBadgeClass(order && order.status);
    badge.textContent = order && order.status ? order.status : "Bilinmiyor";
    statusCell.appendChild(badge);
    row.appendChild(statusCell);

    const actionCell = document.createElement("td");
    actionCell.className = "actions";
    const deleteButton = document.createElement("button");
    deleteButton.type = "button";
    deleteButton.className = "small danger";
    deleteButton.textContent = "Sil";
    deleteButton.addEventListener("click", () => {
      if (orderId === "-") return;
      deleteOrder(orderId);
    });
    actionCell.appendChild(deleteButton);
    row.appendChild(actionCell);

    tbody.appendChild(row);
  });
}

async function loadOrders() {
  setText("ordersStatus", "Siparişler yükleniyor...");
  const res = await fetch(ORDERS_API_URL, { headers: authHeaders() });
  if (handleUnauthorized(res)) return;
  const payload = await safeJson(res);
  if (!res.ok) {
    setText("ordersStatus", "Hata: siparişler getirilemedi");
    addEmptyRow("Siparişler yüklenemedi");
    return;
  }

  const data = payload && payload.data ? payload.data : payload || [];
  renderOrders(data);
  setText("ordersStatus", "");
}

async function deleteOrder(orderId) {
  if (!window.confirm("Sipariş silinsin mi?")) {
    return;
  }

  setText("ordersStatus", "Sipariş siliniyor...");
  const res = await fetch(`${DELETE_ORDER_API_URL}/${orderId}`, {
    method: "DELETE",
    headers: authHeaders(),
  });

  if (handleUnauthorized(res)) return;
  if (!res.ok) {
    setText("ordersStatus", "Hata: sipariş silinemedi");
    return;
  }

  setText("ordersStatus", "");
  await loadOrders();
}

document.addEventListener("DOMContentLoaded", loadOrders);
