'use strict';

const { createApp } = require('./app');
// Initialise the default file-based DB on startup
require('./db/database').getDb();

const PORT = process.env.PORT || 3000;
const app = createApp();

app.listen(PORT, () => {
  console.log(`Todo API listening on port ${PORT}`);
});
