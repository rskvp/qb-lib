

function main() {
    var response = {};
    try {
        console.log(process.env.names());
        console.log(process.env.names()[0] + "=" + process.env[process.env.names()[0]]);

        process.env.TEST = "hello test";
        console.log(process.env.TEST);
    } catch (err) {
        console.error(err);
        response.error = err;
    }
    return response;
}

function testLogin(auth, username, password) {
    var loginData = auth.userSignIn(username, password)
    console.log(JSON.stringify(loginData));
    if (!!loginData.error) {
        console.error(loginData.error);
        return false;
    } else {
        // user logged
        console.log("USER ID: " + loginData.item_id);
        console.log("USER PAYLOAD: " + JSON.stringify(loginData.item_payload));
        console.log("TOKEN: " + loginData.access_token);

        var claims = auth.tokenParse(loginData.access_token);
        console.log("LOGIN CLAIMS: " + JSON.stringify(claims));

        // update protected payload
        var payload = loginData.item_payload;
        payload.info = {
            "age": 35,
            "gender": "male"
        }
        var updated_id = auth.userUpdate(loginData.item_id, payload)
        console.log("UPDATED: " + updated_id);
    }
    return true;
}

function testRegisterAndConfirm(auth, username, password, payload) {
    var data = auth.userSignUp(username, password, payload || {});
    console.log(JSON.stringify(data));
    if (!!data.error) {
        console.error(data.error);
        return false;
    } else {
        try {
            data = auth.userConfirm(data.confirm_token);
            var claims = auth.tokenParse(data.access_token);
            console.log("REGISTERED CLAIMS: " + JSON.stringify(claims));
        } catch (err) {
            console.error(err);
            return false
        }
    }
    return true
}

//----------------------------------------------------------------------------------------------------------------------
//	e x p o r t
//----------------------------------------------------------------------------------------------------------------------

module.exports = {
    run: main
}
