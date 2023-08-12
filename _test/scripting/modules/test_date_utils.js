var date = require("date-utils");

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

function runTest() {
    var response = {};
    try {

        // current date-time
        current = date.wrap();
        response.current = current.format(); // format with default pattern (ISO)

        // js native Date
        jsdate = date.wrap(new Date());
        response.jsdate = jsdate.format(); // format with default pattern (ISO)

        sdate = date.wrap("2020-09-01 12:15");
        response.sdate = sdate.format(); // format with default pattern (ISO)

    } catch (err) {
        console.error("test_data_utils.js", err);
        throw err;
    }
    text = JSON.stringify(response);
    console.log(text);
    return text;
}

//----------------------------------------------------------------------------------------------------------------------
//	e x p o r t s
//----------------------------------------------------------------------------------------------------------------------

module.exports = {

    run: runTest
}