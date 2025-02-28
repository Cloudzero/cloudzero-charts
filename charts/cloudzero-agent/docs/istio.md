# Installing `cloudzero-agent` in Istio Enabled Clusters

When installing the `cloudzero-agent` Helm chart in a Kubernetes cluster that uses Istio, additional steps may be required to ensure proper functioning. This is because Istio automatically configures sidecars to use **mutual TLS (mTLS)**, which can interfere with the agent’s existing TLS communication.

To successfully deploy `cloudzero-agent`, either **exclude some communication from redirecting through envoy**, **disable Istio sidecar injection for a subset of the agent workloads**, or **disable mTLS for a subset of the agent workloads**.

---

## **Option 1: Disable Sidecar Injection for `cloudzero-agent` webhook-server Pods Only**
To disable the sidecar injection **only for a subset of the pods deployed by the cluster**, update the chart input values with the following annotation:

```yaml
insightsController:
  server:
    podAnnotations:
      sidecar.istio.io/inject: "false"
```

---

## **Option 2: Disable Envoy for Webhook Ports Only**
To prevent only requests to a single port on the webhook-server pods from being routed through envoy:

```yaml
insightsController:
  server:
    podAnnotations:
      traffic.sidecar.istio.io/excludeInboundPorts: "8443"
```
In this case, the pods will still have an istio sidecar injected. For more details, see [here](https://istio.io/latest/docs/reference/config/annotations/#SidecarTrafficExcludeInboundPorts).

---
## **Option 3: Disable mTLS for `cloudzero-agent`**
To disable mTLS for the `cloudzero-agent` service, apply the following `PeerAuthentication` resource:

```yaml
apiVersion: security.istio.io/v1beta1
kind: PeerAuthentication
metadata:
  name: cloudzero-agent-mtls
  namespace: <cloudzero-namespace>
spec:
  selector:
    matchLabels:
      app.kubernetes.io/component: webhook-server
  mtls:
    mode: DISABLE
```

### **Steps to Apply:**
1. Replace `<your-namespace>` with the namespace where `cloudzero-agent` is deployed.
2. Apply the resource:
   ```sh
   kubectl apply -f cloudzero-agent-mtls.yaml
   ```
3. Deploy the `cloudzero-agent` chart as instructed in the chart README.md

This configuration **disables mTLS for `cloudzero-agent` webhook-server pods only**, while keeping it enabled for the rest of the cluster.

---
