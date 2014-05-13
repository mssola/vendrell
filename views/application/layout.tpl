<!DOCTYPE html>
<html>
<head>
    <meta http-equiv="content-type" content="text/html; charset=UTF-8" />
    <title>Entrenaments</title>
    <link href="/css/style.css" rel="stylesheet" type="text/css" />
</head>
<body>
    {{if .LoggedIn}}
    <div id="header">
    </div>
    {{end}}

    <div class="main">
        {{ yield }}
    </div>
</body>
</html>
