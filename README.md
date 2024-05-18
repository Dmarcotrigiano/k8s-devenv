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

