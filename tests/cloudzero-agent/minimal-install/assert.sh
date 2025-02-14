#!/bin/bash

# Define variables
TEST_POD_NAME="test-pod"
NAMESPACE="default"
CHART_NAMESPACE="cloudzero"
LABEL_SELECTOR="app.kubernetes.io/component=webhook-server"
LOG_LINE="[/validate/pod - CREATE] - Allowed: true]"

# Create a test pod
kubectl run $TEST_POD_NAME --image=busybox --restart=Never --namespace=$NAMESPACE -- /bin/sh -c "sleep 3600"

# Get the name of the existing pod based on labels
EXISTING_POD_NAME=$(kubectl get pods --namespace=$NAMESPACE -l $LABEL_SELECTOR -o jsonpath="{.items[0].metadata.name}")

# Get the logs from the existing pod
POD_LOGS=$(kubectl logs $EXISTING_POD_NAME --namespace=$NAMESPACE)

# Validate that the specific log line is in the logs
if echo "$POD_LOGS" | grep -q "$LOG_LINE"; then
    echo "Log line found. Test passed."
    exit 0
else
    echo "Log line not found. Test failed."
    exit 1
fi