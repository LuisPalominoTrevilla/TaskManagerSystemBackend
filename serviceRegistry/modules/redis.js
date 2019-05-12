const redis = require('redis');
const client = redis.createClient();

module.exports = {
    addServiceToRegistry(serviceName, host, port) {
        const server = `http://${host}:${port}`;
        return new Promise((resolve, reject) => {
            client.SADD(serviceName, server, (err, res) => {
                if (err) reject({ error: 500, message: 'There was an error while registering the service' });
                else resolve();
            });
        });
    },

    isServiceRegistered(serviceName, host, port) {
        const server = `http://${host}:${port}`;
        return new Promise((resolve, reject) => {
            client.SISMEMBER(serviceName, server, (err, res) => {
                if (err) return reject({ error: 500, message: err.message });
                resolve(res === 1);
            });
        })
    }
};