// On click listener for login button
$("#login-action-btn").click(function(evt){
	let user = {username: $("#uname").val(), password: $("#pass").val()};
	loginUser(user);
});

function loginUser(user) {
	if (user.username === "") {
		//display username missing error
		$("#uname").addClass("is-invalid");
		return;
	}

	if (user.password === "") {
		$("#pass").addClass("is-invalid");

		//display missing password error
		return;
	}

	$("#uname").removeClass("is-invalid");
	$("#pass").removeClass("is-invalid");

	$("#invalid-feedback-login").hide();

	//Perform login web service
	let req = new XMLHttpRequest();

	req.withCredentials = true;
	// On load, what was therequest status?
	req.addEventListener("load", function(evt){

		// 200 = OK
		// Registered - redirect to login.html
		// Anything else = not OK
		// Was not registered
		if(req.status === 200){
			let twoFactor = JSON.parse(req.response)['2fa'];
			sessionStorage.setItem("2fa", twoFactor);
			sessionStorage.setItem("username", user.username);
			sessionStorage.setItem("isAdmin", 0);
			window.location.href = "twofactorauth.html";
		} else if (req.status == 401){
			$("#invalid-feedback-login").show(); //show invalid creds.
		} else {
			alert("We encountered an error."); //other server error. No handler on UI
		}
	});

	// Open the request to the register route, and then send the data in the body
	req.open("POST", "/api/v1/login");

	console.log("user object: " + user);
	req.send(JSON.stringify(user));
}