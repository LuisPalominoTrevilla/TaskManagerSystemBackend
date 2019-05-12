const express = require('express');
const Constants = require('../constants.json');
const registry = require('../modules/registry');
const router = express.Router();

router.post('/', (req, res) => {
    let ip = req.connection.remoteAddress.split(':').slice(-1)[0];
    if (ip === '1') {
        ip = '127.0.0.1';
    }
    const { port, service, healthCheck } = req.body;
    if (port === undefined || service === undefined || healthCheck === undefined) {
        return res.status(400).send({ error: 400, message: 'Missing params' });
    }
    const validService = service === Constants.names.accountsMicroservice ||
        service === Constants.names.tasksMicroservice ||
        service === Constants.names.habitsMicroservice ||
        service === Constants.names.reportsMicroservice;

    if (!validService) {
        return res.status(502).send({ error: 502, message: 'Service is neither required nor existent' });
    }
    registry.registerService(service, ip, port, healthCheck)
        .then(() => {
            res.status(200).send({ message: 'Service registered correctly' });
        })
        .catch(err => {
            res.status(err.error).send(err);
        });
});

module.exports = router;