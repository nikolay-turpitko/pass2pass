{{if .URL}}
[{{.Name}}]({{.URL}})
{{range $k, $v := extractMeta "extra" .Extra | mapKeys toLower}}
{{if eq $k "note" "notes" "description" }}_**{{$k}}:** {{$v}}_{{end}}
{{end}}
{{end}}
