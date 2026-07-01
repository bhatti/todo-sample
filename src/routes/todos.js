'use strict';

const { Router } = require('express');
const { requireFields } = require('../middleware/validate');
const { getUserById } = require('../db/userDao');
const {
  createTodo,
  getTodoById,
  getTodosByUserId,
  updateTodo,
  deleteTodo,
} = require('../db/todoDao');

const router = Router({ mergeParams: true });

// Helper: resolve user or respond 404
function resolveUser(req, res) {
  const user = getUserById(Number(req.params.userId));
  if (!user) {
    res.status(404).json({ error: 'User not found' });
    return null;
  }
  return user;
}

// GET /api/users/:userId/todos
router.get('/', (req, res) => {
  const user = resolveUser(req, res);
  if (!user) return;
  const todos = getTodosByUserId(user.id);
  res.json(todos);
});

// POST /api/users/:userId/todos
router.post('/', requireFields(['title']), (req, res, next) => {
  try {
    const user = resolveUser(req, res);
    if (!user) return;
    const { title, description } = req.body;
    const todo = createTodo({ user_id: user.id, title, description });
    res.status(201).json(todo);
  } catch (err) {
    next(err);
  }
});

// GET /api/users/:userId/todos/:id
router.get('/:id', (req, res) => {
  const user = resolveUser(req, res);
  if (!user) return;

  const todo = getTodoById(Number(req.params.id));
  if (!todo || todo.user_id !== user.id) {
    return res.status(404).json({ error: 'Todo not found' });
  }
  res.json(todo);
});

// PUT /api/users/:userId/todos/:id
router.put('/:id', (req, res, next) => {
  try {
    const user = resolveUser(req, res);
    if (!user) return;

    const existing = getTodoById(Number(req.params.id));
    if (!existing || existing.user_id !== user.id) {
      return res.status(404).json({ error: 'Todo not found' });
    }

    const { title, description, completed } = req.body;
    const todo = updateTodo(Number(req.params.id), { title, description, completed });
    res.json(todo);
  } catch (err) {
    next(err);
  }
});

// DELETE /api/users/:userId/todos/:id
router.delete('/:id', (req, res) => {
  const user = resolveUser(req, res);
  if (!user) return;

  const existing = getTodoById(Number(req.params.id));
  if (!existing || existing.user_id !== user.id) {
    return res.status(404).json({ error: 'Todo not found' });
  }

  deleteTodo(Number(req.params.id));
  res.status(204).end();
});

module.exports = router;
