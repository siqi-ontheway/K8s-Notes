package controller

import (
	"context"
	"fmt"
	"time"

	"kluster/pkg/apis/siqi.dev/v1alpha1"
	klientset "kluster/pkg/client/clientset/versioned"
	skeme "kluster/pkg/client/clientset/versioned/scheme"
	kinf "kluster/pkg/client/informers/externalversions/siqi.dev/v1alpha1"
	klister "kluster/pkg/client/listers/siqi.dev/v1alpha1"
	"kluster/pkg/do"

	"github.com/kanisterio/kanister/pkg/poll"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
)

type controller struct {
	client        kubernetes.Interface            /* Client set to store and pass the secrete token */
	klient        klientset.Interface             /* Customized crd kluster klient */
	kLister       klister.KlusterLister           /* Component of informer to get the resources from cache */
	klusterSynced cache.InformerSynced            /* To get Status that if the cache is successfully synced, passed from reflector */
	queue         workqueue.RateLimitingInterface /* FIFO queue so we can add objects to queue when Add/delete functions are called */
	recorder      record.EventRecorder            /* Event recorder for the cr */
}

var clusterID string

// Create new controllers
func NewController(client kubernetes.Interface, klient klientset.Interface, klusterInformer kinf.KlusterInformer) *controller {
	runtime.Must(skeme.AddToScheme(scheme.Scheme))
	eveBroadCaster := record.NewBroadcaster()
	eveBroadCaster.StartStructuredLogging(0)
	eveBroadCaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{
		Interface: client.CoreV1().Events(""),
	})
	recorder := eveBroadCaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: "Kluster"})

	c := &controller{
		client:        client,
		klient:        klient,
		kLister:       klusterInformer.Lister(),
		klusterSynced: klusterInformer.Informer().HasSynced,
		queue:         workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "kluster"),
		recorder:      recorder,
	}

	// Register functions in informer to handle add/delete events
	klusterInformer.Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc:    c.handleAdd,
			DeleteFunc: c.handleDel,
		},
	)
	return c
}

// Run controllers: sync cache and keep running workers until the channel is closed.
func (c *controller) Run(workers int, ch <-chan struct{}) error {
	defer runtime.HandleCrash()

	// As long as the Shutdown is called, the processItem method will return false
	defer c.queue.ShutDown()
	klog.Infof("start controller")

	// Make sure informer cache has been synced
	if !cache.WaitForCacheSync(ch, c.klusterSynced) {
		klog.Errorf("failed to wait for caches to sync")
		return fmt.Errorf("failed to wait for caches to sync")
	}

	for i := 0; i < workers; i++ {
		go wait.Until(c.worker, 1*time.Second, ch)
	}

	// go routine is non-blocking.
	// This step is to make the channel keep waiting
	// if we do not put something into the channel.
	<-ch
	return nil
}

// The worker will keep running processItem until it returns false
func (c *controller) worker() {
	for c.processItem() {

	}
}

// Poll from the queue and then sync deployment
// Production: Refer to claimWorker method in kubenetes library
func (c *controller) processItem() bool {
	item, shutdown := c.queue.Get()
	if shutdown {
		return false
	}

	err := func(obj interface{}) error {
		defer c.queue.Done(obj)

		if err := c.syncHandler(obj); err != nil {
			return fmt.Errorf("error syncing '%s': %s", obj, err.Error())
		}

		c.queue.Forget(obj)
		klog.Infof("Successfully synced")
		return nil
	}(item)

	if err != nil {
		runtime.HandleError(err)
		return true
	}

	return true
}

