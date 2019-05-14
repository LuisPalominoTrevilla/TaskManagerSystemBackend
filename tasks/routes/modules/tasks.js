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

router.get('/', (req, res) => {
    const validQueries = new Set(['title','description','reminder','userId','completed']);
    const filter = {};
    for (const field in req.query) {
        if (validQueries.has(field)) {
            let value = req.query[field];
            filter[field] = value;
        }
    }
    taskDB.findMany(filter)
        .then(result => {
            res.status(200).json(result);
        })
        .catch(err => {
            res.status(404).send(err);
        });
});

router.get('/:taskId', (req, res) => {
    const taskId = req.params.taskId;
    taskDB.findById(taskId)
        .then(result => {
            const taskFound = {
                ...result,
                id: result.taskId,
                reminder: result.reminder !== 0,
                completed: result.completed !== 0
            };
            delete taskFound.taskId;
            res.status(200).send(taskFound);
        })
        .catch(err => {
            res.status(404).send(err);
        });
});

router.post('/', (req, res) => {
    const newTask = {};
    const { title, description, dueDate, reminderDate, userId } = req.body;

    const completeParams = title !== undefined && description !== undefined &&
        dueDate !== undefined && userId !== undefined && !strings.isEmpty(description) &&
        !strings.isEmpty(title) && !strings.isEmpty(userId) && req.files.image !== undefined;
    const validDates = moment(dueDate).format() !== 'Invalid date' &&
        (reminderDate === undefined || moment(reminderDate).format() != 'Invalid date');
    const correctDueDate = moment(dueDate).isAfter(moment());
    const correctReminderDate = reminderDate === undefined || moment(reminderDate).isBefore(dueDate) && moment(reminderDate).isAfter(moment());

    if (!completeParams) {
        return res.status(400).json({ error: 400, message: 'Required parameters missing' });
    }
    if (!validDates) {
        return res.status(400).json({ error: 400, message: 'Date field is not valid' });
    }
    if (!correctReminderDate) {
        return res.status(400).json({ error: 400, message: 'reminderDate cannot be after the due date or before the current time' });
    }
    if (!correctDueDate) {
        return res.status(400).json({ error: 400, message: 'dueDate cannot be before the current time' });
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
                newTask.completed = false;
                return taskDB.insertOne(newTask);
            })
            .then((queryResults) => {
                newTask.id = queryResults.insertId;
                res.status(200).send(newTask);
            })
            .catch(err => {
                console.log(err);
                res.status(500).send({ error: 500, message: 'An error occurred while trying to create the task' });
            });
    });
});

router.post('/:taskId/complete', (req, res) => {
    const taskId = req.params.taskId;
    let existentTask;
    taskDB.findById(taskId)
        .then(result => {
            existentTask = result;
            if (result.completed === 1) {
                throw { error: 409, message: 'Task has already been completed' };
            }
            return taskDB.update({ completed: 1}, { taskId });
        })
        .then(() => {
            existentTask.completed = true;
            existentTask.reminder = existentTask.reminder !== 0;
            existentTask.id = existentTask.taskId;
            delete existentTask.taskId;
            res.status(200).send(existentTask);
        })
        .catch(err => {
            res.status(err.error).send(err);
        });
});

router.put('/:taskId', (req, res) => {
    const taskId = req.params.taskId;
    let existentTask;
    const fieldsToUpdate = {};

    const { title, description, dueDate, reminder, reminderDate, completed } = req.body;

    const validDates = (dueDate !== undefined || moment(dueDate).format() !== 'Invalid date') &&
        (reminderDate === undefined || moment(reminderDate).format() != 'Invalid date');
    const correctDueDate = dueDate === undefined || moment(dueDate).isAfter(moment());
    const correctReminderDate = reminderDate === undefined ||
        (dueDate === undefined || moment(reminderDate).isBefore(dueDate)) &&
        moment(reminderDate).isAfter(moment());

    if (!validDates) {
        return res.status(400).json({ error: 400, message: 'Date field is not valid' });
    }
    if (!correctReminderDate) {
        return res.status(400).json({ error: 400, message: 'reminderDate cannot be after the due date or before the current time' });
    }
    if (!correctDueDate) {
        return res.status(400).json({ error: 400, message: 'dueDate cannot be before the current time' });
    }

    if (title !== undefined && !strings.isEmpty(title)) {
        fieldsToUpdate.title = title;
    }
    if (description !== undefined && !strings.isEmpty(description)) {
        fieldsToUpdate.description = description;
    }
    if (reminder !== undefined && reminder === 'false') {
        fieldsToUpdate.reminder = 0;
        fieldsToUpdate.reminderDate = null;
    } else if (reminderDate !== undefined) {
        fieldsToUpdate.reminder = 1;
        fieldsToUpdate.reminderDate = moment(reminderDate).format('YYYY-MM-DD HH:mm:ss');
    }
    if (dueDate !== undefined) {
        fieldsToUpdate.dueDate = moment(dueDate).format('YYYY-MM-DD HH:mm:ss');
    }
    if (completed !== undefined && (completed === 'true' || completed === 'false')) {
        fieldsToUpdate.completed = (completed === 'true') ? 1 : 0;
    }

    taskDB.findById(taskId)
        .then(result => {
            existentTask = result;
            return taskDB.update(fieldsToUpdate, { taskId });
        })
        .then(() => {
            const updatedtask = {
                ...existentTask,
                ...fieldsToUpdate
            };
            updatedtask.reminder = updatedtask.reminder !== 0;
            updatedtask.completed = updatedtask.completed !== 0;
            updatedtask.id = updatedtask.taskId;
            delete updatedtask.taskId;
            res.status(200).send(updatedtask);
        })
        .catch(err => {
            res.status(err.error).send(err);
        });
});

router.delete('/:taskId', (req, res) => {
    const taskId = req.params.taskId;
    taskDB.findById(taskId)
        .then(result => {
            const imageFile = result.imageUrl.split('/').slice(-1)[0];
            S3Client.deleteFile(imageFile);
            taskDB.deleteOne(taskId)
                .then(() => {
                    res.status(200).send('OK');
                })
                .catch(err => {
                    console.log(err);
                    res.status(500).send({ error: 500, message: 'An error occurred while trying to delete the task' });
                });
        })
        .catch(err => {
            return res.status(404).send(err);
        });
});

module.exports = router;