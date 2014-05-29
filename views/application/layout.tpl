<!DOCTYPE html>
<html>
<head>
    <meta http-equiv="content-type" content="text/html; charset=UTF-8" />
    <title>Entrenaments</title>
    <meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1, user-scalable=0"/>
    <link href="/css/mobile.css" media="only screen and (max-width: 480px)" rel="stylesheet" type="text/css" />
    <link href="/css/style.css" media="only screen and (min-width: 480px)" rel="stylesheet" type="text/css" />
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
            <a href="/players">Jugadors</a>
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
