const express = require('express');
const moment = require('moment-timezone');
const request = require('request');
const router = express.Router();

moment.tz.setDefault('America/Mexico_City');

router.get('/', (req, res) => {
    res.send('Hello from reports');
});

module.exports = router;