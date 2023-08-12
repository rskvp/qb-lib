var path = require("path");
var fs = require("fs");


module.exports = {

    run: function () {
        try {
            var filename = path.resolve("./data/testfile.txt");

            var response = {
                exists: fs.existsSync(filename),
                readFileSync: fs.readFileSync(filename),
                isDirectory: fs.statSync(filename).isDirectory(),
                isFile: fs.statSync(filename).isFile(),
                size: fs.statSync(filename).size,
                writeFileSync: fs.writeFileSync(filename, "HELLO FROM JAVASCRIPT"),
                readFileSync2: fs.readFileSync(filename),
            };

            var s = JSON.stringify(response);
            console.debug(s);
            return s;
        } catch (err) {
            console.error("test_fs.js", err);
            return err;
        }
    }
}