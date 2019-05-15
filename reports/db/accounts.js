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
            db.query(sql, (err, tasks) => {
                if (err) return reject({ error: 500, message: err.message });
                resolve(tasks);
            });
        });
    }

    getDelayedTasks(userId) {
        return new Promise((resolve, reject) => {
            const today = moment().format('YYYY-MM-DD HH:mm:ss');
            let sql = `select * from tasks where dueDate < ${mysql.escape(today)} AND userId = ${mysql.escape(userId)} AND completed = 0`;
            db.query(sql, (err, tasks) => {
                if (err) return reject({ error: 500, message: err.message });
                resolve(tasks);
            });
        });
    }

    getGoodHabits(userId) {
        return new Promise((resolve, reject) => {
            let sql = `select * from habits where score > 50 AND userId = ${mysql.escape(userId)}`;
            db.query(sql, (err, habits) => {
                if (err) return reject({ error: 500, message: err.message });
                resolve(habits);
            });
        });
    }

    getBadHabits(userId) {
        return new Promise((resolve, reject) => {
            let sql = `select * from habits where score < 0 AND userId = ${mysql.escape(userId)}`;
            db.query(sql, (err, habits) => {
                if (err) return reject({ error: 500, message: err.message });
                resolve(habits);
            });
        });
    }
}

module.exports = accountsModel;