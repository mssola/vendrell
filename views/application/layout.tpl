<!DOCTYPE html>
<html>
<head>
    <meta http-equiv="content-type" content="text/html; charset=UTF-8" />
    <title>Entrenaments</title>
    <link href="/css/style.css" rel="stylesheet" type="text/css" />
    {{if .JS}}
    <script src="/js/vendrell.js"></script>
    {{end}}
</head>
<body>
    {{if .LoggedIn}}
    <header>
      <div class="wrapper">
        <div class="right">
          <div class="inner">
            <a href="/">Inici</a>
            {{if .Download}}
            <a href="{{.Download}}">Baixar-se CSV</a>
            {{end}}
            <a href="/players/new">Crear jugador</a>
            <form id="logout" action="/logout" method="POST">
                <input class="btn yellow-btn" type="submit" value="Surt" />
            </form>
          </div>
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
