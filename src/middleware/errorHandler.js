'use strict';

/**
 * Centralised Express error handler.
 * Maps known SQLite constraint errors to HTTP 409, everything else to 500.
 */
// eslint-disable-next-line no-unused-vars
function errorHandler(err, req, res, next) {
  // SQLite UNIQUE constraint violation
  if (err.code === 'SQLITE_CONSTRAINT_UNIQUE' || err.message?.includes('UNIQUE constraint failed')) {
    return res.status(409).json({ error: err.message });
  }

  console.error(err);
  res.status(500).json({ error: 'Internal server error' });
}

module.exports = { errorHandler };
