'use strict';

const { createTestDb } = require('../helpers/db');
const UserModel = require('../../src/models/user');

let db;
let User;

beforeEach(async () => {
  db = await createTestDb();
  User = UserModel(db);
});

afterEach(() => {
  db.close();
});

describe('User.create', () => {
  test('creates a user and returns the record with id', () => {
    const user = User.create({ username: 'alice', email: 'alice@example.com' });
    expect(user.id).toBeDefined();
    expect(user.username).toBe('alice');
    expect(user.email).toBe('alice@example.com');
    expect(user.created_at).toBeTruthy();
  });

  test('throws on duplicate username', () => {
    User.create({ username: 'alice', email: 'alice@example.com' });
    expect(() => {
      User.create({ username: 'alice', email: 'other@example.com' });
    }).toThrow();
  });

  test('throws on duplicate email', () => {
    User.create({ username: 'alice', email: 'alice@example.com' });
    expect(() => {
      User.create({ username: 'bob', email: 'alice@example.com' });
    }).toThrow();
  });
});

describe('User.findById', () => {
  test('returns user when found', () => {
    const created = User.create({ username: 'bob', email: 'bob@example.com' });
    const found = User.findById(created.id);
    expect(found).not.toBeUndefined();
    expect(found.id).toBe(created.id);
  });

  test('returns undefined when not found', () => {
    const found = User.findById(9999);
    expect(found).toBeUndefined();
  });
});

describe('User.list', () => {
  test('returns empty array when no users', () => {
    expect(User.list()).toEqual([]);
  });

  test('returns all users', () => {
    User.create({ username: 'alice', email: 'alice@example.com' });
    User.create({ username: 'bob', email: 'bob@example.com' });
    const users = User.list();
    expect(users).toHaveLength(2);
  });
});

describe('User.update', () => {
  test('updates username', () => {
    const created = User.create({ username: 'alice', email: 'alice@example.com' });
    const updated = User.update(created.id, { username: 'alice2' });
    expect(updated.username).toBe('alice2');
    expect(updated.email).toBe('alice@example.com');
  });

  test('returns current record when no fields provided', () => {
    const created = User.create({ username: 'alice', email: 'alice@example.com' });
    const result = User.update(created.id, {});
    expect(result.id).toBe(created.id);
  });
});

describe('User.remove', () => {
  test('deletes user and returns true', () => {
    const created = User.create({ username: 'alice', email: 'alice@example.com' });
    const deleted = User.remove(created.id);
    expect(deleted).toBe(true);
    expect(User.findById(created.id)).toBeUndefined();
  });

  test('returns false for non-existent user', () => {
    expect(User.remove(9999)).toBe(false);
  });

  test('cascades deletion to todos', async () => {
    const TodoModel = require('../../src/models/todo');
    const testDb = await createTestDb();
    const TestUser = UserModel(testDb);
    const TestTodo = TodoModel(testDb);

    const user = TestUser.create({ username: 'cascade', email: 'c@c.com' });
    TestTodo.create(user.id, { title: 'task 1' });
    TestUser.remove(user.id);

    const remaining = TestTodo.listByUser(user.id);
    expect(remaining).toHaveLength(0);
    testDb.close();
  });
});
