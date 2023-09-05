# K8s-Notes

## Goals
 - Trying your hand and writing your own kubernetes controllers,
 - Learning how logging, debugging, etc. works on k8s

## Running 

- Before running this controller, you need to make sure there is a k8s cluster running. You can use "minikube start" or enable docker k8s cluster for the demo.
- In order to create ingress for deployment, you need to run: kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v0.46.0/deploy/static/provider/do/deploy.yaml
- To test controller, you can run:
    - kubectl create ns eksposetest
    - kubectl create deployment nginx -n eksposetest --image nginx
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

- The log file after calling the k8s apis

![image](./img/log.png)

## References
https://cloud.tencent.com/developer/article/1493250
https://github.com/kubernetes/sample-controller
https://github.com/zq2599/blog_demos/tree/master
https://blog.csdn.net/boling_cavalry/article/details/128753781
https://medium.com/speechmatics/how-to-write-kubernetes-custom-controllers-in-go-8014c4a04235
https://cloudark.medium.com/kubernetes-custom-controllers-b6c7d0668fdf#:~:text=You%20can%20write%20custom%20controllers%20that%20handle,you%20can%20add%20new%20custom%20resources%20within
https://youtu.be/lzoWSfvE2yA?si=gkFn6-qzXi2l7DuG
