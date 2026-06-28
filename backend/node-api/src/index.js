const express = require('express');
const app = express();
const dotenv = require('dotenv');
dotenv.config({ path: process.env.DOTENV_PATH || '../../.env'});

const userRouter = require('./routes/userRoutes');
const authRouter = require('./routes/auth');

PORT = process.env.NODE_API_PORT;

app.use(express.json());

app.use('/api/users', userRouter);
app.use('/api/auth', authRouter);

app.listen(PORT, () => {
    console.log("Server Started")
});