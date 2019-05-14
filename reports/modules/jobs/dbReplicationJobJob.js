const accountsReplication = require('./serviceReplication');

module.exports = function() {
    accountsReplication('accounts');
    // accounts.insertOne({email: 'luispalominotQ@hotmail.com', name: 'luis', password: 'luis123'});
}