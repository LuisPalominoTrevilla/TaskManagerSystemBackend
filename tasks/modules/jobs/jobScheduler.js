const CronJob = require('cron').CronJob;
const taskReminder = require('./taskReminderJob');

module.exports = {
    startJobs() {
        new CronJob('*/5 * * * *', taskReminder, null, true, 'America/Mexico_City');
    }
}