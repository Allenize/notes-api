package web

import (
	"net/http"
	"strings"
)

const notezLogoB64 = "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"

const pageHTML = `<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<link rel="icon" type="image/png" href="data:image/png;base64,LOGO_SRC_TOKEN">
<link rel="preconnect" href="https://fonts.googleapis.com">
<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
<link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&family=Poppins:wght@600;700&display=swap" rel="stylesheet">
<title>Notez</title>
<style>
  :root {
    --bg: #F6F1E9;
    --surface: #FFFFFF;
    --surface-hover: #FFF8EC;
    --primary: #FFD93D;
    --secondary: #FF9A00;
    --text: #4F200D;
    --text-muted: #96725a;
    --border: #EAE0CF;
    --sidebar-bg: #4F200D;
    --sidebar-text: #F6F1E9;
    --sidebar-text-muted: rgba(246,241,233,0.62);
    --sidebar-active-bg: #FFD93D;
    --sidebar-active-text: #4F200D;
    --radius: 16px;
    --radius-sm: 10px;
    --shadow: 0 2px 10px rgba(79,32,13,0.07), 0 1px 2px rgba(79,32,13,0.05);
    --shadow-hover: 0 10px 28px rgba(79,32,13,0.14);
    --danger: #d1483f;
  }
  * { box-sizing: border-box; }
  html, body { height: 100%; margin: 0; }
  body {
    font-family: 'Inter', -apple-system, sans-serif;
    background: var(--bg);
    color: var(--text);
  }
  h1, h2, h3, .brand-font { font-family: 'Poppins', 'Inter', sans-serif; }
  button { font-family: inherit; cursor: pointer; }
  input, textarea, select { font-family: inherit; }
  a { text-decoration: none; color: inherit; }

  .app-shell { display: none; height: 100vh; overflow: hidden; }
  .app-shell.visible { display: flex; }

  /* Sidebar */
  .sidebar {
    width: 250px;
    flex-shrink: 0;
    background: var(--sidebar-bg);
    color: var(--sidebar-text);
    display: flex;
    flex-direction: column;
    transition: width 0.2s ease;
    overflow: hidden;
  }
  .sidebar.collapsed { width: 76px; }
  .sidebar-header {
    display: flex; align-items: center; justify-content: space-between;
    padding: 22px 20px; flex-shrink: 0;
  }
  .sidebar-header img { height: 28px; display: block; }
  .sidebar.collapsed .sidebar-header img { display: none; }
  .sidebar-toggle-btn {
    background: rgba(246,241,233,0.1); border: none; color: var(--sidebar-text);
    width: 30px; height: 30px; border-radius: 8px; display: flex; align-items: center; justify-content: center;
    flex-shrink: 0;
  }
  .sidebar-toggle-btn:hover { background: rgba(246,241,233,0.2); }

  .sidebar-nav { flex: 1; overflow-y: auto; padding: 4px 12px; }
  .nav-item {
    display: flex; align-items: center; gap: 12px;
    padding: 10px 12px; border-radius: var(--radius-sm);
    color: var(--sidebar-text-muted); font-size: 14px; font-weight: 500;
    margin-bottom: 2px; white-space: nowrap;
  }
  .nav-item svg { flex-shrink: 0; width: 18px; height: 18px; }
  .nav-item:hover { background: rgba(246,241,233,0.08); color: var(--sidebar-text); }
  .nav-item.active { background: var(--sidebar-active-bg); color: var(--sidebar-active-text); }
  .sidebar.collapsed .nav-label, .sidebar.collapsed .nav-section-title, .sidebar.collapsed .user-name,
  .sidebar.collapsed .nav-badge { display: none; }

  .nav-section { margin-top: 18px; }
  .nav-section-title {
    font-size: 11px; text-transform: uppercase; letter-spacing: 0.06em;
    color: var(--sidebar-text-muted); padding: 0 12px; margin-bottom: 6px; font-weight: 600;
  }
  .nav-badge {
    margin-left: auto; background: rgba(246,241,233,0.15); font-size: 11px;
    padding: 1px 7px; border-radius: 999px;
  }
  .nav-item.active .nav-badge { background: rgba(79,32,13,0.15); }

  .sidebar-user {
    display: flex; align-items: center; gap: 10px; padding: 16px 20px;
    border-top: 1px solid rgba(246,241,233,0.1); flex-shrink: 0; cursor: pointer;
  }
  .user-avatar {
    width: 32px; height: 32px; border-radius: 50%; background: var(--primary); color: var(--text);
    display: flex; align-items: center; justify-content: center; font-weight: 700; font-size: 13px; flex-shrink: 0;
  }
  .user-name { font-size: 13px; font-weight: 600; overflow: hidden; text-overflow: ellipsis; }

  /* Main area */
  .main-area { flex: 1; display: flex; flex-direction: column; min-width: 0; }
  .topnav {
    display: flex; align-items: center; gap: 16px; padding: 16px 28px;
    background: var(--surface); border-bottom: 1px solid var(--border); flex-shrink: 0;
  }
  .search-wrap { flex: 1; max-width: 420px; position: relative; }
  .search-wrap svg { position: absolute; left: 14px; top: 50%; transform: translateY(-50%); color: var(--text-muted); }
  .search-wrap input {
    width: 100%; border: 1px solid var(--border); border-radius: 999px;
    padding: 10px 16px 10px 40px; font-size: 14px; background: var(--bg); color: var(--text);
  }
  .search-wrap input:focus { outline: none; border-color: var(--secondary); background: var(--surface); }
  .spacer { flex: 1; }

  .btn {
    border: none; border-radius: var(--radius-sm); padding: 10px 18px; font-size: 14px; font-weight: 600;
    display: inline-flex; align-items: center; gap: 8px; transition: transform 0.1s ease, box-shadow 0.15s ease;
  }
  .btn:active { transform: scale(0.97); }
  .btn-primary { background: var(--primary); color: var(--text); }
  .btn-primary:hover { box-shadow: var(--shadow-hover); }
  .btn-secondary { background: var(--secondary); color: white; }
  .btn-secondary:hover { box-shadow: var(--shadow-hover); }
  .btn-ghost { background: transparent; color: var(--text); border: 1px solid var(--border); }
  .btn-ghost:hover { background: var(--surface-hover); }
  .btn-danger { background: var(--danger); color: white; }
  .btn-danger:hover { opacity: 0.9; }
  .btn-sm { padding: 7px 12px; font-size: 13px; }

  .icon-btn {
    width: 36px; height: 36px; border-radius: 50%; border: none; background: transparent;
    display: flex; align-items: center; justify-content: center; color: var(--text);
  }
  .icon-btn:hover { background: var(--surface-hover); }

  .topnav-user {
    width: 36px; height: 36px; border-radius: 50%; background: var(--primary); color: var(--text);
    display: flex; align-items: center; justify-content: center; font-weight: 700; font-size: 14px; position: relative;
  }
  .user-dropdown {
    position: absolute; top: 46px; right: 0; background: var(--surface); border: 1px solid var(--border);
    border-radius: var(--radius-sm); box-shadow: var(--shadow-hover); min-width: 160px; z-index: 50; display: none; overflow: hidden;
  }
  .user-dropdown.open { display: block; }
  .user-dropdown a, .user-dropdown button {
    display: block; width: 100%; text-align: left; padding: 10px 14px; font-size: 14px; background: none; border: none; color: var(--text);
  }
  .user-dropdown a:hover, .user-dropdown button:hover { background: var(--surface-hover); }

  .content { flex: 1; overflow-y: auto; padding: 28px; }
  .content-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 22px; }
  .content-title { font-size: 22px; font-weight: 700; margin: 0; }
  .content-sub { color: var(--text-muted); font-size: 14px; margin-top: 4px; }

  .stats-row { display: grid; grid-template-columns: repeat(auto-fit, minmax(160px, 1fr)); gap: 16px; margin-bottom: 32px; }
  .stat-card {
    background: var(--surface); border-radius: var(--radius); padding: 20px; box-shadow: var(--shadow);
  }
  .stat-value { font-size: 30px; font-weight: 700; }
  .stat-label { font-size: 13px; color: var(--text-muted); margin-top: 4px; }

  .section { margin-bottom: 32px; }
  .section-title { display: flex; align-items: center; gap: 8px; font-size: 16px; font-weight: 700; margin-bottom: 14px; }
  .section-title svg { color: var(--secondary); }

  .notes-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(260px, 1fr)); gap: 16px; }
  .note-card {
    background: var(--surface); border-radius: var(--radius); padding: 18px; box-shadow: var(--shadow);
    transition: transform 0.15s ease, box-shadow 0.15s ease; cursor: pointer; position: relative; display: flex; flex-direction: column;
  }
  .note-card:hover { transform: translateY(-3px); box-shadow: var(--shadow-hover); }
  .note-card-top { display: flex; align-items: flex-start; justify-content: space-between; gap: 8px; margin-bottom: 8px; }
  .note-title { font-size: 15px; font-weight: 700; line-height: 1.3; }
  .note-preview {
    font-size: 13px; color: var(--text-muted); margin: 6px 0 12px 0; line-height: 1.5;
    display: -webkit-box; -webkit-line-clamp: 3; -webkit-box-orient: vertical; overflow: hidden;
  }
  .note-card-actions { display: flex; gap: 4px; flex-shrink: 0; }
  .note-icon-btn { width: 26px; height: 26px; border-radius: 50%; border: none; background: transparent; color: var(--text-muted); display: flex; align-items: center; justify-content: center; flex-shrink: 0; }
  .note-icon-btn:hover { background: var(--surface-hover); }
  .note-icon-btn.active { color: var(--secondary); }
  .note-meta { margin-top: auto; display: flex; align-items: center; justify-content: space-between; gap: 8px; padding-top: 10px; border-top: 1px solid var(--border); }
  .note-badge { font-size: 11px; font-weight: 600; background: rgba(255,154,0,0.15); color: var(--secondary); padding: 3px 10px; border-radius: 999px; }
  .note-date { font-size: 12px; color: var(--text-muted); }
  .note-tags { display: flex; gap: 5px; flex-wrap: wrap; margin-bottom: 8px; }
  .note-tag { font-size: 11px; background: var(--bg); color: var(--text-muted); padding: 2px 8px; border-radius: 999px; }

  .more-menu-wrap { position: relative; }
  .more-menu {
    position: absolute; top: 30px; right: 0; background: var(--surface); border: 1px solid var(--border);
    border-radius: var(--radius-sm); box-shadow: var(--shadow-hover); min-width: 150px; z-index: 30; display: none; overflow: hidden;
  }
  .more-menu.open { display: block; }
  .more-menu button {
    display: block; width: 100%; text-align: left; padding: 9px 14px; font-size: 13px; background: none; border: none; color: var(--text);
  }
  .more-menu button:hover { background: var(--surface-hover); }
  .more-menu button.danger { color: var(--danger); }

  .empty-state {
    text-align: center; padding: 60px 20px; color: var(--text-muted); background: var(--surface);
    border-radius: var(--radius); border: 1px dashed var(--border);
  }
  .empty-state svg { margin-bottom: 12px; opacity: 0.5; }

  /* Editor modal */
  .modal-overlay {
    position: fixed; inset: 0; background: rgba(79,32,13,0.35); display: none;
    align-items: center; justify-content: center; z-index: 100; padding: 24px;
  }
  .modal-overlay.open { display: flex; }
  .editor-modal {
    background: var(--surface); border-radius: var(--radius); width: 100%; max-width: 720px; max-height: 88vh;
    display: flex; flex-direction: column; overflow: hidden; box-shadow: var(--shadow-hover);
  }
  .editor-header { display: flex; align-items: center; gap: 10px; padding: 16px 20px; border-bottom: 1px solid var(--border); }
  .editor-title-input { flex: 1; border: none; font-size: 18px; font-weight: 700; font-family: 'Poppins', sans-serif; background: none; color: var(--text); }
  .editor-title-input:focus { outline: none; }
  .editor-meta-row { display: flex; gap: 10px; padding: 12px 20px; border-bottom: 1px solid var(--border); flex-wrap: wrap; }
  .editor-meta-row input, .editor-meta-row select {
    border: 1px solid var(--border); border-radius: var(--radius-sm); padding: 7px 12px; font-size: 13px; background: var(--bg); color: var(--text);
  }
  .editor-tabs { display: flex; gap: 4px; padding: 10px 20px 0 20px; }
  .editor-tab { padding: 8px 14px; font-size: 13px; font-weight: 600; color: var(--text-muted); border-bottom: 2px solid transparent; }
  .editor-tab.active { color: var(--text); border-bottom-color: var(--secondary); }
  .editor-body { flex: 1; overflow-y: auto; padding: 16px 20px; }
  #editor-textarea {
    width: 100%; height: 100%; min-height: 260px; border: none; resize: none; font-size: 14px; line-height: 1.6;
    font-family: 'Inter', monospace; color: var(--text); background: none;
  }
  #editor-textarea:focus { outline: none; }
  .editor-preview { font-size: 14px; line-height: 1.6; min-height: 260px; }
  .editor-preview h1, .editor-preview h2, .editor-preview h3 { margin-top: 0.6em; margin-bottom: 0.3em; }
  .editor-preview p { margin: 0.5em 0; }
  .editor-preview code { background: var(--bg); padding: 2px 6px; border-radius: 4px; font-size: 0.9em; }
  .editor-preview pre { background: var(--bg); padding: 12px; border-radius: 8px; overflow-x: auto; }
  .editor-preview pre code { background: none; padding: 0; }
  .editor-preview blockquote { border-left: 3px solid var(--secondary); margin: 0.6em 0; padding-left: 14px; color: var(--text-muted); }
  .editor-preview table { border-collapse: collapse; width: 100%; margin: 0.6em 0; }
  .editor-preview th, .editor-preview td { border: 1px solid var(--border); padding: 6px 10px; text-align: left; font-size: 13px; }
  .editor-preview th { background: var(--bg); }
  .editor-preview ul { padding-left: 22px; }
  .editor-preview ul.checklist { list-style: none; padding-left: 0; }
  .editor-footer {
    display: flex; align-items: center; gap: 14px; padding: 14px 20px; border-top: 1px solid var(--border);
    font-size: 12px; color: var(--text-muted);
  }

  .placeholder-view { text-align: center; padding: 80px 20px; color: var(--text-muted); }

  /* Auth screen */
  .auth-overlay {
    min-height: 100vh; display: flex; align-items: center; justify-content: center; padding: 20px;
  }
  .auth-card {
    background: var(--surface); border-radius: var(--radius); padding: 36px; width: 100%; max-width: 400px; box-shadow: var(--shadow);
  }
  .auth-logo { display: block; height: 34px; margin: 0 auto 28px auto; }
  .auth-tabs { display: flex; gap: 8px; margin-bottom: 22px; background: var(--bg); border-radius: var(--radius-sm); padding: 4px; }
  .auth-tab { flex: 1; text-align: center; padding: 9px; border-radius: 8px; font-size: 14px; font-weight: 600; color: var(--text-muted); }
  .auth-tab.active { background: var(--surface); color: var(--text); box-shadow: 0 1px 3px rgba(79,32,13,0.1); }
  .auth-card input {
    width: 100%; border: 1px solid var(--border); border-radius: var(--radius-sm); padding: 11px 14px;
    font-size: 14px; margin-bottom: 12px; background: var(--bg); color: var(--text);
  }
  .auth-card input:focus { outline: none; border-color: var(--secondary); }
  .auth-error { color: var(--danger); font-size: 13px; margin-bottom: 10px; min-height: 16px; }

  @media (max-width: 860px) {
    .sidebar { position: fixed; z-index: 200; height: 100vh; left: -260px; }
    .sidebar.mobile-open { left: 0; }
    .sidebar.collapsed { width: 250px; left: -260px; }
    .sidebar.collapsed.mobile-open { left: 0; }
    .mobile-only { display: flex !important; }
  }
  .mobile-only { display: none; }
</style>
</head>
<body>

<div id="auth-overlay" class="auth-overlay">
  <div class="auth-card">
    <img class="auth-logo" src="data:image/png;base64,LOGO_SRC_TOKEN" alt="Notez">
    <div class="auth-tabs">
      <div class="auth-tab active" data-tab="login">Log in</div>
      <div class="auth-tab" data-tab="signup">Sign up</div>
    </div>
    <div class="auth-error" id="auth-error"></div>
    <form id="auth-form">
      <input id="auth-username" placeholder="Username" autocomplete="username">
      <input id="auth-password" type="password" placeholder="Password" autocomplete="current-password">
      <button class="btn btn-primary" style="width:100%; justify-content:center;" id="auth-submit" type="submit">Log in</button>
    </form>
  </div>
</div>

<div id="app-shell" class="app-shell">
  <aside id="sidebar" class="sidebar">
    <div class="sidebar-header">
      <img src="data:image/png;base64,LOGO_SRC_TOKEN" alt="Notez">
      <button class="sidebar-toggle-btn" id="sidebar-toggle" title="Collapse sidebar">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M15 18l-6-6 6-6"/></svg>
      </button>
    </div>
    <nav class="sidebar-nav" id="sidebar-nav">
      <a href="#/dashboard" class="nav-item" data-route="dashboard">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="3" width="7" height="9" rx="1.5"/><rect x="14" y="3" width="7" height="5" rx="1.5"/><rect x="14" y="12" width="7" height="9" rx="1.5"/><rect x="3" y="16" width="7" height="5" rx="1.5"/></svg>
        <span class="nav-label">Dashboard</span>
      </a>
      <a href="#/notes" class="nav-item" data-route="notes">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M14 3H6a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V9z"/><path d="M14 3v6h6"/></svg>
        <span class="nav-label">All Notes</span>
      </a>
      <a href="#/favorites" class="nav-item" data-route="favorites">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01z"/></svg>
        <span class="nav-label">Favorites</span>
      </a>
      <div class="nav-section">
        <div class="nav-section-title">Categories</div>
        <div id="categories-list"></div>
      </div>
      <div class="nav-section">
        <div class="nav-section-title">Tags</div>
        <div id="tags-list"></div>
      </div>
      <div class="nav-section">
        <a href="#/archive" class="nav-item" data-route="archive">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="4" width="18" height="4" rx="1"/><path d="M5 8v11a1 1 0 0 0 1 1h12a1 1 0 0 0 1-1V8"/><path d="M10 13h4"/></svg>
          <span class="nav-label">Archive</span>
        </a>
        <a href="#/trash" class="nav-item" data-route="trash">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M3 6h18"/><path d="M8 6V4a1 1 0 0 1 1-1h6a1 1 0 0 1 1 1v2"/><path d="M19 6l-1 14a1 1 0 0 1-1 1H7a1 1 0 0 1-1-1L5 6"/></svg>
          <span class="nav-label">Trash</span>
        </a>
        <a href="#/shared" class="nav-item" data-route="shared">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="18" cy="5" r="3"/><circle cx="6" cy="12" r="3"/><circle cx="18" cy="19" r="3"/><path d="M8.6 13.5l6.8 4M15.4 6.5l-6.8 4"/></svg>
          <span class="nav-label">Shared Notes</span>
        </a>
        <a href="#/settings" class="nav-item" data-route="settings">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 1 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-4 0v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 1 1-2.83-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1 0-4h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 1 1 2.83-2.83l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 4 0v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 1 1 2.83 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 0 4h-.09a1.65 1.65 0 0 0-1.51 1z"/></svg>
          <span class="nav-label">Settings</span>
        </a>
      </div>
    </nav>
    <div class="sidebar-user" id="sidebar-user">
      <div class="user-avatar" id="sidebar-avatar">A</div>
      <div class="user-name" id="sidebar-username">username</div>
    </div>
  </aside>

  <div class="main-area">
    <header class="topnav">
      <button class="icon-btn mobile-only" id="mobile-menu-btn">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="3" y1="6" x2="21" y2="6"/><line x1="3" y1="12" x2="21" y2="12"/><line x1="3" y1="18" x2="21" y2="18"/></svg>
      </button>
      <div class="search-wrap">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg>
        <input id="global-search" placeholder="Search notes...">
      </div>
      <div class="spacer"></div>
      <button class="btn btn-primary" id="new-note-btn">
        <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
        New Note
      </button>
      <div style="position:relative;">
        <div class="topnav-user" id="topnav-avatar">A</div>
        <div class="user-dropdown" id="user-dropdown">
          <a href="#/settings">Settings</a>
          <button id="logout-btn-dropdown">Log out</button>
        </div>
      </div>
    </header>
    <main id="content" class="content"></main>
  </div>
</div>

<div id="editor-overlay" class="modal-overlay">
  <div class="editor-modal">
    <div class="editor-header">
      <input id="editor-title" class="editor-title-input" placeholder="Note title">
      <button class="icon-btn" id="editor-close">
        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
      </button>
    </div>
    <div class="editor-meta-row">
      <input id="editor-category" placeholder="Category" list="category-options" style="max-width:160px;">
      <datalist id="category-options"></datalist>
      <input id="editor-tags" placeholder="Tags, comma separated" style="flex:1;">
    </div>
    <div class="editor-tabs">
      <div class="editor-tab active" data-editor-tab="write">Write</div>
      <div class="editor-tab" data-editor-tab="preview">Preview</div>
    </div>
    <div class="editor-body">
      <textarea id="editor-textarea" placeholder="Start writing in Markdown... **bold**, # heading, - [ ] checklist, &#96;&#96;&#96;code blocks&#96;&#96;&#96;"></textarea>
      <div class="editor-preview" id="editor-preview" style="display:none;"></div>
    </div>
    <div class="editor-footer">
      <span id="save-indicator">All changes saved</span>
      <span id="word-count">0 words</span>
      <span id="reading-time">0 min read</span>
      <div class="spacer"></div>
      <button class="btn btn-ghost btn-sm" id="editor-trash-btn">Move to trash</button>
      <button class="btn btn-primary btn-sm" id="editor-save-btn">Save</button>
    </div>
  </div>
</div>

<script>
NOTEZ_APP_JS_TOKEN
</script>
</body>
</html>
`

