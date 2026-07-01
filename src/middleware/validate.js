'use strict';

/**
 * Returns middleware that checks for required fields in req.body.
 * Responds with 400 if any field is missing or empty.
 * @param {string[]} fields
 */
function requireFields(fields) {
  return (req, res, next) => {
    for (const field of fields) {
      const value = req.body[field];
      if (value === undefined || value === null || value === '') {
        return res.status(400).json({ error: `Field '${field}' is required` });
      }
    }
    next();
  };
}

module.exports = { requireFields };
