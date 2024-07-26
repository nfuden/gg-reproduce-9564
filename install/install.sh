#!/bin/sh
helm repo add glooe https://storage.googleapis.com/gloo-ee-helm
helm repo update
export GLOO_EDGE_HELM_VALUES_FILE="./install/gloo-edge-helm-values.yaml"

if [ -z "$GLOO_EDGE_LICENSE_KEY" ]
then
   echo "Gloo Edge License Key not specified. Please configure the environment variable 'GLOO_EDGE_LICENSE_KEY' with your Gloo Edge License Key."
   exit 1
fi

if [ -z "$GLOO_EDGE_VERSION" ]
then
   export GLOO_EDGE_VERSION="1.17.0"
fi


helm --debug upgrade --install gloo glooe/gloo-ee --namespace gloo-system --create-namespace --set-string license_key=$GLOO_EDGE_LICENSE_KEY -f $GLOO_EDGE_HELM_VALUES_FILE --version $GLOO_EDGE_VERSION