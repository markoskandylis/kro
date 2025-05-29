package commands

import (
	"context"
	"fmt"

	krov1alpha1 "github.com/kro-run/kro/api/v1alpha1"
	kroclient "github.com/kro-run/kro/pkg/client"
	"github.com/kro-run/kro/pkg/graph"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/scale/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func getKubeClient() (client.Client, error) {
	s := runtime.NewScheme()
	_ = scheme.AddToScheme(s)
	_ = krov1alpha1.AddToScheme(s)
	cfg := ctrl.GetConfigOrDie()
	return client.New(cfg, client.Options{Scheme: s})
}

func rgdBuilder(rgd krov1alpha1.ResourceGraphDefinition) (*graph.Graph, error) {
	set, err := kroclient.NewSet(kroclient.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to init kro client: %w", err)
	}
	builder, err := graph.NewBuilder(set.RESTConfig())
	if err != nil {
		return nil, fmt.Errorf("failed to build graph: %w", err)
	}
	graphRuntime, err := builder.NewResourceGraphDefinition(&rgd)
	if err != nil {
		return nil, fmt.Errorf("failed to render graph runtime: %w", err)
	}

	return graphRuntime, nil
}

func getResourceIfExists(c client.Client, name string) (krov1alpha1.ResourceGraphDefinition, bool, error) {
	var rgd krov1alpha1.ResourceGraphDefinition
	err := c.Get(context.Background(), client.ObjectKey{Name: name}, &rgd)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return krov1alpha1.ResourceGraphDefinition{}, false, nil
		}
		return krov1alpha1.ResourceGraphDefinition{}, false, err
	}
	return rgd, true, nil
}
