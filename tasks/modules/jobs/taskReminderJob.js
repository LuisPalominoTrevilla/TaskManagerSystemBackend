const taskDb = require('../../models/task');
const mailer = require('../../modules/mail');
const moment = require('moment');


module.exports = function() {
    taskDb.getReminderTasks(2)
        .then(result => {
            result.forEach(task => {
                mailer.sendTaskReminder(task.userId, task.title, task.description, task.imageUrl, moment(task.dueDate).format('LLL'), moment(task.reminderDate).calendar());
            });
        })
        .catch(err => {
            console.log(err);
        });
}
