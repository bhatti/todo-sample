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

describe('Users API', () => {
  // ── GET /api/users ──────────────────────────────────────────────────────────

  describe('GET /api/users', () => {
    it('returns empty array when no users exist', async () => {
      const res = await request(app).get('/api/users');
      expect(res.status).toBe(200);
      expect(res.body).toEqual([]);
    });
  });

  // ── POST /api/users ─────────────────────────────────────────────────────────

  describe('POST /api/users', () => {
    it('creates a user and returns 201', async () => {
      const res = await request(app)
        .post('/api/users')
        .send({ username: 'alice', email: 'alice@example.com' });

      expect(res.status).toBe(201);
      expect(res.body).toMatchObject({
        id: expect.any(Number),
        username: 'alice',
        email: 'alice@example.com',
        created_at: expect.any(String),
      });
    });

    it('returns 400 when username is missing', async () => {
      const res = await request(app)
        .post('/api/users')
        .send({ email: 'bob@example.com' });
      expect(res.status).toBe(400);
      expect(res.body.error).toMatch(/username/i);
    });

    it('returns 400 when email is missing', async () => {
      const res = await request(app)
        .post('/api/users')
        .send({ username: 'bob' });
      expect(res.status).toBe(400);
      expect(res.body.error).toMatch(/email/i);
    });

    it('returns 409 when username is already taken', async () => {
      await request(app)
        .post('/api/users')
        .send({ username: 'charlie', email: 'charlie@example.com' });

      const res = await request(app)
        .post('/api/users')
        .send({ username: 'charlie', email: 'charlie2@example.com' });
      expect(res.status).toBe(409);
    });

    it('returns 409 when email is already taken', async () => {
      await request(app)
        .post('/api/users')
        .send({ username: 'dave', email: 'dave@example.com' });

      const res = await request(app)
        .post('/api/users')
        .send({ username: 'dave2', email: 'dave@example.com' });
      expect(res.status).toBe(409);
    });
  });

  // ── GET /api/users/:id ──────────────────────────────────────────────────────

  describe('GET /api/users/:id', () => {
    it('returns the user when found', async () => {
      const createRes = await request(app)
        .post('/api/users')
        .send({ username: 'eve', email: 'eve@example.com' });
      const { id } = createRes.body;

      const res = await request(app).get(`/api/users/${id}`);
      expect(res.status).toBe(200);
      expect(res.body.username).toBe('eve');
    });

    it('returns 404 when user not found', async () => {
      const res = await request(app).get('/api/users/999999');
      expect(res.status).toBe(404);
      expect(res.body.error).toMatch(/not found/i);
    });
  });

  // ── PUT /api/users/:id ──────────────────────────────────────────────────────

  describe('PUT /api/users/:id', () => {
    it('updates username', async () => {
      const createRes = await request(app)
        .post('/api/users')
        .send({ username: 'frank', email: 'frank@example.com' });
      const { id } = createRes.body;

      const res = await request(app)
        .put(`/api/users/${id}`)
        .send({ username: 'franklin' });
      expect(res.status).toBe(200);
      expect(res.body.username).toBe('franklin');
      expect(res.body.email).toBe('frank@example.com');
    });

    it('returns 404 when user not found', async () => {
      const res = await request(app)
        .put('/api/users/999999')
        .send({ username: 'ghost' });
      expect(res.status).toBe(404);
    });
  });

  // ── DELETE /api/users/:id ───────────────────────────────────────────────────

  describe('DELETE /api/users/:id', () => {
    it('deletes a user and returns 204', async () => {
      const createRes = await request(app)
        .post('/api/users')
        .send({ username: 'grace', email: 'grace@example.com' });
      const { id } = createRes.body;

      const deleteRes = await request(app).delete(`/api/users/${id}`);
      expect(deleteRes.status).toBe(204);

      const getRes = await request(app).get(`/api/users/${id}`);
      expect(getRes.status).toBe(404);
    });

    it('returns 404 when user not found', async () => {
      const res = await request(app).delete('/api/users/999999');
      expect(res.status).toBe(404);
    });
  });

  // ── GET /api/users (after data) ─────────────────────────────────────────────

  describe('GET /api/users (with data)', () => {
    it('returns all existing users', async () => {
      const res = await request(app).get('/api/users');
      expect(res.status).toBe(200);
      expect(Array.isArray(res.body)).toBe(true);
      // At least the ones created in this suite that weren't deleted
      expect(res.body.length).toBeGreaterThan(0);
    });
  });
});
