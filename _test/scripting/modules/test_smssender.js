var sender = require("sms-sender");

module.exports = {
    run: function () {
        try {
            var response = {};

            var transport = sender.createTransport(
                "smshosting",
                {
                    "enabled": true,
                    "auto-short-url": true,
                    "providers": {
                        "smshosting": {
                            "method": "GET",
                            "endpoint": "https://api.smshosting.it/rest/api/smart/sms/send?authKey={{auth-key}}&authSecret={{auth-secret}}&text={{message}}&to={{to}}&from={{from}}",
                            "params": {
                                "auth-key": "",
                                "auth-secret": "",
                                "message": "",
                                "to": "",
                                "from": "ANGELO"
                            },
                            "headers": {}
                        }
                    }
                }
            );
            // Text Message
            var message = "Check this: https://gianangelogeminiani.me";
            transport.send(message, "+39 347 7857785");

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
}