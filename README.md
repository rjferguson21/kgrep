# kgrep

A powerful command-line tool for searching and filtering Kubernetes resource manifests. `kgrep` makes it easy to find and inspect specific Kubernetes resources within YAML files, Helm charts, or piped input.

## Features

- 🔍 Search for specific Kubernetes resources by kind (Deployment, Service, etc.)
- 📝 Filter resources by name
- 📊 Generate resource summaries
- 🔄 Works with piped input from tools like `helm template` or `curl`
- 🚀 Fast and lightweight

## Requirements

- Go 1.16 or later

## Installation

### Using Go Install

```bash
go install github.com/rjferguson21/kgrep@latest
```

### Using Pre-built Binaries

You can download pre-built binaries for your platform from the [releases page](https://github.com/rjferguson21/kgrep/releases).

For example, to download and install the latest release on Linux (x86_64):

```bash
# Download the latest release
curl -L https://github.com/rjferguson21/kgrep/releases/latest/download/kgrep_Linux_x86_64.tar.gz | tar xz

# Move the binary to your PATH
sudo mv kgrep /usr/local/bin/
```

## Usage

### Basic Examples

```bash
# Search for Services within a YAML file
kgrep --kind Service all.yaml

# Search for a specific Deployment by name
kgrep --kind Deployment --name nginx-deployment manifests.yaml

# Generate a summary of all resources in a file
kgrep --summary manifests.yaml
```

### Advanced Usage

```bash
# Search within Helm chart output
helm template my-chart | kgrep --kind Deployment --name frontend

# Find resources in remote manifests
curl -s https://raw.githubusercontent.com/istio/istio/release-1.13/samples/bookinfo/platform/kube/bookinfo.yaml | kgrep -k Deployment -n reviews-v3

# List all resources of multiple kinds
kgrep --kind "Service,Deployment" cluster-resources.yaml
```

### Command Line Options

- `-k, --kind`: Filter by resource kind (e.g., Deployment, Service)
- `-n, --name`: Filter by resource name
- `-s, --summary`: Display a summary of resources
- `-h, --help`: Show help message

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.