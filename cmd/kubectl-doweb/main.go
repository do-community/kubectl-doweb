package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kamaln7/kubectldoweb"

	"github.com/skratchdot/open-golang/open"
	"github.com/urfave/cli/v2"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

const cloudBase = "https://cloud.digitalocean.com/"

var errHelp = fmt.Errorf("errHelp")

func main() {
	runCLI(os.Args)
}

func newApp() *cli.App {
	return &cli.App{
		Name:  "kubectl-doweb",
		Usage: "a kubectl plugin for opening DigitalOcean resources in a web browser",
		UsageText: `kubectl doweb <type> <name>

EXAMPLES:

   kubectl doweb service main-load-balancer
   kubectl doweb cluster

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
}

func runCLI(args []string) {
	app := newApp()

	err := app.Run(args)
	if err != nil {
		if err == errHelp || err == kubectldoweb.ErrMissingArgument {
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
	typ := c.Args().Get(0)
	name := c.Args().Get(1)

	path, err := kubectldoweb.Run(c.Context, os.Stderr, kubeConfig, namespace, typ, name)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s%s", cloudBase, path)
	return open.Run(url)
}

func defaultKubeconfigPath() string {
	home := homedir.HomeDir()
	if home == "" {
		return ""
	}

	return filepath.Join(home, ".kube", "config")
}
