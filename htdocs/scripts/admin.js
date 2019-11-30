//Constant URI for the API route
const uri = "/api/v1/admins/users"
var selectedUserId = 0;

/* 
 *	Vue for the admin page
 *	Allows the admin user to view a list of other users
 *	Allows that admin to also filter all of the users
 *	TODO: Allow admin to delete users
 */
var admin = new Vue({
	// 	Element containing Vue data
	el: '#admin',

	//	Data to be inserted into the page
	data: {
		users: []
	},

	//	Initializes the data upon page load
	created: function(){
		this.users = this.fetchTable();
	},

	//	All methods useable by the Vue
	methods: {

		deleteUser: function(index, user) {
			confirm("Are you sure you want to delete the selected user?");
			this.users.splice(index, 1);
			deleteUser(user);
		},

		changePassword: function(user) {
			$("#change-password-form").slideDown();
			selectedUserId = user.userid;
		},
		//	Fetch the user data for the users table
		fetchTable: function(){
			$.getJSON(uri, function(json) {
				admin.results = json;
				admin.users = json;
			});
		},

		//	Fetch new, filtered data for the table
		filterTable: function(){
			let id = $('#id').val();
			let un = $('#uname').val();
			let fn = $('#fname').val();
			let ln = $('#lname').val();
			let q = "?";

			//	Add query parameters to the URI
			if(id !== ""){
				q = q + `id=${id}&`;
			}
			if(un !== ""){
				q = q + `un=${un}&`
			}
			if(fn !== ""){
				q = q + `fn=${fn}&`;
			}
			if(ln !== ""){
				q = q + `ln=${ln}`
			}

			let queryURI = uri + q;

			$.getJSON(queryURI, function(json){
				admin.results = json;
				admin.users = json;
			});
		}
	}
});

function deleteUser(user) {
    let req = new XMLHttpRequest();

	req.withCredentials = true;
	// On load, what was therequest status?
	req.addEventListener("load", function(evt){
		if (req.status === 200){
			//display success message
		} else {
			alert("We encountered an error."); //other server error. No handler on UI
		}
	});

	// Open the request to the register route, and then send the data in the body
	req.open("DELETE", "/api/v1/admins/users/" + user.userid); 
	req.send();
}

function changeUserPassword() {
	let password = $("#new-pass").val();

	console.log("password: " + password);

	if (password != "") {
		//send update password request

		let req = new XMLHttpRequest();

		req.withCredentials = true;
		// On load, what was therequest status?
		req.addEventListener("load", function(evt){
			if (req.status === 200){
				$("#change-password-form").slideUp();
			} else {
				alert("We encountered an error."); //other server error. No handler on UI
			}
		});

		// Open the request to the register route, and then send the data in the body
		req.open("PUT", "/api/v1/admins/passwords/" + selectedUserId); 
		req.send(JSON.stringify({password: password}));
	}
}