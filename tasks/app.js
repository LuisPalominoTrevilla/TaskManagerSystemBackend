require('dotenv').config();

const express = require('express');
const helmet = require('helmet');
const cors = require('cors');
const cookieParser = require('cookie-parser');
const logger = require('morgan');
const formData = require('express-form-data');
const jobScheduler = require('./modules/jobs/jobScheduler');
const os = require('os');
const request = require('request');
const db = require('./db/connection');

const endpoints = require('./routes');

jobScheduler.startJobs();

db.testConnection()
    .then(() => console.log('Succesfully connected to database'))
    .catch(err => console.log(err));
db.setTimeZone('America/Mexico_City');

const app = express();

app.use(cors());

request.post(`${process.env.REGISTRY_HOST}${process.env.REGISTRY_ENDPOINT}`, {form: { port: 4002, service: 'tasks', healthCheck: '/healthecheck' }}, (err, res) => {
    if (err) console.log('Service registry not available');
});

const options = {
    uploadDir: os.tmpdir(),
    autoClean: true
};  

app.use(formData.parse(options));
app.use(formData.format());
app.use(formData.stream());
app.use(formData.union());


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
