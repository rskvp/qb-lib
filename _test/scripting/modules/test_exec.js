var exec = require("exec-utils");

module.exports = {
    run: function () {
        var response = {}

        try {
            response = exec.run("ls -lah");
            console.log("response:", JSON.stringify(response));
        } catch (err) {
            console.error(err);
        }

        return JSON.stringify(response);
    }
}
