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

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("There are %d namespaces in the cluster\n", len(namespaces.Items))

	for _, ns := range namespaces.Items {
		fmt.Println(ns.Name)
	}

	nsSpec := &v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "user-akamble"}}
	n, err := clientset.CoreV1().Namespaces().Create(context.TODO(), nsSpec, metav1.CreateOptions{})
	if err != nil {
		fmt.Println("failed to create ns")
		panic(err.Error())
	}
	fmt.Println("created ns successfully", n)
}
