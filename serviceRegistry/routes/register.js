const express = require('express');
const router = express.Router();

router.post('/', (req, res) => {
    const ip = req.connection.remoteAddress;

    console.log(ip);
});

module.exports = router;