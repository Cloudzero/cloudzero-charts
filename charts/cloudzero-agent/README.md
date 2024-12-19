# CloudZero Agent Helm Chart

[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-2.1-4baaaa.svg)](CODE-OF-CONDUCT.md)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
![GitHub release](https://img.shields.io/github/release/Cloudzero/cloudzero-charts.svg)

A Helm chart for deploying Prometheus in agent mode to send cluster metrics to the CloudZero platform.

For the latest release, see [Releases](https://github.com/Cloudzero/cloudzero-charts/releases). You can also [enable release notifications](#enabling-release-notifications).

## Prerequisites

- Kubernetes 1.23+
- Helm 3+
- A CloudZero API key
- Each Kubernetes cluster must have a route to the internet and a rule that allows egress from the agent to the CloudZero collector endpoint at https://api.cloudzero.com on port 443

## Installation

### Easy Install (Most teams will use this!)

To use the chart or a beta version, you must add the repository to Helm. Refer to the [`helm repo`](https://helm.sh/docs/helm/helm_repo/) documentation for command details.

#### 1. Install the Helm Chart

```console
helm repo add cloudzero https://cloudzero.github.io/cloudzero-charts
helm install cloudzero cloudzero/cloudzero-agent \
    --set apiKey=<YOUR_CLOUDZERO_API_KEY>\
    --set clusterName=<YOUR_CLUSTER_NAME>
```

### Advanced Install (Usually when you have specific security requirements.)

#### 1. Create and Configure a Values File

```yaml
# -- values.yaml

# -- clusterName is required to identify this cluster in the CloudZero dashboard.
clusterName: <YOUR_CLUSTER_NAME>

# -- apiKey is the CloudZero apiKey generated in the CloudZero platform.
apiKey: <YOUR_CLOUDZERO_API_KEY>

# -- Other values here...
```

Default values are specified in the chart's `values.yaml` file. Please reference this file for available override values.

### Memory Sizing

Please see the [sizing guide](./docs/sizing-guide.md) in the docs directory.

### Update Helm Chart
If you are updating an existing installation, pull the latest chart information:

```console
helm repo update
```

Next, upgrade the installation to the latest chart version:

```console
helm upgrade --install cloudzero cloudzero/cloudzero-agent
```

## Enabling Release Notifications

To receive a notification when a new version of the chart is [released](https://github.com/Cloudzero/cloudzero-charts/releases), you can [watch the repository](https://docs.github.com/en/account-and-profile/managing-subscriptions-and-notifications-on-github/setting-up-notifications/configuring-notifications#configuring-your-watch-settings-for-an-individual-repository):

1. Navigate to the [repository main page](https://github.com/Cloudzero/cloudzero-charts).
2. Select **Watch > Custom**.
3. Check the **Releases** box.
4. Select **Apply**.

## Useful References

- [Memory Sizing Guide](./docs/sizing-guide.md)
- [Deployment Validation Guide](./docs/deploy-validation.md)

---
