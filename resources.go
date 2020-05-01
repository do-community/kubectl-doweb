package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
)

const nodeIDPrefix = "digitalocean://"
const lbaasAnnotation = "kubernetes.digitalocean.com/load-balancer-id"
const hostnameSuffix = ".k8s.ondigitalocean.com"

type Resource interface {
	CloudPath(context.Context, *restclient.Config, kubernetes.Interface, string, string) (string, error)
}

type Cluster struct{}
type Node struct{}
type Service struct{}
type PersistentVolume struct{}
type PersistentVolumeClaim struct{}

var (
	_ Resource = &Cluster{}
	_ Resource = &Node{}
	_ Resource = &Service{}
	_ Resource = &PersistentVolume{}
	_ Resource = &PersistentVolumeClaim{}
)

var namesMap = map[string]Resource{
	"cluster":                &Cluster{},
	"nodes":                  &Node{},
	"node":                   &Node{},
	"no":                     &Node{},
	"services":               &Service{},
	"service":                &Service{},
	"svc":                    &Service{},
	"persistentvolume":       &PersistentVolume{},
	"persistentvolumes":      &PersistentVolume{},
	"pv":                     &PersistentVolume{},
	"persistentvolumeclaim":  &PersistentVolumeClaim{},
	"persistentvolumeclaims": &PersistentVolumeClaim{},
	"pvc":                    &PersistentVolumeClaim{},
}

func ParseResource(name string) Resource {
	if r, ok := namesMap[name]; ok {
		return r
	}

	return nil
}

func (c *Cluster) CloudPath(ctx context.Context, clientConfig *restclient.Config, clientset kubernetes.Interface, namespace, name string) (string, error) {
	endpoint, err := url.Parse(clientConfig.Host)
	if err != nil {
		return "", err
	}
	if !strings.HasSuffix(endpoint.Host, hostnameSuffix) {
		return "", fmt.Errorf("the cluster does not seem to be a DOKS cluster")
	}

	id := strings.TrimSuffix(endpoint.Host, hostnameSuffix)
	return fmt.Sprintf("kubernetes/clusters/%s", id), nil
}

func (n *Node) CloudPath(ctx context.Context, clientConfig *restclient.Config, clientset kubernetes.Interface, namespace, name string) (string, error) {
	node, err := clientset.CoreV1().Nodes().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	if !strings.HasPrefix(node.Spec.ProviderID, nodeIDPrefix) {
		return "", fmt.Errorf("node %s is not a DigitalOcean-provisioned node", name)
	}

	id := strings.TrimPrefix(node.Spec.ProviderID, nodeIDPrefix)
	return fmt.Sprintf("droplets/%s", id), nil
}

func (s *Service) CloudPath(ctx context.Context, clientConfig *restclient.Config, clientset kubernetes.Interface, namespace, name string) (string, error) {
	svc, err := clientset.CoreV1().Services(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	svcType := svc.Spec.Type
	if svcType != corev1.ServiceTypeLoadBalancer {
		return "", fmt.Errorf("Service %s is of the type %s, not a LoadBalancer", name, svcType)
	}

	id, ok := svc.Annotations[lbaasAnnotation]
	if !ok {
		return "", fmt.Errorf("annotation %s not found on service", lbaasAnnotation)
	}

	return fmt.Sprintf("networking/load_balancers/%s", id), nil
}

func (pv *PersistentVolume) CloudPath(ctx context.Context, clientConfig *restclient.Config, clientset kubernetes.Interface, namespace, name string) (string, error) {
	pvObj, err := clientset.CoreV1().PersistentVolumes().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	pvClass := pvObj.Spec.StorageClassName
	if pvClass != "do-block-storage" {
		return "", fmt.Errorf("PersistentVolume %s is not a DigitalOcean Block Storage Volume. Storage class must be %s but got %s", "do-block-storage", pvClass)
	}

	// TODO: extract os.Stdout
	fmt.Fprintf(os.Stdout, "PersistentVolume name: %s\n", pvObj.Name)
	return "volumes", nil
}

func (pvc *PersistentVolumeClaim) CloudPath(ctx context.Context, clientConfig *restclient.Config, clientset kubernetes.Interface, namespace, name string) (string, error) {
	pvcObj, err := clientset.CoreV1().PersistentVolumeClaims(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	pvcPhase := pvcObj.Status.Phase
	if pvcPhase != corev1.ClaimBound {
		return "", fmt.Errorf("PersistentVolumeClaim %s is not bound to a PersistentVolume. Got phase %s", name, pvcPhase)
	}

	pvClass := *pvcObj.Spec.StorageClassName
	if pvClass != "do-block-storage" {
		return "", fmt.Errorf("PersistentVolume %s is not a DigitalOcean Block Storage Volume. Storage class must be %s but got %s", name, "do-block-storage", pvClass)
	}

	fmt.Fprintf(os.Stdout, "PersistentVolume name: %s\n", pvcObj.Spec.VolumeName)
	return "volumes", nil
}
