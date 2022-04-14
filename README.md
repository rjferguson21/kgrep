# kgrep - search through Kubernetes resources

# Usage
```bash
# Search for services within all.yaml
kgrep --kind Service all.yaml

# Search for a Deployment named foo within helm chart
helm template chart | kgrep --kind Deployment --name foo

# List a summary of resources in a file
curl -s https://raw.githubusercontent.com/istio/istio/release-1.13/samples/bookinfo/platform/kube/bookinfo.yaml | kgrep -s

# Search for a specific Deployment
curl -s https://raw.githubusercontent.com/istio/istio/release-1.13/samples/bookinfo/platform/kube/bookinfo.yaml | kgrep -k Deployment -n reviews-v3
```

# Installation
```bash
go install github.com/rjferguson21/kgrep@latest
```