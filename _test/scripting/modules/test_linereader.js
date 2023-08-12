var linereader = require("line-reader");
var path = require("path");

module.exports = {

    run: function () {
        try {
            var filename = path.resolve("./longfile.txt");

            var response = {
                lines_01: linereader.readLines(filename, 1),
                lines_10: linereader.readLines(filename, 10),
            };

            var count = 0;
            linereader.eachLine(filename, function (text) {
                count++;
                console.info(count, text);
                return false; // continue
            });

            var s = JSON.stringify(response);
            console.debug(s);
            return s;
        } catch (err) {
            console.error("test_linereader.js", err);
            return err;
        }
    }
}