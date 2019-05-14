const requester = require('../requester');
const serviceFactory = require('../../db/servicesFactory');

const async = require('async');

module.exports = function(service) {
    const model = serviceFactory(service);
    return new Promise(resolve => {
        requester.getAccounts()
            .then(retrievedAccounts => {
                async.eachSeries(retrievedAccounts, (account, next) => {
                    console.log('Checking if email', account.email, 'exists');
                    model.exists({ key: model.idField, id: account.email })
                        .then(accountExists => {
                            if (!accountExists) {
                                model.insertOne(account)
                                    .then(() => next())
                                    .catch(() => {
                                        console.log('Error while inserting', account.email);
                                        
                                    });
                            } else {
                                next();
                            }
                        })
                        .catch(() => {
                            next();
                        });
                }, () => resolve());
            })
            .catch(() => {
                resolve();
            });
    });
}