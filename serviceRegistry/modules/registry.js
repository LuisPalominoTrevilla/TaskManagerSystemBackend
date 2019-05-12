const redis = require('../modules/redis');
const healthChecker = require('./healthChecker');

module.exports = {
    registerService(service, ip, port, healthCheck) {
        const server = `http://${ip}:${port}`;
        return new Promise((resolve, reject) => {
            redis.isServiceRegistered(service, server)
                .then(alreadyExists => {
                    if (alreadyExists) {
                        throw { error: 409, message: 'Service is already registered' };
                    }
                    return redis.addServiceToRegistry(service, server);
                })
                .then(() => {
                    return redis.addHealthCheck(server, healthCheck);
                })
                .then(() => {
                    healthChecker(server)
                        .then(resolve)
                        .catch(err => {
                            this.unregisterService(service, server)
                                .finally(() => {
                                    reject(err);
                                });
                        });
                })
                .catch(reject);
        });
    },

    unregisterService(service, server) {
        return new Promise((resolve, reject) => {
            redis.removeService(service, server)
                .then(() => {
                    return redis.removeHealthCheck(server);
                })
                .then(resolve)
                .catch(reject);
        });
    }
}