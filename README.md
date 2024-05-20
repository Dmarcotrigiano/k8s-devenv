# Kubernetes Native Developer Environment

This is a Nx monorepo intended to be used to collocate Go microservices and Kubernetes manifests in a single repository. It is intended to be used as a starting point for a Kubernetes native developer environment.

## Features

- [Nx](https://nx.dev) monorepo with Go support
- [Skaffold](https://skaffold.dev/) for local development
- [ko](https://ko.build/) for building Go binaries and Docker images
- [Helm](https://helm.sh/) for deploying prepackaged manifests to Kubernetes
- [NGINX Ingress Controller](https://kubernetes.github.io/ingress-nginx/) for routing traffic to services from the host machine
- [Ingress DNS](https://minikube.sigs.k8s.io/docs/handbook/addons/ingress-dns/) for resolving DNS names to the Ingress Controller on MacOS. For Linux, you can simply use `/etc/hosts`. 
- [Kafka](https://kafka.apache.org/) for event streaming between microservices
- [Prometheus](https://prometheus.io/) for monitoring
- [Grafana](https://grafana.com/) for visualization


## Getting Started

### Mac OS

1. Install [Homebrew](https://brew.sh/)
2. Install [Docker Desktop](https://www.docker.com/products/docker-desktop)
3. Install [Minikube](https://minikube.sigs.k8s.io/docs/start/)
4. Install [Kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
5. Install [Helm](https://helm.sh/docs/intro/install/)
6. Install [Skaffold](https://skaffold.dev/docs/install/)
7. Start Minikube: 
    ```bash 
    minikube start --driver=docker --kubernetes-version=v1.30.1 --bootstrapper=kubeadm --extra-config=kubelet.authentication-token-webhook=true --extra-config=kubelet.authorization-mode=Webhook --extra-config=scheduler.bind-address=0.0.0.0 --extra-config=controller-manager.bind-address=0.0.0.0
    ```
8. Install Ingress DNS:
    ```bash
    minikube addons enable ingress
    minikube addons enable ingress-dns
    minikube addons disable metrics-server
    ```
9. Run the following command to add the DNS names of your ingress controller to your local DNS (refer to the [Ingress DNS](https://minikube.sigs.k8s.io/docs/handbook/addons/ingress-dns/) documentation for more information):
    ```bash
        sudo bash -c 'cat <<EOF > /etc/resolver/minikube-test
        domain test
        nameserver $(minikube ip)
        search_order 1
        timeout 5
        EOF'
    ```
10. Install the local helm chart which provisions kubernetes resources that are intended to be long-running and not shut down when skaffold is interrupted.

