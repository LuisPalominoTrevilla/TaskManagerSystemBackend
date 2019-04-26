const express = require('express');
const router = express.Router();
const task_module = require('./modules/tasks');

router.get('/', (req, res, next) => {
  res.send("API is alive");
});

router.use('/tasks', task_module);

module.exports = router;
