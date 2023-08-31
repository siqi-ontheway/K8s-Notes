package main

import (
	"context"
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	netV1 "k8s.io/api/networking/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
			AddFunc:    c.handleAdd,
			DeleteFunc: c.handleDel,
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

// The worker will keep running processItem until it returns false
func (c *controller) worker() {
	for c.processItem() {

	}
}

// Poll from the queue and then sync deployment
func (c *controller) processItem() bool {
	item, shutdown := c.queue.Get()
	if shutdown {
		return false
	}
	defer c.queue.Forget(item)
	key, err := cache.MetaNamespaceKeyFunc(item)
	if err != nil {
		fmt.Printf("getting key from cache %s\n", err.Error())
	}
	ns, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		fmt.Printf("splitting key into namespace and name %s\n", err.Error())
		return false
	}

	err = c.syncDeployment(ns, name)
	if err != nil {
		//retry

		fmt.Printf("syncing deployment %s\n", err.Error())
		return false
	}
	return true
}

// Sync deployment: create service and ingress
func (c *controller) syncDeployment(ns, name string) error {
	ctx := context.Background()

	dep, err := c.depLister.Deployments(ns).Get(name)
	if err != nil {
		fmt.Printf("getting deployment from lister %s\n", err.Error())
	}

	// create service
	// We have to modify this, to figure out the port
	// the deployment container is listening on.
	svc := coreV1.Service{
		ObjectMeta: metaV1.ObjectMeta{
			Name:      dep.Name,
			Namespace: ns,
		},
		Spec: coreV1.ServiceSpec{
			// If without a selector, it will have errors
			Selector: depLables(*dep),
			Ports: []coreV1.ServicePort{
				{
					Name: "http",
					Port: 80,
				},
			},
		},
	}
	s, err := c.clientset.CoreV1().Services(ns).Create(ctx, &svc, metaV1.CreateOptions{})
	if err != nil {
		fmt.Printf("creating service %s\n", err.Error())
	}

	//create ingress
	return createIngress(ctx, c.clientset, s)
}

// Create ingress resources after service is created
func createIngress(ctx context.Context, client kubernetes.Interface, svc *coreV1.Service) error {

	pathType := "Prefix"
	ingress := netV1.Ingress{
		TypeMeta: metaV1.TypeMeta{},
		ObjectMeta: metaV1.ObjectMeta{
			Name:      svc.Name,
			Namespace: svc.Namespace,
			Annotations: map[string]string{
				"nginx.ingress.kubernetes.io/rewrite-target": "/",
			},
		},
		Spec: netV1.IngressSpec{
			Rules: []netV1.IngressRule{
				{
					IngressRuleValue: netV1.IngressRuleValue{
						HTTP: &netV1.HTTPIngressRuleValue{
							Paths: []netV1.HTTPIngressPath{
								{
									Path:     fmt.Sprintf("/%s", svc.Name),
									PathType: (*netV1.PathType)(&pathType),
									Backend: netV1.IngressBackend{
										Service: &netV1.IngressServiceBackend{
											Name: svc.Name,
											Port: netV1.ServiceBackendPort{
												Number: 80,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		Status: netV1.IngressStatus{},
	}
	_, err := client.NetworkingV1().Ingresses(svc.Namespace).Create(ctx, &ingress, metaV1.CreateOptions{})
	return err
}

// Add handler: Add obj to queue
func (c *controller) handleAdd(obj interface{}) {
	fmt.Println("Add called")
	// Add obj to queue
	c.queue.Add(obj)
}

// Del handler: Add obj to queue
func (c *controller) handleDel(obj interface{}) {
	fmt.Println("Del called")
	// Add obj to queue
	c.queue.Add(obj)
}

func depLables(dep appsv1.Deployment) map[string]string {
	return dep.Spec.Template.Labels
}
