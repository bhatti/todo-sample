'use strict';

const request = require('supertest');
const { createApp } = require('../../src/app');
const { createTestDb } = require('../helpers/db');

let app;
let db;

beforeEach(async () => {
  db = await createTestDb();
  app = createApp(db);
});

afterEach(() => {
  db.close();
});

describe('POST /users', () => {
  test('201 – creates a user', async () => {
    const res = await request(app)
      .post('/users')
      .send({ username: 'alice', email: 'alice@example.com' });

    expect(res.status).toBe(201);
    expect(res.body.id).toBeDefined();
    expect(res.body.username).toBe('alice');
    expect(res.body.email).toBe('alice@example.com');
  });

  test('400 – missing username', async () => {
    const res = await request(app)
      .post('/users')
      .send({ email: 'alice@example.com' });

    expect(res.status).toBe(400);
    expect(res.body.error).toMatch(/username/i);
  });

  test('400 – missing email', async () => {
    const res = await request(app)
      .post('/users')
      .send({ username: 'alice' });

    expect(res.status).toBe(400);
    expect(res.body.error).toMatch(/email/i);
  });

  test('409 – duplicate username', async () => {
    await request(app).post('/users').send({ username: 'alice', email: 'a@a.com' });
    const res = await request(app).post('/users').send({ username: 'alice', email: 'b@b.com' });
    expect(res.status).toBe(409);
  });

  test('409 – duplicate email', async () => {
    await request(app).post('/users').send({ username: 'alice', email: 'a@a.com' });
    const res = await request(app).post('/users').send({ username: 'bob', email: 'a@a.com' });
    expect(res.status).toBe(409);
  });
});

describe('GET /users', () => {
  test('200 – returns empty array initially', async () => {
    const res = await request(app).get('/users');
    expect(res.status).toBe(200);
    expect(res.body).toEqual([]);
  });

  test('200 – returns all users', async () => {
    await request(app).post('/users').send({ username: 'alice', email: 'a@a.com' });
    await request(app).post('/users').send({ username: 'bob', email: 'b@b.com' });
    const res = await request(app).get('/users');
    expect(res.status).toBe(200);
    expect(res.body).toHaveLength(2);
  });
});

describe('GET /users/:id', () => {
  test('200 – returns user by id', async () => {
    const created = await request(app).post('/users').send({ username: 'alice', email: 'a@a.com' });
    const res = await request(app).get(`/users/${created.body.id}`);
    expect(res.status).toBe(200);
    expect(res.body.username).toBe('alice');
  });

  test('404 – user not found', async () => {
    const res = await request(app).get('/users/9999');
    expect(res.status).toBe(404);
    expect(res.body.error).toMatch(/not found/i);
  });
});

describe('PUT /users/:id', () => {
  test('200 – updates user', async () => {
    const created = await request(app).post('/users').send({ username: 'alice', email: 'a@a.com' });
    const res = await request(app)
      .put(`/users/${created.body.id}`)
      .send({ username: 'alice2' });

    expect(res.status).toBe(200);
    expect(res.body.username).toBe('alice2');
  });

  test('404 – user not found', async () => {
    const res = await request(app).put('/users/9999').send({ username: 'x' });
    expect(res.status).toBe(404);
  });
});

describe('DELETE /users/:id', () => {
  test('204 – deletes user', async () => {
    const created = await request(app).post('/users').send({ username: 'alice', email: 'a@a.com' });
    const res = await request(app).delete(`/users/${created.body.id}`);
    expect(res.status).toBe(204);
  });

  test('204 – cascades todos on delete', async () => {
    const userRes = await request(app).post('/users').send({ username: 'alice', email: 'a@a.com' });
    const userId = userRes.body.id;
    await request(app).post(`/users/${userId}/todos`).send({ title: 'Task 1' });

    await request(app).delete(`/users/${userId}`);

    // User should be gone
    const getRes = await request(app).get(`/users/${userId}`);
    expect(getRes.status).toBe(404);
  });

  test('404 – user not found', async () => {
    const res = await request(app).delete('/users/9999');
    expect(res.status).toBe(404);
  });
});
