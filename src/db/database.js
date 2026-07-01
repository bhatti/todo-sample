'use strict';

const { DatabaseSync } = require('node:sqlite');
const path = require('path');
const { runMigrations } = require('./migrations');

let _db = null;

/**
 * Returns the shared database instance.
 * Initialises it once using the given path (defaults to a file-based DB).
 * @param {string} [dbPath] - Pass ':memory:' for an in-memory DB (tests).
 * @returns {DatabaseSync}
 */
function getDb(dbPath) {
  if (_db) return _db;

  const resolvedPath = dbPath || path.join(__dirname, '..', '..', 'todo.db');
  _db = new DatabaseSync(resolvedPath);

  // Enable foreign-key support
  _db.exec('PRAGMA foreign_keys = ON');

  runMigrations(_db);
  return _db;
}

/**
 * Reset the singleton (used in tests to inject a fresh in-memory DB).
 * @param {DatabaseSync|null} db
 */
function setDb(db) {
  _db = db;
}

module.exports = { getDb, setDb };
