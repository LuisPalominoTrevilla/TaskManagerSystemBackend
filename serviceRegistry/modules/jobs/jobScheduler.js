const CronJob = require('cron').CronJob;
const serviceChecker = require('./serviceCheckerJob');

module.exports = {
    startJobs() {
        new CronJob('* * * * *', serviceChecker, null, true, 'America/Mexico_City');
    }
}