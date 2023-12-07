<br/>
<p align="center">
  <!-- <a href="https://github.com/dhenkel92/ReadME-Generator"> -->
    <!-- <img src="images/logo.png" alt="Logo" width="80" height="80"> -->
  <!-- </a> -->

  <h3 align="center">Kubectl PDB & Eviction plugin</h3>

  <p align="center">
    A kubectl plugin to work with pod disruption budgets.
    <br/>
    <br/>
  </p>
</p>



## Table Of Contents

* [About the Project](#about-the-project)
* [Getting Started](#getting-started)
  * [Installation](#installation)
* [Usage](#usage)
* [Contributing](#contributing)
* [License](#license)
* [Authors](#authors)
* [Acknowledgements](#acknowledgements)

## About The Project

Pod Disruption Budgets (PDBs) help limit the number of concurrent disruptions your application experiences.
This enhances availability while allowing the cluster administrator to manage the cluster nodes.
Unfortunately, kubectl provides limited tooling to interact with PDBs.

This plugin aims to address this issue with the following key features:

- Lists all PDBs matching a given workload.
- Lists all workload pods matching a given PDB.
- Creates new PDBs from the command line.
- Evict a workload from a node.

## Getting Started

To get a local copy up and running follow these simple example steps.

### Installation

### MacOS

You can install the tool via Homebrew and the tap repository can be found [here.](https://github.com/dhenkel92/homebrew-tap)
```
brew update
brew install dhenkel92/homebrew-tap/kubectl-debug-pdb
brew install dhenkel92/homebrew-tap/kubectl-evict
```

In order to get a newer version, just upgrade via Homebrew
```
brew upgrade kubectl-debug-pdb
brew upgrade kubectl-evict
```

### Other distributions

See the [Releases page](https://github.com/dhenkel92/kubectl-debug-pdb/releases) for a list of Peek packages for various distributions.

## Usage

### Kubectl pdb

```
Utility to work with pod disruption budgets

Usage:
  pdb [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  cover       Shows which PDBs are a workload.
  create      Create a new PDB for a given workload.
  help        Help about any command
  pods        List pods covered by a given PDB

Flags:
      --as string                      Username to impersonate for the operation. User could be a regular user or a service account in a namespace.
      --as-group stringArray           Group to impersonate for the operation, this flag can be repeated to specify multiple groups.
      --as-uid string                  UID to impersonate for the operation.
      --cache-dir string               Default cache directory (default "/Users/daniel.henkel/.kube/cache")
      --certificate-authority string   Path to a cert file for the certificate authority
      --client-certificate string      Path to a client certificate file for TLS
      --client-key string              Path to a client key file for TLS
      --cluster string                 The name of the kubeconfig cluster to use
      --context string                 The name of the kubeconfig context to use
      --disable-compression            If true, opt-out of response compression for all requests to the server
  -h, --help                           help for pdb
      --insecure-skip-tls-verify       If true, the server's certificate will not be checked for validity. This will make your HTTPS connections insecure
      --kubeconfig string              Path to the kubeconfig file to use for CLI requests.
  -n, --namespace string               If present, the namespace scope for this CLI request
      --request-timeout string         The length of time to wait before giving up on a single server request. Non-zero values should contain a corresponding time unit (e.g. 1s, 2m, 3h). A value of zero means don't timeout requests. (default "0")
  -s, --server string                  The address and port of the Kubernetes API server
      --tls-server-name string         Server name to use for server certificate validation. If it is not provided, the hostname used to contact the server is used
      --token string                   Bearer token for authentication to the API server
      --user string                    The name of the kubeconfig user to use

Use "pdb [command] --help" for more information about a command.
```

List all pods for a given PDB:
```
kubectl pdb pods <pdb_name>
```

List all PDBs for all pods of a namespace:
```
kubectl pdb cover -n <namespace> [pod_name]
```

Create new PDB:
```
kubectl pdb create <resource_type>/<resource_name> --dry-run -o yaml
```

### Kubectl evict

```
Utility to evict a pod from a node

Usage:
  evict pod [flags]

Flags:
      --dry-run         If true, only print the object that would be sent, without sending it. (default true)
  -h, --help            help for pod
  -o, --output string   Output format. One of: json|yaml (default "json")

Global Flags:
      --as string                      Username to impersonate for the operation. User could be a regular user or a service account in a namespace.
      --as-group stringArray           Group to impersonate for the operation, this flag can be repeated to specify multiple groups.
      --as-uid string                  UID to impersonate for the operation.
      --cache-dir string               Default cache directory (default "/Users/daniel.henkel/.kube/cache")
      --certificate-authority string   Path to a cert file for the certificate authority
      --client-certificate string      Path to a client certificate file for TLS
      --client-key string              Path to a client key file for TLS
      --cluster string                 The name of the kubeconfig cluster to use
      --context string                 The name of the kubeconfig context to use
      --disable-compression            If true, opt-out of response compression for all requests to the server
      --insecure-skip-tls-verify       If true, the server's certificate will not be checked for validity. This will make your HTTPS connections insecure
      --kubeconfig string              Path to the kubeconfig file to use for CLI requests.
  -n, --namespace string               If present, the namespace scope for this CLI request
      --request-timeout string         The length of time to wait before giving up on a single server request. Non-zero values should contain a corresponding time unit (e.g. 1s, 2m, 3h). A value of zero means don't timeout requests. (default "0")
  -s, --server string                  The address and port of the Kubernetes API server
      --tls-server-name string         Server name to use for server certificate validation. If it is not provided, the hostname used to contact the server is used
      --token string                   Bearer token for authentication to the API server
      --user string                    The name of the kubeconfig user to use
```

Evict a pod in dry-run mode (Default: true):
```
kubectl evict pod [-n <namespace>] <pod_name>
```

Trigger real pod eviction:
```
kubectl evict pod [-n <namespace>] <pod_name> --dry-run=false
```

## Contributing

### Creating A Pull Request

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

Distributed under the MIT License. See [LICENSE](https://github.com/dhenkel92/kubectl-debug-pdb/blob/main/LICENSE) for more information.
