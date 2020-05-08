```
NAME:
   kubectl-doweb - a kubectl plugin for opening DigitalOcean resources in a web browser

USAGE:
   kubectl doweb <type> <name>

EXAMPLES:

   kubectl doweb service main-load-balancer
   kubectl doweb cluster

SUPPORTED TYPES:

   cluster, node, service, persistentvolume, persistentvolumeclaim

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --kubeconfig value           absolute path to the kubeconfig file (default: "$HOME/.kube/config")
   --namespace value, -n value  kubernetes object namespace (default: default namespace in kubeconfig)
   --help, -h                   show help (default: false)
```