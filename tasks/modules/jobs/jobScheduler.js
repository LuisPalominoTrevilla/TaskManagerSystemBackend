const CronJob = require('cron').CronJob;
const taskReminder = require('./taskReminderJob');

module.exports = {
    startJobs() {
        new CronJob('*/6 * * * *', taskReminder, null, true, 'America/Mexico_City');
    }
}