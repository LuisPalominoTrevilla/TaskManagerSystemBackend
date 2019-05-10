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

app.get('/healthecheck', (req, res) => {
    res.status(200).json({ status: 'OK' });
});
app.use('/', endpoints);

module.exports = app;
