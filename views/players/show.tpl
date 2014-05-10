
{{if .logged}}

<form action="/players/{{.id}}" method="POST" autocomplete="off" accept-charset="utf-8">
    <input type="text" name="name" autofocus placeholder="Nom" value="{{.name}}" />
    <input type="submit" value="Canviar nom" />
</form>

<form action="/players/{{.id}}/delete" method="POST" autocomplete="off" accept-charset="utf-8">
    <input type="submit" value="Borrar jugador" />
</form>

{{else}}

<form action="/players/{{.id}}/rate" method="POST" autocomplete="off" accept-charset="utf-8">
    <input type="number" name="rating" autofocus placeholder="Puntuació" />
    <input type="submit" value="Puntuar" />
</form>

{{end}}

