package kubernetes

import (
	v1 "k8s.io/api/apps/v1"
)

// ExtractValues converts relevant fields of a Deployment into a flat
// key-value map that can be compared against Helm chart values.
func ExtractValues(dep *v1.Deployment) map[string]interface{} {
	out := make(map[string]interface{})

	if dep == nil {
		return out
	}

	// Top-level replica count
	if dep.Spec.Replicas != nil {
		out["replicaCount"] = int64(*dep.Spec.Replicas)
	}

	containers := dep.Spec.Template.Spec.Containers
	if len(containers) == 0 {
		return out
	}

	// Use the first container as the primary container
	c := containers[0]
	out["image.repository"] = imageRepository(c.Image)
	out["image.tag"] = imageTag(c.Image)

	if c.Resources.Requests != nil {
		if cpu, ok := c.Resources.Requests["cpu"]; ok {
			out["resources.requests.cpu"] = cpu.String()
		}
		if mem, ok := c.Resources.Requests["memory"]; ok {
			out["resources.requests.memory"] = mem.String()
		}
	}

	if c.Resources.Limits != nil {
		if cpu, ok := c.Resources.Limits["cpu"]; ok {
			out["resources.limits.cpu"] = cpu.String()
		}
		if mem, ok := c.Resources.Limits["memory"]; ok {
			out["resources.limits.memory"] = mem.String()
		}
	}

	for _, env := range c.Env {
		out["env."+env.Name] = env.Value
	}

	return out
}

// imageRepository returns the repository portion of an image string.
func imageRepository(image string) string {
	for i := len(image) - 1; i >= 0; i-- {
		if image[i] == ':' {
			return image[:i]
		}
	}
	return image
}

// imageTag returns the tag portion of an image string, defaulting to "latest".
func imageTag(image string) string {
	for i := len(image) - 1; i >= 0; i-- {
		if image[i] == ':' {
			return image[i+1:]
		}
	}
	return "latest"
}
