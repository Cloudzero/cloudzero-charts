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
- A kube-state-metrics exporter running in the cluster, available via Kubernetes Service (see below for details)

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

2. Ensure that required CRDs are installed for certifiacte management. If you have more specific requirements around managing TLS certificates, see the [Certificate Management](https://github.com/Cloudzero/cloudzero-charts/tree/develop/charts/cloudzero-insights-controller#deployment-configurations-and-certificate-management) section in the `cloudzero-insights-controller` subchart.
```console
helm install <RELEASE_NAME> cloudzero/cloudzero-agent \
    --set tags.webhook.issuer.enabled=false \
    --set tags.webhook.certificate.enabled=false \
    --set tags.cert-manager.installCRDs=true
```

3. Fill out all required fields in the `configuration.example.yaml` file in this directory. Rename the file as necessary. Below is an example of a completed configuration file:
```yaml
# -- Account ID of the account the cluster is running in. This must be a string - even if it is a number in your system.
cloudAccountId: YOUR_CLOUD_ACCOUNT_ID
# -- Name of the clusters.
clusterName: YOUR_CLUSTER_NAME
# -- Region the cluster is running in.
region: YOUR_CLOUD_REGION
global:
  # -- CloudZero API key. Required if useExistingSecret is false.
  apiKey: YOUR_CLOUDZERO_API_KEY
  # -- If set, the agent will use the API key in this Secret to authenticate with CloudZero.
  existingSecretName: YOUR_EXISTING_API_KEY_K8S_SECRET

# label and annotation configuration (referred together as 'tags'). See the below 'Labels and Annotations' section for more details.
tags:
  # -- By default, a ValidatingAdmissionWebhook will be deployed that records all created labels and annotations
  enabled: true
  labels:
    # -- This value MUST be set to either true or false. The installation will fail otherwise
    enabled: true
    # -- This value MUST be set to a list of regular expressions which will be used to gather labels from pods, deployments, statefulsets, daemonsets, cronjobs, jobs, nodes, and namespaces
    patterns:
      - '^foo' # -- match all labels whose key starts with "foo"
      - 'bar$' # -- match all labels whose key ends with "bar"
  annotations:
    # -- By default, the gathering of annotations is not enabled. To enable, set this field to true
    enabled: false
    patterns:
      - '.*'
```

4. Install the helm chart using the completed configuration file:
```console
helm install <RELEASE_NAME> cloudzero/cloudzero-agent -f configuration.example.yaml
```

### Update Helm Chart
Alternatively, if you are updating an existing installation, pull the latest chart information first:

```console
helm repo update
```

Next, upgrade the installation to the latest chart version:

```console
helm upgrade --install <RELEASE_NAME> cloudzero/cloudzero-agent -f configuration.example.yaml
```

### Mandatory Values

There are several mandatory values that must be specified for the chart to install properly. Below are the required settings along with strategies for providing custom values during installation:

| Key               | Type   | Default               | Description                                                                                                             |
|-------------------|--------|-----------------------|-------------------------------------------------------------------------------------------------------------------------|
| cloudAccountId    | string | `nil`                 | Account ID in AWS or Subscription ID in Azure or Project Number in GCP where the cluster is running. Must be a string due to Helm limitations.  |
| clusterName       | string | `nil`                 | Name of the cluster. Must be RFC 1123 compliant.                                                                         |
| host              | string | `"api.cloudzero.com"` | CloudZero host to send metrics to.                                                                                      |
| global.apiKey            | string | `nil`                 | The CloudZero API key to use for exporting metrics. Only used if `global.existingSecretName` is not set.                       |
| global.existingSecretName| string | `nil`                 | Name of the secret that contains the CloudZero API key. Required if not providing the API key via `apiKey`.             |
| region            | string | `nil`                 | Region where the cluster is running (e.g., `us-east-1`, `eastus`). For more information, see AWS or Azure documentation. |
| tags.labels.enabled            | string | `nil`                 | If enabled, labels for pods, deployments, statefulsets, daemonsets, cronjobs, jobs, nodes, and namespaces |
| tags.labels.patterns            | string | `nil`                 | An array of regular expressions, which are used to match specific label keys |

#### Overriding Default Values

Default values are specified in the chart's `values.yaml` file. If you need to change any of these values, it is recommended to create a `values-override.yaml` based on the `configuration.example.yaml`  file for your customizations.

##### Using the `--values` Flag

You can use the `--values` (or short form `-f`) flag in your Helm commands to override values in the chart with a new file. Specify the name of the file after the `--values` flag:

```console
helm install <RELEASE_NAME> cloudzero/cloudzero-agent \
    -f values-override.yaml
```

Ensure `values-override.yaml` contains only the values you wish to override from `values.yaml`.

> Note it is possible to save values for different environments, or based on other criteria into seperate values files and multiple files using the `-f` helm parameters.

##### Using the `--set` Flag

You can use the `--set` flag in Helm commands to directly set or override specific values from `values.yaml`. Use dot notation to specify nested values:

```console
helm install <RELEASE_NAME> cloudzero/cloudzero-agent \
    --set global.existingSecretName=<NAME_OF_SECRET> \
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
- By default, only labels from pods and namespaces are exported. To enabled more resources, see the `webhooks.configurations` section of the `values.yaml` file.
- Labels and annotations exports are managed by a subchart, `cloudzero-insights-controller`, which is also maintained in this repository.
- To disambiguate labels/annotations between resources, a prefix representing the resource type is prepended to the label key in the Explorer page. For example, a `foo=bar` node label would be presented as `node:foo: bar`. The exception is pod labels which do not have resource prefixes for backward compatibility with previous versions.
- Annotations are not exported by default; see the `tags.annotations.enabled` setting to enable. To disambiguate annotations from labels, an `annotation` prefix is prepended to the annotation key; i.e., an `foo: bar` annotation on a namespace would be represented in the Explorer as `node:annotation:foo: bar`
- For both labels and annotations, the `enabled` flag applies across all resource types; i.e., setting `true` for `tags.labels.enabled` enabels label exporting for labels on pods, nodes, deployments, statefulsets, namespaces, daemonets, jobs, and cronjobs. Specific resources can be disabled by altering the `webhooks.configurations` configuration.

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

#### Passing Values to Subcharts

Values can be passed to subcharts like [kube-state-metrics](https://github.com/prometheus-community/helm-charts/blob/main/charts/kube-state-metrics/values.yaml) by adding entries in `values-override.yaml` as per their specifications.

A common addition may be to pull the container images from custom image registries/repositories:

`values-override.yaml`
```yaml
kube-state-metrics:
  enabled: true
  image:
    registry: my-custom-registry.io
    repository: my-custom-kube-state-metrics/kube-state-metrics
```

### Custom Scrape Configs

If running without the default `kube-state-metrics` exporter subchart and your existing `kube-state-metrics` deployment does not have the required `prometheus.io/scrape: "true"`, adjust the Prometheus scrape configs as shown:

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
```

### Exporting Pod Labels

Pod labels can be exported as metrics using kube-state-metrics. To customize the labels for export, modify the values-override.yaml file as shown below:

**Example: Exporting only the pod labels named foo and bar:**

```yaml
kube-state-metrics:
  extraArgs:
     - --metric-labels-allowlist=pods=[foo,bar]
```

> This is preferable to including all labels with `*` because the performance and memory impact is reduced. Regular expression matching is not currently supported. See the `kube-state-metrics` [documentation](https://github.com/kubernetes/kube-state-metrics/blob/main/docs/developer/cli-arguments.md) for more details.

⚠️ Important: If you are running an existing `kube-state-metrics` instance, ensure that the labels you want to use are whitelisted. kube-state-metrics version 2.x and above will **_not_** export the `kube_pod_labels` metrics unless they are explicitly allowed. This prevents the use of those labels for cost allocation and other purposes. Make sure you have configured the labels at the appropriate level using the --metric-labels-allowlist parameter:

> eg:  `- --metric-labels-allowlist=namespaces=[*],pods=[*],deployments=[app.kubernetes.io/*,k8s.*]`

## Dependencies

| Repository                                         | Name                     | Version |
|----------------------------------------------------|--------------------------|---------|
| https://prometheus-community.github.io/helm-charts | kube-state-metrics       | 5.15.*  |

## Enabling Release Notifications

To receive a notification when a new version of the chart is [released](https://github.com/Cloudzero/cloudzero-charts/releases), you can [watch the repository](https://docs.github.com/en/account-and-profile/managing-subscriptions-and-notifications-on-github/setting-up-notifications/configuring-notifications#configuring-your-watch-settings-for-an-individual-repository):

1. Navigate to the [repository main page](https://github.com/Cloudzero/cloudzero-charts).
2. Select **Watch > Custom**.
3. Check the **Releases** box.
4. Select **Apply**.


## Troubleshooting

### Issue
I've deployed the chart, but I don't see Kubernetes data in CloudZero.

## Resolution
This can happen for a number of reasons; see below for solutions to the most common problems

### Ensure kube-state-metrics is deployed correctly

1. Review the **Metric Exporters** section.
2. If opting for **Option 1**
  - Is kube-state-metrics installed?
  ```bash
  kubectl get services --all-namespaces | grep kube-state-metrics
  ```
  If the above command does not return any services, install a `kube-state-metrics` exporter, or use **Option 2** in the **Metric Exporters** section.
  
3. If opting for **Option 2**, ensure that `kube-state-metrics.enabled=true` is set as an annotation on the Service.
4. Ensure the cloudzero-agent pod can find the `kube-state-metrics` Service.
   Run the following command:
   ```
   kubectl get services -A -o jsonpath='{range .items[?(@.metadata.annotations.prometheus\.io/scrape=="true")]}{.metadata.name}{" in "}{.metadata.namespace}{"\n"}{end}'
   ```
   If this does not return a `kube-state-metrics` Service, then either annotate the existing Service found in Step 2 with `prometheus.io/scrape: "true"`, or following the instructions in the **Custom Scrape Configs** section above.
5. Ensure connectivity between the `cloudzero-agent` pod and the `kube-state-metrics` Service.
  ```
  SERVER_POD=$(kubectl get pod -l app.kubernetes.io/name=cloudzero-agent -o jsonpath='{.items[0].metadata.name}')
  kubectl exec -it -n <NAMESPACE> $SERVER_POD -- wget -qO- <KSM_SERVICE_NAME>.<KSM_NAMESPACE>.svc.cluster.local:8080/metrics
  ```
  The request should return a 200 response with a list of metrics prefixed with `kube_`, i.e., `kube_pod_info`. If not, ensure that the `kube-state-metrics` deployment is configured correctly.

### Issue
I have Kubernetes data in CloudZero, but I don't see Kubernetes labels as Dimensions.

## Resolution
Note that
1. Only labels on Pods are currently supported, and
2. Labels are "opt-in"; see the **Exporting Pod Labels** section for details.

## Useful References

- [Memory Sizing Guide](./docs/sizing-guide.md)
- [Deployment Validation Guide](./docs/deploy-validation.md)

---
