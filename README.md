# driftwatch

Detects configuration drift between running Kubernetes deployments and their source Helm chart definitions.

---

## Installation

```bash
go install github.com/yourusername/driftwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/driftwatch.git && cd driftwatch && go build ./...
```

---

## Usage

Point `driftwatch` at a namespace and a Helm chart to compare live deployment state against expected configuration:

```bash
driftwatch --namespace production --chart ./charts/myapp
```

Compare against a remote chart release:

```bash
driftwatch --namespace staging --release myapp --repo https://charts.example.com
```

Example output:

```
[DRIFT] deployment/myapp-api
  - image.tag: expected "v1.4.2", got "v1.3.9"
  - resources.limits.memory: expected "512Mi", got "256Mi"

[OK] deployment/myapp-worker
```

### Flags

| Flag | Description |
|------|-------------|
| `--namespace` | Kubernetes namespace to inspect |
| `--chart` | Path or URL to the Helm chart |
| `--release` | Helm release name |
| `--repo` | Helm chart repository URL |
| `--output` | Output format: `text`, `json`, `yaml` (default: `text`) |
| `--kubeconfig` | Path to kubeconfig file (defaults to `~/.kube/config`) |

---

## Requirements

- Go 1.21+
- Kubernetes cluster access (kubeconfig)
- Helm 3

---

## License

MIT © 2024 [Your Name](https://github.com/yourusername)