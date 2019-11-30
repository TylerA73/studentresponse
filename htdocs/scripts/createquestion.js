//Globals
var trueFalseAnswer = true;
var mcMode = false;
var isPosting = false;

var editingMode = false;
var qId = null;

function loaded() {
    if (getQueryVariable("editing") === "true") {
        editingMode = true;
        mcMode = true;
        $("#q-title").val(getQueryVariable("title"));
        $("#response-textfield-container").show();
        $("#q-type-button-container").hide();
        $("#add-response").hide();

        qId = getQueryVariable("qId");

        //fetch the question
        $.getJSON("/api/v1/teachers/questions/" + getQueryVariable("qId"), function(json) {
            console.log("got the json string object: " + JSON.stringify(json));
            $("#create-question-mc").html("Update Question");

            for (var i = 0; i < json.answers.length; i++) {
                let answer = json.answers[i];
                forms.items.push({value: answer.answerText, correct:false, id:answer.answerId})

                if (answer.iscorrect == true) {
                    forms.checked = i;
                }
            }
        });
    }
}

//JQuery button actions
$("#mc-button").click(function () {
    $("#response-textfield-container").fadeIn();
    $("#q-type-button-container").hide();
    mcMode = true;
});

$("#tf-button").click(function () {
    $("#true-false-container").fadeIn();
    $("#q-type-button-container").hide();
    mcMode = false;
});

//Vue models
var forms = new Vue({
    el: '#response-textfield-container',
    data: {
        items: [],
        checked: ''
    },
    methods: {
        addResponse: function () {
            this.items.push({ value: "", correct: false, id:""});

            if (this.items.length == 1) { //select the first radio by default.
                this.checked = 0;
            }
        },

        removeResponse: function (i) {

            this.items.splice(i, 1);

            if (i == this.checked) {
                this.checked = 0;
            }

            if (i < this.checked) { //if a user removes a bullet from one below. 
                this.checked -= 1;
            }
        }
    }
});

function pressTFRadio(radio) {
    if (radio.value === "true") {
        trueFalseAnswer = true;
    } else if (radio.value === "false") {
        trueFalseAnswer = false;
    }
}

function submitAction() {
    //First, let's check if the question title has text.
    if ($("#q-title").val() === "") {
        $("#q-title").addClass("is-invalid");
        return;
    }
    if (mcMode == true) { //Multiple choice
        var newItems = [];
        var newObject = { question: { questiontext: $("#q-title").val(), questiontype: "MC" } };

        for (var i = 0; i < forms.items.length; i++) {
            var item = forms.items[i];

            if (item.value != "") { //if the value has text
                var obj = { answertext: item.value, iscorrect: (i == forms.checked)}

                if (item.id != "") {
                    obj["answerId"] = item.id;
                }

                newItems.push(obj);
            }
        }

        if (newItems.length == 0) { //Error: there is no text in the responses
            alert("Error: None of the responses contain an answer. Please double check and try again.");
            return;
        }

        newObject["answers"] = newItems;

        if (editingMode == false) {
            postNewQuestion(newObject);
        } else {
            updateQuestion(newObject);
        }
    } else { //true/false
        var newObject = {
            question: { questiontext: $("#q-title").val(), questiontype: "MC" },
            answers: [
                {
                    answertext: "True",
                    iscorrect: (trueFalseAnswer == true)
                },
                {
                    answertext: "False",
                    iscorrect: (trueFalseAnswer == false)
                }
            ]
        };
        if (editingMode == false) {
            postNewQuestion(newObject);
        } else {
            updateQuestion(newObject);
        }
    }
}

function updateQuestion(object) {
    if (isPosting == true) {
        return
    }

    let req = new XMLHttpRequest();

	req.withCredentials = true;

    req.addEventListener("load", function(evt){
		if(req.status === 200){
			window.location.href = "teacherquestions.html?classId=" + getQueryVariable("classId");
		} else {
			alert("An internal server error has occured."); //other server error. No handler on UI
		}
	});

	// Open the request to the register route, and then send the data in the body
	req.open("PUT", "/api/v1/teachers/questions/" + qId);

	req.send(JSON.stringify(object));
    isPosting = true;
}

function postNewQuestion(object) {
    if (isPosting == true) {
        return
    }

    let req = new XMLHttpRequest();

	req.withCredentials = true;

    req.addEventListener("load", function(evt){
		if (req.status === 200) {
			window.location.href = "teacherquestions.html?classId=" + getQueryVariable("classId");
		} else {
			alert("An internal server error has occured."); //other server error. No handler on UI
		}
	});

	// Open the request to the register route, and then send the data in the body
	req.open("POST", "/api/v1/teachers/classes/" + getQueryVariable("classId"));

	req.send(JSON.stringify(object));

    isPosting = true;
}