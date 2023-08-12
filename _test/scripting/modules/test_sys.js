var sys = require("sys");


module.exports = {

    run: function () {
        var response = {};
        try {
            response.machine_id = sys.id();
            response.os = sys.getOS();
            response.osVersion = sys.getOSVersion();
            response.info = sys.getInfo();

            var s = JSON.stringify(response);
            console.log(s);
            return s;
        } catch (err) {
            console.error("test_sql.js", err);
            throw err;
        }
    }
}