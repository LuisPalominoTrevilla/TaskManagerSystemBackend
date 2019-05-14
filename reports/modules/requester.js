const request = require('request-promise');

module.exports = {
    getRecords(service) {
        return new Promise((resolve, reject) => {
            request.get(`${process.env.REGISTRY_HOST}/api/v1/${service}`)
                .then(resp => {
                    resolve(JSON.parse(resp));
                })
                .catch(() => {
                    reject();
                })
        });
    }
}