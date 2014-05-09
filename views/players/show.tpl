
<form action="/players/{{.id}}" method="POST" autocomplete="off" accept-charset="utf-8">
    <input type="text" name="name" autofocus placeholder="Nom" value="{{.name}}" />
    <input type="submit" value="Canviar nom" />
</form>

<form action="/players/{{.id}}/delete" method="POST" autocomplete="off" accept-charset="utf-8">
    <input type="submit" value="Borrar jugador" />
</form>

