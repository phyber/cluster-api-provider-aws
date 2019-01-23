/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package machine

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsprovider/v1alpha1"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
	"sigs.k8s.io/cluster-api/pkg/controller/machine"
)

var (
	_ machine.Actuator = (*Actuator)(nil)
)

func contains(s []*clusterv1.Machine, e clusterv1.Machine) bool {
	exists := false
	for _, em := range s {
		if em.Name == e.Name && em.Namespace == e.Namespace {
			exists = true
			break
		}
	}
	return exists
}

func TestGetControlPlaneMachines(t *testing.T) {
	testCases := []struct {
		name        string
		input       *clusterv1.MachineList
		expectedOut []clusterv1.Machine
	}{
		{
			name: "0 machines",
			input: &clusterv1.MachineList{
				Items: []clusterv1.Machine{},
			},
			expectedOut: []clusterv1.Machine{},
		},
		{
			name: "only 2 controlplane machines",
			input: &clusterv1.MachineList{
				Items: []clusterv1.Machine{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "master-0",
							Namespace: "awesome-ns",
						},
						Spec: clusterv1.MachineSpec{
							Versions: clusterv1.MachineVersionInfo{
								Kubelet:      "v1.13.0",
								ControlPlane: "v1.13.0",
							},
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "master-1",
							Namespace: "awesome-ns",
						},
						Spec: clusterv1.MachineSpec{
							Versions: clusterv1.MachineVersionInfo{
								Kubelet:      "v1.13.0",
								ControlPlane: "v1.13.0",
							},
						},
					},
				},
			},
			expectedOut: []clusterv1.Machine{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "master-0",
						Namespace: "awesome-ns",
					},
					Spec: clusterv1.MachineSpec{
						Versions: clusterv1.MachineVersionInfo{
							Kubelet:      "v1.13.0",
							ControlPlane: "v1.13.0",
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "master-1",
						Namespace: "awesome-ns",
					},
					Spec: clusterv1.MachineSpec{
						Versions: clusterv1.MachineVersionInfo{
							Kubelet:      "v1.13.0",
							ControlPlane: "v1.13.0",
						},
					},
				},
			},
		},
		{
			name: "2 controlplane machines and 2 worker machines",
			input: &clusterv1.MachineList{
				Items: []clusterv1.Machine{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "master-0",
							Namespace: "awesome-ns",
						},
						Spec: clusterv1.MachineSpec{
							Versions: clusterv1.MachineVersionInfo{
								Kubelet:      "v1.13.0",
								ControlPlane: "v1.13.0",
							},
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "master-1",
							Namespace: "awesome-ns",
						},
						Spec: clusterv1.MachineSpec{
							Versions: clusterv1.MachineVersionInfo{
								Kubelet:      "v1.13.0",
								ControlPlane: "v1.13.0",
							},
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "worker-0",
							Namespace: "awesome-ns",
						},
						Spec: clusterv1.MachineSpec{
							Versions: clusterv1.MachineVersionInfo{
								Kubelet: "v1.13.0",
							},
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "worker-1",
							Namespace: "awesome-ns",
						},
						Spec: clusterv1.MachineSpec{
							Versions: clusterv1.MachineVersionInfo{
								Kubelet: "v1.13.0",
							},
						},
					},
				},
			},
			expectedOut: []clusterv1.Machine{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "master-0",
						Namespace: "awesome-ns",
					},
					Spec: clusterv1.MachineSpec{
						Versions: clusterv1.MachineVersionInfo{
							Kubelet:      "v1.13.0",
							ControlPlane: "v1.13.0",
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "master-1",
						Namespace: "awesome-ns",
					},
					Spec: clusterv1.MachineSpec{
						Versions: clusterv1.MachineVersionInfo{
							Kubelet:      "v1.13.0",
							ControlPlane: "v1.13.0",
						},
					},
				}},
		},
		{
			name: "only 2 worker machines",
			input: &clusterv1.MachineList{
				Items: []clusterv1.Machine{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "worker-0",
							Namespace: "awesome-ns",
						},
						Spec: clusterv1.MachineSpec{
							Versions: clusterv1.MachineVersionInfo{
								Kubelet: "v1.13.0",
							},
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "worker-1",
							Namespace: "awesome-ns",
						},
						Spec: clusterv1.MachineSpec{
							Versions: clusterv1.MachineVersionInfo{
								Kubelet: "v1.13.0",
							},
						},
					},
				},
			},
			expectedOut: []clusterv1.Machine{},
		},
	}
	testActuator := NewActuator(ActuatorParams{})

	for _, tc := range testCases {
		actual := testActuator.getControlPlaneMachines(tc.input)
		if len(actual) != len(tc.expectedOut) {
			t.Fatalf("[%s] Unexpected number of controlplane machines returned. Got: %d, Want: %d", tc.name, len(actual), len(tc.expectedOut))
		}
		if len(tc.expectedOut) > 1 {
			for _, em := range tc.expectedOut {
				if !contains(actual, em) {
					t.Fatalf("[%s] Expected controlplane machine %q in namespace %q not found", tc.name, em.Name, em.Namespace)
				}
			}
		}
	}
}

