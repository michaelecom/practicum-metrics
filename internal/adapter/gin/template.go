package gin

import "html/template"

// metricsTemplate - HTML шаблон для отображения всех метрик
const metricsTemplate = `<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Metrics Dashboard</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            max-width: 1200px;
            margin: 0 auto;
            padding: 20px;
            background-color: #f5f5f5;
        }
        h1 {
            color: #333;
            border-bottom: 3px solid #4CAF50;
            padding-bottom: 10px;
        }
        h2 {
            color: #555;
            margin-top: 30px;
        }
        table {
            width: 100%;
            border-collapse: collapse;
            background-color: white;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
            margin-bottom: 20px;
        }
        th {
            background-color: #4CAF50;
            color: white;
            padding: 12px;
            text-align: left;
        }
        td {
            padding: 10px 12px;
            border-bottom: 1px solid #ddd;
        }
        tr:hover {
            background-color: #f5f5f5;
        }
        .no-data {
            color: #999;
            font-style: italic;
            padding: 20px;
        }
        .metric-type {
            display: inline-block;
            padding: 4px 8px;
            border-radius: 4px;
            font-size: 12px;
            font-weight: bold;
        }
        .gauge-badge {
            background-color: #2196F3;
            color: white;
        }
        .counter-badge {
            background-color: #FF9800;
            color: white;
        }
    </style>
</head>
<body>
    <h1>📊 Metrics Dashboard</h1>

    <h2><span class="metric-type gauge-badge">GAUGE</span> Gauges</h2>
    {{if .Gauges}}
    <table>
        <thead>
            <tr>
                <th>Name</th>
                <th>Value</th>
            </tr>
        </thead>
        <tbody>
            {{range $name, $value := .Gauges}}
            <tr>
                <td>{{$name}}</td>
                <td>{{$value}}</td>
            </tr>
            {{end}}
        </tbody>
    </table>
    {{else}}
    <p class="no-data">No gauge metrics available</p>
    {{end}}

    <h2><span class="metric-type counter-badge">COUNTER</span> Counters</h2>
    {{if .Counters}}
    <table>
        <thead>
            <tr>
                <th>Name</th>
                <th>Value</th>
            </tr>
        </thead>
        <tbody>
            {{range $name, $value := .Counters}}
            <tr>
                <td>{{$name}}</td>
                <td>{{$value}}</td>
            </tr>
            {{end}}
        </tbody>
    </table>
    {{else}}
    <p class="no-data">No counter metrics available</p>
    {{end}}
</body>
</html>`

// tmpl - предкомпилированный шаблон
var tmpl = template.Must(template.New("metrics").Parse(metricsTemplate))
