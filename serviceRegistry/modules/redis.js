const redis = require('redis');
const client = redis.createClient();

module.exports = {
    addServiceToRegistry(serviceName, server) {
        return new Promise((resolve, reject) => {
            client.SADD(serviceName, server, (err, res) => {
                if (err) reject({ error: 500, message: 'There was an error while registering the service' });
                else resolve();
            });
        });
    },

    isServiceRegistered(serviceName, server) {
        return new Promise((resolve, reject) => {
            client.SISMEMBER(serviceName, server, (err, res) => {
                if (err) return reject({ error: 500, message: err.message });
                resolve(res === 1);
            });
        });
    },

    addHealthCheck(server, healthCheck) {
        return new Promise((resolve, reject) => {
            client.HSET('healthCheck', server, healthCheck, (err, res) => {
                if (err) return reject({ error: 500, message: err.message });
                resolve();
            });
        });
    },

    getHealthCheck(server) {
        return new Promise((resolve, reject) => {
            client.HGET('healthCheck', server, (err, healthCheck) => {
                if (err) return reject({ error: 500, message: err.message });
                resolve(healthCheck);
            });
        });
    },

    removeService(serviceName, server) {
        return new Promise((resolve, reject) => {
            client.SREM(serviceName, server, err => {
                if (err) return reject({ error: 500, message: err.message });
                resolve();
            });
        });
    },

    removeHealthCheck(server) {
        return new Promise((resolve, reject) => {
            client.HDEL('healthCheck', server, err => {
                if (err) return reject({ error: 500, message: err.message });
                resolve();
            });
        });
    }
};