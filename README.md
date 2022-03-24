# kgrep - search through Kubernetes resources

# Usage
```
# Search for services within all.yaml
kgrep --kind Service all.yaml

# Search for a Deployment named foo within helm chart
helm template chart | kgrep --kind Deployment --name foo
```