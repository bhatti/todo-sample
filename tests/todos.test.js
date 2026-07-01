'use strict';

const request = require('supertest');
const { createApp } = require('../src/app');
const { createTestDb, closeTestDb } = require('./helpers/db');

let app;
let db;
let userId;

beforeAll(async () => {
  db = createTestDb();
  app = createApp();

  // Create a user to own todos in this suite
  const res = await request(app)
    .post('/api/users')
    .send({ username: 'testuser', email: 'testuser@example.com' });
  userId = res.body.id;
});

afterAll(() => {
  closeTestDb(db);
});

describe('Todos API', () => {
  // ── GET /api/users/:userId/todos ────────────────────────────────────────────

  describe('GET /api/users/:userId/todos', () => {
    it('returns empty array when user has no todos', async () => {
      const res = await request(app).get(`/api/users/${userId}/todos`);
      expect(res.status).toBe(200);
      expect(res.body).toEqual([]);
    });

    it('returns 404 when user does not exist', async () => {
      const res = await request(app).get('/api/users/999999/todos');
      expect(res.status).toBe(404);
    });
  });

  // ── POST /api/users/:userId/todos ───────────────────────────────────────────

  describe('POST /api/users/:userId/todos', () => {
    it('creates a todo and returns 201', async () => {
      const res = await request(app)
        .post(`/api/users/${userId}/todos`)
        .send({ title: 'Buy milk', description: 'Organic if possible' });

      expect(res.status).toBe(201);
      expect(res.body).toMatchObject({
        id: expect.any(Number),
        user_id: userId,
        title: 'Buy milk',
        description: 'Organic if possible',
        completed: false,
        created_at: expect.any(String),
        updated_at: expect.any(String),
      });
    });

    it('creates a todo without a description', async () => {
      const res = await request(app)
        .post(`/api/users/${userId}/todos`)
        .send({ title: 'Read a book' });

      expect(res.status).toBe(201);
      expect(res.body.title).toBe('Read a book');
      expect(res.body.completed).toBe(false);
    });

    it('returns 400 when title is missing', async () => {
      const res = await request(app)
        .post(`/api/users/${userId}/todos`)
        .send({ description: 'No title here' });
      expect(res.status).toBe(400);
      expect(res.body.error).toMatch(/title/i);
    });

    it('returns 404 when user does not exist', async () => {
      const res = await request(app)
        .post('/api/users/999999/todos')
        .send({ title: 'Ghost todo' });
      expect(res.status).toBe(404);
    });
  });

  // ── GET /api/users/:userId/todos/:id ────────────────────────────────────────

  describe('GET /api/users/:userId/todos/:id', () => {
    it('returns a single todo', async () => {
      const createRes = await request(app)
        .post(`/api/users/${userId}/todos`)
        .send({ title: 'Walk the dog' });
      const { id } = createRes.body;

      const res = await request(app).get(`/api/users/${userId}/todos/${id}`);
      expect(res.status).toBe(200);
      expect(res.body.title).toBe('Walk the dog');
    });

    it('returns 404 when todo does not exist', async () => {
      const res = await request(app).get(`/api/users/${userId}/todos/999999`);
      expect(res.status).toBe(404);
    });

    it('returns 404 when user does not exist', async () => {
      const res = await request(app).get('/api/users/999999/todos/1');
      expect(res.status).toBe(404);
    });
  });

  // ── PUT /api/users/:userId/todos/:id ────────────────────────────────────────

  describe('PUT /api/users/:userId/todos/:id', () => {
    it('updates title and description', async () => {
      const createRes = await request(app)
        .post(`/api/users/${userId}/todos`)
        .send({ title: 'Original title', description: 'Original desc' });
      const { id } = createRes.body;

      const res = await request(app)
        .put(`/api/users/${userId}/todos/${id}`)
        .send({ title: 'Updated title', description: 'Updated desc' });

      expect(res.status).toBe(200);
      expect(res.body.title).toBe('Updated title');
      expect(res.body.description).toBe('Updated desc');
    });

    it('marks a todo as completed', async () => {
      const createRes = await request(app)
        .post(`/api/users/${userId}/todos`)
        .send({ title: 'Finish project' });
      const { id } = createRes.body;

      const res = await request(app)
        .put(`/api/users/${userId}/todos/${id}`)
        .send({ completed: true });

      expect(res.status).toBe(200);
      expect(res.body.completed).toBe(true);
    });

    it('returns 404 when todo does not exist', async () => {
      const res = await request(app)
        .put(`/api/users/${userId}/todos/999999`)
        .send({ title: 'Ghost' });
      expect(res.status).toBe(404);
    });

    it('returns 404 when user does not exist', async () => {
      const res = await request(app)
        .put('/api/users/999999/todos/1')
        .send({ title: 'Ghost' });
      expect(res.status).toBe(404);
    });
  });

  // ── DELETE /api/users/:userId/todos/:id ─────────────────────────────────────

  describe('DELETE /api/users/:userId/todos/:id', () => {
    it('deletes a todo and returns 204', async () => {
      const createRes = await request(app)
        .post(`/api/users/${userId}/todos`)
        .send({ title: 'Temp todo' });
      const { id } = createRes.body;

      const deleteRes = await request(app).delete(`/api/users/${userId}/todos/${id}`);
      expect(deleteRes.status).toBe(204);

      const getRes = await request(app).get(`/api/users/${userId}/todos/${id}`);
      expect(getRes.status).toBe(404);
    });

    it('returns 404 when todo does not exist', async () => {
      const res = await request(app).delete(`/api/users/${userId}/todos/999999`);
      expect(res.status).toBe(404);
    });

    it('returns 404 when user does not exist', async () => {
      const res = await request(app).delete('/api/users/999999/todos/1');
      expect(res.status).toBe(404);
    });
  });

  // ── Cascade delete ──────────────────────────────────────────────────────────

  describe('Cascade delete', () => {
    it('deleting a user also deletes their todos', async () => {
      // Create a new user
      const userRes = await request(app)
        .post('/api/users')
        .send({ username: 'temp_user', email: 'temp@example.com' });
      const tempUserId = userRes.body.id;

      // Create a todo for that user
      const todoRes = await request(app)
        .post(`/api/users/${tempUserId}/todos`)
        .send({ title: 'Will be deleted' });
      const todoId = todoRes.body.id;

      // Delete the user
      await request(app).delete(`/api/users/${tempUserId}`);

      // The user should be gone
      const userCheck = await request(app).get(`/api/users/${tempUserId}`);
      expect(userCheck.status).toBe(404);

      // Verify todo is gone too (query DB directly)
      const { getDb } = require('../src/db/database');
      const row = getDb().prepare('SELECT * FROM todos WHERE id = ?').get(todoId);
      expect(row).toBeUndefined();
    });
  });
});
