Minikube is a lightweight Kubernetes implementation that creates a VM on your local machine and deploys a simple cluster containing only one node. 
 
minikube dashboard: By default, the dashboard is only accessible from within the internal Kubernetes virtual network. The dashboard command creates a temporary proxy to make the dashboard accessible from outside the Kubernetes virtual network.

Pod: A Kubernetes Pod is a group of one or more Containers (such as Docker), tied together for the purposes of administration and networking. One pod has one IP address.

Deployment: A Kubernetes Deployment checks on the health of your Pod and restarts the Pod's Container if it terminates. Deployments are the recommended way to manage the creation and scaling of Pods.

Service: A Service is a method for exposing a network application that is running as one or more Pods in your cluster. The Service abstraction enables this decoupling.

Labels: Each object can have a set of key/value labels defined. Each Key must be unique for a given object.

A Kubernetes cluster consists of two types of resources:

The Control Plane coordinates the cluster. The control plane schedules the containers to run on the cluster's nodes. Node-level components, such as the kubelet, communicate with the control plane using the Kubernetes API, which the control plane exposes. 

Nodes are the workers that run applications. A Node can have multiple pods, and the Kubernetes control plane automatically handles scheduling the pods across the Nodes in the cluster. 

Every Kubernetes Node runs at least:

Kubelet, a process responsible for communication between the Kubernetes control plane and the Node; it manages the Pods and the containers running on a machine.

A container runtime (like Docker) responsible for pulling the container image from a registry, unpacking the container, and running the application.

# Commanly Used kubectl Cmd

scale: kubectl scale deployments/kubernetes-bootcamp --replicas=4

list: kubectl get pods -o wide

describe: kubectl describe deployments/kubernetes-bootcamp

set image: kubectl set image deployments/kubernetes-bootcamp kubernetes-bootcamp=jocatalin/kubernetes-bootcamp:v2