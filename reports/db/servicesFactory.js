const Accounts = require('./accounts');
const Tasks = require('./tasks');
const Habits = require('./habits');

module.exports = function(service) {
    switch(service) {
        case "accounts":
            return new Accounts();
        case "tasks":
            return new Tasks();
        case "habits":
            return new Habits();
    };
}