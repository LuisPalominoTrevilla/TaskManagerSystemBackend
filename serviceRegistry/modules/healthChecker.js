const request = require('request');
const redis = require('./redis');

module.exports = function(server) {
    return new Promise((resolve, reject) => {
        redis.getHealthCheck(server)
        .then(endpoint => {
            request.get(`${server}${endpoint}`, (err, response) => {
                if (err) reject({ error: 502, message: 'healthCheck endpoint did not respond' });
                else resolve();
            });
        })
        .catch(reject);
    });
}