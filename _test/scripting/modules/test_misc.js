var path = require("path");
var fs = require("fs");


module.exports = {

    run: function () {
        try {
            var response = {};
            var filename = path.resolve("./testfile.txt");

            if (fs.existsSync(filename)){
                var body = fs.readFileSync(filename);
                response.body = body;

                if (body.length===0){
                    fs.writeFileSync(filename, "HELLO FROM JAVASCRIPT")
                } else {
                    if (body==="HELLO FROM JAVASCRIPT"){
                        fs.writeFileSync(filename, "HELLO FROM JAVASCRIPT, SECOND STEP")
                    } else{
                        fs.writeFileSync(filename, "HELLO FROM JAVASCRIPT")
                    }
                }
            }

            var s = JSON.stringify(response);
            console.debug(s);
            return s;
        } catch (err) {
            console.error("test_fs.js", err);
            return err;
        }
    }
}