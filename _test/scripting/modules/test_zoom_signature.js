var crypto = require("crypto-utils");

function main() {
    var response = {};
    try {
        var apiSecret = "API-SECRET";
        var apiKey = "API-KEY", role = 0, meetingNumber = "92161605297"

        var timestamp = new Date().getTime() - 30000;
        var msg = crypto.encodeBase64(apiKey + meetingNumber + timestamp + role);
        var hash = crypto.encodeBase64(crypto.encodeSha256(apiSecret, msg));
        var signature = crypto.encodeBase64(apiKey + "." + meetingNumber + "." + timestamp + "." + role + "." + hash)

        response.signature = signature;
        console.log(signature);
    } catch (err) {
        console.error(err);
        response.error = err;
    }
    return response;
}

//----------------------------------------------------------------------------------------------------------------------
//	e x p o r t
//----------------------------------------------------------------------------------------------------------------------

module.exports = {
    run: main
}
