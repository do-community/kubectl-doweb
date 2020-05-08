/*
Copyright 2020 Kamal Nasser All rights reserved.
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

package kubectldoweb

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	restclient "k8s.io/client-go/rest"
)

func newFakeDOCloudPather() *DOCloudPather {
	return &DOCloudPather{
		clientset: fake.NewSimpleClientset(),
		output:    &bytes.Buffer{},
	}
}

func TestDOCloudPather_Cluster(t *testing.T) {
	cp := newFakeDOCloudPather()
	id := "random-uuid"

	tests := []struct {
		name         string
		clientConfig *restclient.Config
		want         string
		wantErr      bool
	}{
		{
			name: "valid DOKS cluster",
			clientConfig: &restclient.Config{
				Host: fmt.Sprintf("https://%s%s", id, hostnameSuffix),
			},
			want:    fmt.Sprintf("kubernetes/clusters/%s", id),
			wantErr: false,
		},
		{
			name: "non-DOKS cluster",
			clientConfig: &restclient.Config{
				Host: fmt.Sprintf("https://%s.example.com", id),
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "invalid host",
			clientConfig: &restclient.Config{
				Host: "http://a b.com/",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cp.clientConfig = tt.clientConfig

			got, err := cp.Cluster(context.TODO())
			if (err != nil) != tt.wantErr {
				t.Errorf("DOCloudPather.Cluster() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DOCloudPather.Cluster() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDOCloudPather_Node(t *testing.T) {
	ctx := context.TODO()
	id := "random-id"
	nodeName := "node-1"

	tests := []struct {
		name     string
		nodeName string
		node     *corev1.Node
		want     string
		wantErr  bool
	}{
		{
			name:     "inexistent node",
			nodeName: nodeName,
			node:     nil,
			want:     "",
			wantErr:  true,
		},
		{
			name:     "valid DOKS node",
			nodeName: nodeName,
			node: &corev1.Node{
				ObjectMeta: metav1.ObjectMeta{
					Name: nodeName,
				},
				Spec: corev1.NodeSpec{
					ProviderID: fmt.Sprintf("%s%s", nodeIDPrefix, id),
				},
			},
			want:    fmt.Sprintf("droplets/%s", id),
			wantErr: false,
		},
		{
			name:     "non-DOKS node",
			nodeName: nodeName,
			node: &corev1.Node{
				ObjectMeta: metav1.ObjectMeta{
					Name: nodeName,
				},
				Spec: corev1.NodeSpec{
					ProviderID: fmt.Sprintf("whomst://%s", id),
				},
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cp := newFakeDOCloudPather()
			if tt.node != nil {
				cp.clientset.CoreV1().Nodes().Create(ctx, tt.node, metav1.CreateOptions{})
			}

			got, err := cp.Node(ctx, tt.nodeName)
			if (err != nil) != tt.wantErr {
				t.Errorf("DOCloudPather.Node() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DOCloudPather.Node() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDOCloudPather_Service(t *testing.T) {
	ctx := context.TODO()
	id := "random-id"
	namespace := "ns"
	svcName := "svc-1"

	tests := []struct {
		name        string
		serviceName string
		service     *corev1.Service
		want        string
		wantErr     bool
	}{
		{
			name:        "inexistent service",
			serviceName: svcName,
			service:     nil,
			want:        "",
			wantErr:     true,
		},
		{
			name:        "valid DOKS LB service",
			serviceName: svcName,
			service: &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name: svcName,
					Annotations: map[string]string{
						lbaasAnnotation: id,
					},
				},
				Spec: corev1.ServiceSpec{
					Type: corev1.ServiceTypeLoadBalancer,
				},
			},
			want:    fmt.Sprintf("networking/load_balancers/%s", id),
			wantErr: false,
		},
		{
			name:        "non-DOKS LB service",
			serviceName: svcName,
			service: &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name: svcName,
				},
				Spec: corev1.ServiceSpec{
					Type: corev1.ServiceTypeLoadBalancer,
				},
			},
			want:    "",
			wantErr: true,
		},
		{
			name:        "non-LB service",
			serviceName: svcName,
			service: &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name: svcName,
				},
				Spec: corev1.ServiceSpec{
					Type: corev1.ServiceTypeClusterIP,
				},
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cp := newFakeDOCloudPather()
			if tt.service != nil {
				cp.clientset.CoreV1().Services(namespace).Create(ctx, tt.service, metav1.CreateOptions{})
			}

			got, err := cp.Service(ctx, namespace, tt.serviceName)
			if (err != nil) != tt.wantErr {
				t.Errorf("DOCloudPather.Service() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DOCloudPather.Service() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDOCloudPather_PersistentVolume(t *testing.T) {
	ctx := context.TODO()
	pvName := "pv-1"

	tests := []struct {
		name       string
		pvName     string
		pv         *corev1.PersistentVolume
		want       string
		wantOutput string
		wantErr    bool
	}{
		{
			name:       "inexistent pv",
			pvName:     pvName,
			pv:         nil,
			want:       "",
			wantOutput: "",
			wantErr:    true,
		},
		{
			name:   "valid DO CSI pv",
			pvName: pvName,
			pv: &corev1.PersistentVolume{
				ObjectMeta: metav1.ObjectMeta{
					Name: pvName,
				},
				Spec: corev1.PersistentVolumeSpec{
					StorageClassName: storageClassName,
				},
			},
			want:       "volumes",
			wantOutput: fmt.Sprintf("PersistentVolume name: %s\n", pvName),
			wantErr:    false,
		},
		{
			name:   "non-DO CSI pv",
			pvName: pvName,
			pv: &corev1.PersistentVolume{
				ObjectMeta: metav1.ObjectMeta{
					Name: pvName,
				},
				Spec: corev1.PersistentVolumeSpec{
					StorageClassName: "whomst",
				},
			},
			want:       "",
			wantOutput: "",
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cp := newFakeDOCloudPather()
			if tt.pv != nil {
				cp.clientset.CoreV1().PersistentVolumes().Create(ctx, tt.pv, metav1.CreateOptions{})
			}

			got, err := cp.PersistentVolume(ctx, tt.pvName)
			if (err != nil) != tt.wantErr {
				t.Errorf("DOCloudPather.PersistentVolume() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DOCloudPather.PersistentVolume() = %v, want %v", got, tt.want)
			}

			gotOutput := cp.output.(*bytes.Buffer).String()
			if gotOutput != tt.wantOutput {
				t.Errorf("DOCloudPather.PersistentVolume() output = %v, want %v", gotOutput, tt.wantOutput)
			}
		})
	}
}

func TestDOCloudPather_PersistentVolumeClaim(t *testing.T) {
	ctx := context.TODO()
	pvName := "pv-1"
	pvcName := "pvc-1"
	namespace := "ns"
	storageClassNameString := string(storageClassName)
	nonDOStorageClassName := "whomst"

	tests := []struct {
		name       string
		pvcName    string
		pvc        *corev1.PersistentVolumeClaim
		want       string
		wantOutput string
		wantErr    bool
	}{
		{
			name:       "inexistent pvc",
			pvcName:    pvcName,
			pvc:        nil,
			want:       "",
			wantOutput: "",
			wantErr:    true,
		},
		{
			name:    "unbound/pending pvc",
			pvcName: pvcName,
			pvc: &corev1.PersistentVolumeClaim{
				ObjectMeta: metav1.ObjectMeta{
					Name: pvcName,
				},
				Spec: corev1.PersistentVolumeClaimSpec{
					StorageClassName: &storageClassNameString,
				},
				Status: corev1.PersistentVolumeClaimStatus{
					Phase: corev1.ClaimPending,
				},
			},
			want:       "",
			wantOutput: "",
			wantErr:    true,
		},
		{
			name:    "pvc with nil storage class",
			pvcName: pvcName,
			pvc: &corev1.PersistentVolumeClaim{
				ObjectMeta: metav1.ObjectMeta{
					Name: pvcName,
				},
				Spec: corev1.PersistentVolumeClaimSpec{},
			},
			want:       "",
			wantOutput: "",
			wantErr:    true,
		},
		{
			name:    "valid DO CSI bound pvc",
			pvcName: pvcName,
			pvc: &corev1.PersistentVolumeClaim{
				ObjectMeta: metav1.ObjectMeta{
					Name: pvcName,
				},
				Spec: corev1.PersistentVolumeClaimSpec{
					VolumeName:       pvName,
					StorageClassName: &storageClassNameString,
				},
				Status: corev1.PersistentVolumeClaimStatus{
					Phase: corev1.ClaimBound,
				},
			},
			want:       "volumes",
			wantOutput: fmt.Sprintf("PersistentVolume name: %s\n", pvName),
			wantErr:    false,
		},
		{
			name:    "non-DO CSI bound pvc",
			pvcName: pvcName,
			pvc: &corev1.PersistentVolumeClaim{
				ObjectMeta: metav1.ObjectMeta{
					Name: pvcName,
				},
				Spec: corev1.PersistentVolumeClaimSpec{
					VolumeName:       pvName,
					StorageClassName: &nonDOStorageClassName,
				},
				Status: corev1.PersistentVolumeClaimStatus{
					Phase: corev1.ClaimBound,
				},
			},
			want:       "",
			wantOutput: "",
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cp := newFakeDOCloudPather()
			if tt.pvc != nil {
				cp.clientset.CoreV1().PersistentVolumeClaims(namespace).Create(ctx, tt.pvc, metav1.CreateOptions{})
			}

			got, err := cp.PersistentVolumeClaim(ctx, namespace, tt.pvcName)
			if (err != nil) != tt.wantErr {
				t.Errorf("DOCloudPather.PersistentVolumeClaim() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DOCloudPather.PersistentVolumeClaim() = %v, want %v", got, tt.want)
			}

			gotOutput := cp.output.(*bytes.Buffer).String()
			if gotOutput != tt.wantOutput {
				t.Errorf("DOCloudPather.PersistentVolumeClaim() output = %v, want %v", gotOutput, tt.wantOutput)
			}
		})
	}
}
