
{{if .Players}}

<div class="force-container">
    <table>
        <tr>
            <th>Jugador</th>
            {{if len .Values}}
                <th>Mínim</th>
                <th>Màxim</th>
                <th>Mitjana</th>
            {{end}}
            {{range $id, $e := .Values}}
                <th>{{inc $id}}</th>
            {{end}}
        </tr>
        {{range .Players}}
            <tr>
                <td><a href="/players/{{.Id}}">{{.Name}}</a></td>
                {{if .Ratings}}
                    <td>{{.Min}}</td>
                    <td>{{.Max}}</td>
                    <td>{{.Avg}}</td>
                {{end}}
                {{range .Ratings}}
                    <td>{{.Value}}</td>
                {{end}}
            </tr>
        {{end}}
    </table>
</div>

{{else}}

<div class="force-empty">
    <span class="empty">No hi ha cap jugador que hagi valorat entrenaments.</span>
</div>

{{end}}

