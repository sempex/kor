package kor

import (
	"context"
	"fmt"
	"slices"

	"github.com/yonahd/kor/pkg/filters"
	rbacv1 "k8s.io/api/rbac/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func RetrieveUsedClusterRoleBindings(clientset kubernetes.Clientset, filterOpts *filters.Options) error {
	usedClusterRoleBindings := make([]string, 0)
	namespaceList, err := clientset.CoreV1().Namespaces().List(context.TODO(), v1.ListOptions{})
	if err != nil {
		return fmt.Errorf("could not get a list of namespaces %v", namespaceList)
	}
	allServiceAccounts := make(map[string][]string)
	clusterRoleBindings, err := clientset.RbacV1().ClusterRoleBindings().List(context.TODO(), v1.ListOptions{})
	if err != nil {
		return fmt.Errorf("could not get a list of ClusterRoleBindings %v", err)
	}

	for _, ns := range namespaceList.Items {
		serviceAccounts, err := clientset.CoreV1().ServiceAccounts(ns.Name).List(context.TODO(), v1.ListOptions{})
		if err != nil {
			return fmt.Errorf("failed to list role bindings in namespace %s: %v", ns.Name, err)
		}

		for _, serviceAccount := range serviceAccounts.Items {
			allServiceAccounts[ns.Name] = append(allServiceAccounts[ns.Name], serviceAccount.Name)
		}
	}

	for _, clusterRoleBinding := range clusterRoleBindings.Items {
		saExists := checkServiceAccount(clusterRoleBinding, allServiceAccounts)
		if saExists {
			usedClusterRoleBindings = append(usedClusterRoleBindings, clusterRoleBinding.Name)
		}
	}
	fmt.Println(usedClusterRoleBindings)
	return nil
}

func checkServiceAccount(clusterrolebinding rbacv1.ClusterRoleBinding, allserviceaccounts map[string][]string) bool {
	serviceAccountDetected := 0
	for _, serviceAccount := range clusterrolebinding.Subjects {
		if serviceAccount.Kind != "ServiceAccount" {
			continue
		}
		namespacedServiceAccount := allserviceaccounts[serviceAccount.Namespace]
		if slices.Contains(namespacedServiceAccount, serviceAccount.Name) {
			serviceAccountDetected++
		}
	}
	return serviceAccountDetected != 0
}

func checkClusterRole(clientset kubernetes.Clientset, clusterrole string) {

}
