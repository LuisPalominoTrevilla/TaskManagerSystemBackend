require('dotenv').config();

const express = require('express');
const helmet = require('helmet');
const cors = require('cors');
const cookieParser = require('cookie-parser');
const logger = require('morgan');
const formData = require('express-form-data');
const jobScheduler = require('./modules/jobs/jobScheduler');
const os = require('os');

const endpoints = require('./routes');
const register = require('./routes/register');

jobScheduler.startJobs();

const app = express();

app.use(cors());

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

app.use('/api/v1/', endpoints);
app.use('/register', register);

module.exports = app;
