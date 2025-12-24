const TOKEN_KEY = "admin_token";

function getToken() {
  return localStorage.getItem(TOKEN_KEY) || "";
}

function setToken(token) {
  if (token) {
    localStorage.setItem(TOKEN_KEY, token);
  }
}

function clearToken() {
  localStorage.removeItem(TOKEN_KEY);
}

function redirectToLogin() {
  window.location.href = "/admin/login";
}

function redirectToCategories() {
  window.location.href = "/admin/categories";
}

function requireAuth() {
  if (!getToken()) {
    redirectToLogin();
  }
}

function redirectIfAuthenticated() {
  if (getToken()) {
    redirectToCategories();
  }
}

function authHeaders() {
  return {
    "Content-Type": "application/json",
    "Authorization": "Bearer " + getToken(),
  };
}

async function safeJson(res) {
  try {
    return await res.json();
  } catch {
    return null;
  }
}

function getId(obj) {
  return obj && (obj._id || obj.id) ? (obj._id || obj.id) : null;
}

function handleUnauthorized(res) {
  if (res && res.status === 401) {
    clearToken();
    redirectToLogin();
    return true;
  }
  return false;
}

function setText(id, text) {
  const el = document.getElementById(id);
  if (el) el.innerText = text || "";
}

function logout() {
  clearToken();
  redirectToLogin();
}

function bindLogout() {
  const logoutBtn = document.getElementById("logoutButton");
  if (logoutBtn) {
    logoutBtn.addEventListener("click", logout);
  }
}

document.addEventListener("DOMContentLoaded", bindLogout);
