const express = require('express');
const moment = require('moment-timezone');
const strings = require('../../utils/strings');
const router = express.Router();

moment.tz.setDefault('America/Mexico_City');

router.get('/', (req, res, next) => {
    res.json({
        message: 'Hello from tasks'
    });
});

router.post('/', (req, res, next) => {
    const newTask = {};
    const { title, description, dueDate, reminderDate, userId } = req.body;

    console.log(req.files.image);

    const completeParams = title !== undefined && description !== undefined &&
        dueDate !== undefined && userId !== undefined && !strings.isEmpty(description) &&
        !strings.isEmpty(title) && !strings.isEmpty(userId);
    const validDates = moment(dueDate).format() !== 'Invalid date' && moment(reminderDate).format() != 'Invalid date';

    if (!completeParams) {
        return res.status(400).json({ error: 400, message: 'Required parameters missing' });
    }
    if (!validDates) {
        return res.status(400).json({ error: 400, message: 'Date field is not valid' });
    }

    newTask.title = title;
    newTask.description = description;
    newTask.dueDate = moment(dueDate).toDate();
    if (reminderDate !== undefined) {
        newTask.reminder = true;
        newTask.reminderDate = moment(reminderDate).toDate();
    } else {
        newTask.reminder = false;
    }
    newTask.userId = userId;

    console.log(JSON.stringify(newTask, null, 2));
    res.send('hola');
});

module.exports = router;