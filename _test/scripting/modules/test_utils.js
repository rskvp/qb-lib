var fu = require("file-utils");
var cu = require("crypto-utils");

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------


//----------------------------------------------------------------------------------------------------------------------
//	e x p o r t s
//----------------------------------------------------------------------------------------------------------------------

module.exports = {

    run: function () {
        var response = {};
        try {

            response.data = fu.fileReadBytes("./data.csv");
            console.log("READ", response.data);
            response.base64 = cu.encodeBase64(response.data);
            console.log("BASE64", response.base64);
            response.decoded = cu.decodeBase64ToText(response.base64);
            console.log("TEXT", response.decoded);

            var s = JSON.stringify(response);
            console.log(s);

            return s;
        } catch (err) {
            console.error("test_utils.js", err);
            throw err;
        }
    }
}