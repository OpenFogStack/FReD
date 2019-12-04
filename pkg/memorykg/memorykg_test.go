package memorykg

import (
	"reflect"
	"sync"
	"testing"
)

func TestKeygroupStorage_Create(t *testing.T) {
	type fields struct {
		keygroups map[string]struct{}
		RWMutex   sync.RWMutex
	}
	type args struct {
		kgname string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := &KeygroupStorage{
				keygroups: tt.fields.keygroups,
				RWMutex:   tt.fields.RWMutex,
			}
			if err := k.Create(tt.args.kgname); (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestKeygroupStorage_Delete(t *testing.T) {
	type fields struct {
		keygroups map[string]struct{}
		RWMutex   sync.RWMutex
	}
	type args struct {
		kgname string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := &KeygroupStorage{
				keygroups: tt.fields.keygroups,
				RWMutex:   tt.fields.RWMutex,
			}
			if err := k.Delete(tt.args.kgname); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestKeygroupStorage_Exists(t *testing.T) {
	type fields struct {
		keygroups map[string]struct{}
		RWMutex   sync.RWMutex
	}
	type args struct {
		kgname string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := &KeygroupStorage{
				keygroups: tt.fields.keygroups,
				RWMutex:   tt.fields.RWMutex,
			}
			if got := k.Exists(tt.args.kgname); got != tt.want {
				t.Errorf("Exists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name  string
		wantK *KeygroupStorage
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotK := New(); !reflect.DeepEqual(gotK, tt.wantK) {
				t.Errorf("New() = %v, want %v", gotK, tt.wantK)
			}
		})
	}
}