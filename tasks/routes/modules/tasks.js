const express = require('express');
const router = express.Router();

router.get('/', (req, res, next) => {
    res.json({
        message: 'Hello from tasks'
    });
});

router.get('/sayHello', (req, res, next) => {
    res.send('Hello there!');
});

module.exports = router;