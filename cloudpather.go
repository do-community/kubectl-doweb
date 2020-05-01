package main

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
)

// CloudPather returns cloud.digitalocean.com paths for K8s resources
type CloudPather interface {
	Cluster(context.Context) (string, error)
	Node(context.Context, string) (string, error)
	Service(context.Context, string, string) (string, error)
	PersistentVolume(context.Context, string) (string, error)
	PersistentVolumeClaim(context.Context, string, string) (string, error)
}

const nodeIDPrefix = "digitalocean://"
const lbaasAnnotation = "kubernetes.digitalocean.com/load-balancer-id"
const hostnameSuffix = ".k8s.ondigitalocean.com"

type DOCloudPather struct {
	clientConfig *restclient.Config
	clientset    kubernetes.Interface
	output       io.Writer
}

var _ CloudPather = &DOCloudPather{}

func (cp *DOCloudPather) Cluster(ctx context.Context) (string, error) {
	endpoint, err := url.Parse(cp.clientConfig.Host)
	if err != nil {
		return "", err
	}
	if !strings.HasSuffix(endpoint.Host, hostnameSuffix) {
		return "", fmt.Errorf("the cluster does not seem to be a DOKS cluster")
	}

	id := strings.TrimSuffix(endpoint.Host, hostnameSuffix)
	return fmt.Sprintf("kubernetes/clusters/%s", id), nil
}

func (cp *DOCloudPather) Node(ctx context.Context, name string) (string, error) {
	node, err := cp.clientset.CoreV1().Nodes().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	if !strings.HasPrefix(node.Spec.ProviderID, nodeIDPrefix) {
		return "", fmt.Errorf("node %s is not a DigitalOcean-provisioned node", name)
	}

	id := strings.TrimPrefix(node.Spec.ProviderID, nodeIDPrefix)
	return fmt.Sprintf("droplets/%s", id), nil
}

func (cp *DOCloudPather) Service(ctx context.Context, namespace, name string) (string, error) {
	svc, err := cp.clientset.CoreV1().Services(namespace).Get(ctx, name, metav1.GetOptions{})
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

func (cp *DOCloudPather) PersistentVolume(ctx context.Context, name string) (string, error) {
	pvObj, err := cp.clientset.CoreV1().PersistentVolumes().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	pvClass := pvObj.Spec.StorageClassName
	if pvClass != "do-block-storage" {
		return "", fmt.Errorf("PersistentVolume %s is not a DigitalOcean Block Storage Volume. Storage class must be %s but got %s", "do-block-storage", pvClass)
	}

	fmt.Fprintf(cp.output, "PersistentVolume name: %s\n", pvObj.Name)
	return "volumes", nil
}

func (cp *DOCloudPather) PersistentVolumeClaim(ctx context.Context, namespace, name string) (string, error) {
	pvcObj, err := cp.clientset.CoreV1().PersistentVolumeClaims(namespace).Get(ctx, name, metav1.GetOptions{})
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

	fmt.Fprintf(cp.output, "PersistentVolume name: %s\n", pvcObj.Spec.VolumeName)
	return "volumes", nil
}
