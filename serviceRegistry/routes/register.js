const express = require('express');
const redis = require('../modules/redis');
const Constants = require('../constants.json');
const router = express.Router();

router.post('/', (req, res) => {
    let ip = req.connection.remoteAddress.split(':').slice(-1)[0];
    if (ip === '1') {
        ip = '127.0.0.1';
    }
    const { port, service } = req.body;
    if (port === undefined || service === undefined) {
        return res.status(400).send({ error: 400, message: 'Missing params' });
    }
    const validService = service === Constants.names.accountsMicroservice ||
        service === Constants.names.tasksMicroservice ||
        service === Constants.names.habitsMicroservice ||
        service === Constants.names.reportsMicroservice;

    if (!validService) {
        return res.status(502).send({ error: 502, message: 'Service is not required or existent' });
    }
    redis.isServiceRegistered(service, ip, port)
        .then(alreadyExists => {
            if (alreadyExists) {
                return res.status(409).send({ status: 409, message: 'Service is already registered' });
            }
            return redis.addServiceToRegistry(service, ip, port);
        })
        .then(() => {
            res.status(200).send({ message: 'Service registered succesfully' });
        })
        .catch(err => {
            console.log(err);
            res.status(500).send({ status: 500, message: 'There was an error while registering the service' });
        });
});

module.exports = router;