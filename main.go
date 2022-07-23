package main

//Inspired from https://github.com/kubernetes/client-go/blob/master/examples/out-of-cluster-client-configuration/main.go

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func createNamespace(cs *kubernetes.Clientset, namespace string) {
	nsSpec := &v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: namespace}}
	n, err := cs.CoreV1().Namespaces().Create(context.TODO(), nsSpec, metav1.CreateOptions{})
	if err != nil {
		fmt.Println("failed to create ns")
		panic(err.Error())
	}
	fmt.Println("created ns successfully", n)
}

func createPod(cs *kubernetes.Clientset, podName, namespace string) {
	podSpec := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "hello-world",
			Namespace: namespace,
			Labels: map[string]string{
				"k8s-app": "kube-dns",
			},
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  "hello",
					Image: "hello-world",
				},
			},
		},
	}

	p, err := cs.CoreV1().Pods(namespace).Create(context.TODO(), podSpec, metav1.CreateOptions{})
	if err != nil {
		fmt.Println("failed to create pod")
		panic(err.Error())
	}
	fmt.Println("created pod successfully", p)
}

func listPods(cs *kubernetes.Clientset, namespace, labelSelector string) {
	fmt.Println("listing pods.....")
	pods, err := cs.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labelSelector})
	if err != nil {
		fmt.Println("failed to list pods with dns selector")
		panic(err.Error())
	}
	for _, pod := range pods.Items {
		fmt.Println(pod.Name)
	}
	fmt.Println("listed pods successfully")
}

func deletePod(cs *kubernetes.Clientset, namespace, name string) {
	err := cs.CoreV1().Pods(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		fmt.Printf("failed to delete %s pod from namespace %s", name, namespace)
		panic(err.Error())
	}
	fmt.Printf("deleted pod %s successfully from %s namespace", name, namespace)
}

func main() {
	var kubeconfig *string
	var namespace string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	flag.StringVar(&namespace, "namespace", "challenge-ns", "namespace to create & to create pods in")

	flag.Parse()
	// create the k8s clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// List namespaces
	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("There are %d namespaces in the cluster; listing....\n", len(namespaces.Items))

	for _, ns := range namespaces.Items {
		fmt.Println(ns.Name)
	}

	// Create namespace
	createNamespace(clientset, namespace)

	// Create pod
	createPod(clientset, "hello-world", namespace)

	// List pods
	listPods(clientset, namespace, "k8s-app=kube-dns")

	// Delete Pod
	deletePod(clientset, namespace, "hello-world")

}
