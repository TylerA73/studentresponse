const uri=  "/api/v1/teachers/questions/" + getQueryVariable("qId");


var res = new Vue({
    // 	Element containing Vue data
    el: '#question-title',
    data() {
        return {
            qtext: getQueryVariable("title"),
            isDisabledAttr: true
        }
    },
    //	Initializes the data upon page load
	created: function(){
		this.questions = this.fetchTable();
	},
    //	All methods useable by the Vue
    methods: {
        //	Fetch the classroom data for the classes table
        fetchTable: function(){
            $.getJSON(uri, function(json) {

                var total = json.question.count;
                var e = document.createElement('span');
                e.innerHTML = '(Number of responses: '+total+')';
                var ref = document.querySelector('#text-title');
                $( "#text-title" ).after( "<p>Number of responses: " + total + "</p>" );

                var all = [["answer","count"]];

                var ansCorrect = 0;
                var ansWrong = 0;
                //Iterate over results
                for (i = 0; i < json.answers.length; i++) { 
                    all.push([json.answers[i].answerText, json.answers[i].count]);
                    if (json.answers[i].iscorrect == true) {
                        ansCorrect = ansCorrect + json.answers[i].count
                    } else {
                        ansWrong = ansWrong + json.answers[i].count
                    }
                }
                console.log(ansCorrect, ansWrong);
                // Load google charts
                google.charts.load('current', {'packages':['corechart']});
                google.charts.setOnLoadCallback(drawChart);

                // Draw the chart and set the chart values
                function drawChart() {
                var data = google.visualization.arrayToDataTable(all);
                var data2 = google.visualization.arrayToDataTable([
                    ['Correct vs Incorrect', 'correct', 'incorrect' ],
                    ['', ansCorrect, ansWrong]
                  ]);


                // Optional; customize the chart
                var options = {
                    'width':500, 
                    'height':500,
                    'backgroundColor': '#69787b',
                    'legendTextStyle': { color: '#FFF' },
                    'titleTextStyle': { color: '#FFF' },
                    'hAxis': {
                        'textStyle':{ color: '#FFF' },
                    },
                    'sliceVisibilityThreshold':0
                };
                var options2 = {
                    'isStacked': 'percent',
                    'width':500, 
                    'height':500,
                    'legend': {position: 'top', maxLines: 3},
                    'hAxis': {
                      'minValue': 0,
                      'textStyle':{ color: '#FFF' },
                    },
                    'backgroundColor': '#69787b',
                    'legendTextStyle': { color: '#FFF' },
                    'titleTextStyle': { color: '#FFF' },
                    'sliceVisibilityThreshold':0
                };

                // Display the chart inside the <div> element with id="piechart"
                var chart = new google.visualization.PieChart(
                    document.getElementById('piechart'));
                chart.draw(data, options);

                chart = new google.visualization.BarChart(
                    document.getElementById("barchart"));
                chart.draw(data, options);

                chart = new google.visualization.BarChart(
                    document.getElementById("stackchart"));
                chart.draw(data2, options2);
                }


            });
        }
    }
});