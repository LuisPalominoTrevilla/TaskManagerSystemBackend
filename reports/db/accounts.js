const Model = require('./model');

const accountsModel = class Accounts extends Model {
    constructor() {
        super('accounts');
        this.idField = 'email';
    }
}

module.exports = accountsModel;