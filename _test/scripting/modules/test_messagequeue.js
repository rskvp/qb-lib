const amqp = require("message-queue");

function connect() {
    try {
        // connect to broker
        const conn = amqp.newConnection({
            "protocol": "amqp",
            "secret": "1234567890", // optional. If assigned messages are encrypted
            "url": "amqp://test:test@localhost:5672/"
        });

        // test broker connection
        if (!!conn && conn.ping()) {

            // declare an exchange (optional, you can send messages also directly at queue)
            try {
                conn.exchangeDeclare({
                    "passive": false,
                    "name": "MyExchange",
                    "kind": "direct",
                    "durable": true,
                    "auto-delete": false,
                    "internal": false,
                    "no-wait": false,
                    "args": null
                });
            } catch (err) {
                console.error("Error declaring Exchange:", err);
            }

            // declare a queue
            try {
                var queue = conn.queueDeclare({
                    "passive": false,
                    "name": "MyQueue",
                    "durable": false,
                    "exclusive": false,
                    "auto-delete": false,
                    "no-wait": false,
                    "args": null
                });
                console.log("Queue Declared:", JSON.stringify(queue));

                // bind exchange and queue
                try{
                    conn.queueBind({
                        "name": queue.name,
                        "exchange": "MyExchange",
                        "key": "only-some-messages",
                        "no-wait": false,
                        "args": null
                    });

                    // consume messages
                    try {
                        var consumer = conn.newListener({
                            "queue": queue.name,
                            "consumer-tag": "only-some-messages",
                            "no-local": false,
                            "auto-ack": false,
                            "exclusive": true,
                            "no-wait": false,
                            "args": null
                        });
                        consumer.listen(function (message) {
                            console.log("receiving:", JSON.stringify(message))
                        });
                    } catch (err) {
                        console.error("Error creating new consumer:", err);
                    }
                }catch(err){
                    console.error("Error Binding queue:", err);
                }
            } catch (err) {
                console.error("Error declaring Queue:", err);
            }
        } else {
            console.error("Connection not available. This should not happen immediatelly after a connection creation.");
        }
    } catch (err) {
        console.error("Broker not available:", err);
    }
}

function run() {
    var response = {};
    try {

        connect();

        var s = JSON.stringify(response);
        console.log(s);
        return s;
    } catch (err) {
        console.error("test_messagequeue.js", err);
        throw err;
    }
}

module.exports = {
    run: run
}