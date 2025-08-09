# Custom Volumes and Volume Mounts

The CloudZero Agent Helm chart supports adding custom volumes and volume mounts to all deployments (agent, aggregator, and webhook) through the `extraVolumes` and `extraVolumeMounts` configuration options.

## Configuration

### Global Configuration

Add custom volumes and volume mounts that will be applied to all deployments:

```yaml
# Global custom volumes that will be added to all deployments
extraVolumes:
  - name: custom-config
    configMap:
      name: my-custom-config
      defaultMode: 0644
  - name: custom-secret
    secret:
      secretName: my-custom-secret
      defaultMode: 0600

# Global custom volume mounts that will be added to all containers
extraVolumeMounts:
  - name: custom-config
    mountPath: /etc/custom-config
    readOnly: true
  - name: custom-secret
    mountPath: /etc/custom-secret
    readOnly: true
```

### Deployment-Specific Configuration

For the webhook deployment, you can also use deployment-specific volumes:

```yaml
insightsController:
  volumes:
    - name: webhook-specific-volume
      configMap:
        name: webhook-config
  volumeMounts:
    - name: webhook-specific-volume
      mountPath: /etc/webhook-config
      readOnly: true
```

## Supported Volume Types

The `extraVolumes` configuration supports all Kubernetes volume types, including:

- **ConfigMap**: Mount configuration files
- **Secret**: Mount sensitive data
- **EmptyDir**: Temporary storage
- **PersistentVolumeClaim**: Persistent storage
- **HostPath**: Host filesystem access (use with caution)
- **NFS**: Network File System
- **CSI**: Container Storage Interface volumes

## Examples

### Example 1: Adding Configuration Files

```yaml
extraVolumes:
  - name: app-config
    configMap:
      name: my-app-config
      items:
        - key: config.yaml
          path: config.yaml

extraVolumeMounts:
  - name: app-config
    mountPath: /etc/app-config
    readOnly: true
```

### Example 2: Adding Secrets

```yaml
extraVolumes:
  - name: tls-certs
    secret:
      secretName: my-tls-secret
      defaultMode: 0400

extraVolumeMounts:
  - name: tls-certs
    mountPath: /etc/ssl/certs
    readOnly: true
```

### Example 3: Adding Temporary Storage

```yaml
extraVolumes:
  - name: temp-storage
    emptyDir:
      sizeLimit: 1Gi

extraVolumeMounts:
  - name: temp-storage
    mountPath: /tmp/app-temp
```

### Example 4: Adding Persistent Storage

```yaml
extraVolumes:
  - name: data-storage
    persistentVolumeClaim:
      claimName: my-pvc

extraVolumeMounts:
  - name: data-storage
    mountPath: /data
```

## Deployment Coverage

The `extraVolumes` and `extraVolumeMounts` configurations are applied to:

1. **Agent Deployment**: Main Prometheus agent container
2. **Aggregator Deployment**: Both collector and shipper containers
3. **Webhook Deployment**: Webhook server container

## Important Notes

1. **Volume Names**: Ensure volume names are unique and don't conflict with existing volumes used by the chart
2. **Mount Paths**: Choose mount paths that don't conflict with existing application paths
3. **Permissions**: Set appropriate file permissions using `defaultMode` for ConfigMaps and Secrets
4. **Resource Limits**: Consider the impact of additional volumes on resource usage
5. **Security**: Be cautious with HostPath volumes and ensure proper security contexts

## Troubleshooting

### Common Issues

1. **Volume Name Conflicts**: If you see errors about duplicate volume names, ensure your custom volume names don't conflict with existing ones
2. **Mount Path Conflicts**: Avoid mounting volumes to paths already used by the application
3. **Permission Issues**: Ensure the container user (65534) has appropriate permissions to access mounted volumes

### Debugging

To verify your volumes are correctly configured, use:

```bash
# Check the rendered templates
helm template my-release cloudzero/cloudzero-agent -f my-values.yaml

# Check deployed resources
kubectl describe deployment my-release-cloudzero-agent-server
kubectl describe deployment my-release-aggregator
kubectl describe deployment my-release-cloudzero-agent-webhook-server
```

## Migration from Deployment-Specific Volumes

If you were previously using deployment-specific volume configurations, you can migrate to the global `extraVolumes` approach:

**Before (webhook-specific):**
```yaml
insightsController:
  volumes:
    - name: my-volume
      configMap:
        name: my-config
  volumeMounts:
    - name: my-volume
      mountPath: /etc/config
```

**After (global):**
```yaml
extraVolumes:
  - name: my-volume
    configMap:
      name: my-config

extraVolumeMounts:
  - name: my-volume
    mountPath: /etc/config
```

The global approach ensures the volume is available across all deployments, providing consistency and reducing configuration duplication.
