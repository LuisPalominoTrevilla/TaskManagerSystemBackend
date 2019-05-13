const nodemailer = require('nodemailer');
const mustache = require('mustache');
const fs = require('fs');
const path = require('path')

let transporter = nodemailer.createTransport({
    host: "smtp.gmail.com",
    port: 587,
    secure: false,
    auth: {
        user: process.env.MAIL_USER,
        pass: process.env.MAIL_PASS
    }
});

module.exports = {
    sendMail(templatePath, templateVars, { to, subject }) {
        let template, html;
        try {
            template = fs.readFileSync(templatePath).toString();
            html = mustache.render(template, templateVars);
        } catch (e) {
            console.log("Error:", e);
            return Promise.reject();
        }
        const mailOptions = {
            from: `"T" <${process.env.MAIL_USER}>`,
            to, subject, html
        };

        return transporter.sendMail(mailOptions)
    },

    sendTaskReminder(email, title, description, imageUrl, dueDate, reminderTime) {
        return new Promise((resolve, reject) => {
            const templateVars = { title, description, imageUrl, dueDate, reminderTime };
            const templatePath = path.join(__dirname, "../templates/reminder.mustache");
            this.sendMail(templatePath, templateVars, { to: email, subject: "Friendly Reminder" })
                .then(resolve)
                .catch(err => {
                    console.log(err);
                    reject(err);
                });
        });
    }
}