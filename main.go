package kubectldoweb

import (
	"fmt"
	"io"

	"golang.org/x/net/context"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type Runner func(ctx context.Context, writer io.Writer, kubeConfig clientcmd.ClientConfig, namespace, typ, name string) (string, error)

const cloudBase = "https://cloud.digitalocean.com/"

var ErrMissingArgument = fmt.Errorf("missing argument")

func Run(ctx context.Context, writer io.Writer, kubeConfig clientcmd.ClientConfig, namespace, typ, name string) (string, error) {
	// if a namespace is not explicitly provided, use the default set in kube config
	if namespace == "" {
		namespace, _, _ = kubeConfig.Namespace()
	}
	if namespace == "" {
		fmt.Println("could not determine namespace using the provided kube config")
		return "", ErrMissingArgument
	}

	clientConfig, err := kubeConfig.ClientConfig()
	if err != nil {
		return "", err
	}

	clientset, err := kubernetes.NewForConfig(clientConfig)
	if err != nil {
		return "", err
	}

	cp := &DOCloudPather{
		clientConfig: clientConfig,
		clientset:    clientset,
		output:       writer,
	}

	fmt.Fprintf(writer, "opening %s %s (namespace %s)\n", typ, name, namespace)
	path, err := cloudPatherByType(ctx, cp, typ, namespace, name)
	if err != nil {
		return "", err
	}

	return path, nil
}

func cloudPatherByType(ctx context.Context, cp CloudPather, typ, namespace, name string) (string, error) {
	// cluster is the only type that doesn't take a name
	if typ == "cluster" {
		return cp.Cluster(ctx)
	}

	if name == "" {
		return "", ErrMissingArgument
	}

	switch typ {
	case "nodes":
		fallthrough
	case "node":
		fallthrough
	case "no":
		return cp.Node(ctx, name)

	case "services":
		fallthrough
	case "service":
		fallthrough
	case "svc":
		return cp.Service(ctx, namespace, name)

	case "persistentvolume":
		fallthrough
	case "persistentvolumes":
		fallthrough
	case "pv":
		return cp.PersistentVolume(ctx, name)

	case "persistentvolumeclaim":
		fallthrough
	case "persistentvolumeclaims":
		fallthrough
	case "pvc":
		return cp.PersistentVolumeClaim(ctx, namespace, name)

	default:
		return "", fmt.Errorf("unknown type %s", typ)
	}
}