const appJS = `
function escapeHtml(s) {
  return s.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;');
}

function renderInline(text) {
  text = escapeHtml(text);
  text = text.replace(/\u0060([^\u0060]+)\u0060/g, '<code>$1</code>');
  text = text.replace(/\*\*([^*]+)\*\*/g, '<strong>$1</strong>');
  text = text.replace(/\*([^*]+)\*/g, '<em>$1</em>');
  return text;
}

function markdownToHtml(src) {
  const lines = src.replace(/\r\n/g, '\n').split('\n');
  let html = '';
  let i = 0;
  let inList = false;
  let inCodeBlock = false;
  let codeBuffer = [];

  function closeList() {
    if (inList) { html += '</ul>'; inList = false; }
  }

  while (i < lines.length) {
    const line = lines[i];

    if (line.trim().startsWith('\u0060\u0060\u0060')) {
      if (!inCodeBlock) {
        closeList();
        inCodeBlock = true;
        codeBuffer = [];
      } else {
        html += '<pre><code>' + escapeHtml(codeBuffer.join('\n')) + '</code></pre>';
        inCodeBlock = false;
      }
      i++;
      continue;
    }
    if (inCodeBlock) {
      codeBuffer.push(line);
      i++;
      continue;
    }

    if (/^\|.*\|\s*$/.test(line) && i + 1 < lines.length && /^\|[\s:|-]+\|\s*$/.test(lines[i + 1])) {
      closeList();
      const headerCells = line.trim().slice(1, -1).split('|').map(c => c.trim());
      html += '<table><thead><tr>' + headerCells.map(c => '<th>' + renderInline(c) + '</th>').join('') + '</tr></thead><tbody>';
      i += 2;
      while (i < lines.length && /^\|.*\|\s*$/.test(lines[i])) {
        const cells = lines[i].trim().slice(1, -1).split('|').map(c => c.trim());
        html += '<tr>' + cells.map(c => '<td>' + renderInline(c) + '</td>').join('') + '</tr>';
        i++;
      }
      html += '</tbody></table>';
      continue;
    }

    const headingMatch = line.match(/^(#{1,3})\s+(.*)$/);
    if (headingMatch) {
      closeList();
      const level = headingMatch[1].length;
      html += '<h' + level + '>' + renderInline(headingMatch[2]) + '</h' + level + '>';
      i++;
      continue;
    }

    const quoteMatch = line.match(/^>\s?(.*)$/);
    if (quoteMatch) {
      closeList();
      html += '<blockquote>' + renderInline(quoteMatch[1]) + '</blockquote>';
      i++;
      continue;
    }

    const checklistMatch = line.match(/^[-*]\s+\[( |x|X)\]\s+(.*)$/);
    if (checklistMatch) {
      if (!inList) { html += '<ul class="checklist">'; inList = true; }
      const checked = checklistMatch[1].toLowerCase() === 'x';
      html += '<li><input type="checkbox" disabled' + (checked ? ' checked' : '') + '> ' + renderInline(checklistMatch[2]) + '</li>';
      i++;
      continue;
    }

    const listMatch = line.match(/^[-*]\s+(.*)$/);
    if (listMatch) {
      if (!inList) { html += '<ul>'; inList = true; }
      html += '<li>' + renderInline(listMatch[1]) + '</li>';
      i++;
      continue;
    }

    closeList();

    if (line.trim() === '') {
      i++;
      continue;
    }

    html += '<p>' + renderInline(line) + '</p>';
    i++;
  }
  closeList();
  if (inCodeBlock) {
    html += '<pre><code>' + escapeHtml(codeBuffer.join('\n')) + '</code></pre>';
  }
  return html;
}

let accessToken = localStorage.getItem('access_token');
let refreshToken = localStorage.getItem('refresh_token');
let currentUsername = localStorage.getItem('username') || '';
let sidebarCollapsed = false;
let mobileMenuOpen = false;
let currentEditingNoteId = null;
let saveDebounceTimer = null;
let allCategoriesCache = [];
let allTagsCache = [];

function saveTokens(pair, username) {
  accessToken = pair.access_token;
  refreshToken = pair.refresh_token;
  localStorage.setItem('access_token', accessToken);
  localStorage.setItem('refresh_token', refreshToken);
  if (username) {
    currentUsername = username;
    localStorage.setItem('username', username);
  }
}

function clearTokens() {
  accessToken = null;
  refreshToken = null;
  localStorage.removeItem('access_token');
  localStorage.removeItem('refresh_token');
  localStorage.removeItem('username');
}

async function api(path, opts) {
  opts = opts || {};
  opts.headers = opts.headers || {};
  if (accessToken) opts.headers['Authorization'] = 'Bearer ' + accessToken;
  if (opts.body && typeof opts.body !== 'string') {
    opts.body = JSON.stringify(opts.body);
    opts.headers['Content-Type'] = 'application/json';
  }

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
  document.getElementById('auth-overlay').style.display = 'flex';
  document.getElementById('app-shell').classList.remove('visible');
}

function showAppView() {
  document.getElementById('auth-overlay').style.display = 'none';
  document.getElementById('app-shell').classList.add('visible');
  const initial = (currentUsername || '?').charAt(0).toUpperCase();
  document.getElementById('sidebar-avatar').textContent = initial;
  document.getElementById('topnav-avatar').textContent = initial;
  document.getElementById('sidebar-username').textContent = currentUsername;
  refreshSidebarLists();
  if (!location.hash || location.hash === '#/') {
    location.hash = '#/dashboard';
  } else {
    router();
  }
}

async function refreshSidebarLists() {
  try {
    const res = await api('/tasks?limit=100&view=all');
    if (!res.ok) return;
    const data = await res.json();
    const notes = data.tasks || [];

    const catSet = {};
    const tagSet = {};
    notes.forEach(function(n) {
      if (n.category) catSet[n.category] = (catSet[n.category] || 0) + 1;
      (n.tags || []).forEach(function(t) {
        tagSet[t] = (tagSet[t] || 0) + 1;
      });
    });

    allCategoriesCache = Object.keys(catSet).sort();
    allTagsCache = Object.keys(tagSet).sort();

    const catList = document.getElementById('categories-list');
    catList.innerHTML = allCategoriesCache.map(function(c) {
      return '<a href="#/category/' + encodeURIComponent(c) + '" class="nav-item" data-route="category" data-value="' + escapeHtml(c) + '">' +
        '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M3 7a2 2 0 0 1 2-2h4l2 2h8a2 2 0 0 1 2 2v8a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/></svg>' +
        '<span class="nav-label">' + escapeHtml(c) + '</span>' +
        '<span class="nav-badge">' + catSet[c] + '</span>' +
        '</a>';
    }).join('') || '<div style="padding:6px 12px; font-size:12px; color: var(--sidebar-text-muted);">No categories yet</div>';

    const tagList = document.getElementById('tags-list');
    tagList.innerHTML = allTagsCache.map(function(t) {
      return '<a href="#/tag/' + encodeURIComponent(t) + '" class="nav-item" data-route="tag" data-value="' + escapeHtml(t) + '">' +
        '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M20.59 13.41L11 3.83A2 2 0 0 0 9.59 3.17H4a1 1 0 0 0-1 1v5.59a2 2 0 0 0 .66 1.41l9.58 9.58a2 2 0 0 0 2.83 0l4.52-4.52a2 2 0 0 0 0-2.83z"/><circle cx="7.5" cy="7.5" r="1"/></svg>' +
        '<span class="nav-label">' + escapeHtml(t) + '</span>' +
        '<span class="nav-badge">' + tagSet[t] + '</span>' +
        '</a>';
    }).join('') || '<div style="padding:6px 12px; font-size:12px; color: var(--sidebar-text-muted);">No tags yet</div>';

    bindNavLinks();
    updateActiveNav();
  } catch (e) {}
}

function bindNavLinks() {
  document.querySelectorAll('.nav-item').forEach(function(el) {
    el.onclick = function(e) {
      if (window.innerWidth <= 860) closeMobileSidebar();
    };
  });
}

function updateActiveNav() {
  const hash = location.hash || '#/dashboard';
  document.querySelectorAll('.nav-item').forEach(function(el) {
    el.classList.remove('active');
  });
  document.querySelectorAll('.nav-item').forEach(function(el) {
    const href = el.getAttribute('href');
    if (href === hash) el.classList.add('active');
  });
}

function stripMarkdown(text) {
  return (text || '')
    .replace(/\u0060\u0060\u0060[\s\S]*?\u0060\u0060\u0060/g, '')
    .replace(/[#>*\u0060|]/g, '')
    .replace(/^\s*-\s*\[[ xX]\]\s*/gm, '')
    .replace(/^\s*-\s*/gm, '')
    .replace(/\n+/g, ' ')
    .trim();
}

function formatDate(iso) {
  const d = new Date(iso);
  return d.toLocaleDateString(undefined, { month: 'short', day: 'numeric', year: 'numeric' });
}

function noteCardHtml(note) {
  const preview = stripMarkdown(note.notes).slice(0, 140);
  const category = note.category ? '<span class="note-badge">' + escapeHtml(note.category) + '</span>' : '<span></span>';
  const tags = (note.tags || []).map(function(t) {
    return '<span class="note-tag">' + escapeHtml(t) + '</span>';
  }).join('');

  return '<div class="note-card" data-note-id="' + note.id + '">' +
    '<div class="note-card-top">' +
    '<div class="note-title">' + escapeHtml(note.title) + '</div>' +
    '<div class="note-card-actions">' +
    '<button class="note-icon-btn fav-btn ' + (note.favorite ? 'active' : '') + '" data-id="' + note.id + '" title="Favorite">' +
    '<svg width="15" height="15" viewBox="0 0 24 24" fill="' + (note.favorite ? 'currentColor' : 'none') + '" stroke="currentColor" stroke-width="2"><path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01z"/></svg>' +
    '</button>' +
    '<button class="note-icon-btn pin-btn ' + (note.pinned ? 'active' : '') + '" data-id="' + note.id + '" title="Pin">' +
    '<svg width="15" height="15" viewBox="0 0 24 24" fill="' + (note.pinned ? 'currentColor' : 'none') + '" stroke="currentColor" stroke-width="2"><path d="M12 17v5"/><path d="M9 3h6l-1 7 4 3H6l4-3z"/></svg>' +
    '</button>' +
    '<div class="more-menu-wrap">' +
    '<button class="note-icon-btn more-btn" data-id="' + note.id + '" title="More">' +
    '<svg width="15" height="15" viewBox="0 0 24 24" fill="currentColor"><circle cx="12" cy="5" r="1.5"/><circle cx="12" cy="12" r="1.5"/><circle cx="12" cy="19" r="1.5"/></svg>' +
    '</button>' +
    '<div class="more-menu" id="more-menu-' + note.id + '">' +
    moreMenuItemsHtml(note) +
    '</div>' +
    '</div>' +
    '</div>' +
    '</div>' +
    (tags ? '<div class="note-tags">' + tags + '</div>' : '') +
    '<div class="note-preview">' + escapeHtml(preview) + (preview.length >= 140 ? '...' : '') + '</div>' +
    '<div class="note-meta">' + category + '<span class="note-date">' + formatDate(note.updated_at) + '</span></div>' +
    '</div>';
}

function moreMenuItemsHtml(note) {
  if (note.trashed) {
    return '<button data-action="restore" data-id="' + note.id + '">Restore</button>' +
      '<button data-action="delete-forever" data-id="' + note.id + '" class="danger">Delete forever</button>';
  }
  if (note.archived) {
    return '<button data-action="edit" data-id="' + note.id + '">Edit</button>' +
      '<button data-action="unarchive" data-id="' + note.id + '">Unarchive</button>' +
      '<button data-action="trash" data-id="' + note.id + '" class="danger">Move to trash</button>';
  }
  return '<button data-action="edit" data-id="' + note.id + '">Edit</button>' +
    '<button data-action="archive" data-id="' + note.id + '">Archive</button>' +
    '<button data-action="trash" data-id="' + note.id + '" class="danger">Move to trash</button>';
}

let currentRoute = { name: 'dashboard', value: null };
let currentSearchTerm = '';

function parseHash() {
  const hash = (location.hash || '#/dashboard').replace(/^#\//, '');
  const parts = hash.split('/');
  if (parts[0] === 'category' || parts[0] === 'tag') {
    return { name: parts[0], value: decodeURIComponent(parts[1] || '') };
  }
  return { name: parts[0] || 'dashboard', value: null };
}

async function router() {
  currentRoute = parseHash();
  updateActiveNav();
  const content = document.getElementById('content');
  content.innerHTML = '<div style="padding:40px; text-align:center; color: var(--text-muted);">Loading...</div>';

  switch (currentRoute.name) {
    case 'dashboard': return renderDashboard();
    case 'notes': return renderNotesView('All Notes', { view: 'all' });
    case 'favorites': return renderNotesView('Favorites', { view: 'all', favorite: 'true' });
    case 'category': return renderNotesView('Category: ' + currentRoute.value, { view: 'all', category: currentRoute.value });
    case 'tag': return renderNotesView('Tag: ' + currentRoute.value, { view: 'all', tag: currentRoute.value });
    case 'archive': return renderNotesView('Archive', { view: 'archived' }, true);
    case 'trash': return renderNotesView('Trash', { view: 'trash' }, true);
    case 'shared': return renderPlaceholder('Shared Notes', 'Sharing notes with others is on the roadmap and not available yet.');
    case 'settings': return renderSettings();
    default: return renderDashboard();
  }
}

function buildQuery(params) {
  const usp = new URLSearchParams(params);
  if (currentSearchTerm) usp.set('q', currentSearchTerm);
  usp.set('limit', '100');
  return usp.toString();
}

async function renderDashboard() {
  const content = document.getElementById('content');

  let stats = { total: 0, favorites: 0, pinned: 0, archived: 0, trashed: 0, categories: [] };
  try {
    const statsRes = await api('/tasks/stats');
    if (statsRes.ok) stats = await statsRes.json();
  } catch (e) {}

  let allNotes = [];
  try {
    const res = await api('/tasks?view=all&limit=100');
    if (res.ok) allNotes = (await res.json()).tasks || [];
  } catch (e) {}

  const pinned = allNotes.filter(function(n) { return n.pinned; }).slice(0, 6);
  const favorites = allNotes.filter(function(n) { return n.favorite; }).slice(0, 6);
  const recent = allNotes.slice().sort(function(a, b) { return new Date(b.created_at) - new Date(a.created_at); }).slice(0, 6);
  const recentlyEdited = allNotes.slice().sort(function(a, b) { return new Date(b.updated_at) - new Date(a.updated_at); }).slice(0, 5);

  let html = '<div class="content-header"><div>' +
    '<h2 class="content-title">Welcome back, ' + escapeHtml(currentUsername) + '</h2>' +
    '<div class="content-sub">Here\'s what\'s happening with your notes</div>' +
    '</div></div>';

  html += '<div class="stats-row">' +
    statCard(stats.total, 'Total Notes') +
    statCard(stats.favorites, 'Favorites') +
    statCard(stats.pinned, 'Pinned') +
    statCard(stats.categories.length, 'Categories') +
    '</div>';

  html += notesSection('Pinned Notes', pinned, sidebarIcon('pin'));
  html += notesSection('Favorite Notes', favorites, sidebarIcon('star'));
  html += notesSection('Recent Notes', recent, sidebarIcon('clock'));

  html += '<div class="section"><div class="section-title">' + sidebarIcon('activity') + ' Activity Timeline</div>';
  if (recentlyEdited.length === 0) {
    html += '<div class="empty-state">No activity yet</div>';
  } else {
    html += '<div style="background:var(--surface); border-radius:var(--radius); box-shadow:var(--shadow); padding: 6px 20px;">';
    html += recentlyEdited.map(function(n) {
      const isNew = n.created_at === n.updated_at;
      return '<div style="display:flex; justify-content:space-between; padding: 12px 0; border-bottom: 1px solid var(--border); font-size:13px;">' +
        '<span><strong>' + escapeHtml(n.title) + '</strong> was ' + (isNew ? 'created' : 'updated') + '</span>' +
        '<span style="color:var(--text-muted);">' + formatDate(n.updated_at) + '</span>' +
        '</div>';
    }).join('');
    html += '</div>';
  }
  html += '</div>';

  content.innerHTML = html;
  bindNoteCardEvents();
}

function statCard(value, label) {
  return '<div class="stat-card"><div class="stat-value">' + value + '</div><div class="stat-label">' + label + '</div></div>';
}

function notesSection(title, notes, iconSvg) {
  if (notes.length === 0) return '';
  return '<div class="section"><div class="section-title">' + iconSvg + ' ' + title + '</div>' +
    '<div class="notes-grid">' + notes.map(noteCardHtml).join('') + '</div></div>';
}

function sidebarIcon(name) {
  const icons = {
    pin: '<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 17v5"/><path d="M9 3h6l-1 7 4 3H6l4-3z"/></svg>',
    star: '<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01z"/></svg>',
    clock: '<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="9"/><path d="M12 7v5l3 3"/></svg>',
    activity: '<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M22 12h-4l-3 9L9 3l-3 9H2"/></svg>'
  };
  return icons[name] || '';
}

async function renderNotesView(title, queryParams, isTrashOrArchive) {
  const content = document.getElementById('content');
  let notes = [];
  try {
    const res = await api('/tasks?' + buildQuery(queryParams));
    if (res.ok) notes = (await res.json()).tasks || [];
  } catch (e) {}

  let html = '<div class="content-header"><div>' +
    '<h2 class="content-title">' + escapeHtml(title) + '</h2>' +
    '<div class="content-sub">' + notes.length + ' note' + (notes.length === 1 ? '' : 's') + '</div>' +
    '</div>';

  if (currentRoute.name === 'trash' && notes.length > 0) {
    html += '<button class="btn btn-danger btn-sm" id="empty-trash-btn">Empty Trash</button>';
  }
  html += '</div>';

  if (notes.length === 0) {
    html += '<div class="empty-state">' +
      '<svg width="40" height="40" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M14 3H6a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V9z"/><path d="M14 3v6h6"/></svg>' +
      '<div>No notes here yet</div></div>';
  } else {
    html += '<div class="notes-grid">' + notes.map(noteCardHtml).join('') + '</div>';
  }

  content.innerHTML = html;
  bindNoteCardEvents();

  const emptyTrashBtn = document.getElementById('empty-trash-btn');
  if (emptyTrashBtn) {
    emptyTrashBtn.onclick = async function() {
      if (!confirm('Permanently delete all ' + notes.length + ' note(s) in trash? This cannot be undone.')) return;
      for (const n of notes) {
        await api('/tasks/' + n.id, { method: 'DELETE' });
      }
      router();
    };
  }
}

function renderPlaceholder(title, message) {
  const content = document.getElementById('content');
  content.innerHTML = '<div class="content-header"><h2 class="content-title">' + escapeHtml(title) + '</h2></div>' +
    '<div class="placeholder-view">' +
    '<svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" style="margin-bottom:16px; opacity:0.4;"><circle cx="12" cy="12" r="9"/><path d="M12 8v4M12 16h.01"/></svg>' +
    '<div>' + escapeHtml(message) + '</div></div>';
}

function renderSettings() {
  const content = document.getElementById('content');
  content.innerHTML = '<div class="content-header"><h2 class="content-title">Settings</h2></div>' +
    '<div style="background:var(--surface); border-radius:var(--radius); box-shadow:var(--shadow); padding:24px; max-width:480px;">' +
    '<div style="display:flex; align-items:center; gap:14px; margin-bottom:20px;">' +
    '<div class="user-avatar" style="width:48px;height:48px;font-size:18px;">' + escapeHtml((currentUsername || '?').charAt(0).toUpperCase()) + '</div>' +
    '<div><div style="font-weight:700; font-size:16px;">' + escapeHtml(currentUsername) + '</div>' +
    '<div style="color:var(--text-muted); font-size:13px;">Signed in</div></div>' +
    '</div>' +
    '<button class="btn btn-ghost" id="settings-logout-btn" style="width:100%; justify-content:center;">Log out</button>' +
    '</div>';
  document.getElementById('settings-logout-btn').onclick = doLogout;
}

function bindNoteCardEvents() {
  document.querySelectorAll('.note-card').forEach(function(card) {
    card.addEventListener('click', function(e) {
      if (e.target.closest('.note-icon-btn') || e.target.closest('.more-menu')) return;
      openEditor(card.dataset.noteId);
    });
  });

  document.querySelectorAll('.fav-btn').forEach(function(btn) {
    btn.addEventListener('click', async function(e) {
      e.stopPropagation();
      const id = btn.dataset.id;
      const isActive = btn.classList.contains('active');
      await api('/tasks/' + id, { method: 'PUT', body: { favorite: !isActive } });
      router();
    });
  });

  document.querySelectorAll('.pin-btn').forEach(function(btn) {
    btn.addEventListener('click', async function(e) {
      e.stopPropagation();
      const id = btn.dataset.id;
      const isActive = btn.classList.contains('active');
      await api('/tasks/' + id, { method: 'PUT', body: { pinned: !isActive } });
      router();
    });
  });

  document.querySelectorAll('.more-btn').forEach(function(btn) {
    btn.addEventListener('click', function(e) {
      e.stopPropagation();
      const id = btn.dataset.id;
      document.querySelectorAll('.more-menu').forEach(function(m) {
        if (m.id !== 'more-menu-' + id) m.classList.remove('open');
      });
      document.getElementById('more-menu-' + id).classList.toggle('open');
    });
  });

  document.querySelectorAll('.more-menu button').forEach(function(btn) {
    btn.addEventListener('click', async function(e) {
      e.stopPropagation();
      const id = btn.dataset.id;
      const action = btn.dataset.action;
      if (action === 'edit') { openEditor(id); return; }
      if (action === 'archive') await api('/tasks/' + id, { method: 'PUT', body: { archived: true } });
      if (action === 'unarchive') await api('/tasks/' + id, { method: 'PUT', body: { archived: false } });
      if (action === 'trash') await api('/tasks/' + id + '/trash', { method: 'POST' });
      if (action === 'restore') await api('/tasks/' + id + '/restore', { method: 'POST' });
      if (action === 'delete-forever') {
        if (!confirm('Permanently delete this note? This cannot be undone.')) return;
        await api('/tasks/' + id, { method: 'DELETE' });
      }
      router();
      refreshSidebarLists();
    });
  });
}

document.addEventListener('click', function() {
  document.querySelectorAll('.more-menu').forEach(function(m) { m.classList.remove('open'); });
});

function wordCountAndReadingTime(text) {
  const words = (text || '').trim().split(/\s+/).filter(Boolean);
  const count = words.length;
  const minutes = Math.max(1, Math.ceil(count / 200));
  return { count: count, minutes: minutes };
}

function updateEditorStats() {
  const text = document.getElementById('editor-textarea').value;
  const stats = wordCountAndReadingTime(text);
  document.getElementById('word-count').textContent = stats.count + ' words';
  document.getElementById('reading-time').textContent = stats.minutes + ' min read';
}

async function openEditor(noteId) {
  currentEditingNoteId = noteId || null;
  const overlay = document.getElementById('editor-overlay');
  const titleInput = document.getElementById('editor-title');
  const categoryInput = document.getElementById('editor-category');
  const tagsInput = document.getElementById('editor-tags');
  const textarea = document.getElementById('editor-textarea');
  const trashBtn = document.getElementById('editor-trash-btn');

  const catOptions = document.getElementById('category-options');
  catOptions.innerHTML = allCategoriesCache.map(function(c) {
    return '<option value="' + escapeHtml(c) + '">';
  }).join('');

  if (noteId) {
    const res = await api('/tasks/' + noteId);
    if (!res.ok) return;
    const note = await res.json();
    titleInput.value = note.title;
    categoryInput.value = note.category || '';
    tagsInput.value = (note.tags || []).join(', ');
    textarea.value = note.notes || '';
    trashBtn.style.display = 'inline-flex';
    document.getElementById('save-indicator').textContent = 'All changes saved';
  } else {
    titleInput.value = '';
    categoryInput.value = '';
    tagsInput.value = '';
    textarea.value = '';
    trashBtn.style.display = 'none';
    document.getElementById('save-indicator').textContent = 'Unsaved note';
  }

  setEditorTab('write');
  updateEditorStats();
  overlay.classList.add('open');
  titleInput.focus();
}

function closeEditor() {
  document.getElementById('editor-overlay').classList.remove('open');
  currentEditingNoteId = null;
  clearTimeout(saveDebounceTimer);
  router();
  refreshSidebarLists();
}

function setEditorTab(tab) {
  document.querySelectorAll('.editor-tab').forEach(function(t) {
    t.classList.toggle('active', t.dataset.editorTab === tab);
  });
  const textarea = document.getElementById('editor-textarea');
  const preview = document.getElementById('editor-preview');
  if (tab === 'preview') {
    preview.innerHTML = markdownToHtml(textarea.value) || '<p style="color:var(--text-muted);">Nothing to preview yet</p>';
    textarea.style.display = 'none';
    preview.style.display = 'block';
  } else {
    textarea.style.display = 'block';
    preview.style.display = 'none';
  }
}

function collectEditorData() {
  const tagsRaw = document.getElementById('editor-tags').value;
  const tags = tagsRaw.split(',').map(function(t) { return t.trim(); }).filter(Boolean);
  return {
    title: document.getElementById('editor-title').value.trim(),
    notes: document.getElementById('editor-textarea').value,
    category: document.getElementById('editor-category').value.trim(),
    tags: tags
  };
}

async function saveEditor(manual) {
  const data = collectEditorData();
  if (!data.title) {
    if (manual) document.getElementById('save-indicator').textContent = 'Title is required';
    return;
  }

  const indicator = document.getElementById('save-indicator');
  indicator.textContent = 'Saving...';

  if (currentEditingNoteId) {
    const res = await api('/tasks/' + currentEditingNoteId, { method: 'PUT', body: data });
    indicator.textContent = res.ok ? 'All changes saved' : 'Failed to save';
  } else {
    const res = await api('/tasks', { method: 'POST', body: data });
    if (res.ok) {
      const created = await res.json();
      currentEditingNoteId = created.id;
      indicator.textContent = 'All changes saved';
      document.getElementById('editor-trash-btn').style.display = 'inline-flex';
    } else {
      indicator.textContent = 'Failed to save';
    }
  }
  refreshSidebarLists();
}

function scheduleAutosave() {
  if (!currentEditingNoteId) return;
  clearTimeout(saveDebounceTimer);
  document.getElementById('save-indicator').textContent = 'Editing...';
  saveDebounceTimer = setTimeout(function() { saveEditor(false); }, 900);
}

function closeMobileSidebar() {
  document.getElementById('sidebar').classList.remove('mobile-open');
  mobileMenuOpen = false;
}

function initEventListeners() {
  document.querySelectorAll('.auth-tab').forEach(function(tab) {
    tab.addEventListener('click', function() {
      document.querySelectorAll('.auth-tab').forEach(function(t) { t.classList.remove('active'); });
      tab.classList.add('active');
      document.getElementById('auth-submit').textContent = tab.dataset.tab === 'login' ? 'Log in' : 'Sign up';
      document.getElementById('auth-error').textContent = '';
    });
  });

  document.getElementById('auth-form').addEventListener('submit', async function(e) {
    e.preventDefault();
    const mode = document.querySelector('.auth-tab.active').dataset.tab;
    const username = document.getElementById('auth-username').value.trim();
    const password = document.getElementById('auth-password').value;
    const errEl = document.getElementById('auth-error');
    errEl.textContent = '';

    if (!username || !password) {
      errEl.textContent = 'Username and password are required';
      return;
    }

    try {
      const res = await fetch('/auth/' + mode, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username: username, password: password })
      });
      const data = await res.json();
      if (!res.ok) {
        errEl.textContent = data.error || 'something went wrong';
        return;
      }
      saveTokens(data, username);
      showAppView();
    } catch (err) {
      errEl.textContent = 'network error, try again';
    }
  });

  document.getElementById('sidebar-toggle').addEventListener('click', function() {
    sidebarCollapsed = !sidebarCollapsed;
    document.getElementById('sidebar').classList.toggle('collapsed', sidebarCollapsed);
  });

  document.getElementById('mobile-menu-btn').addEventListener('click', function() {
    mobileMenuOpen = !mobileMenuOpen;
    document.getElementById('sidebar').classList.toggle('mobile-open', mobileMenuOpen);
  });

  document.getElementById('new-note-btn').addEventListener('click', function() { openEditor(null); });

  document.getElementById('topnav-avatar').addEventListener('click', function(e) {
    e.stopPropagation();
    document.getElementById('user-dropdown').classList.toggle('open');
  });
  document.addEventListener('click', function() {
    document.getElementById('user-dropdown').classList.remove('open');
  });

  document.getElementById('logout-btn-dropdown').addEventListener('click', doLogout);

  let searchDebounce;
  document.getElementById('global-search').addEventListener('input', function(e) {
    clearTimeout(searchDebounce);
    searchDebounce = setTimeout(function() {
      currentSearchTerm = e.target.value.trim();
      router();
    }, 300);
  });

  document.getElementById('editor-close').addEventListener('click', closeEditor);
  document.getElementById('editor-save-btn').addEventListener('click', function() { saveEditor(true); });
  document.getElementById('editor-trash-btn').addEventListener('click', async function() {
    if (currentEditingNoteId) {
      await api('/tasks/' + currentEditingNoteId + '/trash', { method: 'POST' });
    }
    closeEditor();
  });

  document.getElementById('editor-textarea').addEventListener('input', function() {
    updateEditorStats();
    scheduleAutosave();
  });
  document.getElementById('editor-title').addEventListener('input', scheduleAutosave);
  document.getElementById('editor-category').addEventListener('change', scheduleAutosave);
  document.getElementById('editor-tags').addEventListener('input', scheduleAutosave);

  document.querySelectorAll('.editor-tab').forEach(function(tab) {
    tab.addEventListener('click', function() { setEditorTab(tab.dataset.editorTab); });
  });

  document.getElementById('editor-overlay').addEventListener('click', function(e) {
    if (e.target === document.getElementById('editor-overlay')) closeEditor();
  });

  window.addEventListener('hashchange', router);
}

function doLogout() {
  clearTokens();
  showAuthView();
}

initEventListeners();

if (accessToken) {
  showAppView();
} else {
  showAuthView();
}
`

func Handler() http.HandlerFunc {
	html := pageHTML
	html = strings.Replace(html, "LOGO_SRC_TOKEN", notezLogoB64, -1)
	html = strings.Replace(html, "NOTEZ_APP_JS_TOKEN", appJS, 1)
	htmlBytes := []byte(html)
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(htmlBytes)
	}
}
