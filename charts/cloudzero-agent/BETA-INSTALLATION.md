# CloudZero Helm Charts Beta Installation Guide

To install a Beta version of the chart, use the `beta` Helm repository. Follow these steps to add the beta repository:

```sh
helm repo add cloudzero-beta https://cloudzero.github.io/cloudzero-charts/beta
helm repo update
```

Here, we name the repository `cloudzero-beta`. Use this name in Helm install commands instead of the standard `cloudzero`. For example:

```sh
helm install <RELEASE_NAME> cloudzero-beta/cloudzero-agent \
    --set existingSecretName=<NAME_OF_SECRET> \
    --set clusterName=<CLUSTER_NAME> \
    --set-string cloudAccountId=<CLOUD_ACCOUNT_ID> \
    --set region=<REGION> \
    # Optionally deploy kube-state-metrics if it doesn't exist in the cluster already
    --set kube-state-metrics.enabled=<true|false>
```

Follow all other installation instructions as defined in the [README](./README.md).