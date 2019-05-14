const express = require('express');
const router = express.Router();
const proxy = require('./modules/proxy');

router.get('/', (req, res, next) => {
  res.send("API is alive");
});

router.use('/tasks', proxy('tasks'));
router.use('/accounts', proxy('accounts'));
router.use('/habits', proxy('habits'));
router.use('/reports', proxy('reports'));

module.exports = router;
