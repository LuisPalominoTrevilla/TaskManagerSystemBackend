const { connection: db } = require('../db/connection');

const task = {
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
    }
}

module.exports = task;