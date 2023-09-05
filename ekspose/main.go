package main

import (
	"flag"
	"time"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
)

func main() {

	klog.InitFlags(nil)

	flag.Set("logtostderr", "false")
	flag.Set("log_file", "/Users/lisiqi/Desktop/K8s-Notes/ekspose/log")
	flag.Parse()

	// Read kubeconfig (yaml) file to build kubenetes configuration from file
	kubeconfig := flag.String("kubeconfig", "/Users/lisiqi/.kube/config", "location to your kubeconfig file")
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		// Handle error if kubeconfig failed to build
		klog.Errorf("erorr %s building config from flags\n", err.Error())
		config, err = rest.InClusterConfig()
		if err != nil {
			klog.Errorf("error %s, getting inclusterconfig", err.Error())
		}
	}

	// Create clientset to manage resources and monitor the state of the cluster
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		// Handle error if client set failed to build
		klog.Errorf("error %s, creating clientset\n", err.Error())
	}

	// Create informers to cache reosurces and call k8s API. It watches for updates to k8s resources (add/delete)
	// They keep in-mem local cache of resources, which can be retrieved by a given index.
	// They refresh the cache using two mechanisms: List and Watch. Here the sync period is every 10 minutes.
	// The reason we use shared factory is so that one informer instance is shared for all namespaces.
	informers := informers.NewSharedInformerFactory(clientset, 10*time.Minute)

	// Create controller that includes params passed from the clientset and the informer (with local cache of resources and lister)
	c := newController(clientset, informers.Apps().V1().Deployments())
	ch := make(chan struct{})

	// Start informers, handled in goroutine chanels
	informers.Start(ch)
	// Run controlelrs, running two workers in parallel to handle events in passed channels
	if err = c.run(2, ch); err != nil {
		klog.Errorf("Error running controller: %s", err.Error())
	}

}
