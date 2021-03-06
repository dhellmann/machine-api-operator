package testing

import (
	"fmt"
	"time"

	mapiv1 "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
	healthcheckingv1alpha1 "github.com/openshift/machine-api-operator/pkg/apis/healthchecking/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/uuid"
)

const (
	// Namespace contains the default namespace for machine-api components
	Namespace = "openshift-machine-api"
	// MachineAnnotationKey contains default machine node annotation
	MachineAnnotationKey = "machine.openshift.io/machine"
)

var (
	// KnownDate contains date that can be used under tests
	KnownDate = metav1.Time{Time: time.Date(1985, 06, 03, 0, 0, 0, 0, time.Local)}
)

// FooBar returns foo:bar map that can be used as default label
func FooBar() map[string]string {
	return map[string]string{"foo": "bar"}
}

// NewSelector returns new LabelSelector
func NewSelector(labels map[string]string) *metav1.LabelSelector {
	return &metav1.LabelSelector{MatchLabels: labels}
}

// NewSelectorFooBar returns new foo:bar label selector
func NewSelectorFooBar() *metav1.LabelSelector {
	return NewSelector(FooBar())
}

// NewMinAvailableMachineDisruptionBudget returns new MachineDisruptionBudget with min available parameter
func NewMinAvailableMachineDisruptionBudget(minAvailable int32) *healthcheckingv1alpha1.MachineDisruptionBudget {
	return &healthcheckingv1alpha1.MachineDisruptionBudget{
		TypeMeta: metav1.TypeMeta{Kind: "MachineDisruptionBudget"},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "foobar",
			Namespace: Namespace,
		},
		Spec: healthcheckingv1alpha1.MachineDisruptionBudgetSpec{
			MinAvailable: &minAvailable,
			Selector:     NewSelectorFooBar(),
		},
	}
}

// NewMaxUnavailableMachineDisruptionBudget returns new MachineDisruptionBudget with max unavailable parameter
func NewMaxUnavailableMachineDisruptionBudget(maxUnavailable int32) *healthcheckingv1alpha1.MachineDisruptionBudget {
	return &healthcheckingv1alpha1.MachineDisruptionBudget{
		TypeMeta: metav1.TypeMeta{Kind: "MachineDisruptionBudget"},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "foobar",
			Namespace: Namespace,
		},
		Spec: healthcheckingv1alpha1.MachineDisruptionBudgetSpec{
			MaxUnavailable: &maxUnavailable,
			Selector:       NewSelectorFooBar(),
		},
	}
}

// NewNode returns new node object that can be used for testing
func NewNode(name string, ready bool) *corev1.Node {
	nodeReadyStatus := corev1.ConditionTrue
	if !ready {
		nodeReadyStatus = corev1.ConditionUnknown
	}

	return &corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: metav1.NamespaceNone,
			Annotations: map[string]string{
				MachineAnnotationKey: fmt.Sprintf("%s/%s", Namespace, "fakeMachine"),
			},
			Labels: map[string]string{},
			UID:    uuid.NewUUID(),
		},
		TypeMeta: metav1.TypeMeta{
			Kind: "Node",
		},
		Status: corev1.NodeStatus{
			Conditions: []corev1.NodeCondition{
				{
					Type:               corev1.NodeReady,
					Status:             nodeReadyStatus,
					LastTransitionTime: KnownDate,
				},
			},
		},
	}
}

// NewMachine returns new machine object that can be used for testing
func NewMachine(name string, nodeName string) *mapiv1.Machine {
	m := &mapiv1.Machine{
		TypeMeta: metav1.TypeMeta{Kind: "Machine"},
		ObjectMeta: metav1.ObjectMeta{
			Annotations:     make(map[string]string),
			Name:            name,
			Namespace:       Namespace,
			Labels:          FooBar(),
			UID:             uuid.NewUUID(),
			OwnerReferences: []metav1.OwnerReference{{Kind: "MachineSet"}},
		},
		Spec: mapiv1.MachineSpec{},
	}
	if nodeName != "" {
		m.Status = mapiv1.MachineStatus{
			NodeRef: &corev1.ObjectReference{
				Name:      nodeName,
				Namespace: metav1.NamespaceNone,
			},
		}
	}
	return m
}

// NewMachineHealthCheck returns new MachineHealthCheck object that can be used for testing
func NewMachineHealthCheck(name string) *healthcheckingv1alpha1.MachineHealthCheck {
	return &healthcheckingv1alpha1.MachineHealthCheck{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: Namespace,
		},
		TypeMeta: metav1.TypeMeta{
			Kind: "MachineHealthCheck",
		},
		Spec: healthcheckingv1alpha1.MachineHealthCheckSpec{
			Selector: *NewSelectorFooBar(),
		},
		Status: healthcheckingv1alpha1.MachineHealthCheckStatus{},
	}
}

// NewUnhealthyConditionsConfigMap returns new config map object with unhealthy conditions
func NewUnhealthyConditionsConfigMap(name string, data string) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: Namespace,
		},
		TypeMeta: metav1.TypeMeta{
			Kind: "ConfigMap",
		},
		Data: map[string]string{
			"conditions": data,
		},
	}
}
