const express = require('express');
const router = express.Router();
const userController = require('../controllers/userController');

// Directing traffic to the controller function
router.get('/:id', userController.getUser);

module.exports = router;
