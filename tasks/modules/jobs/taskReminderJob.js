const taskDb = require('../../models/task');


module.exports = function() {
    taskDb.getReminderTasks(3)
        .then(result => {
            console.log(result);
        })
        .catch(err => {
            console.log(err);
        });
}
