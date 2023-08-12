var path = require("path");
var nodemailer = require("nodemailer");

module.exports = {
    run: function () {
        // test using MailTrap
        var filename = path.resolve("./data.csv");
        try {
            var response = {};

            var transport = nodemailer.createTransport({
                host: 'smtp.mailtrap.io',
                port: 2525,
                secure: false,
                auth: {
                    user: '8590e66d494168',
                    pass: '8926a5d72a274f'
                }
            });
            // Text Message
            var message = {
                from: '9f48ff8df3-36d93f@inbox.mailtrap.io', // Sender address
                to: 'angelo.geminiani@gmail.com',         // List of recipients
                subject: 'Design Your Model S | Tesla', // Subject line
                text: 'Have the most fun you can in a car. Get your Tesla today!',
            };
            transport.sendMail(message, function (err, info) {
                if (err) {
                    console.error(err)
                } else {
                    console.log("Text email sent");
                }
            });
            // HTML MESSAGE
            var message_html = {
                from: '9f48ff8df3-36d93f@inbox.mailtrap.io', // Sender address
                to: 'angelo.geminiani@gmail.com',         // List of recipients
                subject: 'Test with HTML', // Subject line
                html: '<strong>Hello HTML</strong><br>This is HTML text',
                attachments: [{
                    "filename": path.basename(filename),
                    "path": filename
                }]
            };
            transport.sendMail(message_html, function (err, info) {
                if (err) {
                    console.error(err)
                } else {
                    console.log("HTML email sent");
                }
            });

            response.transport = transport;
            response.message = message;

            var s = JSON.stringify(response);
            console.debug(s);
            return s;
        } catch (err) {
            console.error(err);
            return err;
        }
    },

    // test with GMail account
    runGMail: function(){
        var filename = path.resolve("./data.csv");
        var response = {};

        var user = "";
        var pass = "";
        var host = "smtp.gmail.com";
        var port = 465; // SSL port works. TLS does not work
        var secure = true;

        var transport = nodemailer.createTransport({
            host: host,
            port: port,
            secure: secure,
            auth: {
                user: user,
                pass: pass
            }
        });

        // HTML MESSAGE
        var message_html = {
            from: '9f48ff8df3-36d93f@inbox.mailtrap.io', // Sender address
            to: 'angelo.geminiani@gmail.com',         // List of recipients
            subject: 'Test with HTML', // Subject line
            html: '<strong>Hello HTML</strong><br>This is HTML text',
            attachments: [{
                "filename": path.basename(filename),
                "path": filename
            }]
        };
        transport.sendMail(message_html, function (err, info) {
            if (err) {
                console.error(err)
            } else {
                console.log("HTML email sent");
            }
        });

        response.transport = transport;
        response.message = message_html;

        var s = JSON.stringify(response);
        console.debug(s);
        return s;
    }
}