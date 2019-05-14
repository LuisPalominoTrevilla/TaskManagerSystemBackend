const Model = require('./model');

const habitsModel = class Accounts extends Model {
    constructor() {
        super('habits');
        this.idField = 'id';
    }
}

module.exports = habitsModel;