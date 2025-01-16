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

### Quick Start
> ⚠️ By default, this chart exports only the app.kubernetes.io/component label from pods and namespaces. No annotations are exported. While this provides a safe default for demo purposes, it may be insufficient for your organization.
>
> **Recommendations:**
> * Configure additional labels to align with your organization's FinOps tagging practices.
* Review the [Labels and Annotations](#labels-and-annotations) section for guidance on exposing an expanded set of labels or annotations to meet your organization’s specific requirements.

#### 1. Add CloudZero Helm Repository
Refer to the [`helm repo`](https://helm.sh/docs/helm/helm_repo/) documentation for command details. To use a beta version, refer to the [beta installation document](./BETA-INSTALLATION.md) for the appropriate channel.

```console
helm repo add cloudzero https://cloudzero.github.io/cloudzero-charts
helm repo update
```

#### 2. Install Helm Chart

```console
helm install <RELEASE_NAME> cloudzero/cloudzero-agent \
    --set apiKey=<CLOUDZERO_API_KEY> \
    --set clusterName=<CLUSTER_NAME> \
    --set-string cloudAccountId=<CLOUD_ACCOUNT_ID> \
    --set region=<REGION>
```

---

### Advanced Install

Please visit the [CloudZero Agent Documentation](https://docs.cloudzero.com/docs/installation-of-cloudzero-agent-for-kubernetes) for advanced configuration.


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
