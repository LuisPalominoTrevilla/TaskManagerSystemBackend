const CronJob = require('cron').CronJob;
const dbReplicationJob = require('./dbReplicationJobJob');

module.exports = {
    startJobs() {
        new CronJob('*/5 * * * *', dbReplicationJob, null, true, 'America/Mexico_City');
    }
}