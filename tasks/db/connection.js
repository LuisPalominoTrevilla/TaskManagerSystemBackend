const mysql = require('mysql');

const connectionPool = mysql.createPool({
    connectionLimit : 10,
    host     : 'generalpurpose.crcsqwhwrqwu.us-east-2.rds.amazonaws.com',
    user     : 'admin',
    password : 'arquitecturasoftware',
    database : 'arquitectura',
    timezone : 'local'
});

module.exports = {
    connection: connectionPool,
    testConnection() {
        return new Promise((resolve, reject) => {
            connectionPool.getConnection((err, connection) => {
                if (err) reject(err);
                else {
                    connection.release();
                    resolve();
                }
            });
        });
    }
};