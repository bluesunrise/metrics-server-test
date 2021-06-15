package main

import (
	"context"
//	"flag"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	io_prometheus_client "github.com/prometheus/client_model/go"
//	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
//	"k8s.io/client-go/tools/clientset"
//	"k8s.io/client-go/util/homedir"
	"net/http"
	"os"
//	"path/filepath"
	"strings"
//	"time"

	//"k8s.io/client-go/kubernetes"
	//"k8s.io/client-go/rest"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"
	v1beta1 "k8s.io/metrics/pkg/client/clientset/versioned/typed/metrics/v1beta1"

	"github.com/prometheus/common/expfmt"
)

// https://github.com/kubernetes/client-go/tree/master/examples
// https://github.com/kubernetes/client-go/tree/master/examples#configuration


func main() {
	kubernetesEndpoint := os.Getenv("METRICS_TEST_ENDPOINT")
	if kubernetesEndpoint == "" {
		kubernetesEndpoint = "localhost:8001"
	}
	// kubeStateMetricsEndpoint := "http://" + kubernetesEndpoint + "/api/v1/namespaces/kube-system/services/kube-state-metrics:http-metrics/proxy/metrics"

	config := rest.Config{
		Host:                kubernetesEndpoint,
		APIPath:             "",
		ContentConfig:       rest.ContentConfig{},
		Username:            "",
		Password:            "",
		BearerToken:         "",
		BearerTokenFile:     "",
		Impersonate:         rest.ImpersonationConfig{},
		AuthProvider:        nil,
		AuthConfigPersister: nil,
		ExecProvider:        nil,
		TLSClientConfig:     rest.TLSClientConfig{},
		UserAgent:           "",
		DisableCompression:  false,
		Transport:           nil,
		WrapTransport:       nil,
		QPS:                 0,
		Burst:               0,
		RateLimiter:         nil,
		Timeout:             0,
		Dial:                nil,
	}

	kClientSet, err := kubernetes.NewForConfig(&config)
	if err != nil {
		panic(err.Error())
	}
	mClientSet, err := metricsv.NewForConfig(&config)
	kapi := kClientSet.CoreV1()
	mapi := mClientSet.MetricsV1beta1()
	ctx := context.TODO()
	GatherNodeMetrics(mapi, ctx)
	GatherPodMetrics(mapi, ctx)
	if 1 == 1 {
		GatherNodeInventory(kapi, ctx)
		GatherPodInventory(kapi, ctx)
		GatherNamespaceInventory(kapi, ctx)
		GatherServiceInventory(kapi, ctx)
		GatherPersistentVolumes(kapi, ctx)
		GatherEndPoints(kapi, ctx)
		GatherReplicationControllers(kapi, ctx)
		GatherConfigMaps(kapi, ctx)
		GatherResourceQuotas(kapi, ctx)
		GatherComponentStatuses(kapi, ctx)
		GatherEvents(kapi, ctx)
		// GatherKubeStateMetrics(kubeStateMetricsEndpoint)
	}
}

func GatherKubeStateMetrics(endPoint string) {
	fmt.Println("State Metrics ....")
	resp, err := http.Get(endPoint)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var parser expfmt.TextParser
	expeditionPayload, err := parser.TextToMetricFamilies(resp.Body)
	if err != nil {
		panic(err)
	}
	for _, line := range expeditionPayload {
		fmt.Printf("%s - %d\n", line.GetName(), line.GetType())
		for _, metric := range line.Metric {
			if line.GetType() == io_prometheus_client.MetricType_GAUGE {
				fmt.Printf("\tGAUGE: %f\n", *metric.GetGauge().Value)
			} else {
				fmt.Printf("\tUNKNOWN METRIC Type: %d\n", *line.Type);
			}
			for _, label := range metric.Label {
				fmt.Printf("\t\tLabel = %s - %s\n", *label.Name, *label.Value)
			}
		}
	}
}

