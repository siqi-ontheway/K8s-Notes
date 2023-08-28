package main

import (
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
	appsinformers "k8s.io/client-go/informers/apps/v1"
	"k8s.io/client-go/kubernetes"
	appslisters "k8s.io/client-go/listers/apps/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type controller struct {
	clientset     kubernetes.Interface            /* Client set to interact with k8s cluster */
	depLister     appslisters.DeploymentLister    /* Component of informer to get the resources from cache */
	depCacheSyncd cache.InformerSynced            /* To get Status that if the cache is successfully synced, passed from reflector */
	queue         workqueue.RateLimitingInterface /* FIFO queue so we can add objects to queue when Add/delete functions are called */
}

// Create new controllers
func newController(clientset kubernetes.Interface, depInformer appsinformers.DeploymentInformer) *controller {
	c := &controller{
		clientset:     clientset,
		depLister:     depInformer.Lister(),
		depCacheSyncd: depInformer.Informer().HasSynced,
		queue:         workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "ekspose"),
	}

	// Register functions in informer to handle add/delete events
	depInformer.Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc:    handleAdd,
			DeleteFunc: handleDel,
		},
	)
	return c
}

// Run controllers: sync cache and keep running workers until the channel is closed.
func (c *controller) run(ch <-chan struct{}) {
	fmt.Println("start controller")

	// Make sure informer cache has been synced
	if !cache.WaitForCacheSync(ch, c.depCacheSyncd) {
		fmt.Printf("waiting for cache to be synced\n")
	}

	go wait.Until(c.worker, 1*time.Second, ch)

	// go routine is non-blocking.
	// This step is to make the cahnnel keep waiting
	// if we do not put something into the cahnnel.
	<-ch
}

func (c *controller) worker() {

}

// Add handler: Add obj to queue
func handleAdd(obj interface{}) {
	fmt.Println("Add called")
}

// Del handler: Add obj to queue
func handleDel(obj interface{}) {
	fmt.Println("Del called")
}
