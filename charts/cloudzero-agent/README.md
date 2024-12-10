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

### Get Helm Repository Info

```console
helm repo add cloudzero https://cloudzero.github.io/cloudzero-charts
helm repo update
```

_See [`helm repo`](https://helm.sh/docs/helm/helm_repo/) for command documentation._

### Install Helm Chart

The chart can be installed directly with Helm or any other common Kubernetes deployment tools.

If installing with Helm directly, execute the following steps:

1. Ensure that the most recent chart version is available:
```console
helm repo update
```

2. Provision a TLS certificate; by default, this chart deploys a `ValidatingWebhookConfiguration` resource, which requires a certificate in order validate requests to the webhook server. See related Kubernetes documentation [here](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/#configure-admission-webhooks-on-the-fly).

There are several options for provisioning this certificate. The default is to use a third party tool, [cert-manager](https://cert-manager.io/). If you would prefer not to use cert-manager, see the [Certificate Management](#certificate-management) section for other options.

To use `cert-manager` for certificate management, first install the `cert-manager` CRDs with the following:
```console
helm install <RELEASE_NAME> cloudzero/cloudzero-agent \
    --set insightsController.webhook.issuer.enabled=false \
    --set insightsController.webhook.certificate.enabled=false \
    --set insightsController.cert-manager.installCRDs=true
```
Alternatively, [install the cert-manager CRDs directly](https://cert-manager.io/docs/installation/helm/) with:

```console
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.16.2/cert-manager.yaml
```
Alternatively, [install the cert-manager CRDs directly](https://cert-manager.io/docs/installation/helm/).

3. Fill out all required fields in the `configuration.example.yaml` file in this directory. Rename the file as necessary. Below is an example of a completed configuration file:
```yaml
# -- Account ID of the account the cluster is running in. This must be a string - even if it is a number in your system.
cloudAccountId: YOUR_CLOUD_ACCOUNT_ID
# -- Name of the clusters.
clusterName: YOUR_CLUSTER_NAME
# -- Region the cluster is running in.
region: YOUR_CLOUD_REGION
# -- CloudZero API key. Required if existingSecretName is false.
apiKey: YOUR_CLOUDZERO_API_KEY
# -- If set, the agent will use the API key in this Secret to authenticate with CloudZero.
existingSecretName: YOUR_EXISTING_API_KEY_K8S_SECRET

# label and annotation configuration (managed in the 'insightsController' section). See the below 'Labels and Annotations' section for more details.
insightsController:
  # -- By default, a ValidatingAdmissionWebhook will be deployed that records all created labels and annotations
  enabled: true
  labels:
    # -- This value MUST be set to either true or false. The installation will fail otherwise
    enabled: true
    # -- This value MUST be set to a list of regular expressions which will be used to gather labels from pods, deployments, statefulsets, daemonsets, cronjobs, jobs, nodes, and namespaces
    patterns:
      - '^foo' # -- Match all labels whose key starts with "foo"
      - 'bar$' # -- Match all labels whose key ends with "bar"
  annotations:
    # -- By default, the gathering of annotations is not enabled. To enable, set this field to true
    enabled: false
    patterns:
      - '.*' # -- match all annotations. This is not recommended.
cert-manager:
  # -- Your cluster may already have cert-manager running, in which case this value can be set to false.
  enabled: true
```

4. Install the helm chart using the completed configuration file:
```console
helm install <RELEASE_NAME> cloudzero/cloudzero-agent -f configuration.example.yaml
```

### Mandatory Values

There are several mandatory values that must be specified for the chart to install properly. Below are the required settings along with strategies for providing custom values during installation:

| Key               | Type   | Default               | Description                                                                                                             |
|-------------------|--------|-----------------------|-------------------------------------------------------------------------------------------------------------------------|
| cloudAccountId    | string | `nil`                 | Account ID in AWS or Subscription ID in Azure or Project Number in GCP where the cluster is running. Must be a string due to Helm limitations.  |
| clusterName       | string | `nil`                 | Name of the cluster. Must be RFC 1123 compliant.                                                                         |
| host              | string | `"api.cloudzero.com"` | CloudZero host to send metrics to.                                                                                      |
| apiKey            | string | `nil`                 | The CloudZero API key to use for exporting metrics. Only used if `existingSecretName` is not set.                       |
| existingSecretName| string | `nil`                 | Name of the secret that contains the CloudZero API key. Required if not providing the API key via `apiKey`.             |
| region            | string | `nil`                 | Region where the cluster is running (e.g., `us-east-1`, `eastus`). For more information, see AWS or Azure documentation. |
| insightsController.labels.enabled            | string | `nil`                 | If enabled, labels for pods, deployments, statefulsets, daemonsets, cronjobs, jobs, nodes, and namespaces |
| insightsController.labels.patterns            | string | `nil`                 | An array of regular expressions, which are used to match specific label keys |

#### Overriding Default Values

Default values are specified in the chart's `values.yaml` file. If you need to change any of these values, it is recommended to create a `values-override.yaml` based on the `configuration.example.yaml`  file for your customizations.

##### Using the `--values` Flag

You can use the `--values` (or short form `-f`) flag in your Helm commands to override values in the chart with a new file. Specify the name of the file after the `--values` flag:

```console
helm install <RELEASE_NAME> cloudzero/cloudzero-agent \
    -f values-override.yaml
```

Ensure `values-override.yaml` contains only the values you wish to override from `values.yaml`.

> Note it is possible to save values for different environments, or based on other criteria into separate values files and multiple files using the `-f` helm parameters.

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

### Labels and Annotations

This chart allows the exporting of labels and annotations from the following resources:
- `Pod`
- `Deployment`
- `StatefulSet`
- `Daemonset`
- `Job`
- `CronJob`
- `Node`
- `Namespace`

Additional Notes:
- Labels and annotations exports are managed in the `insightsController` section of the `values.yaml` file.
- By default, only labels from pods and namespaces are exported. To enable more resources, see the `insightsController.labels.resources` and `insightsController.annotations.resources` section of the `values.yaml` file.
- To disambiguate labels/annotations between resources, a prefix representing the resource type is prepended to the label key in the [CloudZero Explorer](https://app.cloudzero.com/explorer). For example, a `foo=bar` node label would be presented as `node:foo: bar`. The exception is pod labels which do not have resource prefixes for backward compatibility with previous versions.
- Annotations are not exported by default; see the `insightsController.annotations.enabled` setting to enable. To disambiguate annotations from labels, an `annotation` prefix is prepended to the annotation key; i.e., an `foo: bar` annotation on a namespace would be represented in the Explorer as `node:annotation:foo: bar`
- For both labels and annotations, the `patterns` array applies across all resource types; i.e., setting `['^foo']` for `insightsController.labels.patterns` will match label keys that start with `foo` for all resource types set to `true` in `insightsController.labels.resources`.

### Certificate Management

The default behavior of the chart is to deploy [cert-manager](https://github.com/cert-manager/cert-manager/tree/master) to create and manage the certificate, but there are several alternate options if using `cert-manager` is not possible:

#### Option 1: use the `cloudzero-certificate` chart

The `cloudzero-certificate` chart, which is maintained in this repo, creates a certificate and stores it in a Secret. The `cloudzero-agent` resources can then access this Secret. To install this helm chart, first decide what the release name of the `cloudzero-agent` will be (shown as `EXAMPLE_CLOUDZERO_AGENT_RELEASE_NAME` here), as the release name will be used as the DNS name in the certificate. Then, do the following:

1. Create a Secret using the `cloudzero-certificate` chart and get the CA bundle value:
```console
helm repo update

# Install the chart, which creates a Secret with a TLS certificate
helm upgrade --install --namespace <YOUR_NAMESPACE> <YOUR_RELEASE_NAME> cloudzero --set cloudzeroAgentReleaseName=<EXAMPLE_CLOUDZERO_AGENT_RELEASE_NAME>

# Get the CA bundle value by running:
CA_BUNDLE=$(kubectl get secret -n <YOUR_NAMESPACE> <YOUR_RELEASE_NAME>-cloudzero-certificate -o jsonpath='{.data.ca\.crt}')

# Confirm that CA_BUNDLE is set:
echo $CA_BUNDLE
```

2. Next, set the following in your `configuration.example.yaml` file in addition to the existing values:
```yaml
insightsController:
  server:
    tls:
      nameOverride: <YOUR_RELEASE_NAME>-cloudzero-certificate # This should be the name of the secret created in the previous step
  webhooks:
    caBundle: $CA_BUNDLE # This should be the value of the CA_BUNDLE variable in the previous step
  certificate:
    enabled: false
  issuer:
    enabled: false
cert-manager:
  enabled: false
```

3. Finally, continue on with the rest of the installation in the [Installation](#installation) section. The only new requirement is that when installing the `cloudzero-agent` helm chart, you must use the same value for the release name as was set in `cloudzeroAgentReleaseName` in step #1. For example:

```console
helm install <EXAMPLE_CLOUDZERO_AGENT_RELEASE_NAME> cloudzero/cloudzero-agent -f configuration.example.yaml
```

#### Option 2: bring your own certificate

The `cloudzero-agent` chart can also use any certificate provided by an external source. The common name of the certificate should be `<RELEASE_NAME>.<RELEASE_NAMESPACE>.cluster.local`. A Kubernetes Secret should be created with the keys:
```
ca.crt: <base64 encoded caBundle>
tls.crt: <base64 encoded certificate>
tls.key: <base64 encoded key>
```
Finally, set the following in your `configuration.example.yaml` file in addition to the existing values:
```yaml
insightsController:
  server:
    tls:
      nameOverride: <YOUR_TLS_SECRET_NAME>-cloudzero-certificate # This should be the name of the secret created in the previous step
  webhooks:
    caBundle: $CA_BUNDLE # This should be the value of the `ca.crt` value created in the Secret
  certificate:
    enabled: false
  issuer:
    enabled: false
cert-manager:
  enabled: false
```

### Secret Management

The chart requires a CloudZero API key to send metric data. Admins can retrieve API keys [here](https://app.cloudzero.com/organization/api-keys).

The API key is typically stored in a Secret in the cluster. The `cloudzero-agent` chart will create a Secret if the API key is provided via the `apiKey` argument. Alternatively, the API key can be stored in a Secret external to the chart; this Secret name would then be set as the `existingSecretName` argument. If creating a Secret external to the chart, ensure the Secret is in the same namespace as the chart and follows this format:

**Example User-Created Secret Content**
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

### Update Helm Chart
If you are updating an existing installation, pull the latest chart information:

```console
helm repo update
```

Next, upgrade the installation to the latest chart version:

```console
helm upgrade --install <RELEASE_NAME> cloudzero/cloudzero-agent -f configuration.example.yaml
```

#### Passing Values to Subcharts

Values can be passed to subcharts like [kube-state-metrics](https://github.com/prometheus-community/helm-charts/blob/main/charts/kube-state-metrics/values.yaml) by adding entries in `values-override.yaml` as per their specifications.

A common addition may be to pull the container images from custom image registries/repositories:

`values-override.yaml`
```yaml
kubeStateMetrics:
  enabled: true
  image:
    registry: my-custom-registry.io
    repository: my-custom-kube-state-metrics/kube-state-metrics
```

## Dependencies

| Repository                                         | Name                     | Version |
|----------------------------------------------------|--------------------------|---------|
| https://prometheus-community.github.io/helm-charts | kube-state-metrics       | 5.15.*  |

Note that while `kube-state-metrics` is listed as a dependency, it is referred to as `cloudzero-state-metrics` within the helm chart. This is to enforce the idea that this KSM deployment is used exclusively by the `cloudzero-agent`

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
