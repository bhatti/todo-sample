'use strict';

const initSqlJs = require('sql.js');
const fs = require('fs');
const path = require('path');

const DB_PATH = process.env.DATABASE_PATH || path.join(__dirname, '..', 'todo.db');

/**
 * sql.js wrapper that exposes a synchronous API similar to better-sqlite3.
 * This lets us keep the model layer clean and synchronous.
 */
class SqlJsDatabase {
  constructor(sqlJs, data) {
    this._db = new sqlJs.Database(data || null);
    this._path = null;
  }

  /**
   * Execute DDL/DML statements (no return value needed).
   * @param {string} sql
   */
  exec(sql) {
    this._db.run(sql);
    return this;
  }

  /**
   * Prepare and run a SELECT, returning the first matching row or undefined.
   * @param {string} sql
   * @returns {{ get: (...params) => object|undefined }}
   */
  prepare(sql) {
    const db = this._db;
    return {
      get: (...params) => {
        const stmt = db.prepare(sql);
        stmt.bind(params);
        if (stmt.step()) {
          const row = stmt.getAsObject();
          stmt.free();
          return row;
        }
        stmt.free();
        return undefined;
      },
      all: (...params) => {
        const stmt = db.prepare(sql);
        stmt.bind(params);
        const rows = [];
        while (stmt.step()) {
          rows.push(stmt.getAsObject());
        }
        stmt.free();
        return rows;
      },
      run: (...params) => {
        const stmt = db.prepare(sql);
        stmt.bind(params);
        stmt.step();
        stmt.free();
        const changes = db.getRowsModified();
        return { changes };
      },
    };
  }

  /**
   * Persist the database to disk if a path is set.
   */
  _save() {
    if (this._path && this._path !== ':memory:') {
      const data = this._db.export();
      fs.writeFileSync(this._path, Buffer.from(data));
    }
  }

  close() {
    this._db.close();
  }
}

// sql.js uses WASM and must be initialised asynchronously once.
let _sqlJsPromise = null;

function getSqlJs() {
  if (!_sqlJsPromise) {
    _sqlJsPromise = initSqlJs();
  }
  return _sqlJsPromise;
}

/**
 * Opens and initialises the SQLite database synchronously.
 * NOTE: Because sql.js init is async, callers should use openDb() which
 * returns a promise. The app.js / server.js wait for it.
 * @param {string} [dbPath] - ':memory:' or a file path.
 * @returns {Promise<SqlJsDatabase>}
 */
async function openDb(dbPath) {
  const resolvedPath = dbPath || DB_PATH;
  const sqlJs = await getSqlJs();

  let fileData = null;
  if (resolvedPath !== ':memory:' && fs.existsSync(resolvedPath)) {
    fileData = fs.readFileSync(resolvedPath);
  }

  const wrapper = new SqlJsDatabase(sqlJs, fileData);
  wrapper._path = resolvedPath;

  // Enable FK constraints (sql.js supports PRAGMA).
  wrapper.exec('PRAGMA foreign_keys = ON;');

  // Schema migrations.
  wrapper.exec(`
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
      completed   INTEGER NOT NULL DEFAULT 0,
      priority    TEXT    NOT NULL DEFAULT 'medium',
      status      TEXT    NOT NULL DEFAULT 'pending',
      due_at      TEXT,
      created_at  TEXT    NOT NULL DEFAULT (datetime('now')),
      updated_at  TEXT    NOT NULL DEFAULT (datetime('now'))
    );
  `);

  return wrapper;
}

module.exports = { openDb };
