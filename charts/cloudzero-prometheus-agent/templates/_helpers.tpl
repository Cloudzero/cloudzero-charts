{{/* Define the url to which metrics should be sent */}}
{{- define "cloudzero-prometheus-agent.remoteWriteUrl" -}}
https://{{ .Values.global.cloudzeroHost }}/v1/container-metrics?cluster_name={{.Values.global.cluster_name}}&cloud_account_id={{.Values.global.cloud_account_id}}
{{- end -}}

{{ define "cloudzero-prometheus-agent.secretName" -}}
{{ .Values.global.secretNameOverride | default (printf "%s-api-key" .Release.Name) }}
{{- end}}
