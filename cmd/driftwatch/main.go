package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/example/driftwatch/internal/drift"
	"github.com/example/driftwatch/internal/helm"
	"github.com/example/driftwatch/internal/kubernetes"
	"github.com/example/driftwatch/internal/report"
)

func main() {
	var (
		namespace  = flag.String("namespace", "", "Kubernetes namespace (empty = all namespaces)")
		deployment = flag.String("deployment", "", "Deployment name (required)")
		chartPath  = flag.String("chart", "", "Path to Helm chart values file (required)")
		outputFmt  = flag.String("output", "text", "Output format: text or json")
		kubeconfig = flag.String("kubeconfig", "", "Path to kubeconfig file (defaults to in-cluster)")
	)
	flag.Parse()

	if *deployment == "" || *chartPath == "" {
		fmt.Fprintln(os.Stderr, "error: --deployment and --chart flags are required")
		flag.Usage()
		os.Exit(1)
	}

	ctx := context.Background()

	kClient, err := kubernetes.NewClient(*kubeconfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating kubernetes client: %v\n", err)
		os.Exit(1)
	}

	dep, err := kClient.GetDeployment(ctx, *namespace, *deployment)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error fetching deployment: %v\n", err)
		os.Exit(1)
	}

	liveValues := kubernetes.ExtractValues(dep)

	hLoader := helm.NewLoader()
	chartValues, err := hLoader.LoadFromPath(*chartPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading chart values: %v\n", err)
		os.Exit(1)
	}

	detector := drift.NewDetector()
	result := detector.Compare(chartValues, liveValues)

	formatter := report.NewFormatter(os.Stdout)
	switch *outputFmt {
	case "json":
		if err := formatter.WriteJSON(result); err != nil {
			fmt.Fprintf(os.Stderr, "error writing json output: %v\n", err)
			os.Exit(1)
		}
	default:
		if err := formatter.WriteText(result); err != nil {
			fmt.Fprintf(os.Stderr, "error writing text output: %v\n", err)
			os.Exit(1)
		}
	}

	if result.HasDrift {
		os.Exit(2)
	}
}