func TestMachineEqual(t *testing.T) {
	testCases := []struct {
		name          string
		inM1          clusterv1.Machine
		inM2          clusterv1.Machine
		expectedEqual bool
	}{
		{
			name: "machines are equal",
			inM1: clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "machine1",
					Namespace: "my-awesome-ns",
				},
			},
			inM2: clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "machine1",
					Namespace: "my-awesome-ns",
				},
			},
			expectedEqual: true,
		},
		{
			name: "machines are not equal: names are different",
			inM1: clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "machine 1",
					Namespace: "my-awesome-ns",
				},
			},
			inM2: clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "machine 2",
					Namespace: "my-awsesome-ns",
				},
			},
			expectedEqual: false,
		},
		{
			name: "machines are not equal: namespace are different",
			inM1: clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "machine1",
					Namespace: "my-awesome-ns",
				},
			},
			inM2: clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "machine1",
					Namespace: "your-awsesome-ns",
				},
			},
			expectedEqual: false,
		},
	}

	for _, tc := range testCases {
		actualEqual := machinesEqual(&tc.inM1, &tc.inM2)
		if tc.expectedEqual {
			if !actualEqual {
				t.Fatalf("[%s] Expected Machine1 [Name:%q, Namespace:%q], Equal Machine2 [Name:%q, Namespace:%q]",
					tc.name, tc.inM1.Name, tc.inM1.Namespace, tc.inM2.Name, tc.inM2.Namespace)
			}
		} else {
			if actualEqual {
				t.Fatalf("[%s] Expected Machine1 [Name:%q, Namespace:%q], NOT Equal Machine2 [Name:%q, Namespace:%q]",
					tc.name, tc.inM1.Name, tc.inM1.Namespace, tc.inM2.Name, tc.inM2.Namespace)
			}
		}
	}
}

func TestImmutableStateChange(t *testing.T) {
	testCases := []struct {
		name                string
		machineConfig       v1alpha1.AWSMachineProviderSpec
		instanceDescription v1alpha1.Instance
		expectedValue       bool
	}{
		{
			name: "instance type is unchanged",
			machineConfig: v1alpha1.AWSMachineProviderSpec{
				InstanceType: "t2.micro",
			},
			instanceDescription: v1alpha1.Instance{
				Type: "t2.micro",
			},
			expectedValue: false,
		},
		{
			name: "instance type is changed",
			machineConfig: v1alpha1.AWSMachineProviderSpec{
				InstanceType: "m5.large",
			},
			instanceDescription: v1alpha1.Instance{
				Type: "t2.micro",
			},
			expectedValue: true,
		},
		{
			name: "iam profile is unchanged",
			machineConfig: v1alpha1.AWSMachineProviderSpec{
				IAMInstanceProfile: "test-profile",
			},
			instanceDescription: v1alpha1.Instance{
				IAMProfile: "test-profile",
			},
			expectedValue: false,
		},
		{
			name: "iam profile is changed",
			machineConfig: v1alpha1.AWSMachineProviderSpec{
				IAMInstanceProfile: "test-profile-updated",
			},
			instanceDescription: v1alpha1.Instance{
				IAMProfile: "test-profile",
			},
			expectedValue: true,
		},
		{
			name: "keyname is unchanged",
			machineConfig: v1alpha1.AWSMachineProviderSpec{
				KeyName: "SSHKey",
			},
			instanceDescription: v1alpha1.Instance{
				KeyName: aws.String("SSHKey"),
			},
			expectedValue: false,
		},
		{
			name: "keyname is changed",
			machineConfig: v1alpha1.AWSMachineProviderSpec{
				KeyName: "SSHKey2",
			},
			instanceDescription: v1alpha1.Instance{
				KeyName: aws.String("SSHKey"),
			},
			expectedValue: true,
		},
		{
			name: "instance with public ip is unchanged",
			machineConfig: v1alpha1.AWSMachineProviderSpec{
				PublicIP: aws.Bool(true),
			},
			instanceDescription: v1alpha1.Instance{
				// This IP chosen from RFC5737 TEST-NET-1
				PublicIP: aws.String("192.0.2.1"),
			},
			expectedValue: false,
		},
		{
			name: "instance with public ip is changed",
			machineConfig: v1alpha1.AWSMachineProviderSpec{
				PublicIP: aws.Bool(false),
			},
			instanceDescription: v1alpha1.Instance{
				// This IP chosen from RFC5737 TEST-NET-1
				PublicIP: aws.String("192.0.2.1"),
			},
			expectedValue: true,
		},
		{
			name: "instance without public ip is unchanged",
			machineConfig: v1alpha1.AWSMachineProviderSpec{
				PublicIP: aws.Bool(false),
			},
			instanceDescription: v1alpha1.Instance{
				// This IP chosen from RFC5737 TEST-NET-1
				PublicIP: aws.String(""),
			},
			expectedValue: false,
		},
		{
			name: "instance without public ip is changed",
			machineConfig: v1alpha1.AWSMachineProviderSpec{
				PublicIP: aws.Bool(true),
			},
			instanceDescription: v1alpha1.Instance{
				// This IP chosen from RFC5737 TEST-NET-1
				PublicIP: aws.String(""),
			},
			expectedValue: true,
		},
		{
			name: "subnetid is unchanged",
			machineConfig: v1alpha1.AWSMachineProviderSpec{
				Subnet: &v1alpha1.AWSResourceReference{
					ID: aws.String("subnet-abcdef"),
				},
			},
			instanceDescription: v1alpha1.Instance{
				SubnetID: "subnet-abcdef",
			},
			expectedValue: false,
		},
		{
			name: "subnetid is changed",
			machineConfig: v1alpha1.AWSMachineProviderSpec{
				Subnet: &v1alpha1.AWSResourceReference{
					ID: aws.String("subnet-123456"),
				},
			},
			instanceDescription: v1alpha1.Instance{
				SubnetID: "subnet-abcdef",
			},
			expectedValue: true,
		},
	}

	for _, tc := range testCases {
		changed := immutableStateChanged(&tc.machineConfig, &tc.instanceDescription)

		// true case: values changed
		if tc.expectedValue != changed {
			t.Fatalf("[%s] Expected Machine Config [%+v], NOT Equal Instance Description [%+v]",
				tc.name, tc.machineConfig, tc.instanceDescription)
		}
	}
}
