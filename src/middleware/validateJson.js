'use strict';

/**
 * Handles malformed JSON bodies.
 */
function validateJson(err, req, res, next) {
  if (err.type === 'entity.parse.failed') {
    return res.status(400).json({ error: 'Invalid JSON body' });
  }
  return next(err);
}

module.exports = validateJson;
