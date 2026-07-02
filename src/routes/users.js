'use strict';

const { Router } = require('express');
const UserModel = require('../models/user');

/**
 * Users router.
 * @param {import('better-sqlite3').Database} db
 */
function usersRouter(db) {
  const router = Router();
  const User = UserModel(db);

  /**
   * POST /users — create a user
   * Body: { username: string, email: string }
   */
  router.post('/', (req, res) => {
    const { username, email } = req.body || {};

    if (!username || typeof username !== 'string' || username.trim() === '') {
      return res.status(400).json({ error: 'username is required' });
    }
    if (!email || typeof email !== 'string' || email.trim() === '') {
      return res.status(400).json({ error: 'email is required' });
    }

    try {
      const user = User.create({ username: username.trim(), email: email.trim() });
      return res.status(201).json(user);
    } catch (err) {
      if (err.message && err.message.includes('UNIQUE constraint failed')) {
        return res.status(409).json({ error: err.message });
      }
      throw err;
    }
  });

  /**
   * GET /users — list all users
   */
  router.get('/', (req, res) => {
    const users = User.list();
    return res.json(users);
  });

  /**
   * GET /users/:id — get a single user
   */
  router.get('/:id', (req, res) => {
    const id = Number(req.params.id);
    if (!Number.isInteger(id)) {
      return res.status(400).json({ error: 'Invalid user id' });
    }
    const user = User.findById(id);
    if (!user) {
      return res.status(404).json({ error: 'User not found' });
    }
    return res.json(user);
  });

  /**
   * PUT /users/:id — update a user
   * Body: { username?: string, email?: string }
   */
  router.put('/:id', (req, res) => {
    const id = Number(req.params.id);
    if (!Number.isInteger(id)) {
      return res.status(400).json({ error: 'Invalid user id' });
    }

    const existing = User.findById(id);
    if (!existing) {
      return res.status(404).json({ error: 'User not found' });
    }

    const { username, email } = req.body || {};

    try {
      const updated = User.update(id, { username, email });
      return res.json(updated);
    } catch (err) {
      if (err.message && err.message.includes('UNIQUE constraint failed')) {
        return res.status(409).json({ error: err.message });
      }
      throw err;
    }
  });

  /**
   * DELETE /users/:id — delete a user (cascades todos)
   */
  router.delete('/:id', (req, res) => {
    const id = Number(req.params.id);
    if (!Number.isInteger(id)) {
      return res.status(400).json({ error: 'Invalid user id' });
    }
    const deleted = User.remove(id);
    if (!deleted) {
      return res.status(404).json({ error: 'User not found' });
    }
    return res.status(204).send();
  });

  return router;
}

module.exports = usersRouter;
