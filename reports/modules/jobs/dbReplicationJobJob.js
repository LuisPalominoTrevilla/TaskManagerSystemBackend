const serviceReplication = require('./serviceReplication');

module.exports = function() {
    serviceReplication('accounts')
        .then(() => {
            serviceReplication('tasks')
        });
    // accounts.insertOne({email: 'luispalominotQ@hotmail.com', name: 'luis', password: 'luis123'});
}