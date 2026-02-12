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

`kgrep` reads Kubernetes manifests from stdin and filters by kind and/or name.

### Kind/Name Syntax

The fastest way to filter resources is with the positional `Kind/Name` syntax:

```bash
# Search for all Deployments
helm template my-chart | kgrep Deployment

# Search for a specific Service by name
cat manifests.yaml | kgrep Service/nginx

# Search for any resource named "nginx" (wildcard kind)
kubectl get all -o yaml | kgrep '*/nginx'

# Explicit wildcard name (same as just "Service")
cat manifests.yaml | kgrep 'Service/*'
```

### Using Flags

Flags provide more control and support regex patterns:

```bash
# Filter by kind and name
helm template my-chart | kgrep --kind Deployment --name frontend

# Use regex to match multiple kinds
kubectl get all -o yaml | kgrep --kind "Service|Deployment"

# Use regex to match name patterns
cat manifests.yaml | kgrep --kind Pod --name "nginx.*"

# Generate a summary of resources
helm template my-chart | kgrep -s
```

### Examples

```bash
# Search within Helm chart output
helm template my-chart | kgrep Deployment/frontend

# Find resources in remote manifests
curl -s https://raw.githubusercontent.com/.../bookinfo.yaml | kgrep Deployment/reviews-v3

# List all Deployments with summary output
kubectl get all -o yaml | kgrep Deployment -s
```

### Command Line Options

```
kgrep [kind[/name]] [flags]
```

- `kind[/name]`: Positional filter using Kind/Name syntax (use `*` as wildcard)
- `-k, --kind`: Filter by resource kind (supports regex)
- `-n, --name`: Filter by resource name (supports regex)
- `-s, --summary`: Display a summary of resources
- `-h, --help`: Show help message

Flags override positional arguments when both are provided.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.