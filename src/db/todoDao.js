'use strict';

const { getDb } = require('./database');

/**
 * Create a new todo for a user.
 * @param {{ user_id: number, title: string, description?: string }} data
 * @returns {object} The created todo row.
 */
function createTodo({ user_id, title, description }) {
  const db = getDb();
  const stmt = db.prepare(
    'INSERT INTO todos (user_id, title, description) VALUES (?, ?, ?)'
  );
  const result = stmt.run(user_id, title, description || null);
  return getTodoById(result.lastInsertRowid);
}

/**
 * Get a single todo by ID.
 * @param {number} id
 * @returns {object|undefined}
 */
function getTodoById(id) {
  const db = getDb();
  const row = db.prepare('SELECT * FROM todos WHERE id = ?').get(id);
  return row ? normalizeTodo(row) : undefined;
}

/**
 * Get all todos for a user.
 * @param {number} userId
 * @returns {object[]}
 */
function getTodosByUserId(userId) {
  const db = getDb();
  const rows = db.prepare('SELECT * FROM todos WHERE user_id = ? ORDER BY id').all(userId);
  return rows.map(normalizeTodo);
}

/**
 * Update a todo.
 * @param {number} id
 * @param {{ title?: string, description?: string, completed?: boolean }} data
 * @returns {object|undefined}
 */
function updateTodo(id, { title, description, completed }) {
  const db = getDb();
  const existing = getTodoById(id);
  if (!existing) return undefined;

  const newTitle = title !== undefined ? title : existing.title;
  const newDescription = description !== undefined ? description : existing.description;
  const newCompleted = completed !== undefined ? (completed ? 1 : 0) : (existing.completed ? 1 : 0);

  db.prepare(
    `UPDATE todos
     SET title = ?, description = ?, completed = ?, updated_at = datetime('now')
     WHERE id = ?`
  ).run(newTitle, newDescription, newCompleted, id);

  return getTodoById(id);
}

/**
 * Delete a todo.
 * @param {number} id
 * @returns {boolean}
 */
function deleteTodo(id) {
  const db = getDb();
  const result = db.prepare('DELETE FROM todos WHERE id = ?').run(id);
  return result.changes > 0;
}

/**
 * Convert SQLite integer (0/1) to boolean for the completed field.
 * @param {object} row
 * @returns {object}
 */
function normalizeTodo(row) {
  return { ...row, completed: row.completed === 1 };
}

module.exports = { createTodo, getTodoById, getTodosByUserId, updateTodo, deleteTodo };
