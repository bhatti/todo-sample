'use strict';

/**
 * User model – thin data-access wrapper around the SQLite `users` table.
 * @param {import('better-sqlite3').Database} db
 */
function UserModel(db) {
  /**
   * Create a new user.
   * @param {{ username: string, email: string }} data
   * @returns {{ id: number, username: string, email: string, created_at: string, updated_at: string }}
   */
  function create({ username, email }) {
    const stmt = db.prepare(
      `INSERT INTO users (username, email) VALUES (?, ?) RETURNING *`
    );
    return stmt.get(username, email);
  }

  /**
   * Return all users.
   * @returns {Array}
   */
  function list() {
    return db.prepare('SELECT * FROM users ORDER BY id').all();
  }

  /**
   * Find a user by id.
   * @param {number} id
   * @returns {object|undefined}
   */
  function findById(id) {
    return db.prepare('SELECT * FROM users WHERE id = ?').get(id);
  }

  /**
   * Update a user.
   * @param {number} id
   * @param {{ username?: string, email?: string }} data
   * @returns {object|undefined}
   */
  function update(id, { username, email }) {
    const fields = [];
    const values = [];

    if (username !== undefined) {
      fields.push('username = ?');
      values.push(username);
    }
    if (email !== undefined) {
      fields.push('email = ?');
      values.push(email);
    }

    if (fields.length === 0) return findById(id);

    fields.push("updated_at = datetime('now')");
    values.push(id);

    const stmt = db.prepare(
      `UPDATE users SET ${fields.join(', ')} WHERE id = ? RETURNING *`
    );
    return stmt.get(...values);
  }

  /**
   * Delete a user (cascades to todos).
   * @param {number} id
   * @returns {boolean} true if a row was deleted
   */
  function remove(id) {
    const result = db.prepare('DELETE FROM users WHERE id = ?').run(id);
    return result.changes > 0;
  }

  return { create, list, findById, update, remove };
}

module.exports = UserModel;
