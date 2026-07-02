'use strict';

/**
 * Centralised Express error-handling middleware.
 * Must be registered last (after all routes).
 */
// eslint-disable-next-line no-unused-vars
function errorHandler(err, req, res, next) {
  console.error(err);

  // SQLite unique constraint violation
  if (err.code === 'SQLITE_CONSTRAINT_UNIQUE' || (err.message && err.message.includes('UNIQUE constraint failed'))) {
    return res.status(409).json({ error: err.message || 'Unique constraint violation' });
  }

  // SQLite foreign-key violation
  if (err.code === 'SQLITE_CONSTRAINT_FOREIGNKEY' || (err.message && err.message.includes('FOREIGN KEY constraint failed'))) {
    return res.status(400).json({ error: 'Foreign key constraint failed' });
  }

  const status = err.status || err.statusCode || 500;
  const message = err.message || 'Internal server error';
  return res.status(status).json({ error: message });
}

module.exports = errorHandler;
