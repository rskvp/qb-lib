var http = require("http");

/**
 * MAIN EXPORTED FUNCTION
 * @returns {string}
 */
function run() {
    var response = {}

    try {
        // server
        var app = http.newServer();
        app.static({
            "enabled": true,
            "prefix": "/",
            "root": "./www",
            "index": "",
            "compress": true
        });

        // middleware
        app.use("/*", onMiddleware);

        // routing
        app.all("/api/:name", onRoute);

        // websocket
        app.websocket("/websocket", onWebsocket);

        // file upload
        app.all("/upload/*", onUpload);

        // OPEN SERVER
        app.listen([{
            "addr": "80",
            "tls": false
        }, {
            "addr": "443",
            "tls": true,
            "ssl_cert": "./cert/ssl-cert.pem",
            "ssl_key": "./cert/ssl-cert.key"
        }]);

        console.log("SERVER OPENED");

        // client
        /*  */
        var client = http.newClient();
        console.log("client", "GET", "http://localhost/");
        var res = client.get("http://localhost/");
        console.log(JSON.stringify(res));
        console.log("http://localhost/", res.text());
        res = client.get("http://localhost/api/hello");
        console.log("http://localhost/api/hello", res.text());
        res = client.get("http://localhost/api/foo");
        console.log("http://localhost/api/foo", res.text());

        // -----------------------
        // console.log("JOIN SERVER");
        app.join();
        // -----------------------

    } catch (err) {
        console.error(err);
    }

    return JSON.stringify(response);
}

function onUpload(req, res) {
    console.log("handled onUpload");

    // get multipart form with data and files
    var form = req.multipart();
    if (!!form) {

    } else {
        console.error("NOTHING FROM UPLOAD REQUEST")
    }
}

function onMiddleware(req, res, next) {
    console.log("MIDDLEWARE");
    next();
}

function onRoute(req, res) {
    var data = {
        "message": "This is a response to: " + req.originalUrl,
        "query": req.query,
        "params": req.params,
        "name": req.param("name"),
        "path": req.path,
        "originalUrl": req.originalUrl,
        "baseUrl": req.baseUrl
    }
    res.json(data);
}

function onWebsocket(req, res) {
    try {
        var data = {
            "params": req.params,
            "uid": req.param("uid"),
            "response_message": "Hello websocket"
        }
        console.log("websocket request", req.text());
        console.log("data", JSON.stringify(data));

        // return data
        res.json(data);
    } catch (err) {
        console.error("onWebsocket", err);
    }
}

module.exports = {
    run: run
}
