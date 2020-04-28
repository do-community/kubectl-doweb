package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/skratchdot/open-golang/open"
	"github.com/urfave/cli/v2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

const cloudBase = "https://cloud.digitalocean.com/"

func main() {
	app := &cli.App{
		Name:   "kubectl-dobrowse",
		Usage:  "a kubectl plugin for opening DigitalOcean resources in a web browser",
		Action: rootCmd,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "kubeconfig",
				Usage: "absolute path to the kubeconfig file",
				Value: defaultKubeconfigPath(),
			},
			&cli.StringFlag{
				Name:        "namespace",
				Usage:       "kubernetes object namespace",
				Value:       "",
				Aliases:     []string{"n"},
				DefaultText: "default namespace in kubeconfig",
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func rootCmd(c *cli.Context) error {
	if c.Args().Len() != 2 {
		fmt.Printf("usage: kubectl dobrowse <type> <name>\n\n\texample: kubectl dobrowse service main-load-balancer")
		return nil
	}

	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: c.String("kubeconfig")},
		&clientcmd.ConfigOverrides{})
	clientConfig, err := kubeConfig.ClientConfig()

	resourceType, resourceName := c.Args().Get(0), c.Args().Get(1)
	resourceNamespace := c.String("namespace")
	if resourceNamespace == "" {
		resourceNamespace, _, _ = kubeConfig.Namespace()
	}
	fmt.Printf("opening %s %s (namespace %s)\n", resourceType, resourceName, resourceNamespace)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	clientset, err := kubernetes.NewForConfig(clientConfig)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	path, err := getCloudPath(clientset, resourceNamespace, resourceType, resourceName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	open.Run(fmt.Sprintf("%s%s", cloudBase, path))

	return nil
}

func getCloudPath(clientset *kubernetes.Clientset, resourceNamespace, resourceType, resourceName string) (string, error) {
	switch resourceType {
	case "node":
		fallthrough
	case "nodes":
		fallthrough
	case "no":
		node, err := clientset.CoreV1().Nodes().Get(context.TODO(), resourceName, metav1.GetOptions{})
		if err != nil {
			return "", err
		}

		if !strings.HasPrefix(node.Spec.ProviderID, "digitalocean://") {
			return "", fmt.Errorf("node %s is not a DigitalOcean-provisioned node", resourceName)
		}

		id := strings.TrimPrefix(node.Spec.ProviderID, "digitalocean://")
		return fmt.Sprintf("droplets/%s", id), nil

	case "service":
		fallthrough
	case "services":
		fallthrough
	case "svc":
		svc, err := clientset.CoreV1().Services(resourceNamespace).Get(context.TODO(), resourceName, metav1.GetOptions{})
		if err != nil {
			return "", err
		}

		svcType := svc.Spec.Type
		if svcType != corev1.ServiceTypeLoadBalancer {
			return "", fmt.Errorf("Service %s is of the type %s, not a LoadBalancer", resourceName, svcType)
		}

		id, ok := svc.Annotations["kubernetes.digitalocean.com/load-balancer-id"]
		if !ok {
			return "", fmt.Errorf("annotation kubernetes.digitalocean.com/load-balancer-id not found on service")
		}

		return fmt.Sprintf("networking/load_balancers/%s", id), nil

	case "persistentvolume":
		fallthrough
	case "persistentvolumes":
		fallthrough
	case "pv":
		pv, err := clientset.CoreV1().PersistentVolumes().Get(context.TODO(), resourceName, metav1.GetOptions{})
		if err != nil {
			return "", err
		}

		pvClass := pv.Spec.StorageClassName
		if pvClass != "do-block-storage" {
			return "", fmt.Errorf("PersistentVolume %s is not a DigitalOcean Block Storage Volume. Storage class must be %s but got %s", "do-block-storage", pvClass)
		}

		fmt.Printf("PersistentVolume name: %s\n", pv.Name)
		return "volumes", nil

	case "persistentvolumeclaim":
		fallthrough
	case "persistentvolumeclaims":
		fallthrough
	case "pvc":
		pvc, err := clientset.CoreV1().PersistentVolumeClaims(resourceNamespace).Get(context.TODO(), resourceName, metav1.GetOptions{})
		if err != nil {
			return "", err
		}

		pvcPhase := pvc.Status.Phase
		if pvcPhase != corev1.ClaimBound {
			return "", fmt.Errorf("PersistentVolumeClaim %s is not bound to a PersistentVolume. Got phase %s", resourceName, pvcPhase)
		}

		pvClass := *pvc.Spec.StorageClassName
		if pvClass != "do-block-storage" {
			return "", fmt.Errorf("PersistentVolume %s is not a DigitalOcean Block Storage Volume. Storage class must be %s but got %s", resourceName, "do-block-storage", pvClass)
		}

		fmt.Printf("PersistentVolume name: %s\n", pvc.Spec.VolumeName)
		return "volumes", nil

	default:
		return "", fmt.Errorf("unknown resource type %s", resourceType)
	}
}

func defaultKubeconfigPath() string {
	home := homedir.HomeDir()
	if home == "" {
		return ""
	}

	return filepath.Join(home, ".kube", "config")
}
