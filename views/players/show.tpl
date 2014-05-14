
{{if .LoggedIn}}

<div class="dialog">
    <div class="dialog-header">
        <h1>Canviar nom</h1>
    </div>
    <div class="dialog-body">
        <form action="/players/{{.Id}}" method="POST" autocomplete="off" accept-charset="utf-8">
            <input class="text" type="text" name="name" autofocus placeholder="Nom" value="{{.Name}}" />
            <input class="yellow-btn" type="submit" value="Canviar nom" />
        </form>
    </div>
</div>

<form action="/players/{{.Id}}/delete" method="POST" autocomplete="off" accept-charset="utf-8">
    <input type="submit" value="Borrar jugador" />
</form>

{{else}}

<div class="dialog">
    <div class="dialog-header">
        <h1>Canviar nom</h1>
    </div>
    <div class="dialog-body">
        <form action="/players/{{.Id}}/rate" method="POST" autocomplete="off" accept-charset="utf-8">
            <input class="text" type="number" name="rating" autofocus placeholder="Puntuació" />
            <input class="yellow-btn" type="submit" value="Puntuar" />
            <div class="clearout"></div>
        </form>
    </diV>
</diV>

{{end}}

