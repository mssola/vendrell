<!DOCTYPE html>
<html>
<head>
    <meta http-equiv="content-type" content="text/html; charset=UTF-8" />
    <title>Entrenaments</title>
    <link href="/css/style.css" rel="stylesheet" type="text/css" />
</head>
<body>
    {{if .LoggedIn}}
    <header>
      <div class="wrapper">
        <div class="right">
          <a href="/">Inici</a>
          <a href="/players/new">Crear jugador</a>
          <form id="logout" action="/logout" method="POST">
              <input class="btn yellow-btn" type="submit" value="Surt" />
          </form>
        </div>
      </div>
    </header>
    {{end}}

    <div id="content">
      <div class="wrapper">
        {{ yield }}
      </div>
    </div>
</body>
</html>
