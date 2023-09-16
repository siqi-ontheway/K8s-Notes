# K8s-Notes

## Goals
 - Trying your hand and writing your own kubernetes controllers,
 - Learning how logging, debugging, etc. works on k8s

## Description

A custom Kubernetes controller to watch for the creation and deletion of deployments and create/delete services and ingresses accordingly. The controller is containerized and deployed with necessary manifests related to roles and service account.

## Files

- controller related files are in "ekspose" folder.
- k8s tutorial learning notes are in "Tutorial.md".

## Deployment
```
cd ekspose
docker build -t siqili/ekspose:0.0.1 .
docker push siqili/ekspose:0.0.1 
```

## Run & Test

- Before running this controller, you need to make sure there is a k8s cluster running. You can use "minikube start" or enable docker k8s cluster for the demo.
- In order to create ingress for deployment, you need to run: kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v0.46.0/deploy/static/provider/do/deploy.yaml
- In order to deploy the controller, you can run:
    - kubectl create -f deployment.yaml
- To monitor the log you can run:
    - to get the pod: kubectl get pods -n ekspose
    - to get the log: kubectl logs -n ekspose
- To test controller, you can run:
    - kubectl create ns ekspose
    - kubectl create deployment nginx -n ekspose --image nginx
    - kubectl delete deployment -n eksposetest nginx
- To list, you can run:
    - kubectl get all -n ingress-nginx
    - kubectl get nodes
    - kubectl get deployments
    - kubectl get service -n eksposetest
    - kubectl get ingress -n eksposetest
- To clear, you can run: 
    - kubectl delete all --all -n ekspose


## Demo

Below is a demo. 

- The output ofter calling the k8s apis

![image](./img/k8sapi.png)

- The kubectl logs after calling the k8s apis

![image](./img/log.png)

## Code-generator

1. deep copy obejcts
2. clientset
3. informers
4. lister
'''
execDir=/Users/lisiqi/go/pkg/mod/k8s.io/code-generator@v0.28.1
'''

riomous-MacBook-Pro:kluster lisiqi$ "${execDir}"/generate-groups.sh all kluster/pkg/client kluster/pkg/apis siqi.dev:v1alpha1 --go-header-file /Users/lisiqi/go/pkg/mod/k8s.io/gengo@v0.0.0-20220902162205-c0856e24416d/boilerplate/boilerplate.go.txt

kubectl create secret generic dosecret --from-literal token=DOTOKEN

interface {
    CreateCluster
    DeleteCluster
    GetKubeConfig
}


kubectl patch klusters.siqi.dev/kluster-0 -p '{"metadata":{"finalizers":[]}}' --type=merge

kubectl delete klusters.siqi.dev kluster-0
kubectl create -f manifests/kluster0.yaml 
kubectl get klusters.siqi.dev

kluster.go has kluster spec, using which our cluster is created in digital ocean. To do that, so we create do.go which takes kubenetes clientset (to get the token) and klient.spec and make a request to digital ocean api

kluster status is a sub resource
, which is useful once reflected in printer column
controller-gen code add the status to cr by comments


docker build -t siqili/kluster:0.1.0 .
docker push siqili/kluster:0.1.0
kubectl create deployment kluster --image siqili/kluster:0.1.0 --dry-run=client -oyaml > install/deploy.yaml
kubectl create -f install/deploy.yaml

kubectl create serviceaccount kluster-sa --dry-run=client -oyaml > install/sa.yaml
kubectl create clusterrole kluster-cr --resource Kluster --verb list,watch --dry-run=client -oyaml > install/clusterrole.yaml
kubectl create clusterrolebinding kluster-crb --clusterrole k
luster-cr --serviceaccount default:kluster-sa --dry-run=client -oyaml > install/crb.yaml
## References
https://cloud.tencent.com/developer/article/1493250
https://github.com/kubernetes/sample-controller
https://github.com/zq2599/blog_demos/tree/master
https://blog.csdn.net/boling_cavalry/article/details/128753781
https://medium.com/speechmatics/how-to-write-kubernetes-custom-controllers-in-go-8014c4a04235
https://cloudark.medium.com/kubernetes-custom-controllers-b6c7d0668fdf#:~:text=You%20can%20write%20custom%20controllers%20that%20handle,you%20can%20add%20new%20custom%20resources%20within
https://youtu.be/lzoWSfvE2yA?si=gkFn6-qzXi2l7DuG