func GatherNodeInventory(api v1.CoreV1Interface, context context.Context) {
	fmt.Println("Node Inventory ....")
	nodes, _ := api.Nodes().List(context, metav1.ListOptions{})
	for _, node := range nodes.Items {
		fmt.Printf("node = %s\n", node.Name)
		for key, element := range node.Labels {
			fmt.Printf("\tlabel = %s:%s\n", key, element)
		}
		fmt.Println("Node Capacity:")
        spew.Dump(node.Status.Capacity)
		fmt.Println("Node Allocatable:")
		spew.Dump(node.Status.Allocatable)
		fmt.Println("Conditions:")
		spew.Dump(node.Status.Conditions)
		fmt.Println("Addresses:")
		spew.Dump(node.Status.Addresses)
		fmt.Println("Volumes Attached:")
		spew.Dump(node.Status.VolumesAttached)
		fmt.Println("Volumes IN USE:")
		spew.Dump(node.Status.VolumesAttached)
		fmt.Println("Info:")
		spew.Dump(node.Status.NodeInfo)
		fmt.Println("Config:")
		spew.Dump(node.Status.Config)
		fmt.Println("Images:")
		spew.Dump(node.Status.Images)
		fmt.Println("Phase:")
		spew.Dump(node.Status.Phase)
	}
}

func GatherPodInventory(api v1.CoreV1Interface, context context.Context) {
	// nodes
	fmt.Println("Pod Inventory ....")
	pods, _ := api.Pods("").List(context, metav1.ListOptions{})
	for _, pod := range pods.Items {
		fmt.Printf("pod = %s\n", pod.Name)
		for key, element := range pod	.Labels {
			fmt.Printf("\tlabel = %s:%s\n", key, element)
		}
		// TODO:
		//fmt.Println("Conditions:")
		//spew.Dump(pod.Status.Conditions)
		//fmt.Println("Phase:")
		//spew.Dump(pod.Status.Phase)
	}
}

func GatherNamespaceInventory(api v1.CoreV1Interface, context context.Context) {
	fmt.Println("Namespace Inventory ....")
	namespaces, _ := api.Namespaces().List(context, metav1.ListOptions{})
	for _, ns := range namespaces.Items {
		fmt.Printf("namespace = %s\n", ns.Name)
	}
}

func GatherServiceInventory(api v1.CoreV1Interface, context context.Context) {
	fmt.Println("Service Inventory ....")
	services, _ := api.Services("").List(context, metav1.ListOptions{})
	for _, s := range services.Items {
		fmt.Printf("service = %s\n", s.Name)
	}
}

func GatherPersistentVolumes(api v1.CoreV1Interface, context context.Context) {
	volumes, _ := api.PersistentVolumes().List(context, metav1.ListOptions{})
	for _, v := range volumes.Items {
		fmt.Printf("volume = %s\n", v.Name)
	}
	claims, _ := api.PersistentVolumeClaims("").List(context, metav1.ListOptions{})
	for _, v := range claims.Items {
		fmt.Printf("volume claim = %s\n", v.Name)
	}
}

func GatherEndPoints(api v1.CoreV1Interface, context context.Context) {
	fmt.Println("End Points ....")
	endpoints, _ := api.Endpoints("").List(context, metav1.ListOptions{})
	for _, ep := range endpoints.Items {
		fmt.Printf("endPoint = %s\n", ep.Name)
	}
}

func GatherReplicationControllers(api v1.CoreV1Interface, context context.Context) {
	fmt.Println("Replication Controllers ....")
	replicationControllers, _ := api.ReplicationControllers("").List(context, metav1.ListOptions{})
	for _, rp := range replicationControllers.Items {
		fmt.Printf("replication Controller = %s\n", rp.Name)
	}
}

func GatherConfigMaps(api v1.CoreV1Interface, context context.Context) {
	fmt.Println("Config Maps ....")
	configMaps, _ := api.ConfigMaps("").List(context, metav1.ListOptions{})
	for _, cm := range configMaps.Items {
		fmt.Printf("Config Map = %s\n", cm.Name)
	}
}

