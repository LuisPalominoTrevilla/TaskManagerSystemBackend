const CronJob = require('cron').CronJob;

module.exports = {
    startJobs() {
        new CronJob('*/4 * * * *', () => {}, null, true, 'America/Mexico_City');
    }
}