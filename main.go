package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/skratchdot/open-golang/open"
	"github.com/urfave/cli/v2"
	"golang.org/x/net/context"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

const cloudBase = "https://cloud.digitalocean.com/"

var errHelp = fmt.Errorf("errHelp")

func main() {
	run(os.Args)
}

func run(args []string) {
	app := &cli.App{
		Name:  "kubectl-dobrowse",
		Usage: "a kubectl plugin for opening DigitalOcean resources in a web browser",
		UsageText: `kubectl dobrowse <type> <name>

EXAMPLES:

   kubectl dobrowse service main-load-balancer
   kubectl dobrowse cluster

SUPPORTED TYPES:

   cluster, node, service, persistentvolume, persistentvolumeclaim`,
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

	err := app.Run(args)
	if err != nil {
		if err == errHelp {
			app.Run([]string{"", "help"})
		} else {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(1)
	}
}

func rootCmd(c *cli.Context) error {
	if c.Args().Len() < 1 {
		return errHelp
	}

	kubeConfigPath := c.String("kubeconfig")
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeConfigPath},
		&clientcmd.ConfigOverrides{},
	)

	namespace := c.String("namespace")
	if namespace == "" {
		namespace, _, _ = kubeConfig.Namespace()
	}
	if namespace == "" {
		fmt.Println("could not determine namespace")
		return errHelp
	}
	typ := c.Args().Get(0)
	name := c.Args().Get(1)

	clientConfig, err := kubeConfig.ClientConfig()
	if err != nil {
		return err
	}

	clientset, err := kubernetes.NewForConfig(clientConfig)
	if err != nil {
		return err
	}

	cp := &DOCloudPather{
		clientConfig: clientConfig,
		clientset:    clientset,
		output:       os.Stdout,
	}

	fmt.Printf("opening %s %s (namespace %s)\n", typ, name, namespace)
	path, err := cloudPatherWithType(context.Background(), cp, typ, namespace, name)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s%s", cloudBase, path)

	open.Run(url)
	return nil
}

func cloudPatherWithType(ctx context.Context, cp CloudPather, typ, namespace, name string) (string, error) {
	// cluster is the only type that doesn't take a name
	if typ == "cluster" {
		return cp.Cluster(ctx)
	}

	if name == "" {
		return "", errHelp
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

func defaultKubeconfigPath() string {
	home := homedir.HomeDir()
	if home == "" {
		return ""
	}

	return filepath.Join(home, ".kube", "config")
}
