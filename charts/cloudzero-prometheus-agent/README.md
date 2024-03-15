# Cloudzero Prometheus Agent

[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-2.1-4baaaa.svg)](CODE-OF-CONDUCT.md)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
![GitHub release](https://img.shields.io/github/release/cloudzero/template-cloudzero-open-source.svg)

A chart for deploying Prometheus in agent mode for the purpose of sending cluster metrics to the CloudZero platform.

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

```console
helm install [RELEASE_NAME] cloudzero/cloudzero-prometheus-agent \
    --set api_key=CLOUDZERO_API_KEY \
    --set global.cluster_name=CLUSTER_NAME \
    --set global.cloud_account_id=CLOUD_ACCOUNT_ID \
```

## Dependencies

By default this chart installs additional, dependent charts:

- [prometheus-community/kube-prometheus-stack](https://github.com/prometheus-community/helm-charts/tree/main/charts/kube-prometheus-stack)
- [prometheus-community/kube-state-metrics](https://github.com/prometheus-community/helm-charts/tree/main/charts/kube-state-metrics)
- [prometheus-community/prometheus-node-exporter](https://github.com/prometheus-community/helm-charts/tree/main/charts/prometheus-node-exporter)
- [grafana/grafana](https://github.com/grafana/helm-charts/tree/main/charts/grafana)
