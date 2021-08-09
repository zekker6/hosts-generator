package kubernetes

import (
	"context"
	"github.com/pkg/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"sort"
)

type Client struct {
	kubeconfig string
}

func NewKubernetesClient(kubeconfig string) *Client {
	return &Client{kubeconfig: kubeconfig}
}

func (k Client) Get() ([]string, error) {
	config, err := clientcmd.BuildConfigFromFlags("", k.kubeconfig)
	if err != nil {
		return []string{}, errors.Wrap(err, "failed to create client from config")
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return []string{}, errors.Wrap(err, "failed to create client from config")
	}

	ingresses, err := clientset.NetworkingV1().Ingresses("").List(context.Background(), v1.ListOptions{})
	if err != nil {
		return []string{}, errors.Wrap(err, "failed to list ingress resources")
	}

	hosts := make([]string, 0, len(ingresses.Items))
	for _, ing := range ingresses.Items {
		for _, r := range ing.Spec.Rules {
			hosts = append(hosts, r.Host)
		}
	}

	sort.Strings(hosts)

	return hosts, nil

}
