// Package helm provides utilities for loading Helm chart definitions
// and extracting their default values for use in drift detection.
//
// The primary entry point is [Loader], which can load charts from the
// local filesystem (directory or .tgz archive) via [Loader.LoadFromPath],
// or wrap a pre-parsed values map via [Loader.LoadFromValues].
//
// Loaded values are returned as a [ChartValues] struct, which is then
// consumed by the drift detector in internal/drift to compare against
// live Kubernetes deployment configurations.
//
// Typical usage:
//
//	l := helm.NewLoader()
//	cv, err := l.LoadFromPath("./charts/my-app")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(cv.ChartName, cv.Version)
package helm
