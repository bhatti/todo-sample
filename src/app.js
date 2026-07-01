'use strict';

const express = require('express');
const usersRouter = require('./routes/users');
const todosRouter = require('./routes/todos');
const { errorHandler } = require('./middleware/errorHandler');

/**
 * Create and configure the Express application.
 * Does NOT call app.listen() so it can be imported cleanly in tests.
 * @returns {import('express').Application}
 */
function createApp() {
  const app = express();

  app.use(express.json());

  // Health check
  app.get('/health', (req, res) => {
    res.json({ status: 'ok' });
  });

  // API routes
  app.use('/api/users', usersRouter);
  app.use('/api/users/:userId/todos', todosRouter);

  // Centralised error handler (must be last)
  app.use(errorHandler);

  return app;
}

module.exports = { createApp };
