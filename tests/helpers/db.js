'use strict';

const { openDb } = require('../../src/db');

/**
 * Returns a fresh in-memory SQLite database for isolated tests.
 * @returns {Promise<object>}
 */
async function createTestDb() {
  return openDb(':memory:');
}

module.exports = { createTestDb };
