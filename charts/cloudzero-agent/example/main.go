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
	startTime := time.Now()
	totalJobs := 100
	jobs := make(chan int, totalJobs)
	var wg sync.WaitGroup

	for w := 1; w <= 20; w++ {
		wg.Add(1)
		go worker(clientset, w, jobs, &wg)
	}

	for job := 1; job <= totalJobs; job++ {
		jobs <- job
	}

	close(jobs)
	wg.Wait()
	fmt.Println("Total time ", time.Since(startTime))
}

func scaleDeployments(clientset *kubernetes.Clientset, namespace string) {
	log.Println("Scaling down in namespace", namespace)
	deploymentsClient := clientset.AppsV1().Deployments(namespace)

	// List all deployments
	deployments, err := deploymentsClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Failed to list deployments: %v", err)
	}

	for _, deployment := range deployments.Items {
		// Set replicas to 0
		deployment.Spec.Replicas = int32Ptr(0)
		_, err := deploymentsClient.Update(context.TODO(), &deployment, metav1.UpdateOptions{})
		if err != nil {
			log.Printf("Failed to scale deployment %s in namespace %s.: %v", deployment.Name, namespace, err)
		} else {
			fmt.Printf("Scaled deployment %s to 0 replicas in namespace %s.\n", deployment.Name, namespace)
		}
	}
}

// scaleStatefulSets scales all statefulsets in the given namespace down to 0 replicas
func scaleStatefulSets(clientset *kubernetes.Clientset, namespace string) {
	statefulSetsClient := clientset.AppsV1().StatefulSets(namespace)

	// List all statefulsets
	statefulSets, err := statefulSetsClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Failed to list statefulsets: %v", err)
	}

	for _, statefulSet := range statefulSets.Items {
		// Set replicas to 0
		statefulSet.Spec.Replicas = int32Ptr(0)
		_, err := statefulSetsClient.Update(context.TODO(), &statefulSet, metav1.UpdateOptions{})
		if err != nil {
			log.Printf("Failed to scale statefulset %s: %v", statefulSet.Name, err)
		} else {
			fmt.Printf("Scaled statefulset %s to 0 replicas in namespace %s.\n", statefulSet.Name, namespace)
		}
	}
}

// Helper function to return a pointer to an int32
func int32Ptr(i int32) *int32 {
	return &i
}

func worker(clientset *kubernetes.Clientset, w int, jobs chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		// processJobs(clientset, w, job)
		scaleDown(clientset, job)
	}
}

func scaleDown(clientset *kubernetes.Clientset, job int) {
	scaleDeployments(clientset, "test-namespace-"+strconv.Itoa(job))
	scaleStatefulSets(clientset, "test-namespace-"+strconv.Itoa(job))

}

func processJobs(clientset *kubernetes.Clientset, w int, job int) {
	fmt.Println("Worker", w, "started  job", job)
	numResources := 100
	namespaceName := "test-namespace-" + strconv.Itoa(job)
	replicas := int32(2)

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
					"deploy-label": "num-" + strconv.Itoa(job),
				},
				Annotations: map[string]string{
					"deploy-annotation": "num-" + strconv.Itoa(job),
				},
			},
			Spec: appsv1.DeploymentSpec{
				Replicas: &replicas,
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"app":       deploymentName,
						"pod-label": "num-" + strconv.Itoa(job),
					},
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"app":       deploymentName,
							"pod-label": "num-" + strconv.Itoa(job),
						},
						Annotations: map[string]string{
							"pod-annotation": "num-" + strconv.Itoa(job),
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
					"sts-label": "num-" + strconv.Itoa(job),
				},
				Annotations: map[string]string{
					"sts-annotation": "num-" + strconv.Itoa(job),
				},
			},
			Spec: appsv1.StatefulSetSpec{
				ServiceName: statefulSetName + "-service",
				Replicas:    &replicas,
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"app":       statefulSetName,
						"pod-label": "num-" + strconv.Itoa(job),
					},
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"app":       statefulSetName,
							"pod-label": "num-" + strconv.Itoa(job),
						},
						Annotations: map[string]string{
							"pod-annotation": "num-" + strconv.Itoa(job),
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
	fmt.Println("Worker", w, "finished job", job)
}
