{{/*
Expand the name of the chart.
*/}}
{{- define "obs-tools-usage.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "obs-tools-usage.fullname" -}}
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
Create chart name and version as used by the chart label.
*/}}
{{- define "obs-tools-usage.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "obs-tools-usage.labels" -}}
helm.sh/chart: {{ include "obs-tools-usage.chart" . }}
{{ include "obs-tools-usage.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "obs-tools-usage.selectorLabels" -}}
app.kubernetes.io/name: {{ include "obs-tools-usage.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "obs-tools-usage.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "obs-tools-usage.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Service name for a specific service
*/}}
{{- define "obs-tools-usage.serviceName" -}}
{{- printf "%s-%s" (include "obs-tools-usage.fullname" .) .serviceName }}
{{- end }}

{{/*
Service labels for a specific service
*/}}
{{- define "obs-tools-usage.serviceLabels" -}}
{{ include "obs-tools-usage.labels" . }}
app.kubernetes.io/component: {{ .serviceName }}
{{- end }}

{{/*
Service selector labels for a specific service
*/}}
{{- define "obs-tools-usage.serviceSelectorLabels" -}}
{{ include "obs-tools-usage.selectorLabels" . }}
app.kubernetes.io/component: {{ .serviceName }}
{{- end }}

{{/*
Image name for a specific service
*/}}
{{- define "obs-tools-usage.image" -}}
{{- $registry := .Values.image.registry | default .Values.global.imageRegistry }}
{{- if $registry }}
{{- printf "%s/%s:%s" $registry .Values.image.repository .Values.image.tag }}
{{- else }}
{{- printf "%s:%s" .Values.image.repository .Values.image.tag }}
{{- end }}
{{- end }}

{{/*
Environment variables for database connections
*/}}
{{- define "obs-tools-usage.postgresEnv" -}}
- name: DB_HOST
  value: {{ include "obs-tools-usage.fullname" . }}-postgresql
- name: DB_PORT
  value: "5432"
- name: DB_USER
  value: "postgres"
- name: DB_PASSWORD
  valueFrom:
    secretKeyRef:
      name: {{ include "obs-tools-usage.fullname" . }}-postgresql
      key: postgres-password
- name: DB_NAME
  value: "product_service"
- name: DB_SSL_MODE
  value: "disable"
{{- end }}

{{- define "obs-tools-usage.redisEnv" -}}
- name: REDIS_HOST
  value: {{ include "obs-tools-usage.fullname" . }}-redis-master
- name: REDIS_PORT
  value: "6379"
- name: REDIS_PASSWORD
  value: ""
- name: REDIS_DB
  value: "0"
{{- end }}

{{- define "obs-tools-usage.mariadbEnv" -}}
- name: DB_HOST
  value: {{ include "obs-tools-usage.fullname" . }}-mariadb
- name: DB_PORT
  value: "3306"
- name: DB_USER
  value: "payment"
- name: DB_PASSWORD
  valueFrom:
    secretKeyRef:
      name: {{ include "obs-tools-usage.fullname" . }}-mariadb
      key: mariadb-password
- name: DB_NAME
  value: "payment_service"
- name: DB_SSL_MODE
  value: "false"
{{- end }}

{{- define "obs-tools-usage.kafkaEnv" -}}
- name: KAFKA_BROKERS
  value: {{ include "obs-tools-usage.fullname" . }}-kafka:9092
{{- end }}

{{/*
Service URLs
*/}}
{{- define "obs-tools-usage.productServiceUrl" -}}
{{- printf "http://%s-product-service:8080" (include "obs-tools-usage.fullname" .) }}
{{- end }}

{{- define "obs-tools-usage.basketServiceUrl" -}}
{{- printf "http://%s-basket-service:8081" (include "obs-tools-usage.fullname" .) }}
{{- end }}

{{- define "obs-tools-usage.paymentServiceUrl" -}}
{{- printf "http://%s-payment-service:8082" (include "obs-tools-usage.fullname" .) }}
{{- end }}
