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

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/do-community/kubectldoweb"

	"github.com/skratchdot/open-golang/open"
	"github.com/urfave/cli/v2"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

const cloudBase = "https://cloud.digitalocean.com/"

type opener func(input string) error

var errHelp = fmt.Errorf("errHelp")

func main() {
	runCLI(os.Args)
}

func newApp(rootCmd cli.ActionFunc) *cli.App {
	return &cli.App{
		Name:  "kubectl-doweb",
		Usage: "a kubectl plugin for opening DigitalOcean resources in a web browser",
		UsageText: `kubectl doweb <type> <name>

EXAMPLES:

   kubectl doweb service main-load-balancer
   kubectl doweb cluster

SUPPORTED TYPES:

   cluster, node (no), service (svc), persistentvolume (pv), persistentvolumeclaim (pvc)`,
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
	app := newApp(newRootCmd(kubectldoweb.Run, open.Run))

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

func newRootCmd(runner kubectldoweb.Runner, opnr opener) cli.ActionFunc {
	return func(c *cli.Context) error {
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

		path, err := runner(c.Context, os.Stderr, kubeConfig, namespace, typ, name)
		if err != nil {
			return err
		}

		url := fmt.Sprintf("%s%s", cloudBase, path)
		return opnr(url)
	}
}

func defaultKubeconfigPath() string {
	home := homedir.HomeDir()
	if home == "" {
		return ""
	}

	return filepath.Join(home, ".kube", "config")
}
