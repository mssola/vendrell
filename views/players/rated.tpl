
{{if .error}}

<span>No hem pogut agafar la teva puntuació!</span>
<br />
<span>Tingues en compte que la teva puntuació ha de ser un nombre entre 0 i 10.</span>

<br />

<span>Pots tornar a provar de puntuar <a href="/players/{{.id}}">aquí</a>.</span>

{{else}}

<span>Gràcies per puntuar aquest entrenament!</span>

{{end}}

