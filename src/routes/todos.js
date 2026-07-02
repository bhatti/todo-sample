'use strict';

const { Router } = require('express');
const UserModel = require('../models/user');
const TodoModel = require('../models/todo');

/**
 * Todos router (nested under /users/:userId/todos).
 * @param {import('better-sqlite3').Database} db
 */
function todosRouter(db) {
  const router = Router({ mergeParams: true });
  const User = UserModel(db);
  const Todo = TodoModel(db);

  // Middleware: verify the parent user exists.
  function requireUser(req, res, next) {
    const userId = Number(req.params.userId);
    if (!Number.isInteger(userId)) {
      return res.status(400).json({ error: 'Invalid user id' });
    }
    const user = User.findById(userId);
    if (!user) {
      return res.status(404).json({ error: 'User not found' });
    }
    req.userId = userId;
    return next();
  }

  router.use(requireUser);

  /**
   * POST /users/:userId/todos — create a todo
   * Body: { title: string, description?: string, priority?: string, status?: string, due_at?: string }
   */
  router.post('/', (req, res) => {
    const { title, description, priority, status, due_at } = req.body || {};

    if (!title || typeof title !== 'string' || title.trim() === '') {
      return res.status(400).json({ error: 'title is required' });
    }

    const todo = Todo.create(req.userId, {
      title: title.trim(),
      description,
      priority,
      status,
      due_at,
    });
    return res.status(201).json(todo);
  });

  /**
   * GET /users/:userId/todos — list todos for user
   */
  router.get('/', (req, res) => {
    const todos = Todo.listByUser(req.userId);
    return res.json(todos);
  });

  /**
   * GET /users/:userId/todos/:id — get a single todo
   */
  router.get('/:id', (req, res) => {
    const id = Number(req.params.id);
    if (!Number.isInteger(id)) {
      return res.status(400).json({ error: 'Invalid todo id' });
    }
    const todo = Todo.findByIdAndUser(req.userId, id);
    if (!todo) {
      return res.status(404).json({ error: 'Todo not found' });
    }
    return res.json(todo);
  });

  /**
   * PUT /users/:userId/todos/:id — update a todo
   */
  router.put('/:id', (req, res) => {
    const id = Number(req.params.id);
    if (!Number.isInteger(id)) {
      return res.status(400).json({ error: 'Invalid todo id' });
    }

    const existing = Todo.findByIdAndUser(req.userId, id);
    if (!existing) {
      return res.status(404).json({ error: 'Todo not found' });
    }

    const { title, description, priority, status, due_at, completed } = req.body || {};

    if (title !== undefined && (typeof title !== 'string' || title.trim() === '')) {
      return res.status(400).json({ error: 'title must be a non-empty string' });
    }

    const data = {};
    if (title !== undefined) data.title = title.trim();
    if (description !== undefined) data.description = description;
    if (priority !== undefined) data.priority = priority;
    if (status !== undefined) data.status = status;
    if (due_at !== undefined) data.due_at = due_at;
    if (completed !== undefined) data.completed = completed ? 1 : 0;

    const updated = Todo.update(req.userId, id, data);
    return res.json(updated);
  });

  /**
   * PATCH /users/:userId/todos/:id — toggle completed
   */
  router.patch('/:id', (req, res) => {
    const id = Number(req.params.id);
    if (!Number.isInteger(id)) {
      return res.status(400).json({ error: 'Invalid todo id' });
    }
    const todo = Todo.toggleCompleted(req.userId, id);
    if (!todo) {
      return res.status(404).json({ error: 'Todo not found' });
    }
    return res.json(todo);
  });

  /**
   * DELETE /users/:userId/todos/:id — delete a todo
   */
  router.delete('/:id', (req, res) => {
    const id = Number(req.params.id);
    if (!Number.isInteger(id)) {
      return res.status(400).json({ error: 'Invalid todo id' });
    }
    const deleted = Todo.remove(req.userId, id);
    if (!deleted) {
      return res.status(404).json({ error: 'Todo not found' });
    }
    return res.status(204).send();
  });

  return router;
}

module.exports = todosRouter;
