var http = require("http");
var auth0 = require("auth0");

var count = 0;
var config = {
    "secrets": {
        "auth": "this-is-token-to-authenticate",
        "access": "hsdfuhksdhf5435khjsd",
        "refresh": "hsdfuhqswe34qwksdhfkhjsd"
    },
    "cache-storage": {
        "driver": "arango",
        "dsn": "test_user:test_password@tcp(localhost:8529)/test)"
    },
    "auth-storage": {
        "driver": "arango",
        "dsn": "test_user:test_password@tcp(localhost:8529)/test)"
    }
};

/**
 * Launch a web server instance
 * @returns {string}
 */
function main() {
    var response = {}

    console.log("RUNNING HTTP TEST");

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
        app.all("/auth/signin", onUserSignIn);
        app.all("/auth/signup", onUserSignUp);
        app.all("/auth/update", onUserUpdate);

        app.all("/upload/*", onUpload);

        // websocket
        app.websocket("/websocket", onWebsocket);

        app.listen([{
            "addr": ":80",
            "tls": false
        }, {
            "addr": ":443",
            "tls": true,
            "ssl_cert": "./cert/ssl-cert.pem",
            "ssl_key": "./cert/ssl-cert.key"
        }]);


        // -----------------------
        // console.log("JOIN SERVER");
        app.join();
        // -----------------------

    } catch (err) {
        console.error(err);
    }

    console.log("EXITING SERVER. bye bye");

    return JSON.stringify(response);
}

function onUpload(req, res) {
    console.log("handled onUpload");
    var response = {
        "error": undefined,
        "response": []
    };
    try {
        // get multipart form with data and files
        var form = req.multipart();
        if (!!form) {
            console.log("params: " + JSON.stringify(form.data));
            console.log("files: " + form.files.length);
            for (var i = 0; i < form.files.length; i++) {
                var file = form.files[i]
                var pathLevel = 3; // save file under "./YEAR/MONTH/DAY/" path
                var max = 2100000; // max size = 2Mb
                var path = file.save("./uploads/", pathLevel, max); // save file under "./uploads/YEAR/MONTH/DAY/" path
                // add file system path
                response.response.push(path);
                console.log("SAVED TO: " + path);
            }
        } else {
            console.error("NOTHING FROM UPLOAD REQUEST")
            response.error = "NOTHING FROM UPLOAD REQUEST";
        }
    } catch (err) {
        console.log(err);
        response.error = err.message;
    }
    if (!!response.error) {
        res.status(403);
    } else {
        res.json(response);
    }
}

function onUserSignIn(req, res) {
    console.log("USER SIGN-IN");

}

function onUserSignUp(req, res) {
    console.log("USER SIGN-UP");

}

function onUserUpdate(req, res) {
    console.log("USER UPDATE");
    var auth = req.getAuth();
    if (!auth.token) {
        res.status(401); // Unauthorized
    }
    var token = auth.token;

}

function onMiddleware(req, res, next) {
    count++;
    console.log(count + " - MIDDLEWARE");
    var authorization = req.getAuth();
    console.log("Authorization: " + JSON.stringify(authorization));
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
    var auth = req.getAuth();
    if (!auth.token) {
        res.status(401); // Unauthorized
    } else {
        res.send(data, "application/json");
    }
}

function onWebsocket(req, res) {
    try {
        var data = {
            "params": req.params,
            "uid": req.param("uid")
        }
        console.log("websocket request", req.text());
        console.log("data", JSON.stringify(data));

        // return data
        res.json(data);
    } catch (err) {
        console.error("onWebsocket", err);
    }
}

//----------------------------------------------------------------------------------------------------------------------
//	e x p o r t
//----------------------------------------------------------------------------------------------------------------------

module.exports = {
    run: main
}
