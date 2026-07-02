'use strict';

const { createAppAsync } = require('./app');

const PORT = process.env.PORT || 3000;

createAppAsync().then((app) => {
  app.listen(PORT, () => {
    console.log(`Todo API listening on port ${PORT}`);
  });
}).catch((err) => {
  console.error('Failed to start server:', err);
  process.exit(1);
});
