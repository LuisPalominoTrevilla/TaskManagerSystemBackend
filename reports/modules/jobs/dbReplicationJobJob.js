const accountsReplication = require('./accountsReplication');

module.exports = function() {
    accountsReplication('accounts');
    // accounts.insertOne({email: 'luispalominotQ@hotmail.com', name: 'luis', password: 'luis123'});
}