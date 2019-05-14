const express = require('express');
const moment = require('moment-timezone');
const serviceFactory = require('../../db/servicesFactory');
const router = express.Router();

moment.tz.setDefault('America/Mexico_City');

router.get('/', (req, res) => {
    res.send('Hello from reports');
});

router.get('/users/:userId', (req, res) => {
    const accounts = serviceFactory('accounts');
    const userId = req.params.userId;
    accounts.getTasksDueToday(userId)
        .then(tasks => {
            res.status(200).json(tasks);
        })
        .catch(err => {
            res.status(err.error).send(err);
        })
});

module.exports = router;