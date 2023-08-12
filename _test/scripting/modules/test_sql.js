var sql = require("sql");


module.exports = {

    run: function () {
        var response = {
            "driver": "mysql",
            "dsn": "admin:admin@tcp(localhost:3306)/test",
            "connected": false,
        };
        try {
            var db;
            try {
                db = sql.open(response.driver, response.dsn);
                response.connected = true;

                // query
                console.log("query");
                response.query = db.query("select * from table1")
                // insert
                console.log("insert");
                response.insert = db.insert("table1", {
                    "age": Math.floor(Math.random() * 70),
                    "first_name": "Juan " + (new Date())
                })
                // update
                console.log("update");
                response.update = db.update("table1", "id", 1, {
                    "age": Math.floor(Math.random() * 70),
                    "first_name": "UPDATED " + (new Date())
                })
                // delete
                console.log("delete");
                response.delete = db.delete("table1", "id > 10");
                // count
                console.log("count");
                response.count = db.count("table1");
                // countDistinct
                console.log("countDistinct");
                response.countDistinct = db.countDistinct("table1", "age", "first_name IS NOT NULL");

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
            console.error("test_sql.js", err);
            throw err;
        }
    }
}