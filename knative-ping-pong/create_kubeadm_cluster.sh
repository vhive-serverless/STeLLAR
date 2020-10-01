#!/bin/bash
kubeadm init --ignore-preflight-errors=all --cri-socket /run/containerd/containerd.sock --pod-network-cidr=192.168.0.0/16

# Install Calico network add-on
kubectl apply -f https://docs.projectcalico.org/manifests/calico.yaml

# Untaint master (allow pods to be scheduled on master) 
kubectl taint nodes --all node-role.kubernetes.io/master-


# Install KNative in the cluster
kubectl apply --filename https://github.com/knative/serving/releases/download/v0.17.0/serving-crds.yaml
kubectl apply --filename https://github.com/knative/serving/releases/download/v0.17.0/serving-core.yaml

# Configure network
kubectl apply --filename https://raw.githubusercontent.com/Kong/kubernetes-ingress-controller/0.9.x/deploy/single/all-in-one-dbless.yaml
kubectl patch configmap/config-network \
  --namespace knative-serving \
  --type merge \
  --patch '{"data":{"ingress.class":"kong"}}'

PUBLIC_IP=$(curl ifconfig.me)
kubectl patch svc kong-proxy -n kong -p '{"spec": {"type": "LoadBalancer", "externalIPs":["'${PUBLIC_IP}'"]}}'

kn service create helloworld-go --image gcr.io/knative-samples/helloworld-go --env TARGET="Go Sample v1"

kubectl --namespace kong get service kong-proxy

NODE_PORT=$(kubectl --namespace kong get service kong-proxy -o go-template='{{(index .spec.ports 0).nodePort}}')
cat proxy_template.yaml | sed 's/EXTERNALIP/'$(curl ifconfig.me)'/g' | sed 's/NODEPORT/'${NODE_PORT}'/g' > proxy.yaml
