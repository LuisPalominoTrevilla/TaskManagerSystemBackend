const { connection: db } = require('../db/connection');
const mysql = require('mysql');
const moment = require('moment-timezone');
moment.tz.setDefault('America/Mexico_City');

const task = {
    findById(taskId) {
        let sql = `SELECT * FROM tasks WHERE taskId=${mysql.escape(taskId)}`;
        return new Promise((resolve, reject) => {
            db.query(sql, (err, result) => {
                if (err) return reject({ error: 500, message: err.message });
                if (result.length === 0) reject({ error: 404, message: 'There is no task associated with the taskId provided' });
                else resolve(result[0]);
            })
        });
    },

    findMany(filter) {
        let sql = 'SELECT * FROM tasks WHERE ';
        for (const field in filter) {
            sql += `${field} = ${mysql.escape(filter[field])}, `;
        }
        sql = sql.slice(0, -2);
        return new Promise((resolve, reject) => {
            db.query(sql, (err, results) => {
                if (err) return reject({ error: 500, message: err.message });
                resolve(results)
            })
        });
    },

    insertOne(newTask) {
        const sql = 'INSERT INTO tasks (title, description, dueDate, reminder, reminderDate, imageUrl, userId) VALUES ?';
        const values = [
            [newTask.title, newTask.description, newTask.dueDate, newTask.reminder, newTask.reminderDate, newTask.imageUrl, newTask.userId],
        ]
        return new Promise((resolve, reject) => {
            db.query(sql, [values], (err, result) => {
                if (err) return reject(err);
                resolve(result);
            });
        });
    },

    update(fields, where) {
        let sql = 'UPDATE tasks SET ';
        for (const field in fields) {
            sql += `${field} = ${mysql.escape(fields[field])}, `;
        }
        sql = `${sql.slice(0, -2)} WHERE `;
        for (const field in where) {
            sql += `${field} = ${mysql.escape(where[field])}, `;
        }
        sql = sql.slice(0, -2);
        return new Promise((resolve, reject) => {
            db.query(sql, (err, result) => {
                if (err) return reject({ error: 500, message: err.message });
                resolve(result);
            });
        });
    },

    deleteOne(taskId) {
        const sql = `DELETE FROM tasks WHERE taskId=${mysql.escape(taskId)}`;
        return new Promise((resolve, reject) => {
            db.query(sql, (err, result) => {
                if (err) return reject(err);
                resolve(result);
            })
        });
    },

    getReminderTasks(minutesAfter = 0) {
        const datetimeBefore = moment().add(minutesAfter, 'minute').format('YYYY-MM-DD HH:mm:ss');
        const datetimeAfter = moment().subtract(minutesAfter, 'minute').format('YYYY-MM-DD HH:mm:ss');
        const sql = `SELECT * FROM tasks WHERE reminder = 1 AND reminderDate < ${mysql.escape(datetimeBefore)} AND
            reminderDate >= ${mysql.escape(datetimeAfter)} AND
            completed = 0 ORDER BY reminderDate ASC`;
        return new Promise((resolve, reject) => {
            db.query(sql, (err, result) => {
                if (err) return reject(err);
                resolve(result);
            })
        });
    }
}

module.exports = task;