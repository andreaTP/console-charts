package kube

import (
	"fmt"

	"github.com/lightbend/gotests/util"

	"k8s.io/client-go/kubernetes"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NOTE: Following two utilities use kubectl, that means they won't work when running tests from inside the cluster.
// A change to parse yaml and use client-go is needed to make that work.

func ApplyYaml(namespace string, filepath string) error {
	if _, err := util.Cmd("kubectl", "-n", namespace, "apply", "-f", filepath).Run(); err != nil {
		return err
	}

	return nil
}

func DeleteYaml(namespace string, filepath string) error {
	if _, err := util.Cmd("kubectl", "-n", namespace, "delete", "-f", filepath).Run(); err != nil {
		return err
	}

	return nil
}

func DeletePvc(namespace string, name string) error {
	if _, err := util.Cmd("kubectl", "-n", namespace, "delete", "pvc", name).Run(); err != nil {
		return err
	}

	return nil
}

func CreateNamespace(k8sClient *kubernetes.Clientset, name string) error {
	namespace := &apiv1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
	_, err := k8sClient.CoreV1().Namespaces().Create(namespace)
	return err
}

func IsDeploymentAvailable(k8sClient *kubernetes.Clientset, namespace string, name string) error {
	dep, err := k8sClient.AppsV1().Deployments(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	if len(dep.Status.Conditions) == 0 {
		return fmt.Errorf("deployment is pending")
	}
	if dep.Status.Conditions[0].Type != "Available" {
		return fmt.Errorf("deployment not available: %v", dep.Status.Conditions[0].Type)
	}
	return nil
}