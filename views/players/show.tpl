
{{if .LoggedIn}}

<div id="left">

{{if .One.Ratings}}

    <table>
        <tr>
            <th>Mínim</th>
            <th>Màxim</th>
            <th>Mitjana</th>
        </tr>
        <tr>
            <td>{{.One.Min}}</td>
            <td>{{.One.Max}}</td>
            <td>{{.One.Avg}}</td>
        </tr>
    </table>

    <table>
        <tr>
            <th>Entrenament</th>
            <th>Puntuació</th>
            <th>Data</th>
        </tr>
        {{range $idx, $rating := .One.Ratings}}
            <tr>
                <td>{{inc $idx}}</td>
                <td>{{$rating.Value}}</td>
                <td>{{fmtDate $rating.Created_at}}</td>
            </tr>
        {{end}}
    </table>

{{else}}

<span class="empty">Aquest jugador encara no ha valorat cap entrenament.</span>

{{end}}

</div>

<div id="right">
    <div class="dialog">
        <div class="dialog-header">
            <h1>Canviar nom</h1>
        </div>
        <div class="dialog-body">
            <form action="/players/{{.One.Id}}" method="POST" autocomplete="off" accept-charset="utf-8">
                <input id="name" class="text" type="text" name="name" placeholder="Nom" value="{{.One.Name}}" />
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
            <form action="/players/{{.One.Id}}/delete" method="POST" autocomplete="off" accept-charset="utf-8">
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
        <h1>Puntuar</h1>
    </div>
    <div class="dialog-body">
        <form action="/players/{{.One.Id}}/rate" method="POST" autocomplete="off" accept-charset="utf-8">
            <label for="rating">Si haguéssis de puntuar l'entrenament d'avui
            del 0 al 10 segons l'esforç que t'ha suposat, com el
            qualificaries ?</label>
            <input class="text" type="number" id="rating" name="rating" autofocus placeholder="Puntuació (0-10)" />
            <input class="btn yellow-btn" type="submit" value="Enviar" />
            <div class="clearout"></div>
        </form>
    </div>

    <div class="borg">
        <table>
            <tr>
                <th colspan="2">Com a referència</th>
            </tr>
            <tr><td>0</td><td class="second">Res</td></tr>
            <tr><td>0.5</td><td class="second">Molt, molt fàcil</td></tr>
            <tr><td>1</td><td class="second">Molt fàcil</td></tr>
            <tr><td>2</td><td class="second">Fàcil</td></tr>
            <tr><td>3</td><td class="second">Moderat</td></tr>
            <tr><td>4</td><td class="second">Una mica dur</td></tr>
            <tr><td>5</td><td class="second">Dur</td></tr>
            <tr><td>6</td><td class="second"></td></tr>
            <tr><td>7</td><td class="second">Molt dur</td></tr>
            <tr><td>8</td><td class="second"></td></tr>
            <tr><td>9</td><td class="second">Molt, molt dur</td></tr>
            <tr><td>10</td><td class="second">Impossible</td></tr>
        </table>
    </div>
</div>


{{end}}

