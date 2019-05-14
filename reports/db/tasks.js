const Model = require('./model');

const tasksModal = class Tasks extends Model {
    constructor() {
        super('tasks');
        this.idField = 'taskId';
    }
}

module.exports = tasksModal;