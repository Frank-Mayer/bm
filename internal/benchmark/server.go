package benchmark

import (
	"github.com/gorilla/websocket"

	"net"
	"net/http"
)

func handleAppScript(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/javascript")
	w.Header().Set("Cache-Control", "no-cache")
	if _, err := w.Write(appJs); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func handleIndex(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Cache-Control", "no-cache")
	if _, err := w.Write(html); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func handleStyle(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/css")
	w.Header().Set("Cache-Control", "no-cache")
	if _, err := w.Write(style); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func handleWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "could not upgrade connection", http.StatusInternalServerError)
		return
	}
	defer conn.Close()
	for {
		data := <-dataChan
		if err := conn.WriteJSON(data); err != nil {
			conn.Close()
			break
		}
	}
}

func StartServer() string {
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/ws", handleWs)
	http.HandleFunc("/app.js", handleAppScript)
	http.HandleFunc("/style.css", handleStyle)
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}
	go func() {
		if err := http.Serve(listener, nil); err != nil {
			panic(err)
		}
	}()
	return listener.Addr().String()
}

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	appJs = []byte(appJsStr)
	style = []byte(styleStr)
	html  = []byte(htmlStr)
)

const appJsStr = `const ws = new WebSocket("ws://" + location.host + "/ws");
function createChart(name, ...labels) {
    const data = {
        labels: [],
        datasets: labels.map(label => {
            const h = hash(label) % 360;
            return {
                label: label,
                data: [],
                borderColor: "hsl(" + h + ", 100%, 50%)",
                backgroundColor: "hsla(" + h + ", 100%, 50%, 0.3)",
            }       
        }),
    };
    const config = {
        type: "line",
        data: data,
        options: {
            responsive: false,
            plugins: {
                legend: {
                    display: true,
                },
                title: {
                    display: true,
                    text: name,
                },
            },
            scales: {
                xAxes: [{
                    display: false,
                    ticks: {
                        display: false,
                    },
                }],
                yAxes: [{
                    display: false,
                    ticks: {
                        display: false,
                    },
                }]
            },
        },
    };
    const canvas = document.createElement("canvas");
    canvas.classList.add("chart");
    canvas.classList.add(name);
    document.body.appendChild(canvas);
    const ctx = canvas.getContext("2d");
    return new Chart(ctx, config);
}
const cpuChart = createChart("CPU", "CPU usage");
const memChart = createChart("Memory", "Virtual memory size", "Resident memory size");
let i = 0;
let t = 0;
ws.onmessage = function (event) {
    if (i >= 100) {
        cpuChart.data.labels.shift();
        for (let i = 0; i < cpuChart.data.datasets.length; i++) {
            cpuChart.data.datasets[i].data.shift();
        }
        memChart.data.labels.shift();
        for (let i = 0; i < memChart.data.datasets.length; i++) {
            memChart.data.datasets[i].data.shift();
        }
    } else {
        i++;
    }
    const data = JSON.parse(event.data);
    cpuChart.data.labels.push(timeFormat(data.time));
    cpuChart.data.datasets[0].data.push(data.cpu);
    cpuChart.update();
    memChart.data.labels.push(timeFormat(data.time));
    memChart.data.datasets[0].data.push(data.vms);
    memChart.data.datasets[1].data.push(data.rss);
    memChart.update();
    document.title = data.title;
    if (t) {
        window.clearTimeout(t);
    }
    t = window.setTimeout(() => {
        document.body.classList.add("idle");
    }, 5000);
};
ws.onclose = function () {
    if (t) {
        window.clearTimeout(t);
    }
    document.body.classList.add("idle");
};
function hash(s) {
    for (var i = 0, h = 0; i < s.length; i++) {
        h = Math.imul(31, h) + s.charCodeAt(i) | 0;
    }
    return Math.abs(h);
}
function timeFormat(x) {
    console.log(x);
    return x;
}`

const styleStr = `:root {
    font-family: "Segoe UI", Tahoma, Geneva, Verdana, sans-serif;
    color-scheme: light dark;
}
body {
    margin: 0;
    padding: 5vh 5vw;
}
.chart {
    width: 90vmin;
    height: 50vh;
}
body.idle::before {
    content: "Idle";
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    display: block;
    text-align: center;
}
`

const htmlStr = `<!DOCTYPE html>
<html>
<head>
    <title>benchmark</title>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link rel="stylesheet" href="/style.css">
</head>
<body>
    <script src="https://cdn.jsdelivr.net/npm/chart.js@4.4.1/dist/chart.umd.min.js"></script>
    <script defer src="/app.js"></script>
</body>
</html>`
