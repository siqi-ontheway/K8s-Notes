# K8s-Notes

## Goals
 - Trying your hand and writing your own kubernetes controllers,
 - Learning how logging, debugging, etc. works on k8s


minikube start
kubectl create ns eksposetest
kubectl create deployment nginx -n eksposetest --image nginx
kubectl delete deployment -n eksposetest nginx
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v0.46.0/deploy/static/provider/do/deploy.yaml
kubectl delete all --all -n ekspose

kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v0.46.0/deploy/static/provider/do/deploy.yaml
kubectl get all -n ingress-nginx
kubectl get nodes
kubectl get deployments
kubectl create deployment nginx -n eksposetest --image nginx
kubectl get deployments
kubectl create deployment hello-node --image=registry.k8s.io/e2e-test-images/agnhost:2.39 -- /agnhost netexec --http-port=8080
kubectl get service -n testone
kubectl get ingress -n testone
kubectl get service -n ingress-nginx

## References
https://cloud.tencent.com/developer/article/1493250
https://github.com/kubernetes/sample-controller
https://github.com/zq2599/blog_demos/tree/master
https://blog.csdn.net/boling_cavalry/article/details/128753781
https://medium.com/speechmatics/how-to-write-kubernetes-custom-controllers-in-go-8014c4a04235
https://cloudark.medium.com/kubernetes-custom-controllers-b6c7d0668fdf#:~:text=You%20can%20write%20custom%20controllers%20that%20handle,you%20can%20add%20new%20custom%20resources%20within
https://youtu.be/lzoWSfvE2yA?si=gkFn6-qzXi2l7DuG
