const uri = "/api/v1/teachers/classes/" + getQueryVariable("classId");

var questiontable = new Vue({
	// 	Element containing Vue data
	el: '#questions',

	//	Data to be inserted into the page
	data: {
		questions: [],
		list: true
	},

	//	Initializes the data upon page load
	created: function(){
		this.questions = this.fetchTable();
		addQr();
	},

	//	All methods useable by the Vue
	methods: {
		//	Fetch the classroom data for the classes table
		fetchTable: function(){
			$.getJSON(uri, function(json) {
				questiontable.results = json;
				questiontable.questions = json;
				console.log(JSON.stringify(json));
			});
        },
        
        editQuestion: function(index) {
            let question = this.questions[index];
            window.location.href = "createquestion.html?qId=" + question.questionId + "&title=" + question.questionText + "&classId=" + question.classcode + "&editing=true";
        },
        statsQuestion: function(index) {
            let question = this.questions[index];
            window.location.href = "statsquestion.html?qId=" + question.questionId + "&title=" + question.questionText;
        },
        deleteQuestion: function(index) {
            confirm("Are you sure that you want to delete this question?");
            let question = this.questions[index];
            deleteQuestionWithId(question.questionId, index)

        }
	}
});

function deleteQuestionWithId(qId, index) {
    let req = new XMLHttpRequest();

	req.withCredentials = true;
	// On load, what was therequest status?
	req.addEventListener("load", function(evt){
		if(req.status === 200){
            questiontable.questions.splice(index, 1); //If deletion was successful, remove question from the table.
		} else {
			alert("We encountered an error."); //other server error. No handler on UI
		}
	});

	// Open the request to the register route, and then send the data in the body
	req.open("DELETE", "/api/v1/teachers/questions/" + qId); 

	req.send();
}

function pressNewQuestion() {
    window.location.href = "/createquestion.html?classId=" + getQueryVariable("classId");
}

function addQr(){
	let param = window.location.href.split("?")[1];
	let classId = param.split("=")[1];
	$("#qr").html(`<img src='/api/v1/teachers/classes/${classId}/qrjoin'/>`)
}
