package kubernetes

import (
	"context"
	"testing"

	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func makeDeployment(namespace, name string, replicas int32) *v1.Deployment {
	return &v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace},
		Spec: v1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": name},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"app": name}},
				Spec:       corev1.PodSpec{},
			},
		},
	}
}

func TestGetDeployment_Found(t *testing.T) {
	cs := fake.NewSimpleClientset(makeDeployment("default", "myapp", 2))
	client := NewClientFromInterface(cs)

	dep, err := client.GetDeployment(context.Background(), "default", "myapp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dep.Name != "myapp" {
		t.Errorf("expected name myapp, got %s", dep.Name)
	}
}

func TestGetDeployment_NotFound(t *testing.T) {
	cs := fake.NewSimpleClientset()
	client := NewClientFromInterface(cs)

	_, err := client.GetDeployment(context.Background(), "default", "missing")
	if err == nil {
		t.Fatal("expected error for missing deployment, got nil")
	}
}

func TestListDeployments(t *testing.T) {
	cs := fake.NewSimpleClientset(
		makeDeployment("default", "app-a", 1),
		makeDeployment("default", "app-b", 3),
		makeDeployment("other", "app-c", 1),
	)
	client := NewClientFromInterface(cs)

	items, err := client.ListDeployments(context.Background(), "default")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 2 {
		t.Errorf("expected 2 deployments, got %d", len(items))
	}
}

func TestListDeployments_AllNamespaces(t *testing.T) {
	cs := fake.NewSimpleClientset(
		makeDeployment("default", "app-a", 1),
		makeDeployment("other", "app-c", 1),
	)
	client := NewClientFromInterface(cs)

	items, err := client.ListDeployments(context.Background(), "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 2 {
		t.Errorf("expected 2 deployments across namespaces, got %d", len(items))
	}
}
