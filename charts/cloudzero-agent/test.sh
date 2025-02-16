#!/bin/bash

# GENERATE_CERTIFICATE=false

# # Check if the caBundle in the ValidatingWebhookConfiguration is the same for all webhooks
# caBundles=()
# caBundles+=($(kubectl get validatingwebhookconfiguration cloudzero-agent-webhook-server-webhook-deployments -o jsonpath='{.webhooks[0].clientConfig.caBundle}'))
# caBundles+=($(kubectl get validatingwebhookconfiguration cloudzero-agent-webhook-server-webhook-namespaces -o jsonpath='{.webhooks[0].clientConfig.caBundle}'))
# caBundles+=($(kubectl get validatingwebhookconfiguration cloudzero-agent-webhook-server-webhook-nodes -o jsonpath='{.webhooks[0].clientConfig.caBundle}'))
# caBundles+=($(kubectl get validatingwebhookconfiguration cloudzero-agent-webhook-server-webhook-pods -o jsonpath='{.webhooks[0].clientConfig.caBundle}'))

# CA_BUNDLE=${caBundles[0]}
# for caBundle in "${caBundles[@]}"; do
#     if [[ "$caBundle" != "$CA_BUNDLE" ]]; then
#         echo "Mismatch found between ValidatingWebhookConfiguration caBundle values."
#         GENERATE_CERTIFICATE=true
#     fi
# done

# SECRET_NAME=cloudzero-agent-webhook-server-tls
# NAMESPACE=bat

# EXISTING_TLS_CRT=$(kubectl get secret $SECRET_NAME -n $NAMESPACE -o jsonpath='{.data.tls\.crt}')
# EXISTING_TLS_KEY=$(kubectl get secret $SECRET_NAME -n $NAMESPACE -o jsonpath='{.data.tls\.key}')

# # Check if the existing TLS certificate is expired
# if [[ -n "$EXISTING_TLS_CRT" ]]; then
#     EXPIRATION_DATE=$(echo "$EXISTING_TLS_CRT" | base64 -d | openssl x509 -noout -enddate)
#     echo $EXPIRATION_DATE
#     # if [[ $(date -d "$EXPIRATION_DATE" +%s) -lt $(date +%s) ]]; then
#     #     echo "The existing TLS certificate has expired."
#     #     GENERATE_CERTIFICATE=true
#     # fi
#     # Check if the SANs in the certificate match the service name
#     SAN=$(echo "$EXISTING_TLS_CRT" | base64 -d | openssl x509 -text -noout | grep DNS | sed 's/.*DNS://')
#     if [[ "$SAN" != "cloudzero-agent-webhook-server-svc.bat.svc" ]]; then
#         echo "The SANs in the certificate do not match the service name."
#         GENERATE_CERTIFICATE=true
#     fi
#     # Check that caBundle and tls.crt are the same
#     if [[ "$CA_BUNDLE" != $EXISTING_TLS_CRT ]]; then
#         echo "The caBundle in the ValidatingWebhookConfiguration does not match the tls.crt in the TLS Secret."
#         GENERATE_CERTIFICATE=true
#     fi
# fi

# # Check if the TLS Secret already has certificate information
# if [[ -z "$EXISTING_TLS_CRT" ]] || [[ -z "$EXISTING_TLS_KEY" ]] || [[ $GENERATE_CERTIFICATE == "true" ]] ; then
#     echo "The TLS Secret and/or at least one webhook configuration contains empty certificate information, or the certificate is invalid/expired. Creating a new certificate..."
# else
#     echo "The TLS Secret and all webhook configurations contain non-empty certificate information. Will not create a new certificate and will not patch resources."
#     exit 0
# fi

kubectl get validatingwebhookconfigurations cloudzero-agent-webhook-server-webhook-namespaces -o jsonpath='{.webhooks[0].clientConfig.caBundle}' | md5sum
kubectl get secret cloudzero-agent-webhook-server-tls -n foobar -o jsonpath='{.data.tls\.crt}' | md5sum
kubectl get secret cloudzero-agent-webhook-server-tls -n foobar -o jsonpath='{.data.tls\.key}' | md5sum
# echo $GENERATE_CERTIFICATE