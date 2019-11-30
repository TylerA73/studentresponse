//Constant URI for the API route
const uri = "/api/v1/teachers/classes";

/* 
 *	Vue for the instructor page
 *	Allows instructor to view their classrooms, and create new ones
 *	TODO: Link classroom buttons to an edit page
 */
var instructor = new Vue({
	// 	Element containing Vue data
	el: '#instructor',

	//	Data to be inserted into the page
	data: {
		classrooms: [],
		list: true
	},

	//	Initializes the data upon page load
	created: function(){
		this.classrooms = this.fetchTable();
	},

	//	All methods useable by the Vue
	methods: {

		goToClass: function (classroom) {
			window.location.href = "/teacherquestions.html?classId=" + classroom.classcode;
		},

		//	Fetch the classroom data for the classes table
		fetchTable: function(){
			$.getJSON(uri, function(json) {
				instructor.results = json;
				instructor.classrooms = json;
			});
		},

		//	Post a new class to the database
		postClass: function(){
			let className = $("#class-name").val();
			let classroom = {
				classname: className
			}

			let req = new XMLHttpRequest();
			req.addEventListener("load", function(evt){
				if(req.status === 200){
					alert("New class added");
				}else{
					alert("We encountered a problem");
				}
			});

			req.open("POST", uri);
			req.send(JSON.stringify(classroom));
			$.getJSON(queryURI, fetchTable());

			this.list = false;
		}
	}
});