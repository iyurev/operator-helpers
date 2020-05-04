package k8s_resources

import (
	"fmt"
	"k8s.io/api/apps/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func makeResources(cpuLimit, memLimit, cpuReq, memReq string) corev1.ResourceRequirements {
	res := corev1.ResourceRequirements{
		Limits: map[corev1.ResourceName]resource.Quantity{
			corev1.ResourceLimitsCPU:    resource.MustParse(cpuLimit),
			corev1.ResourceLimitsMemory: resource.MustParse(memLimit),
		},
		Requests: map[corev1.ResourceName]resource.Quantity{
			corev1.ResourceRequestsCPU:    resource.MustParse(cpuReq),
			corev1.ResourceRequestsMemory: resource.MustParse(memReq),
		},
	}
	return res
}

func MakeCommonLabels(appName, deploymentName string) map[string]string {
	return map[string]string{
		"app":        appName,
		"deployment": deploymentName,
	}
}

func CommonServiceAccount(saName, namespace string) (corev1.ServiceAccount, error) {
	if saName == "" {
		return corev1.ServiceAccount{}, fmt.Errorf("Empty service account name!!!")
	}
	if namespace == "" {
		return corev1.ServiceAccount{}, fmt.Errorf("Empty namespace name!!!")
	}
	sa := corev1.ServiceAccount{
		TypeMeta: v1.TypeMeta{
			Kind:       "ServiceAccount",
			APIVersion: "v1",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      saName,
			Namespace: namespace,
		},
		Secrets:                      nil,
		ImagePullSecrets:             nil,
		AutomountServiceAccountToken: nil,
	}
	return sa, nil
}

func commonDeployment(name, image string, port int) (v1beta1.Deployment, error) {
	return v1beta1.Deployment{}, nil
}
