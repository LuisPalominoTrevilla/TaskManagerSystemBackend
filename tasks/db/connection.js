const mysql = require('mysql');

const connectionPool = mysql.createPool({
    connectionLimit : 10,
    host     : process.env.DB_HOST,
    user     : process.env.DB_USER,
    password : process.env.DB_PASSWORD,
    database : process.env.DB_NAME,
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
    },
    setTimeZone(tz) {
        connectionPool.query(`SET time_zone = '${tz}';`);
    }
};