func GatherEvents(api v1.CoreV1Interface, context context.Context) {
	fmt.Println("Events ....")
	events, _ := api.Events("").List(context, metav1.ListOptions{})
	for _, e := range events.Items {
		fmt.Printf("Events = %s\n", e.Name)
	}
}

func GatherResourceQuotas(api v1.CoreV1Interface, context context.Context) {
	fmt.Println("Resource Quotas ....")
	quotas, _ := api.ResourceQuotas("").List(context, metav1.ListOptions{})
	for _, q := range quotas.Items {
		fmt.Printf("Quota = %s\n", q.Name)
	}
}

func GatherComponentStatuses(api v1.CoreV1Interface, context context.Context) {
	fmt.Println("Component Statuses ....")
	statuses, _ := api.ComponentStatuses().List(context, metav1.ListOptions{})
	for _, status := range statuses.Items {
		fmt.Printf("Component Status = %s\n", status.Name)
	}
}

func GatherNodeMetrics(api v1beta1.MetricsV1beta1Interface, context context.Context) {
	nodes, _ := api.NodeMetricses().List(context, metav1.ListOptions{})
	for _, node := range nodes.Items {
		fmt.Println(node.Name, node.Namespace, node.ClusterName, node.Kind)
		fmt.Println("\t", node.Usage.Cpu())
		fmt.Println("\t", node.Usage.Memory())
		fmt.Println("\t", node.Usage.Storage())
		fmt.Println("\t", node.Usage.StorageEphemeral())
		fmt.Println("\t", node.Usage.Pods())
	}
}

func GatherPodMetrics(api v1beta1.MetricsV1beta1Interface, context context.Context) {
	pods, _ := api.PodMetricses("").List(context, metav1.ListOptions{})
	for _, pod := range pods.Items {
		fmt.Println(pod.Name, pod.Namespace, pod.Timestamp)
		for _, container := range pod.Containers {
			fmt.Println("\t" + container.Name)
			fmt.Println("\t\t", container.Usage.Cpu())
			fmt.Println("\t\t", container.Usage.Memory().Value())
			fmt.Println("\t\t", container.Usage.Storage().Value())
			fmt.Println("\t\t", container.Usage.StorageEphemeral().Value())
			fmt.Println("\t\t", container.Usage.Pods().Value())
		}
	}
}

func tester() {
	fmt.Printf("%q\n", strings.SplitAfterN("wordpress-39933f", "-", 2))
	fmt.Printf("%q\n", strings.SplitAfterN("word-press-39933f", "-", 2))
	fmt.Printf("%q\n", strings.Split("wordpress-39933f", "-"))
	fmt.Printf("%q\n", strings.Split("word-press-39933f", "-"))
}

//func main() {
//	var kubeconfig string = "./config"
//	if home := homedir.HomeDir(); home != "" {
//		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
//	} else {
//		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
//	}
//	//flag.Parse()
//
//	// use the current context in kubeconfig
//
//		// Examples for error handling:
//		// - Use helper functions like e.g. errors.IsNotFound()
//		// - And/or cast to StatusError and use its properties like e.g. ErrStatus.Message
//		namespace := "default"
//		pod := "example-xxxxx"
//		_, err = clientset.CoreV1().Pods(namespace).Get(context.TODO(), pod, metav1.GetOptions{})
//		if errors.IsNotFound(err) {
//			fmt.Printf("Pod %s in namespace %s not found\n", pod, namespace)
//		} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
//			fmt.Printf("Error getting pod %s in namespace %s: %v\n",
//				pod, namespace, statusError.ErrStatus.Message)
//		} else if err != nil {
//			panic(err.Error())
//		} else {
//			fmt.Printf("Found pod %s in namespace %s\n", pod, namespace)
//		}
//
//		time.Sleep(10 * time.Second)
//	}
//}
