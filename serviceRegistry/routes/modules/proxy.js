const express = require('express');
const moment = require('moment-timezone');
const request = require('request');
const registry = require('../../modules/registry');

moment.tz.setDefault('America/Mexico_City');

module.exports = function(microservice) {
    const router = express.Router();

    router.get('/*', (req, res) => {
        registry.getAvailableHost(microservice)
            .then(host => {
                const requestedUrl = `${host}/${microservice}${req.url}`;
                request.get(requestedUrl).pipe(res);
            })
            .catch(err => {
                let errorMessage = err;
                if (errorMessage.error === undefined) {
                    errorMessage = { error: 500, message: 'An error occurred on our side' };
                }
                res.status(errorMessage.error).send(errorMessage);
            });
    });

    router.post('/*', (req, res) => {
        registry.getAvailableHost(microservice)
            .then(host => {
                const requestedUrl = `${host}/${microservice}${req.url}`;
                const sendData = {};

                let contentType = req.headers['content-type'];
                if (contentType === undefined) {
                    res.status(406).send({ error: 406, message: 'Content type is not supported' });
                    return;  
                }
                contentType = contentType.split(';')[0];

                switch(contentType) {
                    case 'multipart/form-data':
                        sendData.formData = req.body;
                    break;
                    case 'application/json':
                        sendData.json = req.body;
                    break;
                    case 'application/x-www-form-urlencoded':
                        sendData.form = req.body;
                    break;
                }
                request.post(requestedUrl, sendData)
                    .pipe(res);
            })
            .catch(err => {
                let errorMessage = err;
                if (errorMessage.error === undefined) {
                    errorMessage = { error: 500, message: 'An error occurred on our side' };
                }
                res.status(errorMessage.error).send(errorMessage);
            });
    });
    
    router.put('/*', (req, res) => {
        registry.getAvailableHost(microservice)
            .then(host => {
                const requestedUrl = `${host}/${microservice}${req.url}`;
                const sendData = {};

                let contentType = req.headers['content-type'];
                if (contentType === undefined) {
                    res.status(406).send({ error: 406, message: 'Content type is not supported' });
                    return;  
                }
                contentType = contentType.split(';')[0];

                switch(contentType) {
                    case 'multipart/form-data':
                        sendData.formData = req.body;
                    break;
                    case 'application/json':
                        sendData.json = req.body;
                    break;
                    case 'application/x-www-form-urlencoded':
                        sendData.form = req.body;
                    break;
                }
                request.put(requestedUrl, sendData)
                    .pipe(res);
            })
            .catch(err => {
                let errorMessage = err;
                if (errorMessage.error === undefined) {
                    errorMessage = { error: 500, message: 'An error occurred on our side' };
                }
                res.status(errorMessage.error).send(errorMessage);
            });
    });

    router.delete('/*', (req, res) => {
        registry.getAvailableHost(microservice)
            .then(host => {
                const requestedUrl = `${host}/${microservice}${req.url}`;
                request.delete(requestedUrl).pipe(res);
            })
            .catch(err => {
                let errorMessage = err;
                if (errorMessage.error === undefined) {
                    errorMessage = { error: 500, message: 'An error occurred on our side' };
                }
                res.status(errorMessage.error).send(errorMessage);
            });
    });

    return router;
}


