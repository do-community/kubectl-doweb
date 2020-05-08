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
	"testing"

	"golang.org/x/net/context"
)

type NoopCloudPather struct{}

var _ CloudPather = &NoopCloudPather{}

func (_ *NoopCloudPather) Cluster(ctx context.Context) (string, error) {
	return "cluster", nil
}

func (_ *NoopCloudPather) Node(ctx context.Context, name string) (string, error) {
	return "no", nil
}

func (_ *NoopCloudPather) Service(ctx context.Context, namespace, name string) (string, error) {
	return "svc", nil
}

func (_ *NoopCloudPather) PersistentVolume(ctx context.Context, name string) (string, error) {
	return "pv", nil
}

func (_ *NoopCloudPather) PersistentVolumeClaim(ctx context.Context, namespace, name string) (string, error) {
	return "pvc", nil
}

func Test_cloudPatherWithType(t *testing.T) {
	cp := &NoopCloudPather{}
	tests := []struct {
		typ     string
		name    string
		want    string
		wantErr bool
	}{
		{
			typ:     "cluster",
			name:    "",
			want:    "cluster",
			wantErr: false,
		},
		{
			typ:     "cluster",
			name:    "name",
			want:    "cluster",
			wantErr: false,
		},

		{
			typ:     "no",
			name:    "",
			want:    "",
			wantErr: true,
		},
		{
			typ:     "no",
			name:    "name",
			want:    "no",
			wantErr: false,
		},
		{
			typ:     "node",
			name:    "name",
			want:    "no",
			wantErr: false,
		},
		{
			typ:     "nodes",
			name:    "name",
			want:    "no",
			wantErr: false,
		},

		{
			typ:     "svc",
			name:    "",
			want:    "",
			wantErr: true,
		},
		{
			typ:     "svc",
			name:    "name",
			want:    "svc",
			wantErr: false,
		},
		{
			typ:     "services",
			name:    "name",
			want:    "svc",
			wantErr: false,
		},
		{
			typ:     "service",
			name:    "name",
			want:    "svc",
			wantErr: false,
		},

		{
			typ:     "pv",
			name:    "",
			want:    "",
			wantErr: true,
		},
		{
			typ:     "pv",
			name:    "name",
			want:    "pv",
			wantErr: false,
		},
		{
			typ:     "persistentvolume",
			name:    "name",
			want:    "pv",
			wantErr: false,
		},
		{
			typ:     "persistentvolumes",
			name:    "name",
			want:    "pv",
			wantErr: false,
		},

		{
			typ:     "pvc",
			name:    "",
			want:    "",
			wantErr: true,
		},
		{
			typ:     "pvc",
			name:    "name",
			want:    "pvc",
			wantErr: false,
		},
		{
			typ:     "persistentvolumeclaim",
			name:    "name",
			want:    "pvc",
			wantErr: false,
		},
		{
			typ:     "persistentvolumeclaims",
			name:    "name",
			want:    "pvc",
			wantErr: false,
		},

		{
			typ:     "whomst",
			name:    "",
			want:    "",
			wantErr: true,
		},
		{
			typ:     "whomst",
			name:    "name",
			want:    "",
			wantErr: true,
		},
		{
			typ:     "",
			name:    "",
			want:    "",
			wantErr: true,
		},
		{
			typ:     "",
			name:    "name",
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.typ, func(t *testing.T) {
			got, err := cloudPatherByType(context.Background(), cp, tt.typ, "namespace", tt.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("cloudPatherWithType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("cloudPatherWithType() = %v, want %v", got, tt.want)
			}
		})
	}
}
