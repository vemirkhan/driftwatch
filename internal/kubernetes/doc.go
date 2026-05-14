// Package kubernetes provides utilities for interacting with a Kubernetes
// cluster in the context of drift detection.
//
// It exposes two main concerns:
//
//  1. Client — a thin wrapper around the official client-go Kubernetes
//     clientset, providing methods to fetch and list Deployments.
//
//  2. ExtractValues — converts a live Deployment object into a flat
//     key-value map whose keys mirror common Helm chart value paths
//     (e.g. "replicaCount", "image.repository", "resources.requests.cpu").
//     This normalised representation is consumed by the drift detector in
//     internal/drift to compare running state against chart-defined values.
//
// Example usage:
//
//	client, err := kubernetes.NewClient("/home/user/.kube/config")
//	if err != nil { ... }
//
//	dep, err := client.GetDeployment(ctx, "default", "my-release-app")
//	if err != nil { ... }
//
//	liveVals := kubernetes.ExtractValues(dep)
package kubernetes
