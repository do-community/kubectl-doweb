# 🗺 kubectl-doweb

kubectl-doweb is a kubectl plugin for opening DigitalOcean resources in a web browser. For example, if you have a LoadBalancer Service in your DOKS cluster, this plugin will let you open the corresponding page of the Load Balancer in the DigitalOcean Control Panel.

Usage: `kubectl doweb <resource type> <resource name>`

## Supported Resources

| Resource              | Target Page                                                  |
| --------------------- | ------------------------------------------------------------ |
| Cluster               | Overview page of the DOKS cluster                            |
| Node                  | The Droplet page of a specific worker node                   |
| Service               | LoadBalancer services only. Opens the underlying DigitalOcean Load Balancer |
| PersistentVolume      | Prints the Volume name and opens the Volumes page            |
| PersistentVolumeClaim | If bound, prints the Volume name and opens the Volumes page  |

## Installation

1. Download the latest release from [the releases page](https://github.com/do-community/kubectl-dobrowse/releases).
2. Extract the downloaded archive and move the binary `kubectl-doweb` to any directory in your `$PATH` such as `~/bin` if it exists or `/usr/local/bin`. This allows `kubectl` to detect it as a plugin.
3. To use, run `kubectl doweb`.

You can verify that it is installed properly by running `kubectl plugin list` and looking for `kubectl-dobrowse` in the output.

## Usage

Run `kubectl doweb <type> <name>`,`<type>` being the resource type and `<name>` being the resource name. The supported types are: `cluster, node (no), service (svc), persistentvolume (pv), persistentvolumeclaim (pvc)`.

The default namespace is used. To set a different namespace, use the `--namespace` or `-n` option.

Examples:

* `kubectl doweb cluster`
* `kubectl doweb node pool-c0yaq2bd6-95th`
* `kubectl doweb --namespace nginx-ingress service nginx-ingress`
* `kubectl doweb pvc kibana-data-01`

kubectl-doweb attempts to use the kube config file found in `$HOME/.kube/config`. To set a different path, use the `--kubeconfig` option.

---

```
NAME:
   kubectl-doweb - a kubectl plugin for opening DigitalOcean resources in a web browser

USAGE:
   kubectl doweb <type> <name>

EXAMPLES:

   kubectl doweb service main-load-balancer
   kubectl doweb cluster

SUPPORTED TYPES:

   cluster, node (no), service (svc), persistentvolume (pv), persistentvolumeclaim (pvc)

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --kubeconfig value           absolute path to the kubeconfig file (default: "$HOME/.kube/config")
   --namespace value, -n value  kubernetes object namespace (default: default namespace in kubeconfig)
   --help, -h                   show help (default: false)
```
