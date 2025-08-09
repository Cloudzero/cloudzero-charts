# Using CSI Secrets Store with CloudZero Agent

The CloudZero Agent Helm chart supports using the [Kubernetes CSI Secrets Store driver](https://secrets-store-csi-driver.sigs.k8s.io/) to retrieve secrets from external secret management systems like AWS Secrets Manager, Azure Key Vault, or Google Secret Manager.

## Prerequisites

1. **CSI Secrets Store Driver**: The CSI Secrets Store driver must be installed in your cluster.
   ```bash
   helm repo add secrets-store-csi-driver https://kubernetes-sigs.github.io/secrets-store-csi-driver/charts
   helm install csi-secrets-store secrets-store-csi-driver/secrets-store-csi-driver --namespace kube-system
   ```

2. **Provider**: Install the appropriate provider for your cloud platform:
   - **AWS**: [AWS Provider](https://github.com/aws/secrets-store-csi-driver-provider-aws)
   - **Azure**: [Azure Provider](https://azure.github.io/secrets-store-csi-driver-provider-azure/)
   - **GCP**: [GCP Provider](https://github.com/GoogleCloudPlatform/secrets-store-csi-driver-provider-gcp)

3. **SecretProviderClass**: Create a SecretProviderClass resource that defines how to retrieve your CloudZero API key.

## Configuration

### Basic Setup

Configure the CloudZero Agent to use CSI Secrets Store by adding extra volumes and volume mounts:

```yaml
# Required basic configuration
clusterName: "my-cluster"
cloudAccountId: "123456789012"  # Must be a string
region: "us-east-1"

# Use existing secret name to satisfy validation
existingSecretName: "cloudzero-api-key-secret"
apiKey: null

# Configure CSI Secrets Store volume
extraVolumes:
  - name: cloudzero-api-key-csi
    csi:
      driver: secrets-store.csi.k8s.io
      readOnly: true
      volumeAttributes:
        secretProviderClass: "cloudzero-secrets-provider"

extraVolumeMounts:
  - name: cloudzero-api-key-csi
    mountPath: /etc/csi-secrets/
    readOnly: true
```

### SecretProviderClass Examples

#### AWS Secrets Manager

```yaml
apiVersion: secrets-store.csi.x-k8s.io/v1
kind: SecretProviderClass
metadata:
  name: cloudzero-secrets-provider
  namespace: default
spec:
  provider: aws
  parameters:
    objects: |
      - objectName: "cloudzero-api-key"
        objectType: "secretsmanager"
        jmesPath:
          - path: "api_key"
            objectAlias: "value"
```

#### Azure Key Vault

```yaml
apiVersion: secrets-store.csi.x-k8s.io/v1
kind: SecretProviderClass
metadata:
  name: cloudzero-secrets-provider
  namespace: default
spec:
  provider: azure
  parameters:
    usePodIdentity: "false"
    useVMManagedIdentity: "true"
    userAssignedIdentityID: "your-managed-identity-client-id"
    keyvaultName: "your-keyvault-name"
    objects: |
      array:
        - |
          objectName: cloudzero-api-key
          objectType: secret
          objectAlias: value
    tenantId: "your-tenant-id"
```

#### Google Secret Manager

```yaml
apiVersion: secrets-store.csi.x-k8s.io/v1
kind: SecretProviderClass
metadata:
  name: cloudzero-secrets-provider
  namespace: default
spec:
  provider: gcp
  parameters:
    secrets: |
      - resourceName: "projects/your-project-id/secrets/cloudzero-api-key/versions/latest"
        path: "value"
```

## Advanced Configuration

### Multiple Secrets

You can mount multiple secrets from different sources:

```yaml
extraVolumes:
  - name: cloudzero-api-key-csi
    csi:
      driver: secrets-store.csi.k8s.io
      readOnly: true
      volumeAttributes:
        secretProviderClass: "cloudzero-secrets-provider"
  - name: additional-secrets
    csi:
      driver: secrets-store.csi.k8s.io
      readOnly: true
      volumeAttributes:
        secretProviderClass: "additional-secrets-provider"

extraVolumeMounts:
  - name: cloudzero-api-key-csi
    mountPath: /etc/csi-secrets/
    readOnly: true
  - name: additional-secrets
    mountPath: /etc/additional-secrets/
    readOnly: true
```

### Authentication with nodePublishSecretRef

For providers that require additional authentication:

```yaml
extraVolumes:
  - name: cloudzero-api-key-csi
    csi:
      driver: secrets-store.csi.k8s.io
      readOnly: true
      volumeAttributes:
        secretProviderClass: "cloudzero-secrets-provider"
      nodePublishSecretRef:
        name: secrets-store-creds
```

## Deployment

Deploy the CloudZero Agent with CSI Secrets Store configuration:

```bash
helm install cloudzero-agent cloudzero/cloudzero-agent -f csi-secrets-store-values.yaml
```

## Verification

After deployment, verify that the secrets are mounted correctly:

```bash
# Check that the CSI volume is mounted
kubectl describe pod -l app.kubernetes.io/name=cloudzero-agent

# Verify secret content (be careful with sensitive data)
kubectl exec -it deployment/cloudzero-agent-server -- ls -la /etc/csi-secrets/
```

## Troubleshooting

1. **Volume Mount Issues**: Ensure the CSI Secrets Store driver is running and the SecretProviderClass is in the same namespace as the CloudZero Agent.

2. **Authentication Failures**: Verify that your cloud provider authentication is configured correctly (IAM roles, managed identities, service accounts).

3. **Secret Not Found**: Check that the secret exists in your secret store and the SecretProviderClass references it correctly.

4. **Permission Denied**: Ensure the pod's service account has the necessary permissions to access the secret store.

## Security Considerations

- Use least-privilege access when configuring cloud provider permissions
- Regularly rotate your CloudZero API keys
- Monitor access to your secret stores
- Consider using pod identity or workload identity for authentication instead of static credentials

## References

- [CSI Secrets Store Driver Documentation](https://secrets-store-csi-driver.sigs.k8s.io/)
- [AWS Provider Documentation](https://github.com/aws/secrets-store-csi-driver-provider-aws)
- [Azure Provider Documentation](https://azure.github.io/secrets-store-csi-driver-provider-azure/)
- [GCP Provider Documentation](https://github.com/GoogleCloudPlatform/secrets-store-csi-driver-provider-gcp)
