package web

import (
	"net/http"
	"strings"
)

const faviconB64 = "PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHZpZXdCb3g9IjAgMCA2NCA2NCI+CiAgPGNpcmNsZSBjeD0iMzIiIGN5PSIzMiIgcj0iMzAiIGZpbGw9IiMyODVBNDgiLz4KICA8Y2lyY2xlIGN4PSIzMiIgY3k9IjMyIiByPSIyNSIgZmlsbD0iI0VBRjdFRiIvPgogIDxjaXJjbGUgY3g9IjMyIiBjeT0iMzIiIHI9IjIyIiBmaWxsPSIjNDA4QTcxIi8+CiAgPGcgc3Ryb2tlPSIjRUFGN0VGIiBzdHJva2Utd2lkdGg9IjIiIHN0cm9rZS1saW5lY2FwPSJyb3VuZCI+CiAgICA8bGluZSB4MT0iMzIiIHkxPSIzMiIgeDI9IjMyIiB5Mj0iMTAiLz4KICAgIDxsaW5lIHgxPSIzMiIgeTE9IjMyIiB4Mj0iNTIiIHkyPSIyMCIvPgogICAgPGxpbmUgeDE9IjMyIiB5MT0iMzIiIHgyPSI1MiIgeTI9IjQ0Ii8+CiAgICA8bGluZSB4MT0iMzIiIHkxPSIzMiIgeDI9IjMyIiB5Mj0iNTQiLz4KICAgIDxsaW5lIHgxPSIzMiIgeTE9IjMyIiB4Mj0iMTIiIHkyPSI0NCIvPgogICAgPGxpbmUgeDE9IjMyIiB5MT0iMzIiIHgyPSIxMiIgeTI9IjIwIi8+CiAgPC9nPgogIDxjaXJjbGUgY3g9IjMyIiBjeT0iMzIiIHI9IjMiIGZpbGw9IiNCMEU0Q0MiLz4KPC9zdmc+Cg=="

