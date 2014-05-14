
{{if .Error}}

<div class="dialog">
    <div class="dialog-header">
        <h1>Error !</h1>
    </div>
    <div class="dialog-body dialog-body-text">
        <span>No hem pogut agafar la teva puntuació!</span>
        <br />
        <span>Tingues en compte que la teva puntuació ha de ser un nombre entre 0 i 10.</span>

        <br />
        <br />

        <span>Pots tornar a provar de puntuar <a href="/players/{{.Id}}">aquí</a>.</span>
    </div>
</div>

{{else}}

<div class="dialog">
    <div class="dialog-header">
        <h1>Ho tenim !</h1>
    </div>
    <div class="dialog-body dialog-body-text">
        <span>Gràcies per puntuar aquest entrenament!</span>
        <br />
        <br />
        <span>Pots tornar a la teva pàgina <a href="/players/{{.Id}}">aquí</a>.</span>
    </div>
</div>

{{end}}

