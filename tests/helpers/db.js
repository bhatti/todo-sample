'use strict';

const { DatabaseSync } = require('node:sqlite');
const { runMigrations } = require('../../src/db/migrations');
const { setDb } = require('../../src/db/database');

/**
 * Create a fresh in-memory SQLite database and inject it as the singleton.
 * Call this before each test suite (beforeAll) and reset after (afterAll).
 * @returns {DatabaseSync}
 */
function createTestDb() {
  const db = new DatabaseSync(':memory:');
  db.exec('PRAGMA foreign_keys = ON');
  runMigrations(db);
  setDb(db);
  return db;
}

/**
 * Close and tear down the test database.
 * @param {DatabaseSync} db
 */
function closeTestDb(db) {
  if (db) db.close();
  setDb(null);
}

module.exports = { createTestDb, closeTestDb };
