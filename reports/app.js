require('dotenv').config();

const express = require('express');
const helmet = require('helmet');
const cors = require('cors');
const cookieParser = require('cookie-parser');
const logger = require('morgan');
const jobScheduler = require('./modules/jobs/jobScheduler');
const db = require('./db/connection');
const request = require('request');

const endpoints = require('./routes');

jobScheduler.startJobs();

db.testConnection()
    .then(() => console.log('Succesfully connected to database'))
    .catch(err => console.log(err));
db.setTimeZone('America/Mexico_City');

const app = express();

app.use(cors());

request.post(`${process.env.REGISTRY_HOST}${process.env.REGISTRY_ENDPOINT}`, {form: { port: 4003, service: 'reports', healthCheck: '/healthecheck' }}, (err, res) => {
    if (err) console.log('Service registry not available');
});

app.use(logger('dev'));
app.use(express.json());
app.use(express.urlencoded({ extended: false }));
app.use(cookieParser());
app.use(helmet());

app.get('/healthecheck', (req, res) => {
    res.status(200).json({ status: 'OK' });
});
app.use('/', endpoints);

module.exports = app;
