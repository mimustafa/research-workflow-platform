{{- define "research-workflow.name" -}}
{{- .Chart.Name | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "research-workflow.fullname" -}}
{{- printf "%s-%s" .Release.Name (include "research-workflow.name" .) | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "research-workflow.labels" -}}
app.kubernetes.io/name: {{ include "research-workflow.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
helm.sh/chart: {{ printf "%s-%s" .Chart.Name .Chart.Version | quote }}
{{- end }}
