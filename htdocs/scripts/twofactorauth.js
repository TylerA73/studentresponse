//Constant URI for the API route
const uri = "/api/v1/2fa/";
var onCompletePage = "/instructor.html";

$(document).ready(function(){
	if (getQueryVariable("admin") === "true") {
		onCompletePage = "/admin.html";
	}
	var authentication = new Vue({
	// 	Element containing Vue data
	el: '#authentication',

	//	Data to be inserted into the page
	data: {
		setup: true
	},

	created: function() {
		this.setup = this.checkAuthType();
	},

	//	All methods useable by the Vue
	methods: {

		checkAuthType: function(){
			let auth = sessionStorage.getItem("2fa");
			if(auth === "setup"){
				$("#qr").html("<img src='/api/v1/2fa/qr'/>");
				return true;
			}else{
				return false;
			}
		},

		skipClick: function(){
			
			let url = uri + "pass";
			let req = new XMLHttpRequest();
			req.addEventListener("load", function(evt){
				if(req.status === 200){
					window.location.href = onCompletePage
				}else{
					alert("There was an error");
				}
			});
			req.open("POST", url);
			req.send();
		},

		setupProceedClick: function(){

			let url = uri + "challenge";
			let code = {
				"code" : String($("#setup-code").val())
			}
			console.log(JSON.stringify(code));
			if(code.code.length !== 6){
				alert("Must be a 6 digit code")
			}else{
				let req = new XMLHttpRequest();
				req.addEventListener("load", function(evt){
					if(req.status === 200){
						window.location.href = onCompletePage
					}else{
						alert("There was an error");
					}
				});
				req.open("POST", url);
				req.send(JSON.stringify(code));
			}
		},

		challengeProceedClick: function(){
			let url = uri + "challenge";
			let code = {
				"code" : String($("#challenge-code").val())
			}
			console.log(JSON.stringify(code));
			if(code.code.length !== 6){
				alert("Must be a 6 digit code")
			}else{
				let req = new XMLHttpRequest();
				req.addEventListener("load", function(evt){
					if(req.status === 200){
						window.location.href = onCompletePage
					}else{
						alert("There was an error");
					}
				});
				req.open("POST", url);
				req.send(JSON.stringify(code));
			}
		}
	}
});
});

