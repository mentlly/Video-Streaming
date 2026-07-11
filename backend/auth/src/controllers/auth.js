const db = require('../config/db');
const bcrypt = require('bcryptjs');
const jwt = require('jsonwebtoken');

const JWT_SECRET = process.env.JWT_SECRET

exports.register = async (req, res) => {
    const { email, password } = req.body;

    if (!email) {
        return res.status(400).json({ error: 'Email is required' });
    } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
        return res.status(400).json({ error: 'Enter a valid email'})
    }

    if (!password) {
        return res.status(400).json({ error: 'Password is required '});
    } else if (password.length < 8) {
        return res.status(400).json({ error: 'Password should be greater than or equal to 8 characters'});
    } else if (password.length > 16) {
        return res.status(400).json({ error: 'Password should be less than or equal to 16 characters'});
    }

    try {
        // Check if user already exists
        const userCheck = await db.query('SELECT * FROM users WHERE email = $1', [email]);
        if (userCheck.rows.length > 0) {
            return res.status(409).json({ error: 'Email is already registered' });
        }

        // Hash the password securely
        const saltRounds = 10;
        const passwordHash = await bcrypt.hash(password, saltRounds);

        // Insert user into database
        await db.query(
            'INSERT INTO users (email, password_hash) VALUES ($1, $2)',
            [email, passwordHash]
        );

        return res.status(201).json({ message: 'User registered successfully' });
    } catch (error) {
        console.error('Registration Error:', error);
        return res.status(500).json({ error: 'Internal server error during registration' });
    }
};

exports.login = async (req, res) => {
    const { email, password } = req.body;

    if (!email) {
        return res.status(400).json({ error: 'Email is required' });
    } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
        return res.status(400).json({ error: 'Enter a valid email'})
    }

    if (!password) {
        return res.status(400).json({ error: 'Password is required '});
    } else if (password.length < 8) {
        return res.status(400).json({ error: 'Password should be greater than or equal to 8 characters'});
    } else if (password.length > 16) {
        return res.status(400).json({ error: 'Password should be less than or equal to 16 characters'});
    }

    try {
        // Find user by email
        const result = await db.query('SELECT * FROM users WHERE email = $1', [email]);
        const user = result.rows[0];

        if (!user) {
            return res.status(401).json({ error: 'Invalid email or password' });
        }

        // Compare incoming password with stored hash
        const isMatch = await bcrypt.compare(password, user.password_hash);
        if (!isMatch) {
            return res.status(401).json({ error: 'Invalid email or password' });
        }

        // Generate JWT token (expires in 1 hours)
        const token = jwt.sign(
            { userId: user.account_id },
            JWT_SECRET,
            { expiresIn: '1h' }
        );

        res.cookie('jwt_token', token, 
            { 
                httpOnly: true,
                secure: false,
                sameSite: 'lax', 
                maxAge: 24 * 60 * 60 * 1000
            }
        )

        // Return token to client
        return res.status(200).json({
            message: 'Login successful',
        });
    } catch (error) {
        console.error('Login Error:', error);
        return res.status(500).json({ error: 'Internal server error during login' });
    }
};