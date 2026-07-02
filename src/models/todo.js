'use strict';

/**
 * Todo model – data-access wrapper around the SQLite `todos` table.
 * @param {import('better-sqlite3').Database} db
 */
function TodoModel(db) {
  /**
   * Create a todo for a user.
   * @param {number} userId
   * @param {{ title: string, description?: string, priority?: string, status?: string, due_at?: string }} data
   * @returns {object}
   */
  function create(userId, { title, description, priority, status, due_at }) {
    const stmt = db.prepare(
      `INSERT INTO todos (user_id, title, description, priority, status, due_at)
       VALUES (?, ?, ?, ?, ?, ?) RETURNING *`
    );
    return stmt.get(
      userId,
      title,
      description || null,
      priority || 'medium',
      status || 'pending',
      due_at || null
    );
  }

  /**
   * List all todos for a user.
   * @param {number} userId
   * @returns {Array}
   */
  function listByUser(userId) {
    return db
      .prepare('SELECT * FROM todos WHERE user_id = ? ORDER BY id')
      .all(userId);
  }

  /**
   * Find a todo by id, scoped to a user.
   * @param {number} userId
   * @param {number} id
   * @returns {object|undefined}
   */
  function findByIdAndUser(userId, id) {
    return db
      .prepare('SELECT * FROM todos WHERE id = ? AND user_id = ?')
      .get(id, userId);
  }

  /**
   * Update a todo.
   * @param {number} userId
   * @param {number} id
   * @param {object} data
   * @returns {object|undefined}
   */
  function update(userId, id, data) {
    const allowed = ['title', 'description', 'priority', 'status', 'due_at', 'completed'];
    const fields = [];
    const values = [];

    for (const key of allowed) {
      if (data[key] !== undefined) {
        fields.push(`${key} = ?`);
        values.push(data[key]);
      }
    }

    if (fields.length === 0) return findByIdAndUser(userId, id);

    fields.push("updated_at = datetime('now')");
    values.push(id, userId);

    const stmt = db.prepare(
      `UPDATE todos SET ${fields.join(', ')} WHERE id = ? AND user_id = ? RETURNING *`
    );
    return stmt.get(...values);
  }

  /**
   * Toggle the completed flag.
   * @param {number} userId
   * @param {number} id
   * @returns {object|undefined}
   */
  function toggleCompleted(userId, id) {
    const todo = findByIdAndUser(userId, id);
    if (!todo) return undefined;
    const newVal = todo.completed ? 0 : 1;
    return db
      .prepare(
        `UPDATE todos SET completed = ?, updated_at = datetime('now')
         WHERE id = ? AND user_id = ? RETURNING *`
      )
      .get(newVal, id, userId);
  }

  /**
   * Delete a todo.
   * @param {number} userId
   * @param {number} id
   * @returns {boolean}
   */
  function remove(userId, id) {
    const result = db
      .prepare('DELETE FROM todos WHERE id = ? AND user_id = ?')
      .run(id, userId);
    return result.changes > 0;
  }

  return { create, listByUser, findByIdAndUser, update, toggleCompleted, remove };
}

module.exports = TodoModel;
