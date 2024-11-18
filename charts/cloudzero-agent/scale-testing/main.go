package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// Create Kubernetes client
	kubeconfigPath := "/Users/daniel.mepham/.kube/config"
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	config.QPS = 10000000
	config.Burst = 100000000
	if err != nil {
		log.Fatalf("Failed to create in-cluster config: %v", err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to create Kubernetes client: %v", err)
	}

	// Number of namespaces, deployments, and statefulsets
	numNamespaces := 999
	numResources := 100
	replicas := int32(2)
	numThreads := 20

	// WaitGroup for managing goroutines
	var wg sync.WaitGroup
	worker := func(namespace int) {
		defer wg.Done()
		namespaceName := "test-namespace-" + strconv.Itoa(namespace)

		// Ensure the namespace exists
		_, err := clientset.CoreV1().Namespaces().Get(context.TODO(), namespaceName, metav1.GetOptions{})
		if err != nil {
			_, err = clientset.CoreV1().Namespaces().Create(context.TODO(), &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespaceName,
				},
			}, metav1.CreateOptions{})
			if err != nil {
				log.Fatalf("Failed to create namespace %s: %v", namespaceName, err)
			}
			fmt.Printf("Created namespace: %s\n", namespaceName)
		}

		// Create Deployments
		for j := 1; j <= numResources; j++ {
			deploymentName := "deployment-" + strconv.Itoa(j)
			deployment := &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      deploymentName,
					Namespace: namespaceName,
					Labels: map[string]string{
						"deploy-label": "num-" + strconv.Itoa(namespace),
					},
					Annotations: map[string]string{
						"deploy-annotation": "num-" + strconv.Itoa(namespace),
					},
				},
				Spec: appsv1.DeploymentSpec{
					Replicas: &replicas,
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							"app":       deploymentName,
							"pod-label": "num-" + strconv.Itoa(namespace),
						},
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								"app":       deploymentName,
								"pod-label": "num-" + strconv.Itoa(namespace),
							},
							Annotations: map[string]string{
								"pod-annotation": "num-" + strconv.Itoa(namespace),
							},
						},
						Spec: corev1.PodSpec{
							NodeSelector: map[string]string{"fake": "test"},
							Containers: []corev1.Container{
								{
									Name:  "nginx",
									Image: "nginx:latest",
								},
							},
						},
					},
				},
			}
			_, err := clientset.AppsV1().Deployments(namespaceName).Create(context.TODO(), deployment, metav1.CreateOptions{})
			if err != nil {
				log.Printf("Failed to create Deployment %s in namespace %s: %v", deploymentName, namespaceName, err)
			}
		}

		// Create StatefulSets
		for j := 1; j <= numResources; j++ {
			statefulSetName := "statefulset-" + strconv.Itoa(j)
			statefulSet := &appsv1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:      statefulSetName,
					Namespace: namespaceName,
					Labels: map[string]string{
						"sts-label": "num-" + strconv.Itoa(namespace),
					},
					Annotations: map[string]string{
						"sts-annotation": "num-" + strconv.Itoa(namespace),
					},
				},
				Spec: appsv1.StatefulSetSpec{
					ServiceName: statefulSetName + "-service",
					Replicas:    &replicas,
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							"app":       statefulSetName,
							"pod-label": "num-" + strconv.Itoa(namespace),
						},
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								"app":       statefulSetName,
								"pod-label": "num-" + strconv.Itoa(namespace),
							},
							Annotations: map[string]string{
								"pod-annotation": "num-" + strconv.Itoa(namespace),
							},
						},
						Spec: corev1.PodSpec{
							NodeSelector: map[string]string{"fake": "test"},
							Containers: []corev1.Container{
								{
									Name:  "nginx",
									Image: "nginx:latest",
								},
							},
						},
					},
				},
			}
			_, err := clientset.AppsV1().StatefulSets(namespaceName).Create(context.TODO(), statefulSet, metav1.CreateOptions{})
			if err != nil {
				log.Printf("Failed to create StatefulSet %s in namespace %s: %v", statefulSetName, namespaceName, err)
			}
		}

		fmt.Printf("Completed namespace: %s\n", namespaceName)
		time.Sleep(500 * time.Millisecond) // Throttle requests slightly
	}

	// Launch workers
	for i := 1; i <= numThreads; i++ {
		wg.Add(1)
		namespaceNum
		go worker(i)
	}

	// Wait for all workers to finish
	wg.Wait()

	fmt.Println("Finished creating Deployments and StatefulSets.")
}
