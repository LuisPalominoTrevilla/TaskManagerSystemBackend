const redis = require('redis');
const client = redis.createClient();

module.exports = {
    addServiceToRegistry(serviceName, host, port) {
        const server = `http://${host}:${port}`;
        return new Promise((resolve, reject) => {
            client.SADD(serviceName, server, (err, res) => {
                if (err) reject(err);
                else resolve();
            });
        });
    },

    isServiceRegistered(serviceName, host, port) {
        const server = `http://${host}:${port}`;
        return new Promise((resolve, reject) => {
            client.SISMEMBER(serviceName, server, (err, res) => {
                if (err) return reject(err);
                resolve(res === 1);
            });
        })
    }
};