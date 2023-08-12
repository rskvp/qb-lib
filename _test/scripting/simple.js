(function(){
    console.log("simple.js", "RUNNING...");
    const simple = require("simple");
    const result = simple.echo("HELLO ECHO");
    console.log("simple.echo(\"HELLO ECHO\")", result);
    return result
})();