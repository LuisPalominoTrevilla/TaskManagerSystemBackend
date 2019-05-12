const redis = require('../../modules/redis');
const healthChecker = require('../../modules/healthChecker');
const registry = require('../../modules/registry');
const Constants = require('../../constants.json');
const async = require('async');

module.exports = function() {
    const serviceNames = Constants.names;
    const services = [];
    for (const i in serviceNames) {
        services.push(serviceNames[i]);
    }
    async.each(services, (service, next) => {
        redis.getServicesFromRegistry(service)
            .then(members => {
                async.each(members, (member, nextMember) => {
                    console.log('Checking health from member ', member);
                    healthChecker(member)
                        .then(nextMember)
                        .catch(() => {
                            console.log('Member is down, unregistering...');
                            registry.unregisterService(service, member)
                                .finally(() => nextMember());
                        });
                }, next);
            })
            .catch(() => next());
    });
}