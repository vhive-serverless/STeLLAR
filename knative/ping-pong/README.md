# KNative Ping Pong
This toy example shows how a client can call one Knative service - `f1`, which calls another Knative service `f2`. `f2` then responds to `f1`, which sends `f2`'s response back to the client. The workflow can be seen in the diagram below.

![Flow](images/flow.pdf)

We create a Kubernetes cluster with [kubeadm](https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/create-cluster-kubeadm/) and use [containerd](https://github.com/containerd/containerd) as the container runtime for Kubernetes.

## Quick Start
In order to do the experiment from scratch:
1. ```./install.sh```
2. Start containerd in another terminal - `containerd`.
2. ```source /etc/profile && ./create_kubeadm_cluster.sh```
3. ```./call.sh``` - call `f1`

## Details
The `install.sh` script installs runc, containerd, Kubernetes and Knative CLI - `kn`.

The example uses the `plamenppetrov/f1` and `plamenppetrov/f2` images on DockerHub for `f1` and `f2`, respectively. If you would like to build your own images for `f1` and `f2`, you can use the Dockerfiles provided. In this case, make sure to change the `.yaml` files for `f1` and `f2` accordingly. 

We patch the load balancer to use the public IP of the master node. Each Knative service can be called by issuing an HTTP request to the load balancer, with the `Host` header field set to the virtual domain name of the service. `call.sh` is an example of this.
