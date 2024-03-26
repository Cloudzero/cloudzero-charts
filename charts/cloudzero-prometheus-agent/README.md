# cloudzero-agent

[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-2.1-4baaaa.svg)](CODE-OF-CONDUCT.md)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
![GitHub release](https://img.shields.io/github/release/cloudzero/template-cloudzero-open-source.svg)

A chart for using Prometheus in agent mode to send cluster metrics to the CloudZero platform.

## Prerequisites

- Kubernetes 1.19+
- Helm 3+
- A CloudZero API key

## Installation

### Get Helm Repository Info

```console
helm repo add cloudzero cloudzero/cloudzero-prometheus-agent
helm repo update
```

_See [`helm repo`](https://helm.sh/docs/helm/helm_repo/) for command documentation._

### Install Helm Chart

The chart can be installed directly with Helm or any other common Kubernetes deployment tool.

```console
helm install [RELEASE_NAME] cloudzero/cloudzero-prometheus-agent \
    --set api_key=CLOUDZERO_API_KEY \
    --set cloudzero.cluster_name=CLUSTER_NAME \
    --set cloudzero.cloud_account_id=CLOUD_ACCOUNT_ID \
```

### Secret Management

The chart requires an CloudZero API key in order to send metric data to the CloudZero platform. Admins can retrieve API keys [here](https://app.cloudzero.com/organization/api-keys).

The Deployment running Prometheus ingests the API key via a Secret; this Secret can be created by the chart (default), or an existing secret containing the API key can be specified.

If using a Secret external to this chart for the API key, ensure the Secret is created in the same namespace as the chart and that the Secret data follows the format:
   
```yaml
data:
  value: <API_KEY>
```

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| cloudzero.cloud_account_id | string | `nil` | Account ID of the account the cluster is running in. |
| cloudzero.cluster_name | string | `nil` | Name of the clusters. |
| cloudzero.host | string | `"api.cloudzero.com"` | CloudZero host to send metrics to. |
| credentials.createSecret | bool | `true` | If true, a secret containing the CloudZero API key will be created using the `api_key` value. |
| credentials.secretName | string | `""` | The name of the secret that contains the CloudZero API key. Required if createSecret is false. |

## Requirements

| Repository | Name | Version |
|------------|------|---------|
| https://prometheus-community.github.io/helm-charts | kube-state-metrics | 5.15.* |
| https://prometheus-community.github.io/helm-charts | prometheus-node-exporter | 4.24.* |