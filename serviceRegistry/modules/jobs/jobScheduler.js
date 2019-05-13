const CronJob = require('cron').CronJob;
const serviceChecker = require('./serviceCheckerJob');

module.exports = {
    startJobs() {
        new CronJob('*/5 * * * *', serviceChecker, null, true, 'America/Mexico_City');
    }
}