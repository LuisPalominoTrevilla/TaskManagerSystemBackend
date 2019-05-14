const request = require('request-promise');

module.exports = {
    getAccounts() {
        return new Promise((resolve, reject) => {
            request.get(`${process.env.REGISTRY_HOST}/api/v1/accounts`)
                .then(resp => {
                    resolve(JSON.parse(resp));
                })
                .catch(() => {
                    reject();
                })
        });
    }
}