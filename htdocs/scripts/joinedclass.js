let prestring = ""
let globalClassCode = ""
let className = ""

//A component for a single question
Vue.component('single-question', {
    props: ['question','response', 'index'],
    data: function () {
        vm = this
        return {
            answers: this.question.answers,
            show: true
        }
    },
    methods: {
        checkAns: function (event) {
            vm = this
            if (vm.question.answers === undefined || vm.question.answers.length == 0){
                vm.show = false
            } else{
                vm.show = true
            }
        }
    },
    template: `<div class="question mb-3">
                <div class="card">
                    <h5 class="card-header">Question #{{index + 1}}</h5>
                    <div class="card-body">
                        <p class="card-text">{{question.questionText}}</p>
                        <a v-if="show" class="btn btn-primary" v-on:click="checkAns"data-toggle="collapse" v-bind:href="'#answerArea' + question.questionId" role="button" aria-expanded="false" aria-controls="answerArea">Open/Close Answers</a>
                    </div>
                <answer-area v-if="show" v-bind:id="'answerArea' + question.questionId" v-bind:answers="question.answers" v-bind:response="response"></answer-area>
                </div>
            </div>`   
})

//A component for the answer-area
Vue.component('answer-area', {
    data: function() {
        return {
            show: false,
            radioval: "",
            responseId: 0 
        }
    },
    props: ['answers', 'response'],
    methods: {
        saveRequest: function (event) {
            vm = this
            let inResp = false
            if (this.radioval !== "" && vm.responseId === 0) {
                $.ajax({
                    type: "POST",
                    url: prestring + "/api/v1/questions/" + this.radioval.questionId + "/answers/" + this.radioval.answerId,
                    dataType: 'json',
                    success: function (json) {
                        vm.responseId = json.responseId
                        vm.response.push(json)
                    }
                })
            } else if (this.radioval !== "" && vm.responseId !== 0) {
                $.ajax({
                    type: "PUT",
                    url: prestring + "/api/v1/responses/" + vm.responseId + "/answers/" + this.radioval.answerId,
                    dataType: 'json',
                    success: function (json) {
                        
                    }
                })
            } else {
                alert("Please select an answer before submitting")
            }
            toastr["success"]("Response saved")
        }
    },
    template: `<div class="collapse" id="">
                <div class="card card-body answer-card">
                    <div class="radio" v-for="answer in answers" v-bind:key="answer.answerId" v-bind:answer="answer">
                    <label><input type="radio" v-model="radioval" v-bind:response=0 v-bind:test="answer.questionId" v-bind:value="answer" v-bind:name="'answerRadio' + answer.questionId" v-bind:key="answer.answerId">{{answer.answerText}}</label>
                    <hr>
                    </div>
                    <button v-on:click="saveRequest" class="btn btn-primary answerBtn">Save Answer</button>
                </div>
            </div>`

})

//Vue instance object
new Vue ({
    el: '#question-area',
    data: {
        classCode: "",
        className: "",
        questions: [],
        response: [],
        index: 0,
        timer: ''
    },
    //Once the instance creates itself
    created: function(){
        curWin = window.location.href
        let vm = this
        prestring = curWin.substring(0, curWin.lastIndexOf("/"));
        let classCode = this.getQueryVariable("classCode")
        if (classCode === "" || classCode === undefined && sessionStorage.getItem('classJson') !== null) {
            classJson = sessionStorage.getItem('classJson')
            let parsed = JSON.parse(classJson)
            classCode = parsed.classcode
            globalClassCode = parsed.classcode
            className = parsed.classname
            vm.classCode = parsed.classcode
            vm.className = parsed.classname
            this.getData(vm, classCode)
        } else if (classCode === "" || classCode === undefined && sessionStorage.getItem('classJson') === null) {
            window.location.replace(prestring)
            alert("Class does not exist")
        } else {
            $.getJSON(prestring + "/api/v1/classes/" + classCode, function (json) {
                globalClassCode = json.classcode
                className = json.classname
                vm.classCode = json.classcode
                vm.className = json.classname
                vm.getData(vm, classCode)
            }).fail(function () {
                window.location.replace(prestring)
                alert("Class does not exist")
            })
        }
    },
    //Methods that can be called
    methods: {
        //Gets initial data
        getData: function(vm, classCode) {
            $.getJSON(prestring + "/api/v1/classes/" + classCode + "/questions", function (json) {
                vm.questions = json
                vm.questions.forEach(question => {
                    $.ajax({
                        url: prestring + "/api/v1/questions/" + question.questionId + "/answers",
                        dataType: 'json',
                        success: function (json) {
                            Vue.set(question, 'answers', json);
                        }
                    })
                })
            }).fail(function(){
                toastr["error"]("No questions available yet")
            })
        },
        //taken from extras.js to use inside Vue instance
        //Gets variable out of uri
        getQueryVariable: function(variable){
            var query = window.location.search.substring(1);
            var vars = query.split('&');
            for (var i = 0; i < vars.length; i++) {
                var pair = vars[i].split('=');
                if (decodeURIComponent(pair[0]) == variable) {
                    return decodeURIComponent(pair[1]);
                }
            }
        },
        //Refreshes the questions
        refreshData: function () {
            vm = this
            classCode = this.classCode
            $.getJSON(prestring + "/api/v1/classes/" + classCode + "/questions", function (json) {
                json.forEach(question => {
                    var inarr = false
                    for (var i = 0; i < vm.questions.length; i++) {
                        if (vm.questions[i].questionId == question.questionId) {
                            inarr = true
                        }
                    }
                    if (!inarr) {
                        (function(vm){$.ajax({
                            url: prestring + "/api/v1/questions/" + question.questionId + "/answers",
                            dataType: 'json',
                            success: function (json) {
                                question["answers"] = json
                                vm.questions.push(question)
                            }
                        })
                    })(vm)
                    }
                })
            }).fail(function () {
                toastr["error"]("No questions available yet")
            })
        },
    }
})