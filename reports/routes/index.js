const express = require('express');
const router = express.Router();
const reports = require('./modules/reports');

router.get('/', (req, res, next) => {
  res.send("API is alive");
});

router.use('/reports', reports);

module.exports = router;
