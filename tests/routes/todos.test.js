'use strict';

const request = require('supertest');
const { createApp } = require('../../src/app');
const { createTestDb } = require('../helpers/db');

let app;
let db;
let userId;

beforeEach(async () => {
  db = await createTestDb();
  app = createApp(db);
  const res = await request(app)
    .post('/users')
    .send({ username: 'alice', email: 'alice@example.com' });
  userId = res.body.id;
});

afterEach(() => {
  db.close();
});

describe('POST /users/:userId/todos', () => {
  test('201 – creates a todo', async () => {
    const res = await request(app)
      .post(`/users/${userId}/todos`)
      .send({ title: 'Buy milk' });

    expect(res.status).toBe(201);
    expect(res.body.id).toBeDefined();
    expect(res.body.title).toBe('Buy milk');
    expect(res.body.user_id).toBe(userId);
    expect(res.body.completed).toBe(0);
  });

  test('404 – unknown userId', async () => {
    const res = await request(app)
      .post('/users/9999/todos')
      .send({ title: 'Ghost' });

    expect(res.status).toBe(404);
    expect(res.body.error).toMatch(/not found/i);
  });

  test('400 – missing title', async () => {
    const res = await request(app)
      .post(`/users/${userId}/todos`)
      .send({ description: 'no title here' });

    expect(res.status).toBe(400);
    expect(res.body.error).toMatch(/title/i);
  });

  test('201 – stores optional fields', async () => {
    const res = await request(app)
      .post(`/users/${userId}/todos`)
      .send({ title: 'Study', description: 'Read Go book', priority: 'high', status: 'in_progress', due_at: '2026-12-31' });

    expect(res.status).toBe(201);
    expect(res.body.description).toBe('Read Go book');
    expect(res.body.priority).toBe('high');
    expect(res.body.status).toBe('in_progress');
    expect(res.body.due_at).toBe('2026-12-31');
  });
});

describe('GET /users/:userId/todos', () => {
  test('200 – returns empty list', async () => {
    const res = await request(app).get(`/users/${userId}/todos`);
    expect(res.status).toBe(200);
    expect(res.body).toEqual([]);
  });

  test('200 – returns user todos', async () => {
    await request(app).post(`/users/${userId}/todos`).send({ title: 'Task 1' });
    await request(app).post(`/users/${userId}/todos`).send({ title: 'Task 2' });
    const res = await request(app).get(`/users/${userId}/todos`);
    expect(res.status).toBe(200);
    expect(res.body).toHaveLength(2);
  });

  test('404 – unknown userId', async () => {
    const res = await request(app).get('/users/9999/todos');
    expect(res.status).toBe(404);
  });
});

describe('GET /users/:userId/todos/:id', () => {
  test('200 – returns todo', async () => {
    const created = await request(app)
      .post(`/users/${userId}/todos`)
      .send({ title: 'Find me' });
    const res = await request(app).get(`/users/${userId}/todos/${created.body.id}`);
    expect(res.status).toBe(200);
    expect(res.body.title).toBe('Find me');
  });

  test('404 – todo not found', async () => {
    const res = await request(app).get(`/users/${userId}/todos/9999`);
    expect(res.status).toBe(404);
  });

  test('404 – userId not found', async () => {
    const res = await request(app).get('/users/9999/todos/1');
    expect(res.status).toBe(404);
  });
});

describe('PUT /users/:userId/todos/:id', () => {
  test('200 – updates todo', async () => {
    const created = await request(app)
      .post(`/users/${userId}/todos`)
      .send({ title: 'Old' });
    const res = await request(app)
      .put(`/users/${userId}/todos/${created.body.id}`)
      .send({ title: 'New', description: 'Updated' });

    expect(res.status).toBe(200);
    expect(res.body.title).toBe('New');
    expect(res.body.description).toBe('Updated');
  });

  test('404 – todo not found', async () => {
    const res = await request(app).put(`/users/${userId}/todos/9999`).send({ title: 'x' });
    expect(res.status).toBe(404);
  });

  test('400 – empty title', async () => {
    const created = await request(app)
      .post(`/users/${userId}/todos`)
      .send({ title: 'Task' });
    const res = await request(app)
      .put(`/users/${userId}/todos/${created.body.id}`)
      .send({ title: '   ' });
    expect(res.status).toBe(400);
  });
});

describe('PATCH /users/:userId/todos/:id', () => {
  test('200 – toggles completed from 0 to 1', async () => {
    const created = await request(app)
      .post(`/users/${userId}/todos`)
      .send({ title: 'Task' });

    expect(created.body.completed).toBe(0);

    const res = await request(app).patch(`/users/${userId}/todos/${created.body.id}`);
    expect(res.status).toBe(200);
    expect(res.body.completed).toBe(1);
  });

  test('200 – toggles completed from 1 to 0', async () => {
    const created = await request(app)
      .post(`/users/${userId}/todos`)
      .send({ title: 'Task' });

    await request(app).patch(`/users/${userId}/todos/${created.body.id}`); // 0 → 1
    const res = await request(app).patch(`/users/${userId}/todos/${created.body.id}`); // 1 → 0
    expect(res.status).toBe(200);
    expect(res.body.completed).toBe(0);
  });

  test('404 – todo not found', async () => {
    const res = await request(app).patch(`/users/${userId}/todos/9999`);
    expect(res.status).toBe(404);
  });
});

describe('DELETE /users/:userId/todos/:id', () => {
  test('204 – deletes todo', async () => {
    const created = await request(app)
      .post(`/users/${userId}/todos`)
      .send({ title: 'Delete me' });
    const res = await request(app).delete(`/users/${userId}/todos/${created.body.id}`);
    expect(res.status).toBe(204);
  });

  test('404 – todo not found', async () => {
    const res = await request(app).delete(`/users/${userId}/todos/9999`);
    expect(res.status).toBe(404);
  });
});
