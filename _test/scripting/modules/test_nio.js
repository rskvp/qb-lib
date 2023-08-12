var nio = require("nio");
var files = require("file-utils");
var crypto = require("crypto-utils");

function onFileUpload(name, params) {
    console.log("onFileUpload", name, params);
    // "this" is the server
    console.log("Clients connected: ", this.count());

    console.log("file name: ", params[0]);
    console.log("file data raw: ", params[1]);
    console.log("file data decoded: ", crypto.decodeBase64ToText(params[1]));

    return "File Uploaded!";
}

module.exports = {

    run: function () {
        var response = {};
        try {
            console.log("SERVER STARTING");
            var server = nio.newServer(10001);
            server.open();
            console.log("SERVER STARTED");
            server.listen("file_upload", onFileUpload);

            console.log("CLIENT OPENING");
            var client = nio.newClient("localhost:10001");
            client.secure(true);
            client.open();
            console.log("CLIENT OPENED");
            console.log("CLIENT SENDING MESSAGE");
            var resp = client.send("file_upload", "data.csv", files.fileReadBytes("./data.csv"));
            console.log("SERVER RESPONSE", resp);

            var s = JSON.stringify(response);
            console.log(s);

            return s;
        } catch (err) {
            console.error("test_nio.js", err);
            throw err;
        }
    }
}