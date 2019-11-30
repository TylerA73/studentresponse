$("#login-action-btn").click(function (evt) {
	let user = $("#uname-login").val();
	let pass = $("#pwd").val();
	if ((user === "") || (pass === "")) {
		$("#invalid-feedback-login").show();
	} else {
		console.log(user + " " + pass);
		$("#content").html("<p>Logged in</p>");
	}


	//success
	// $("#login").fadeOut();
	// $("#main").fadeOut();
});


/*
// Login Button click event
$("#login-btn").click(function(evt){
  console.log("Login clicked");
  // document.querySelector("#main").style.display = "none";
  //document.querySelector("#login").style.display = "block";
  $("#login").slideDown("fast");
});






$.get("login.html", function(data){
  var login = new Vue({
    el: '#login',
    data: {
      content: data
    }
  });
  $("#login-action-button").click(function(evt) {
    console.log("perform login");

    //if login failes
    $("#invalid-feedback-login").show();

    //success
    // $("#login").fadeOut();
    // $("#main").fadeOut();
});
*/

new Vue({
	el: '#join-class-area',
	data: {
		show: false
	},
	methods: {
		showJoin: function (event) {
			if (this.show == false) {
				this.show = true
			}
		},
		//Checks if class exists
		joinClass: function (event) {
			let classCode = $('#class-input').val();
			if (classCode === "") {
				toastr["error"]("Please enter a class")
				
			} else {
				$.getJSON(window.location.href + "api/v1/classes/" + classCode, function(json) {
					sessionStorage.setItem('classJson', JSON.stringify(json))
					window.location.href = "joinedclass.html"
				}).fail(function () {
					toastr["error"]("Class does not exist")
				})
			}
		}
	}
})