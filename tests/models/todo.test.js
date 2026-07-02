'use strict';

const { createTestDb } = require('../helpers/db');
const UserModel = require('../../src/models/user');
const TodoModel = require('../../src/models/todo');

let db;
let User;
let Todo;
let testUser;

beforeEach(async () => {
  db = await createTestDb();
  User = UserModel(db);
  Todo = TodoModel(db);
  testUser = User.create({ username: 'alice', email: 'alice@example.com' });
});

afterEach(() => {
  db.close();
});

describe('Todo.create', () => {
  test('creates a todo for an existing user', () => {
    const todo = Todo.create(testUser.id, { title: 'Buy milk' });
    expect(todo.id).toBeDefined();
    expect(todo.title).toBe('Buy milk');
    expect(todo.user_id).toBe(testUser.id);
    expect(todo.completed).toBe(0);
    expect(todo.priority).toBe('medium');
    expect(todo.status).toBe('pending');
  });

  test('throws FK error for non-existent user', () => {
    expect(() => {
      Todo.create(9999, { title: 'Ghost task' });
    }).toThrow();
  });

  test('stores optional fields', () => {
    const todo = Todo.create(testUser.id, {
      title: 'Study',
      description: 'Read books',
      priority: 'high',
      status: 'in_progress',
      due_at: '2026-12-31',
    });
    expect(todo.description).toBe('Read books');
    expect(todo.priority).toBe('high');
    expect(todo.status).toBe('in_progress');
    expect(todo.due_at).toBe('2026-12-31');
  });
});

describe('Todo.listByUser', () => {
  test('returns empty array for user with no todos', () => {
    expect(Todo.listByUser(testUser.id)).toEqual([]);
  });

  test('returns all todos for the user', () => {
    Todo.create(testUser.id, { title: 'Task 1' });
    Todo.create(testUser.id, { title: 'Task 2' });
    const todos = Todo.listByUser(testUser.id);
    expect(todos).toHaveLength(2);
  });

  test('does not return todos from other users', () => {
    const otherUser = User.create({ username: 'bob', email: 'bob@example.com' });
    Todo.create(otherUser.id, { title: "Bob's task" });
    expect(Todo.listByUser(testUser.id)).toHaveLength(0);
  });
});

describe('Todo.findByIdAndUser', () => {
  test('returns todo when found', () => {
    const todo = Todo.create(testUser.id, { title: 'Find me' });
    const found = Todo.findByIdAndUser(testUser.id, todo.id);
    expect(found).not.toBeUndefined();
    expect(found.id).toBe(todo.id);
  });

  test('returns undefined for wrong user', () => {
    const otherUser = User.create({ username: 'bob', email: 'bob@example.com' });
    const todo = Todo.create(otherUser.id, { title: "Bob's" });
    expect(Todo.findByIdAndUser(testUser.id, todo.id)).toBeUndefined();
  });

  test('returns undefined for non-existent todo', () => {
    expect(Todo.findByIdAndUser(testUser.id, 9999)).toBeUndefined();
  });
});

describe('Todo.update', () => {
  test('updates title and description', () => {
    const todo = Todo.create(testUser.id, { title: 'Old title' });
    const updated = Todo.update(testUser.id, todo.id, { title: 'New title', description: 'Details' });
    expect(updated.title).toBe('New title');
    expect(updated.description).toBe('Details');
  });

  test('updates completed flag', () => {
    const todo = Todo.create(testUser.id, { title: 'Task' });
    const updated = Todo.update(testUser.id, todo.id, { completed: 1 });
    expect(updated.completed).toBe(1);
  });

  test('returns current record when no fields provided', () => {
    const todo = Todo.create(testUser.id, { title: 'Task' });
    const result = Todo.update(testUser.id, todo.id, {});
    expect(result.id).toBe(todo.id);
  });
});

describe('Todo.toggleCompleted', () => {
  test('toggles from 0 to 1', () => {
    const todo = Todo.create(testUser.id, { title: 'Task' });
    expect(todo.completed).toBe(0);
    const toggled = Todo.toggleCompleted(testUser.id, todo.id);
    expect(toggled.completed).toBe(1);
  });

  test('toggles from 1 to 0', () => {
    const todo = Todo.create(testUser.id, { title: 'Task' });
    Todo.toggleCompleted(testUser.id, todo.id); // 0 → 1
    const toggled = Todo.toggleCompleted(testUser.id, todo.id); // 1 → 0
    expect(toggled.completed).toBe(0);
  });

  test('returns undefined for non-existent todo', () => {
    expect(Todo.toggleCompleted(testUser.id, 9999)).toBeUndefined();
  });
});

describe('Todo.remove', () => {
  test('deletes todo and returns true', () => {
    const todo = Todo.create(testUser.id, { title: 'Delete me' });
    expect(Todo.remove(testUser.id, todo.id)).toBe(true);
    expect(Todo.findByIdAndUser(testUser.id, todo.id)).toBeUndefined();
  });

  test('returns false for non-existent todo', () => {
    expect(Todo.remove(testUser.id, 9999)).toBe(false);
  });
});
