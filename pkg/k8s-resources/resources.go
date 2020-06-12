package k8s_resources

import (
	"context"
	"fmt"
	"k8s.io/api/apps/v1beta1"
	v1 "k8s.io/api/core/v1"
	rbac "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"
)

const (
	KindNamespace   string = "Namespace"
	RoleRefAPIGroup        = "rbac.authorization.k8s.io"
	ClusterRoleKind        = "ClusterRole"
)

func MakeResources(cpuLimit, memLimit, cpuReq, memReq string) v1.ResourceRequirements {
	res := v1.ResourceRequirements{
		Limits: map[v1.ResourceName]resource.Quantity{
			v1.ResourceLimitsCPU:    resource.MustParse(cpuLimit),
			v1.ResourceLimitsMemory: resource.MustParse(memLimit),
		},
		Requests: map[v1.ResourceName]resource.Quantity{
			v1.ResourceRequestsCPU:    resource.MustParse(cpuReq),
			v1.ResourceRequestsMemory: resource.MustParse(memReq),
		},
	}
	return res
}

func MakeCommonLabels(appName string) map[string]string {
	return map[string]string{
		"app":        appName,
		"deployment": appName,
	}
}

func CommonServiceAccount(saName string) (v1.ServiceAccount, error) {
	if saName == "" {
		return v1.ServiceAccount{}, fmt.Errorf("Empty service account name!!!")
	}
	sa := v1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ServiceAccount",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: saName,
		},
		Secrets:                      nil,
		ImagePullSecrets:             nil,
		AutomountServiceAccountToken: nil,
	}
	return sa, nil
}

func CommonDeployment(name, image string, port, replicas int32, res v1.ResourceRequirements) (v1beta1.Deployment, error) {
	meta := metav1.ObjectMeta{
		Name:   name,
		Labels: MakeCommonLabels(name),
	}
	container := v1.Container{
		Name:      name,
		Image:     image,
		Resources: res,
		Ports: []v1.ContainerPort{
			v1.ContainerPort{
				Name:          name,
				Protocol:      v1.ProtocolTCP,
				ContainerPort: port,
			},
		},
	}
	deployment := v1beta1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "v1",
		},
		ObjectMeta: meta,
		Spec: v1beta1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: MakeCommonLabels(name),
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: meta,
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						container,
					},
				},
			},
		},
	}
	return deployment, nil
}

func CommonNamespace(name string) (v1.Namespace, error) {
	if name == "" {
		return v1.Namespace{}, fmt.Errorf("Empty namespace name")
	}
	ns := v1.Namespace{
		TypeMeta: metav1.TypeMeta{
			Kind:       KindNamespace,
			APIVersion: "",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: map[string]string{"name": name},
		},
	}
	return ns, nil
}

func SaClusterRoleBindingToNamespace(namespace, name, sa, clusterrole string) rbac.ClusterRoleBinding {
	rb := rbac.ClusterRoleBinding{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "",
			Kind:       "",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Subjects: []rbac.Subject{
			rbac.Subject{
				Kind:      rbac.ServiceAccountKind,
				Name:      sa,
				Namespace: namespace,
			},
		},
		RoleRef: rbac.RoleRef{
			APIGroup: RoleRefAPIGroup,
			Kind:     ClusterRoleKind,
			Name:     clusterrole,
		},
	}
	return rb
}

func CreateSecret(client kubernetes.Clientset, name string) (*v1.Secret, error) {
	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
	newSecret, err := client.CoreV1().Secrets("").Create(context.TODO(), secret, metav1.CreateOptions{})
	if err != nil && errors.IsAlreadyExists(err) {
		fmt.Printf("Secret %s is already exists", name)
		return newSecret, nil
	}
	if err != nil {
		return &v1.Secret{}, err
	}
	return newSecret, nil
}
