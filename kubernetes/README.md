# What is this?

This is a list of things to do in order to get a kubernetes pod up and running with multiple replicas of the servers.

## Dependencies (for mac)
- docker
- VirtualBox
- minikube
- kubectl
- helm

# Steps (after installing all the dependencies):
1. `minikube start`
2. `eval $(minikube docker-env)`
2. `minikube addons enable ingress`
3. `helm init`
4. `helm dep update`
4. `minikube ip`
5. `helm install ./ -n go-rest-service-example`

# Generating a new helm chart:
`helm create <service-name>`
