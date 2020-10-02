# KNative Ping Pong
This toy example shows how to connect a consumer KNative service to a producer KNative service. The producer service is accepting HTTP requests.
The consumer service connects to the producer and just prints the status code of the HTTP response.

The example uses containerd as the container runtime for Kubernetes.

## Quick Start
In order to do the experiment from scratch:
1. ```./install.sh```
2. ```source /etc/profile && ./create_kubeadm_cluster.sh```
3. ```./call_proxy.sh```

## Details
The `install.sh` script installs runc, containerd, Kubernetes and Knative CLI - `kn`. You can customize this script or create your own if you want to use
different versions.

If you would like to build your own proxy/consumer image, you can use the `Dockerfile`. Alternatively, feel free to use `plamenppetrov/knative-ping-pong`.
