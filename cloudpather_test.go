package main

import (
	"context"
	"testing"
)

func TestDOCloudPather_Cluster(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		cp      *DOCloudPather
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cp.Cluster(tt.args.ctx)
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
	type args struct {
		ctx  context.Context
		name string
	}
	tests := []struct {
		name    string
		cp      *DOCloudPather
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cp.Node(tt.args.ctx, tt.args.name)
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
	type args struct {
		ctx       context.Context
		namespace string
		name      string
	}
	tests := []struct {
		name    string
		cp      *DOCloudPather
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cp.Service(tt.args.ctx, tt.args.namespace, tt.args.name)
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
	type args struct {
		ctx  context.Context
		name string
	}
	tests := []struct {
		name    string
		cp      *DOCloudPather
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cp.PersistentVolume(tt.args.ctx, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("DOCloudPather.PersistentVolume() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DOCloudPather.PersistentVolume() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDOCloudPather_PersistentVolumeClaim(t *testing.T) {
	type args struct {
		ctx       context.Context
		namespace string
		name      string
	}
	tests := []struct {
		name    string
		cp      *DOCloudPather
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cp.PersistentVolumeClaim(tt.args.ctx, tt.args.namespace, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("DOCloudPather.PersistentVolumeClaim() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DOCloudPather.PersistentVolumeClaim() = %v, want %v", got, tt.want)
			}
		})
	}
}
