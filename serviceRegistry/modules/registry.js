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
                .then(resolve)
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
    },

    getAvailableHost(service) {
        return new Promise((resolve, reject) => {
            redis.getRandomService(service)
                .then(host => {
                    if (host === null) {
                        throw { error: 502, message: 'There are no available services for the moment' };
                    }
                    healthChecker(host)
                        .then(() => {
                            resolve(host);
                        })
                        .catch(() => {
                            this.unregisterService(service, host)
                                .finally(() => {
                                    return this.getAvailableHost(service);
                                })
                                .then(resolve)
                                .catch(reject);
                        });
                })
                .catch(reject);
        });
    }
}