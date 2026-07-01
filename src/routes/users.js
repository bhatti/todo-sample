'use strict';

const { Router } = require('express');
const { requireFields } = require('../middleware/validate');
const {
  createUser,
  getUserById,
  getAllUsers,
  updateUser,
  deleteUser,
} = require('../db/userDao');

const router = Router();

// GET /api/users — list all users
router.get('/', (req, res) => {
  const users = getAllUsers();
  res.json(users);
});

// POST /api/users — create a user
router.post('/', requireFields(['username', 'email']), (req, res, next) => {
  try {
    const { username, email } = req.body;
    const user = createUser({ username, email });
    res.status(201).json(user);
  } catch (err) {
    next(err);
  }
});

// GET /api/users/:id — get a single user
router.get('/:id', (req, res) => {
  const user = getUserById(Number(req.params.id));
  if (!user) return res.status(404).json({ error: 'User not found' });
  res.json(user);
});

// PUT /api/users/:id — update a user
router.put('/:id', (req, res, next) => {
  try {
    const { username, email } = req.body;
    const user = updateUser(Number(req.params.id), { username, email });
    if (!user) return res.status(404).json({ error: 'User not found' });
    res.json(user);
  } catch (err) {
    next(err);
  }
});

// DELETE /api/users/:id — delete a user
router.delete('/:id', (req, res) => {
  const deleted = deleteUser(Number(req.params.id));
  if (!deleted) return res.status(404).json({ error: 'User not found' });
  res.status(204).end();
});

module.exports = router;
