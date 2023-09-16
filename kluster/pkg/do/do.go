package do

import (
	"context"
	"fmt"
	"strings"

	"kluster/pkg/apis/siqi.dev/v1alpha1"

	"github.com/digitalocean/godo"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var Token string

// Create digital ocean cluster
func Create(c kubernetes.Interface, spec v1alpha1.KlusterSpec) (string, error) {
	token, err := getToken(c, spec.TokenSecret)
	if err != nil {
		return "", err
	}
	Token = token

	client := godo.NewFromToken(token)
	fmt.Println(client)
	request := &godo.KubernetesClusterCreateRequest{
		Name:        spec.Name,
		RegionSlug:  spec.Region,
		VersionSlug: spec.Version,
		NodePools: []*godo.KubernetesNodePoolCreateRequest{
			&godo.KubernetesNodePoolCreateRequest{
				Size:  spec.NodePools[0].Size,
				Name:  spec.NodePools[0].Name,
				Count: spec.NodePools[0].Count,
			},
		},
	}

	cluster, _, err := client.Kubernetes.Create(context.Background(), request)
	if err != nil {
		return "", err
	}

	return cluster.ID, nil
}

// Get digital ocean cluster status
func ClusterState(id string) (string, error) {
	client := godo.NewFromToken(Token)
	cluster, _, err := client.Kubernetes.Get(context.Background(), id)
	if err != nil {
		return "", err
	}
	return string(cluster.Status.State), err
}

// Get token from secretes of existing clusters
func getToken(client kubernetes.Interface, sec string) (string, error) {
	namespace := strings.Split(sec, "/")[0]
	name := strings.Split(sec, "/")[1]
	s, err := client.CoreV1().Secrets(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	return string(s.Data["token"]), nil
}

// Delete digital ocean cluster
func Delete(id string) error {
	client := godo.NewFromToken(Token)
	_, err := client.Kubernetes.Delete(context.TODO(), id)
	if err != nil {
		return err
	}

	return nil
}
