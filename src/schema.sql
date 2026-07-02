-- Todo API – SQLite schema reference
-- This file is for documentation only; migrations run automatically via src/db.js

PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS users (
  id         INTEGER PRIMARY KEY AUTOINCREMENT,
  username   TEXT    NOT NULL UNIQUE,
  email      TEXT    NOT NULL UNIQUE,
  created_at TEXT    NOT NULL DEFAULT (datetime('now')),
  updated_at TEXT    NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS todos (
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id     INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  title       TEXT    NOT NULL,
  description TEXT,
  completed   INTEGER NOT NULL DEFAULT 0,   -- 0 = false, 1 = true
  priority    TEXT    NOT NULL DEFAULT 'medium',  -- low | medium | high | urgent
  status      TEXT    NOT NULL DEFAULT 'pending', -- pending | in_progress | done | cancelled
  due_at      TEXT,                          -- ISO 8601 datetime string, optional
  created_at  TEXT    NOT NULL DEFAULT (datetime('now')),
  updated_at  TEXT    NOT NULL DEFAULT (datetime('now'))
);
