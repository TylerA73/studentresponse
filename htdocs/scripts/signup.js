/* Event listener for the user registration button */
$("#register-btn").click(function(event){
	let first = $("#fname").val();
	let last = $("#lname").val();
	let username = $("#uname").val();
	let pass = $("#pass").val();

	if(first === "" || last === "" || username === "" || pass === ""){
		alert("Please enter information");
	}else{
		if(pass.length < 5){
			alert("Password must be at least 5 characters")
		}else{

			// Build the json
			let user = {
				username: username,
				firstName: first,
				lastName: last,
				password: pass

			}

			// Pass the json to the API route
			registerUser(JSON.stringify(user));
		}
	}
});

/* Calls the API to register the user */
function registerUser(user){

	let req = new XMLHttpRequest();

	// On load, what was therequest status?
	req.addEventListener("load", function(evt){

		// 200 = OK
		// Registered - redirect to login.html
		// Anything else = not OK
		// Was not registered
		if(req.status === 200){
			alert("Registered.");
			window.location.href = "login.html";
		} else if (req.status === 409) {
			alert("That username already exists.");
		} else {
			alert("We encountered an error.");
		}
	});

	// Open the request to the register route, and then send the data in the body
	req.open("POST", "/api/v1/register");
	req.send(user);


}
