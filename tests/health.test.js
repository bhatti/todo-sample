'use strict';

const request = require('supertest');
const { createApp } = require('../src/app');
const { createTestDb, closeTestDb } = require('./helpers/db');

let app;
let db;

beforeAll(() => {
  db = createTestDb();
  app = createApp();
});

afterAll(() => {
  closeTestDb(db);
});

describe('GET /health', () => {
  it('returns 200 with status ok', async () => {
    const res = await request(app).get('/health');
    expect(res.status).toBe(200);
    expect(res.body).toEqual({ status: 'ok' });
  });
});
