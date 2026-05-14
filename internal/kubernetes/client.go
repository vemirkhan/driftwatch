package kubernetes

import (
	"context"
	"fmt"

	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Client wraps a Kubernetes clientset for deployment inspection.
type Client struct {
	clientset kubernetes.Interface
}

// NewClient creates a Client using the provided kubeconfig path.
// If kubeconfig is empty, it falls back to in-cluster configuration.
func NewClient(kubeconfig string) (*Client, error) {
	var cfg *rest.Config
	var err error

	if kubeconfig != "" {
		cfg, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	} else {
		cfg, err = rest.InClusterConfig()
	}
	if err != nil {
		return nil, fmt.Errorf("building kube config: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("creating clientset: %w", err)
	}

	return &Client{clientset: clientset}, nil
}

// NewClientFromInterface creates a Client from an existing kubernetes.Interface (useful for testing).
func NewClientFromInterface(cs kubernetes.Interface) *Client {
	return &Client{clientset: cs}
}

// GetDeployment retrieves a single Deployment by namespace and name.
func (c *Client) GetDeployment(ctx context.Context, namespace, name string) (*v1.Deployment, error) {
	dep, err := c.clientset.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("getting deployment %s/%s: %w", namespace, name, err)
	}
	return dep, nil
}

// ListDeployments returns all Deployments in the given namespace.
// Pass an empty string to list across all namespaces.
func (c *Client) ListDeployments(ctx context.Context, namespace string) ([]v1.Deployment, error) {
	list, err := c.clientset.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("listing deployments in namespace %q: %w", namespace, err)
	}
	return list.Items, nil
}
