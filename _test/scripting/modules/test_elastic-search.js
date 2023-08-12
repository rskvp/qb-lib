var elastic = require("elastic-search");

function main() {
    var response = {};
    try {
        var DB_AUTH = "root:!root"
        var config = {
            "db_internal": {
                "driver": "arango",
                "dsn": DB_AUTH + "@tcp(localhost:8529)/test)"
            },
            "db_external": {
                "driver": "arango",
                "dsn": DB_AUTH + "@tcp(localhost:8529)/test)"
            }
        };
        var engine = elastic.newEngine(config);
        engine.open();
        try {
            // put some data to be indexed
            engine.put("test", "010", "Marino Todisco Hello boy, this is some elastic data to search for. Text can be long long long and long again!!! Rimini 23012020");

            // search for some indexed data
            var data = engine.get("", "Give me elastic or long data!!");
            response.data = data

            console.log("DATA: " + JSON.stringify(data));
            for(var i=0;i<data.length;i++){
                console.log(i + "\t", JSON.stringify(data[i]));
            }
        } finally {
            engine.close();
        }
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
