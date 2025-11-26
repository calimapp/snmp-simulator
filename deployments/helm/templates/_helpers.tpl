

{{/* Common labels */}}
{{- define "common.labels" -}}
helm.sh/chart: {{ .Chart.Name }}
app.kubernetes.io/name: {{ .Chart.Name }}
app.kubernetes.io/version: {{ .Chart.AppVersion }}
app.kubernetes.io/part-of: {{ .Chart.Name }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}
