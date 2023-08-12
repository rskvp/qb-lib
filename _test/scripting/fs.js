(function(){
    let test = require("./test_fs.js");

    console.log("fs.js", "RUNNING...");
    let result = test.run();
    // console.reset();
    return result
})();