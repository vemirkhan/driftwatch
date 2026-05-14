package kubernetes

import (
	"testing"

	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func ptr32(i int32) *int32 { return &i }

func TestExtractValues_Basic(t *testing.T) {
	dep := &v1.Deployment{
		Spec: v1.DeploymentSpec{
			Replicas: ptr32(3),
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{Image: "nginx:1.21"},
					},
				},
			},
		},
	}

	vals := ExtractValues(dep)

	if vals["replicaCount"] != int64(3) {
		t.Errorf("expected replicaCount=3, got %v", vals["replicaCount"])
	}
	if vals["image.repository"] != "nginx" {
		t.Errorf("expected image.repository=nginx, got %v", vals["image.repository"])
	}
	if vals["image.tag"] != "1.21" {
		t.Errorf("expected image.tag=1.21, got %v", vals["image.tag"])
	}
}

func TestExtractValues_NoTag(t *testing.T) {
	dep := &v1.Deployment{
		Spec: v1.DeploymentSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{Image: "myapp"},
					},
				},
			},
		},
	}

	vals := ExtractValues(dep)
	if vals["image.tag"] != "latest" {
		t.Errorf("expected image.tag=latest, got %v", vals["image.tag"])
	}
}

func TestExtractValues_Resources(t *testing.T) {
	dep := &v1.Deployment{
		Spec: v1.DeploymentSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Image: "app:v2",
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									"cpu":    resource.MustParse("100m"),
									"memory": resource.MustParse("128Mi"),
								},
							},
						},
					},
				},
			},
		},
	}

	vals := ExtractValues(dep)
	if vals["resources.requests.cpu"] != "100m" {
		t.Errorf("expected resources.requests.cpu=100m, got %v", vals["resources.requests.cpu"])
	}
}

func TestExtractValues_Nil(t *testing.T) {
	vals := ExtractValues(nil)
	if len(vals) != 0 {
		t.Errorf("expected empty map for nil deployment, got %v", vals)
	}
}
