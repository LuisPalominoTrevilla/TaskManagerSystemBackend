const Accounts = require('./accounts');
const Tasks = require('./tasks');

module.exports = function(service) {
    switch(service) {
        case "accounts":
            return new Accounts();
        case "tasks":
            return new Tasks();
    };
}