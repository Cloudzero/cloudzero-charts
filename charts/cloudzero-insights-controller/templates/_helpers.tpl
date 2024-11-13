{{/*
Expand the name of the chart.
*/}}
{{- define "insights-controller.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "insights-controller.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Generate certificates for the webhook server
*/}}
{{- define "insights-controller.generate-certs" -}}
{{- $altNames := list ( printf "DNS:webhook-server.%s" ( .Release.Namespace )) -}}
{{- $ca := genCA "insights-controller-ca" 365 -}}
{{- $cert := genSignedCert ( include "insights-controller.name" . ) nil $altNames 365 $ca -}}
tls.crt: {{ $cert.Cert | b64enc }}
tls.key: {{ $cert.Key | b64enc }}
{{- end -}}


{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "insights-controller.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "insights-controller.labels" -}}
helm.sh/chart: {{ include "insights-controller.chart" . }}
{{ include "insights-controller.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "insights-controller.selectorLabels" -}}
app.kubernetes.io/name: {{ include "insights-controller.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "insights-controller.annotations" -}}
{{- if .Values.webhook.annotations }}
{{ toYaml .Values.webhook.annotations }}
{{- end }}
{{- if and .Values.webhook.certificate.enabled .Values.webhook.issuer.enabled }}
cert-manager.io/inject-ca-from: {{ .Values.webhook.caInjection | default (printf "%s/%s" .Release.Namespace (include "insights-controller.certificateName" .)) }}
{{- end }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "insights-controller.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "insights-controller.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Name for the webhook server deployment
*/}}
{{- define "insights-controller.deploymentName" -}}
{{- printf "%s-server" (include "insights-controller.fullname" .) }}
{{- end }}

{{/*
Name for the webhook server service
*/}}
{{- define "insights-controller.serviceName" -}}
{{- printf "%s-svc" (include "insights-controller.fullname" .) }}
{{- end }}

{{/*
Name for the validating webhook configuration resource
*/}}
{{- define "insights-controller.validatingWebhookConfigName" -}}
{{- printf "%s-webhook" (include "insights-controller.fullname" .) }}
{{- end }}

{{/*
Name for the validating webhook
*/}}
{{- define "insights-controller.validatingWebhookName" -}}
{{- printf "%s.%s.svc" (include "insights-controller.validatingWebhookConfigName" .) .Release.Namespace }}
{{- end }}

{{/*
Name for the certificate resource
*/}}
{{- define "insights-controller.certificateName" -}}
{{- printf "%s-certificate" (include "insights-controller.fullname" .) }}
{{- end }}

{{/*
Name for the certificate secret
*/}}
{{- define "insights-controller.tlsSecretName" -}}
{{- printf "%s-tls" (include "insights-controller.fullname" .) }}
{{- end }}

{{/*
Name for the cloudzero API key secret name
*/}}
{{- define "insights-controller.cloudzeroSecretName" -}}
{{- if and .Values.global .Values.global.existingSecretName }}
{{- .Values.global.existingSecretName }}
{{- else}}
{{- .Values.existingSecretName | default (printf "%s-api-key" .Release.Name) }}
{{- end }}
{{- end }}

{{/*
Name for the cloudzero specific configuration file
*/}}
{{- define "insights-controller.cloudzeroConfigMapName" -}}
{{- if and .Values.global .Values.global.cloudzeroConfigMapName }}
{{- tpl .Values.global.cloudzeroConfigMapName . }}
{{- end }}
{{- end }}

{{/*
Name for the webhook server configuration file
*/}}
{{- define "insights-controller.configMapName" -}}
{{- printf "%s-configuration" (include "insights-controller.fullname" .) }}
{{- end }}

{{/*
Mount path for the cloudzero configuration file
*/}}
{{- define "insights-controller.cloudzeroConfigurationMountPath" -}}
{{- printf "/etc/%s-cloudzero-config" .Chart.Name  }}
{{- end }}

{{/*
Mount path for the insights server configuration file
*/}}
{{- define "insights-controller.serverConfigurationMountPath" -}}
{{- printf "/etc/%s-insights-server-config" .Chart.Name  }}
{{- end }}

{{/*
Name for the issuer resource
*/}}
{{- define "insights-controller.issuerName" -}}
{{- printf "%s-issuer" (include "insights-controller.fullname" .) }}
{{- end }}

{{/*
Name for the job resource
*/}}
{{- define "insights-controller.initJobName" -}}
{{- printf "%s-init" (include "insights-controller.fullname" .) }}
{{- end }}

{{/*
Name for the clusterrole resource
*/}}
{{- define "insights-controller.clusterRoleName" -}}
{{- printf "%s" (include "insights-controller.fullname" .) }}
{{- end }}

{{/*
Name for the clusterrolebinding resource
*/}}
{{- define "insights-controller.clusterRoleBindingName" -}}
{{- printf "%s" (include "insights-controller.fullname" .) }}
{{- end }}

{{/*
Name for the clusterrolebinding resource
*/}}
{{- define "insights-controller.cloudAccountId" -}}
{{- if .Values.global.cloudAccountId }}
{{ tpl .Values.global.cloudAccountId . }}
{{- else}}
{{ .Values.cloudAccountId . }}
{{- end }}
{{- end }}
