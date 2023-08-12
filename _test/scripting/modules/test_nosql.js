var arango = require("nosql");


module.exports = {

    run: function () {
        var response = {
            "driver": "arango",
            "dsn": "root:admin@http(localhost:8529)/test",
            "connected": false,
        };
        try {
            var db;
            try {
                // open database
                db = arango.open(response.driver, response.dsn);
                response.connected = true;

                // ensure collections exists
                console.log("COLLECTION: ", db.collection("nodes").name());
                console.log("COLLECTION: ", db.collection("nodes_terminals").name());

                db.collection("nodes_terminals").ensureIndex("persist", ["node_key"], false);

                // query
                response.query = db.query("FOR doc IN nodes_terminals\n" +
                    "FILTER doc.node_key == @node_key\n" +
                    "SORT doc._key\n" +
                    "RETURN doc._key", {"node_key": "abcd"})
                console.log("query", response.query);

                // insert
                response.insert = db.insert("nodes_terminals", {
                    "node_key": "abcd",
                    "timestamp": (new Date()).getTime()
                })
                console.log("insert", JSON.stringify(response.insert));

                // upsert
                response.upsert = db.upsert("nodes_terminals", {
                    "_key": "test_001",
                    "timestamp": (new Date()).getTime()
                })
                console.log("upsert", JSON.stringify(response.upsert));

                // update
                response.update = db.update("nodes_terminals", {
                    "_key": "test_001",
                    "age": Math.floor(Math.random() * 70),
                    "first_name": "UPDATED " + (new Date())
                })
                console.log("update", JSON.stringify(response.update));

                // delete
                response.delete = db.delete("nodes_terminals", response.insert);
                console.log("delete", JSON.stringify(response.delete));

                // count
                response.count = db.count("FOR doc IN nodes_terminals\n" +
                    "FILTER doc.node_key == @node_key\n" +
                    "SORT doc._key\n" +
                    "RETURN doc", {"node_key": "abcd"}
                );
                console.log("count", response.count);

            } finally {
                // finally close
                if (!!db) {
                    db.close();
                }
            }

            var s = JSON.stringify(response);
            console.log(s);
            return s;
        } catch (err) {
            console.error("test_nosql.js", err);
            throw err;
        }
    }
}