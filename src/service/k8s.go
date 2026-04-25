package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func authN() (*kubernetes.Clientset, error) {
	// authenticate against the kubernetes cluster using either the kubeconfig file or the in-cluster service account
	// return the clientset or the error
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		homeDir, err := os.UserHomeDir()
		if err == nil && homeDir != "" {
			kubeconfig = filepath.Join(homeDir, ".kube", "config")
		}
	}

	var cfg *rest.Config
	var cfgErr error

	if kubeconfig != "" {
		cfg, cfgErr = clientcmd.BuildConfigFromFlags("", kubeconfig)
	}

	if cfg == nil {
		cfg, cfgErr = rest.InClusterConfig()
	}

	if cfgErr != nil {
		return nil, fmt.Errorf("failed to build kubernetes config: %w", cfgErr)
	}

	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes clientset: %w", err)
	}

	return clientset, nil
}

func SpawnJob(jobID, image_uri string) error {
	clientset, err := authN()
	if err != nil {
		return fmt.Errorf("failed to authenticate against the kubernetes cluster: %w", err)
	}

	jobs := clientset.BatchV1().Jobs("default")
	var backOffLimit int32 = 0

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobID,
			Namespace: "default",
			Labels: map[string]string{
				"app":    "orchestrator",
				"type":   "job",
				"job-id": jobID,
			},
		},
		Spec: batchv1.JobSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:    jobID,
							Image:   image_uri,
							Command: []string{"echo", "Hello, World!"},
						},
					},
					RestartPolicy: v1.RestartPolicyNever,
				},
			},
			BackoffLimit: &backOffLimit,
		},
	}

	_, err = jobs.Create(context.TODO(), job, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create job: %w", err)
	}

	return err
}

func GetJobStatus(jobID string) (string, error) {
	clientset, err := authN()
	if err != nil {
		return "", fmt.Errorf("failed to authenticate against the kubernetes cluster: %w", err)
	}

	jobs := clientset.BatchV1().Jobs("default")
	job, err := jobs.Get(context.TODO(), jobID, metav1.GetOptions{})

	if err != nil {
		return "", fmt.Errorf("failed to get job: %w", err)
	}

	if job.Status.Active > 0 {
		return "running", nil
	} else if job.Status.Succeeded > 0 {
		return "succeeded", nil
	} else if job.Status.Failed > 0 {
		return "failed", nil
	} else {
		return "pending", nil
	}

}