// Handle add and delete event sync
func (c *controller) syncHandler(item interface{}) error {
	key, err := cache.MetaNamespaceKeyFunc(item)
	if err != nil {
		klog.Errorf("getting key from cache %s\n", err.Error())
	}
	ns, name, err := cache.SplitMetaNamespaceKey(key)

	if err != nil {
		klog.Errorf("splitting key into namespace and name %s\n", err.Error())
		return err
	}

	// Check if the object has been deleted from k8s cluster
	kluster, err := c.kLister.Klusters(ns).Get(name)
	if err != nil {

		// If error is that the object is not found in k8s cluster, i.e. the event is deletion
		if apierrors.IsNotFound(err) {
			klog.Infof("kluster %s was deleted\n", name)

			return deleteDOCluster(clusterID)
		}
		klog.Errorf("error %s, Getting the kluster resource from lister", err.Error())
		return err
	}

	klog.Infof("kluster spec that we have is %+v\n", kluster.Spec)

	clusterID, err = do.Create(c.client, kluster.Spec)
	klog.Infof("clusterID is %+s\n", clusterID)
	if err != nil {
		klog.Errorf("error %s, creating the cluster\n", err.Error())
		c.retry(err, item)
		return err
	}

	c.recorder.Event(kluster, corev1.EventTypeNormal, "ClusterCreation", "DO API was called to create the cluster")

	err = c.updateStatus(clusterID, "creating", kluster)
	if err != nil {
		klog.Errorf("error %s, updating the status of the kluster %s\n", err.Error(), kluster.Name)
		return err
	}

	// Query DO API to make sure the cluster is created
	err = c.waitForCluster(kluster.Spec, clusterID)
	if err != nil {
		klog.Errorf("Cluster is already deleted")
		return err
	}

	err = c.updateStatus(clusterID, "running", kluster)
	if err != nil {
		// In prod env, we need to retry if a kluster is not created successfully
		klog.Errorf("error %s, updating the status of the kluster %s after waiting\n", err.Error(), kluster.Name)
		return err
	}

	c.recorder.Event(kluster, corev1.EventTypeNormal, "ClusterCreationCompleted", "DO cluster creation was completed")

	return nil
}

// Retry for five times if failed to sync deployment
func (c *controller) retry(err error, key interface{}) {
	if err == nil {
		// Item is successfully processed.
		c.queue.Forget(key)
		return
	}
	// If item is not successfully processed,
	// check how many times you have retried until 5 times
	if c.queue.NumRequeues(key) < 5 {
		klog.Infof("Error syncing: %v\n", err)
		c.queue.AddRateLimited(key)
		return
	}

	// If you have reached the 5 limit times of retry, forget this item
	c.queue.Forget(key)
	// report error
	runtime.HandleError(err)
	klog.Errorf("Dropping pod %q out of the queue: %v", key, err)
}

// Delete actual cluster from digital ocean
func deleteDOCluster(klusterID string) error {
	err := do.Delete(klusterID)
	klog.Infof("ClusterID: %s", klusterID)
	if err != nil {
		klog.Errorf("error %s, destroying the cluster\n", err.Error())
		return err
	}
	klog.Infof("Cluster was deleted succcessfully")

	return nil
}

// Wait for cluster to finish creating
func (c *controller) waitForCluster(spec v1alpha1.KlusterSpec, clusterID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	return poll.Wait(ctx, func(ctx context.Context) (bool, error) {
		state, err := do.ClusterState(clusterID)
		if err != nil {
			return false, err
		}
		if state == "running" {
			return true, nil
		}

		return false, nil
	})
}

// Update the latest status of a kluster
func (c *controller) updateStatus(id, progress string, kluster *v1alpha1.Kluster) error {
	// get the latest version of kluster, or there would be error when  fetching the object after it is updated
	k, err := c.klient.SiqiV1alpha1().Klusters(kluster.Namespace).Get(context.Background(), kluster.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	k.Status.KlusterID = id
	k.Status.Progress = progress
	_, err = c.klient.SiqiV1alpha1().Klusters(kluster.Namespace).UpdateStatus(context.Background(), k, metav1.UpdateOptions{})
	return err
}

// Add handler: Add obj to queue
func (c *controller) handleAdd(obj interface{}) {
	klog.Infof("Add called")
	// Add obj to queue
	c.queue.Add(obj)
}

// Del handler: Add obj to queue
func (c *controller) handleDel(obj interface{}) {
	klog.Infof("Del called")
	// Add obj to queue
	c.queue.Add(obj)
}
