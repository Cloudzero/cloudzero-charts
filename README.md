# CloudZero Helm Charts

[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-2.1-4baaaa.svg)](CODE-OF-CONDUCT.md)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
![GitHub release](https://img.shields.io/github/release/Cloudzero/cloudzero-charts.svg)

This repository contains helm charts for use by CloudZero users, which can be installed into any cloud-hosted Kubernetes cluster.

## Table of Contents

- [Documentation](#documentation)
- [Installation](#installation)
- [Testing](#testing)
- [Contributing](#contributing)
- [Support + Feedback](#support--feedback)
- [Vulnerability Reporting](#vulnerability-reporting)
- [What is CloudZero?](#what-is-cloudzero)
- [License](#license)

## Documentation

Detailed documentation of each helm chart should be available within each chart README.md. Other documentation that may be helpful:
- [CloudZero Docs](https://docs.cloudzero.com/) for general information on CloudZero.
- [Helm Documentation](https://helm.sh/) for information on what Helm is, and how it is used to install artifacts on Kubernetes clusters.
- [Kubernetes Documentation](https://kubernetes.io/docs/home/) for information on Kubernetes itself.

## Installation

The helm charts in this repository generally assume the use of Helm v3 for installation. More detailed installation instructions are located within the README of each chart. See the [official Helm documentation](https://helm.sh/docs/intro/install/) for instructions on installing Helm v3.

After the `helm` command is available, charts should be installable with the `install` or `upgrade` command:
```bash
helm upgrade --install <RELEASE_NAME> cloudzero/cloudzero-agent \
    --set existingSecretName=<NAME_OF_SECRET> \
    --set clusterName=<CLUSTER_NAME> \
    --set-string cloudAccountId=<CLOUD_ACCOUNT_ID> \
    --set region=<REGION>
```

Installation can also be managed by deployment tools such as ArgoCD or Spinnaker if desired, but installation instructions in this repository assume the use of the Helm CLI.


### Mandatory Values

There are several mandatory values that must be specified for the chart to install properly. Below are the required settings along with strategies for providing custom values during installation:

| Key               | Type   | Default               | Description                                                                                                             |
|-------------------|--------|-----------------------|-------------------------------------------------------------------------------------------------------------------------|
| cloudAccountId    | string | `nil`                 | Account ID in AWS or Subscription ID in Azure where the cluster is running. Must be a string due to Helm limitations.  |
| clusterName       | string | `nil`                 | Name of the cluster. Must be RFC 1123 compliant.                                                                         |
| host              | string | `"api.cloudzero.com"` | CloudZero host to send metrics to.                                                                                      |
| apiKey            | string | `nil`                 | The CloudZero API key to use for exporting metrics. Only used if `existingSecretName` is not set.                       |
| existingSecretName| string | `nil`                 | Name of the secret that contains the CloudZero API key. Required if not providing the API key via `apiKey`.             |
| region            | string | `nil`                 | Region where the cluster is running (e.g., `us-east-1`, `eastus`). For more information, see AWS or Azure documentation. |



#### Overriding Default Values

Default values are specified in the chart's `values.yaml` file. If you need to change any of these values, it is recommended to create a `values-override.yaml` file for your customizations.

##### Using the `--values` Flag

You can use the `--values` (or short form `-f`) flag in your Helm commands to override values in the chart with a new file. Specify the name of the file after the `--values` flag:

```console
helm install <RELEASE_NAME> cloudzero/cloudzero-agent \
    --set existingSecretName=<NAME_OF_SECRET> \
    --set clusterName=<CLUSTER_NAME> \
    --set-string cloudAccountId=<CLOUD_ACCOUNT_ID> \
    --set region=<REGION> \
    -f values-override.yaml
```

Ensure `values-override.yaml` contains only the values you wish to override from `values.yaml`.

> Note it is possible to save values for different environments, or based on other criteria into seperate values files and multiple files using the `-f` helm parameters.

##### Using the `--set` Flag

You can use the `--set` flag in Helm commands to directly set or override specific values from `values.yaml`. Use dot notation to specify nested values:

```console
helm install <RELEASE_NAME> cloudzero/cloudzero-agent \
    --set existingSecretName=<NAME_OF_SECRET> \
    --set clusterName=<CLUSTER_NAME> \
    --set-string cloudAccountId=<CLOUD_ACCOUNT_ID> \
    --set region=<REGION> \
    --set server.resources.limits.memory=2048Mi \
    -f values-override.yaml
```

### Secret Management

The chart requires a CloudZero API key to send metric data. Admins can retrieve API keys [here](https://app.cloudzero.com/organization/api-keys).

The API key can be supplied as an existing secret (default) or created by the chart. Ensure the Secret is in the same namespace as the chart and follows this format:

**values-override.yaml**
```yaml
data:
  value: <API_KEY>
```

Example of creating a secret:
```console
kubectl create secret -n example-namespace generic example-secret-name --from-literal=value=<example-api-key-value>
```

The secret can then be used with `existingSecretName`.

### Memory Sizing

Please see the [sizing guide](./docs/sizing-guide.md) in the docs directory.

### Metric Exporters

This chart depends on metrics from [kube-state-metrics](https://github.com/kubernetes/kube-state-metrics) and [node-exporter](https://github.com/prometheus/node_exporter) projects as subcharts.

By default, these subcharts are disabled to allow scraping from existing instances. Configure the `cloudzero-agent` to use existing service endpoint addresses in `values.yaml`:

```yaml
validator:
  serviceEndpoints:
     kubeStateMetrics: <kube-state-metrics>.<example-namespace>.svc.cluster.local:8080
     prometheusNodeExporter: <node-exporter>.<example-namespace>.svc.cluster.local:9100
```

Alternatively, deploy them automatically by enabling settings in `values-override.yaml`:

```yaml
kube-state-metrics:
  enabled: true
prometheus-node-exporter:
  enabled: true
```

#### Passing Values to Subcharts

Values can be passed to subcharts like [kube-state-metrics](https://github.com/prometheus-community/helm-charts/blob/main/charts/kube-state-metrics/values.yaml) and [prometheus-node-exporter](https://github.com/prometheus-community/helm-charts/blob/main/charts/prometheus-node-exporter/values.yaml) by adding entries in `values-override.yaml` as per their specifications.

A common addition may be to pull the container images from custom image registries/repositories:

`values-override.yaml`
```yaml
kube-state-metrics:
  enabled: true
  image:
    registry: my-custom-registry.io
    repository: my-custom-kube-state-metrics/kube-state-metrics

prometheus-node-exporter:
  enabled: true
  image:
    registry: my-custom-registry.io
    repository: my-custom-prometheus/node-exporter
```

### Custom Scrape Configs

If running without the default exporters, adjust Prometheus scrape configs:

`values-override.yaml`
```yaml
prometheusConfig:
  scrapeJobs:
    kubeStateMetrics:
      enabled: false # this disables the default kube-state-metrics scrape job, which will be replaced by an entry in additionalScrapeJobs
    additionalScrapeJobs:
    - job_name: custom-kube-state-metrics
      honor_timestamps: true
      scrape_interval: 1m
      scrape_timeout: 10s
      metrics_path: /metrics
      static_configs:
        - targets:
          - 'my-kube-state-metrics-service.default.svc.cluster.local:8080'
          - 'my-node-exporter.default.svc.cluster.local:9100'
      relabel_configs:
      - separator: ;
        regex: __meta_kubernetes_service_label_(.+)
        replacement: $1
        action: labelmap
      - source_labels: [__meta_kubernetes_namespace]
        separator: ;
        regex: (.*)
        target_label: namespace
        replacement: $1
        action: replace
      - source_labels: [__meta_kubernetes_service_name]
        separator: ;
        regex: (.*)
        target_label: service
        replacement: $1
        action: replace
      - source_labels: [__meta_kubernetes_pod_node_name]
        separator: ;
        regex: (.*)
        target_label: node
        replacement: $1
        action: replace
      kubernetes_sd_configs:
        - role: endpoints
          kubeconfig_file: ""
          follow_redirects: true
          enable_http2: true
```

### Exporting Pod Labels

Pod labels can be exported as metrics using `kube-state-metrics`. Customize the labels to export in `values-override.yaml`:

Note a subset of relevant pod labels can be included -- for example only exporting the pod labels named `foo` and `bar` -- can be achieved with the following:

```yaml
kube-state-metrics:
  extraArgs:
    - --metric-labels-allowlist=pods=[foo,bar]
```

> This is preferable to including all labels with `*` because the performance and memory impact is reduced. Regular expression matching is not currently supported. See the `kube-state-metrics` [documentation](https://github.com/kubernetes/kube-state-metrics/blob/main/docs/developer/cli-arguments.md) for more details.


## Dependencies

| Repository                                         | Name                     | Version |
|----------------------------------------------------|--------------------------|---------|
| https://prometheus-community.github.io/helm-charts | kube-state-metrics       | 5.15.*  |
| https://prometheus-community.github.io/helm-charts | prometheus-node-exporter | 4.24.*  |

## Useful References

- [Memory Sizing Guide](./docs/sizing-guide.md)
- [Deployment Validation Guide](./docs/deploy-validation.md)

## Contributing

We appreciate feedback and contribution to this repo! Before you get started, please see the following:

- [This repo's contribution guide](GENERAL-CONTRIBUTING.md)

## Support + Feedback

Contact support@cloudzero.com for usage, questions, specific cases. See the [CloudZero Docs](https://docs.cloudzero.com/) for general information on CloudZero.

## Vulnerability Reporting

Please do not report security vulnerabilities on the public GitHub issue tracker. Email [security@cloudzero.com](mailto:security@cloudzero.com) instead.

## What is CloudZero?

CloudZero is the only cloud cost intelligence platform that puts engineering in control by connecting technical decisions to business results.:

- [Cost Allocation And Tagging](https://www.cloudzero.com/tour/allocation) Organize and allocate cloud spend in new ways, increase tagging coverage, or work on showback.
- [Kubernetes Cost Visibility](https://www.cloudzero.com/tour/kubernetes) Understand your Kubernetes spend alongside total spend across containerized and non-containerized environments.
- [FinOps And Financial Reporting](https://www.cloudzero.com/tour/finops) Operationalize reporting on metrics such as cost per customer, COGS, gross margin. Forecast spend, reconcile invoices and easily investigate variance.
- [Engineering Accountability](https://www.cloudzero.com/tour/engineering) Foster a cost-conscious culture, where engineers understand spend, proactively consider cost, and get immediate feedback with fewer interruptions and faster and more efficient innovation.
- [Optimization And Reducing Waste](https://www.cloudzero.com/tour/optimization) Focus on immediately reducing spend by understanding where we have waste, inefficiencies, and discounting opportunities.

Learn more about [CloudZero](https://www.cloudzero.com/) on our website [www.cloudzero.com](https://www.cloudzero.com/)

## License

This project is licenced under the Apache 2.0 [LICENSE](LICENSE).