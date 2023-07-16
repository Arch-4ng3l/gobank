

function send(method) {

    if(method === "login") {
        var number = document.getElementsByName('number')[0].value;
        var password = document.getElementsByName('password')[0].value;
	    var requestdata = {
            number: number,
            password: password
        };
        var req = "/api/login";
    }
    else {
        var firstName = document.getElementsByName('firstName')[0].value;
        var lastName = document.getElementsByName('lastName')[0].value;
        var password = document.getElementsByName('password')[0].value;

        var requestdata = {
            firstName: firstName, 
            lastName: lastName, 
            password: password
        };

        var req = "/api/account";
    }
    fetch(req, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(requestdata),
    })

}

