const Model = require('./model');

const tasksModal = class Tasks extends Model {
    constructor() {
        super('tasks');
        this.idField = 'id';
    }
}

module.exports = tasksModal;