const express = require('express');
const moment = require('moment-timezone');
const strings = require('../../modules/strings');
const path = require('path');
const fs = require('fs');
const fileType = require('file-type');
const S3Client = require('../../modules/aws-s3');
const router = express.Router();
const taskDB = require('../../models/task');

moment.tz.setDefault('America/Mexico_City');

router.get('/', (req, res, next) => {
    res.json({
        message: 'Hello from tasks'
    });
});

router.post('/', (req, res, next) => {
    const newTask = {};
    const { title, description, dueDate, reminderDate, userId } = req.body;

    const completeParams = title !== undefined && description !== undefined &&
        dueDate !== undefined && userId !== undefined && !strings.isEmpty(description) &&
        !strings.isEmpty(title) && !strings.isEmpty(userId);
    const validDates = moment(dueDate).format() !== 'Invalid date' &&
        (reminderDate === undefined || moment(reminderDate).format() != 'Invalid date');

    if (!completeParams) {
        return res.status(400).json({ error: 400, message: 'Required parameters missing' });
    }
    if (!validDates) {
        return res.status(400).json({ error: 400, message: 'Date field is not valid' });
    }

    newTask.title = title;
    newTask.description = description;
    newTask.dueDate = moment(dueDate).format('YYYY-MM-DD HH:mm:ss');
    if (reminderDate !== undefined) {
        newTask.reminder = true;
        newTask.reminderDate = moment(reminderDate).format('YYYY-MM-DD HH:mm:ss');
    } else {
        newTask.reminder = false;
    }
    newTask.userId = userId;

    fs.readFile(req.files.image.path, (err, file) => {
        if (err) return res.status(500).json({ error: 500, message: 'Image could not be read' });

        const { ext: type , mime } = fileType(file);

        if (type !== 'jpg' && type !== 'png' && type !== 'gif') {
            return res.status(400).json({ error: 400, message: 'File provided is not a valid image, valid types are jpeg, png or gif' });
        }

        const imageName = path.basename(req.files.image.path);

        S3Client.uploadFile(imageName, file, type, {'ContentType': mime})
            .then(url => {
                newTask.imageUrl = url;
                return taskDB.insertOne(newTask);
            })
            .then(() => {
                res.status(200).send(newTask);
            })
            .catch(err => {
                console.log(err);
                res.status(500).send({ error: 500, message: 'An error occurred while trying to create the task' });
            });
    });
});

module.exports = router;