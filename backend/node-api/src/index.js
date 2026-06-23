const express = require('express');
const app = express();
const userRouter = require('./routes/userRoutes');
const dotenv = require('dotenv');
dotenv.config({ path: process.env.DOTENV_PATH || '../../.env'});

PORT = process.env.NODE_API_PORT;

app.use(express.json());

app.use('/api/users', userRouter);

app.listen(PORT, () => {
    console.log("Server Started")
});