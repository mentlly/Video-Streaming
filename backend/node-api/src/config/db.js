const { Pool } = require('pg');

const pool = new Pool({
  host: 'postgres',
  user: process.env.DB_USER,
  password: process.env.DB_PASSWORD,
  database: process.env.DB_NAME,
  port: process.env.DB_PORT,
});

pool.query('SELECT NOW()', (err, res) => {
  if (err) {
    console.error('Database connection failed:', err.stack);
  } else {
    console.log('Connected to PostgreSQL successfully at:', res.rows[0].now);
  }
});

const initDb = async () => {
    const queryText = `
        CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        email VARCHAR(255) UNIQUE NOT NULL,
        password_hash TEXT NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        );
    `;
    try {
        await pool.query(queryText);
        console.log("Database user table initialized.")
    } catch(err) {
        console.error("Error in creating user table", err)
    }
};
initDb();

module.exports = pool;