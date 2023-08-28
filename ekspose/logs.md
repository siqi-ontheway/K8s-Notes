View the container logs

Anything that the application would normally send to standard output becomes logs for the container within the Pod. We can retrieve these logs using the kubectl logs command:

kubectl logs "$POD_NAME"