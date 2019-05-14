const Model = require('./model');
const mysql = require('mysql');
const { connection: db } = require('../db/connection');

const moment = require('moment');

moment.tz.setDefault('America/Mexico_City');

const accountsModel = class Accounts extends Model {
    constructor() {
        super('accounts');
        this.idField = 'email';
    }

    getTasksDueToday(userId) {
        return new Promise((resolve, reject) => {
            const today = moment().format('YYYY-MM-DD HH:mm:ss');
            const end = moment(today).endOf('day').format('YYYY-MM-DD HH:mm:ss');
            let sql = `select * from tasks where dueDate >= ${mysql.escape(today)} AND dueDate < ${mysql.escape(end)} AND userId = ${mysql.escape(userId)}`;
            console.log(sql);
            db.query(sql, (err, tasks) => {
                if (err) return reject({ error: 500, message: err.message });
                resolve(tasks);
            });
        });
    }
}

module.exports = accountsModel;