
{{if .LoggedIn}}

<div id="right">
    <div class="dialog">
        <div class="dialog-header">
            <h1>Canviar nom</h1>
        </div>
        <div class="dialog-body">
            <form action="/players/{{.Id}}" method="POST" autocomplete="off" accept-charset="utf-8">
                <input id="name" class="text" type="text" name="name" placeholder="Nom" value="{{.Name}}" />
                <input class="btn yellow-btn" type="submit" value="Canviar" />
                <div class="clearout"></div>
            </form>
        </div>
    </div>

    <div class="dialog">
        <div class="dialog-header">
            <h1>Borrar jugador</h1>
        </div>
        <div class="dialog-body">
            <form action="/players/{{.Id}}/delete" method="POST" autocomplete="off" accept-charset="utf-8">
                <span>Abans de borrar aquest jugador tingues en compte que aquesta
                acció <b>no</b> és reversible.</span>
                <p>Com a seguretat extra, si realment vols borrar aquest usuari
                hauràs d'escriure el seu nom una altra vegada.</p>
                <input id="rm-text" class="text" type="text" name="name" placeholder="Nom" />
                <input id="rm-btn" class="btn red-btn" type="submit" value="Borrar" disabled />
                <div class="clearout"></div>
            </form>
        </div>
    </div>
</div>

{{else}}

<div class="dialog">
    <div class="dialog-header">
        <h1>Canviar nom</h1>
    </div>
    <div class="dialog-body">
        <form action="/players/{{.Id}}/rate" method="POST" autocomplete="off" accept-charset="utf-8">
            <input class="text" type="number" name="rating" autofocus placeholder="Puntuació" />
            <input class="btn yellow-btn" type="submit" value="Puntuar" />
            <div class="clearout"></div>
        </form>
    </diV>
</diV>

{{end}}