const pageHTML = `<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<link rel="icon" type="image/svg+xml" href="data:image/svg+xml;base64,FAVICON_PLACEHOLDER">
<title>Notes</title>
<style>
  :root {
    --bg: #091413;
    --surface: #0f201e;
    --surface-2: #142c29;
    --border: #285A48;
    --accent: #408A71;
    --accent-bright: #B0E4CC;
    --text: #EAF7EF;
    --text-dim: #7fa896;
    --down: #e8836b;
  }
  * { box-sizing: border-box; }
  html, body { height: 100%; }
  body {
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
    background: var(--bg);
    color: var(--text);
    margin: 0;
    min-height: 100vh;
    display: flex;
    justify-content: center;
    padding: clamp(20px, 4vw, 56px) 20px;
  }
  .wrap { width: 100%; max-width: 640px; }
  .header { display: flex; align-items: center; gap: 12px; margin-bottom: 28px; }
  .header svg { width: 30px; height: 30px; flex-shrink: 0; }
  h1 { margin: 0; font-size: 22px; font-weight: 600; letter-spacing: -0.02em; }
  .spacer { flex: 1; }
  button.linklike {
    background: none; border: none; color: var(--text-dim); font-size: 13px;
    cursor: pointer; text-decoration: underline; padding: 0;
  }
  button.linklike:hover { color: var(--text); }

  .card {
    background: var(--surface);
    border: 1px solid var(--border);
    border-radius: 14px;
    padding: 28px;
  }
  .tabs { display: flex; gap: 8px; margin-bottom: 20px; }
  .tab {
    flex: 1; text-align: center; padding: 10px; border-radius: 8px;
    cursor: pointer; font-size: 14px; color: var(--text-dim); border: 1px solid transparent;
  }
  .tab.active { background: var(--surface-2); color: var(--text); border-color: var(--border); }

  input, textarea {
    width: 100%; background: var(--surface-2); border: 1px solid var(--border);
    color: var(--text); border-radius: 8px; padding: 11px 14px; font-size: 14px;
    margin-bottom: 12px; font-family: inherit;
  }
  input::placeholder, textarea::placeholder { color: var(--text-dim); }
  textarea { resize: vertical; min-height: 60px; }

  button.primary {
    width: 100%; background: var(--accent); color: white; border: none;
    border-radius: 8px; padding: 12px; font-size: 14px; font-weight: 600;
    cursor: pointer; transition: background 0.15s ease;
  }
  button.primary:hover { background: var(--accent-bright); color: var(--bg); }
  button.primary:disabled { opacity: 0.5; cursor: not-allowed; }

  .error { color: var(--down); font-size: 13px; margin-bottom: 12px; min-height: 16px; }

  .app-bar { display: flex; align-items: center; gap: 10px; margin-bottom: 20px; }
  .app-bar input { margin-bottom: 0; }
  .filter-row { display: flex; gap: 8px; margin-bottom: 20px; }
  .filter-chip {
    padding: 7px 14px; border-radius: 999px; font-size: 13px; cursor: pointer;
    border: 1px solid var(--border); color: var(--text-dim); background: var(--surface);
  }
  .filter-chip.active { background: var(--accent); color: white; border-color: var(--accent); }

  .task-list { display: flex; flex-direction: column; gap: 10px; }
  .task-card {
    background: var(--surface); border: 1px solid var(--border); border-radius: 12px;
    padding: 16px 18px; display: flex; gap: 12px; align-items: flex-start;
  }
  .task-check {
    width: 20px; height: 20px; border-radius: 50%; border: 2px solid var(--accent);
    flex-shrink: 0; margin-top: 2px; cursor: pointer; display: flex; align-items: center; justify-content: center;
  }
  .task-check.done { background: var(--accent-bright); border-color: var(--accent-bright); }
  .task-check.done::after { content: "✓"; color: var(--bg); font-size: 13px; font-weight: 700; }
  .task-body { flex: 1; min-width: 0; }
  .task-title { font-size: 15px; font-weight: 600; word-break: break-word; }
  .task-title.done { text-decoration: line-through; color: var(--text-dim); }
  .task-notes { font-size: 13px; color: var(--text-dim); margin-top: 4px; word-break: break-word; }
  .task-delete {
    background: none; border: none; color: var(--text-dim); cursor: pointer; font-size: 18px;
    padding: 2px 6px; flex-shrink: 0;
  }
  .task-delete:hover { color: var(--down); }
  .empty-state { text-align: center; color: var(--text-dim); padding: 40px 20px; }
</style>
</head>
<body>
<div class="wrap">
  <div class="header">
    <svg viewBox="0 0 240 240" xmlns="http://www.w3.org/2000/svg">
      <circle cx="120" cy="120" r="110" fill="#285A48"/>
      <circle cx="120" cy="120" r="96" fill="#EAF7EF"/>
      <circle cx="120" cy="120" r="88" fill="#408A71"/>
      <g stroke="#EAF7EF" stroke-width="2.5" stroke-linecap="round">
        <line x1="120" y1="120" x2="120" y2="32"/>
        <line x1="120" y1="120" x2="182" y2="58"/>
        <line x1="120" y1="120" x2="208" y2="120"/>
        <line x1="120" y1="120" x2="182" y2="182"/>
        <line x1="120" y1="120" x2="120" y2="208"/>
        <line x1="120" y1="120" x2="58" y2="182"/>
        <line x1="120" y1="120" x2="32" y2="120"/>
        <line x1="120" y1="120" x2="58" y2="58"/>
      </g>
      <circle cx="120" cy="120" r="8" fill="#B0E4CC"/>
    </svg>
    <h1>Notes</h1>
    <div class="spacer"></div>
    <button class="linklike" id="logout-btn" style="display:none">Log out</button>
  </div>

  <div id="auth-view">
    <div class="card">
      <div class="tabs">
        <div class="tab active" data-tab="login">Log in</div>
        <div class="tab" data-tab="signup">Sign up</div>
      </div>
      <div class="error" id="auth-error"></div>
      <form id="auth-form">
        <input id="auth-username" placeholder="Username" autocomplete="username">
        <input id="auth-password" type="password" placeholder="Password" autocomplete="current-password">
        <button class="primary" id="auth-submit" type="submit">Log in</button>
      </form>
    </div>
  </div>

  <div id="app-view" style="display:none">
    <div class="card" style="margin-bottom: 20px;">
      <input id="new-title" placeholder="Task title">
      <textarea id="new-notes" placeholder="Notes (optional)"></textarea>
      <button class="primary" id="add-task-btn">Add task</button>
    </div>

    <div class="app-bar">
      <input id="search-input" placeholder="Search tasks...">
    </div>
    <div class="filter-row">
      <div class="filter-chip active" data-filter="all">All</div>
      <div class="filter-chip" data-filter="false">Active</div>
      <div class="filter-chip" data-filter="true">Done</div>
    </div>

    <div class="task-list" id="task-list"></div>
  </div>
</div>

<script>
let accessToken = localStorage.getItem('access_token');
let refreshToken = localStorage.getItem('refresh_token');
let currentFilter = 'all';
let currentSearch = '';
let authMode = 'login';

function saveTokens(pair) {
  accessToken = pair.access_token;
  refreshToken = pair.refresh_token;
  localStorage.setItem('access_token', accessToken);
  localStorage.setItem('refresh_token', refreshToken);
}

function clearTokens() {
  accessToken = null;
  refreshToken = null;
  localStorage.removeItem('access_token');
  localStorage.removeItem('refresh_token');
}

async function api(path, opts = {}) {
  opts.headers = opts.headers || {};
  if (accessToken) opts.headers['Authorization'] = 'Bearer ' + accessToken;
  if (opts.body && typeof opts.body !== 'string') opts.body = JSON.stringify(opts.body);
  if (opts.body) opts.headers['Content-Type'] = 'application/json';

  let res = await fetch(path, opts);

  if (res.status === 401 && refreshToken && path !== '/auth/refresh') {
    const refreshRes = await fetch('/auth/refresh', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ refresh_token: refreshToken })
    });
    if (refreshRes.ok) {
      saveTokens(await refreshRes.json());
      opts.headers['Authorization'] = 'Bearer ' + accessToken;
      res = await fetch(path, opts);
    } else {
      clearTokens();
      showAuthView();
      throw new Error('session expired');
    }
  }
  return res;
}

function showAuthView() {
  document.getElementById('auth-view').style.display = 'block';
  document.getElementById('app-view').style.display = 'none';
  document.getElementById('logout-btn').style.display = 'none';
}

function showAppView() {
  document.getElementById('auth-view').style.display = 'none';
  document.getElementById('app-view').style.display = 'block';
  document.getElementById('logout-btn').style.display = 'inline';
  loadTasks();
}

document.querySelectorAll('.tab').forEach(tab => {
  tab.addEventListener('click', () => {
    document.querySelectorAll('.tab').forEach(t => t.classList.remove('active'));
    tab.classList.add('active');
    authMode = tab.dataset.tab;
    document.getElementById('auth-submit').textContent = authMode === 'login' ? 'Log in' : 'Sign up';
    document.getElementById('auth-error').textContent = '';
  });
});

document.getElementById('auth-form').addEventListener('submit', async (e) => {
  e.preventDefault();
  const username = document.getElementById('auth-username').value.trim();
  const password = document.getElementById('auth-password').value;
  const errEl = document.getElementById('auth-error');
  errEl.textContent = '';

  if (!username || !password) {
    errEl.textContent = 'Username and password are required';
    return;
  }

  try {
    const res = await fetch('/auth/' + authMode, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username, password })
    });
    const data = await res.json();
    if (!res.ok) {
      errEl.textContent = data.error || 'something went wrong';
      return;
    }
    saveTokens(data);
    showAppView();
  } catch (e) {
    errEl.textContent = 'network error, try again';
  }
});

document.getElementById('logout-btn').addEventListener('click', () => {
  clearTokens();
  showAuthView();
});

document.getElementById('add-task-btn').addEventListener('click', async () => {
  const title = document.getElementById('new-title').value.trim();
  const notes = document.getElementById('new-notes').value.trim();
  if (!title) return;

  const res = await api('/tasks', { method: 'POST', body: { title, notes } });
  if (res.ok) {
    document.getElementById('new-title').value = '';
    document.getElementById('new-notes').value = '';
    loadTasks();
  }
});

document.querySelectorAll('.filter-chip').forEach(chip => {
  chip.addEventListener('click', () => {
    document.querySelectorAll('.filter-chip').forEach(c => c.classList.remove('active'));
    chip.classList.add('active');
    currentFilter = chip.dataset.filter;
    loadTasks();
  });
});

let searchDebounce;
document.getElementById('search-input').addEventListener('input', (e) => {
  clearTimeout(searchDebounce);
  searchDebounce = setTimeout(() => {
    currentSearch = e.target.value.trim();
    loadTasks();
  }, 300);
});

async function loadTasks() {
  const params = new URLSearchParams();
  if (currentFilter !== 'all') params.set('done', currentFilter);
  if (currentSearch) params.set('q', currentSearch);
  params.set('limit', '100');

  try {
    const res = await api('/tasks?' + params.toString());
    if (!res.ok) return;
    const data = await res.json();
    renderTasks(data.tasks || []);
  } catch (e) {}
}

function renderTasks(tasks) {
  const list = document.getElementById('task-list');
  if (tasks.length === 0) {
    list.innerHTML = '<div class="empty-state">No tasks yet</div>';
    return;
  }
  list.innerHTML = tasks.map(function(t) {
    var checkClass = 'task-check' + (t.done ? ' done' : '');
    var titleClass = 'task-title' + (t.done ? ' done' : '');
    var notesHtml = t.notes ? '<div class="task-notes">' + escapeHtml(t.notes) + '</div>' : '';
    return '<div class="task-card">' +
      '<div class="' + checkClass + '" onclick="toggleDone(\'' + t.id + '\', ' + (!t.done) + ')"></div>' +
      '<div class="task-body">' +
      '<div class="' + titleClass + '">' + escapeHtml(t.title) + '</div>' +
      notesHtml +
      '</div>' +
      '<button class="task-delete" onclick="deleteTask(\'' + t.id + '\')">&times;</button>' +
      '</div>';
  }).join('');
}

function escapeHtml(s) {
  const div = document.createElement('div');
  div.textContent = s;
  return div.innerHTML;
}

async function toggleDone(id, done) {
  await api('/tasks/' + id, { method: 'PUT', body: { done } });
  loadTasks();
}

async function deleteTask(id) {
  await api('/tasks/' + id, { method: 'DELETE' });
  loadTasks();
}

if (accessToken) {
  showAppView();
} else {
  showAuthView();
}
</script>
</body>
</html>
`

func Handler() http.HandlerFunc {
	html := []byte(strings.Replace(pageHTML, "FAVICON_PLACEHOLDER", faviconB64, 1))
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(html)
	}
}
