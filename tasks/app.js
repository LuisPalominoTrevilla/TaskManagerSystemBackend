const express = require('express');
const helmet = require('helmet');
const cookieParser = require('cookie-parser');
const logger = require('morgan');

const endpoints = require('./routes');

const app = express();

app.use(logger('dev'));
app.use(express.json());
app.use(express.urlencoded({ extended: false }));
app.use(cookieParser());
app.use(helmet());

app.use('/api', endpoints);
app.get('/', (req, res) => {
    res.send("Tasks microservice");
});

module.exports = app;
