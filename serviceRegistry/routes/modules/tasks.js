const express = require('express');
const moment = require('moment-timezone');
const router = express.Router();
const request = require('request');

moment.tz.setDefault('America/Mexico_City');

router.post('/*', (req, res) => {
    let contentType = req.headers['content-type'];
    if (contentType === undefined) {

        return;
    }
    contentType = contentType.split(';')[0];
    if (contentType === 'multipart/form-data') {
        request.post(`http://127.0.0.1:4002/tasks${req.url}`, {formData: req.body})
            .pipe(res);
    }
});

router.put('/*', (req, res) => {
    let contentType = req.headers['content-type'];
    if (contentType === undefined) {

        return;
    }
    contentType = contentType.split(';')[0];
    if (contentType === 'application/x-www-form-urlencoded') {
        request.put(`http://127.0.0.1:4002/tasks${req.url}`, {form: req.body})
            .pipe(res);
    } else if (contentType === 'multipart/form-data') {

    } else {
        res.status(406).send({ error: 406, message: 'Content type is not supported' });
    }
    console.log(contentType);
});

module.exports = router;