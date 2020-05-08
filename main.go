/*
Copyright 2020 Kamal Nasser All rights reserved.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
