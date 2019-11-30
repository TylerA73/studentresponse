$("#log-out-btn").click(function(){
	let uri = "/api/v1/logout"
	let req = new XMLHttpRequest();
	req.addEventListener("load", function(){
		if(req.status === 200){
			sessionStorage.removeItem("username");
			sessionStorage.removeItem("2fa");
			sessionStorage.removeItem("isAdmin");
			window.location.href = "/index.html"
		}
	});
	req.open("POST", uri);
	req.send();
});