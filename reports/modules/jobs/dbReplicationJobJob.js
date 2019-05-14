const serviceReplication = require('./serviceReplication');

module.exports = function() {
    serviceReplication('accounts')
        .then(() => {
            return serviceReplication('tasks');
        })
        .then(() => {
            serviceReplication('habits');
        });

}