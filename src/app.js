'use strict';

require('express-async-errors');

const express = require('express');
const { openDb } = require('./db');
const usersRouter = require('./routes/users');
const todosRouter = require('./routes/todos');
const validateJson = require('./middleware/validateJson');
const errorHandler = require('./middleware/errorHandler');

/**
 * Creates and configures the Express application.
 * @param {object} [db] - Optional DB instance (used in tests).
 * @returns {import('express').Application}
 */
function createApp(db) {
  const app = express();

  app.use(express.json());

  // Routes — db must already be resolved before calling createApp.
  app.use('/users', usersRouter(db));
  app.use('/users/:userId/todos', todosRouter(db));

  // Error middlewares (order matters)
  app.use(validateJson);
  app.use(errorHandler);

  return app;
}

/**
 * Creates app with a fresh in-memory or file-backed DB.
 * @param {string} [dbPath]
 * @returns {Promise<import('express').Application>}
 */
async function createAppAsync(dbPath) {
  const db = await openDb(dbPath);
  return createApp(db);
}

module.exports = { createApp, createAppAsync };
