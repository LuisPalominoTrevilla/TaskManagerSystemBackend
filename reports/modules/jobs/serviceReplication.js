const requester = require('../requester');
const serviceFactory = require('../../db/servicesFactory');

const async = require('async');

module.exports = function(service) {
    console.log('replicatting', service);
    const model = serviceFactory(service);
    return new Promise(resolve => {
        requester.getRecords(service)
            .then(retrievedRecords => {
                async.eachSeries(retrievedRecords, (record, next) => {
                    model.exists({ key: model.idField, id: record[model.idField] })
                        .then(recordExists => {
                            if (!recordExists) {
                                model.insertOne(record)
                                    .then(() => next())
                                    .catch(() => {
                                        next();
                                    });
                            } else {
                                model.updateOne(record)
                                    .then(() => next())
                                    .catch(() => next());
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