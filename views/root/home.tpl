
<table>
    <tr>
        <th>Jugador</th>
        <th>Mínim</th>
        <th>Màxim</th>
        <th>Mitjana</th>
        {{range $id, $e := .Values}}
            <th>{{inc $id}}</th>
        {{end}}
    </tr>
    {{range .Players}}
        <tr>
            <td><a href="/players/{{.Id}}">{{.Name}}</a></td>
            <td>{{.Min}}</td>
            <td>{{.Max}}</td>
            <td>{{.Avg}}</td>
            {{range .Values}}
                <td>{{.}}</td>
            {{end}}
        </tr>
    {{end}}
</table>

