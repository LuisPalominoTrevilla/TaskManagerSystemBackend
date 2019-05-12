const redis = require('../modules/redis');

module.exports = {
    registerService(service, ip, port) {
        return new Promise((resolve, reject) => {
            redis.isServiceRegistered(service, ip, port)
                .then(alreadyExists => {
                    if (alreadyExists) {
                        throw { error: 409, message: 'Service is already registered' };
                    }
                    return redis.addServiceToRegistry(service, ip, port);
                })
                .then(resolve)
                .catch(reject);
        });
    }
}