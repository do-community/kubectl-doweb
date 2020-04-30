package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/skratchdot/open-golang/open"
	"github.com/urfave/cli/v2"
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
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func rootCmd(c *cli.Context) error {
	if c.Args().Len() != 2 {
		return fmt.Errorf("usage: kubectl dobrowse <type> <name>\n\n\texample: kubectl dobrowse service main-load-balancer")
	}

	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: c.String("kubeconfig")},
		&clientcmd.ConfigOverrides{},
	)

	namespace := c.String("namespace")
	if namespace == "" {
		namespace, _, _ = kubeConfig.Namespace()
	}
	typ := c.Args().Get(0)
	name := c.Args().Get(1)

	resource := ParseResource(typ)
	if resource == nil {
		return fmt.Errorf("unsupported resource type %s", typ)
	}

	fmt.Printf("opening %s %s (namespace %s)\n", typ, name, namespace)
	return browse(kubeConfig, resource, namespace, name)
}

func browse(kubeConfig clientcmd.ClientConfig, resource Resource, namespace, name string) error {
	clientConfig, err := kubeConfig.ClientConfig()
	if err != nil {
		return err
	}

	clientset, err := kubernetes.NewForConfig(clientConfig)
	if err != nil {
		return err
	}

	path, err := resource.CloudPath(context.Background(), clientset, namespace, name)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s%s", cloudBase, path)
	open.Run(url)

	return nil
}

func defaultKubeconfigPath() string {
	home := homedir.HomeDir()
	if home == "" {
		return ""
	}

	return filepath.Join(home, ".kube", "config")
}
