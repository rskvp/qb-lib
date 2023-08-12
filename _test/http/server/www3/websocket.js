(function () {
    try {
        let vws = window["__vws"];
        if (!!vws) {
            const URL = "ws://127.0.0.1:80/websocket"
            // const URL = "ws://165.22.21.124:80/endpoint"
            // const URL = "ws://localhost:80/endpoint"

            log("info", "Vanilla.Websocket", "version: " + vws.version);
            const client = vws.create(URL);
            log("info", "Vanilla.Websocket client", "host: " + client.host);

            // handle connection error
            client.on("on_error", (err) => {
                log("error", "Vanilla.Websocket client#on_error", err);
            });
            client.on("on_message", (message) => {
                let source = "broadcast server message"
                if (message["uid"] === "test_message") {
                    source = "response to client request"
                }
                console.info("Vanilla.Websocket client#on_message", source, message);
            });
            client.on("on_close", (err) => {
                console.info("Vanilla.Websocket client#on_close");
            });

            // send a message
            send(client);

            $("#btn_send").on("click", (e) => {
                e.preventDefault;
                send(client);
            });

        } else {
            log("error", "Vanilla.Websocket not loaded!", "");
        }
    } catch (e) {
        console.error(e);
    }

    function send(client) {
        // creates a message for the server

        const message = {
            "uid": "test_message",
            "payload": {
                "app_token": "iuhdiu87w23ruh897dfyc2w3r",
                "namespace": "get",
                "function": "array",
                "params": ["Hello Vanilla.Websocket"]
            }
        };

         log("info", "Sending message", message)
        // send message to server
        client.send(message, (full_response) => {
            const response = full_response["response"]||{};
            const error = response["error"];
            if (!!error) {
                log("error", "Vanilla.Websocket client", error);
            } else {
                const data = response["data"];
                log("info", "Response", full_response);
            }
        });
    }

    function log(level, context, message) {
        try {
            let style = "alert-primary";
            switch (level) {
                case "error":
                    console.error(context, message);
                    style = "alert-danger"
                    break;
                default:
                    console.info(context, message);
            }
            message = JSON.stringify(message)
            const html = "<div class='row'><div class=\"alert " + style + "\" role=\"alert\">" + context + message + " </div></div>";
            $(html).appendTo($("#logs"));
        } catch (e) {
            console.error(e);
        }
    }

})();