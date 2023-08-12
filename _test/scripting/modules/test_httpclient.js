var http = require("http");

function main() {
    var response = {}
    try {
        var client = http.newClient();
        client.addHeader("Authorization", "Bearer SOME_TOKEN_HERE");

        testGet(client, "http://gianangelogeminiani.me");

        testPost(client, "http://ivybot.tech:4199/api/program/invoke", {
            "namespace":"ivy",
            "function":"backoffice.echo",
            "params":"[\"sdf\",\"sdfsd\"]"
        });

    } catch (err) {
        console.error(err);
    }

    return response;
}

function testGet(client, url) {
    var res = client.get(url);
    console.log("GET: ", JSON.stringify(res),  "\n" + res.text());
    console.log("HEADER: ", JSON.stringify(res.header));
}

function testPost(client, url, body) {
    var res = client.post(url, body);
    console.log("POST: ", JSON.stringify(res),  "\n" + res.text());
    console.log("Content-Type: ", res.header["Content-Type"]);
    console.log("HEADER: ", JSON.stringify(res.header));
    console.log("BODY: ", res.body);
    console.log("BODY JSON: ", res.json());
}

//----------------------------------------------------------------------------------------------------------------------
//	e x p o r t
//----------------------------------------------------------------------------------------------------------------------

module.exports = {
    run: main
}
