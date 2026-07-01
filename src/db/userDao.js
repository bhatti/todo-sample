'use strict';

const { getDb } = require('./database');

/**
 * Create a new user.
 * @param {{ username: string, email: string }} data
 * @returns {object} The created user row.
 */
function createUser({ username, email }) {
  const db = getDb();
  const stmt = db.prepare(
    'INSERT INTO users (username, email) VALUES (?, ?)'
  );
  const result = stmt.run(username, email);
  return getUserById(result.lastInsertRowid);
}

/**
 * Get a single user by ID.
 * @param {number} id
 * @returns {object|undefined}
 */
function getUserById(id) {
  const db = getDb();
  return db.prepare('SELECT * FROM users WHERE id = ?').get(id);
}

/**
 * Get all users.
 * @returns {object[]}
 */
function getAllUsers() {
  const db = getDb();
  return db.prepare('SELECT * FROM users ORDER BY id').all();
}

/**
 * Update an existing user.
 * @param {number} id
 * @param {{ username?: string, email?: string }} data
 * @returns {object|undefined} The updated user row, or undefined if not found.
 */
function updateUser(id, { username, email }) {
  const db = getDb();
  const existing = getUserById(id);
  if (!existing) return undefined;

  const newUsername = username !== undefined ? username : existing.username;
  const newEmail = email !== undefined ? email : existing.email;

  db.prepare('UPDATE users SET username = ?, email = ? WHERE id = ?')
    .run(newUsername, newEmail, id);

  return getUserById(id);
}

/**
 * Delete a user (cascades to todos).
 * @param {number} id
 * @returns {boolean} True if a row was deleted.
 */
function deleteUser(id) {
  const db = getDb();
  const result = db.prepare('DELETE FROM users WHERE id = ?').run(id);
  return result.changes > 0;
}

module.exports = { createUser, getUserById, getAllUsers, updateUser, deleteUser };
