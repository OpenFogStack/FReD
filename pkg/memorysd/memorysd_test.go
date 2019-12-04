package memorysd

import (
	"reflect"
	"sync"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name  string
		wantS *Storage
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotS := New(); !reflect.DeepEqual(gotS, tt.wantS) {
				t.Errorf("New() = %v, want %v", gotS, tt.wantS)
			}
		})
	}
}

func TestStorage_Create(t *testing.T) {
	type fields struct {
		keygroups map[string]Keygroup
		RWMutex   sync.RWMutex
	}
	type args struct {
		kgname string
		data   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    uint64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				keygroups: tt.fields.keygroups,
				RWMutex:   tt.fields.RWMutex,
			}
			got, err := s.Create(tt.args.kgname, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Create() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_CreateKeygroup(t *testing.T) {
	type fields struct {
		keygroups map[string]Keygroup
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
			s := &Storage{
				keygroups: tt.fields.keygroups,
				RWMutex:   tt.fields.RWMutex,
			}
			if err := s.CreateKeygroup(tt.args.kgname); (err != nil) != tt.wantErr {
				t.Errorf("CreateKeygroup() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStorage_Delete(t *testing.T) {
	type fields struct {
		keygroups map[string]Keygroup
		RWMutex   sync.RWMutex
	}
	type args struct {
		kgname string
		id     uint64
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
			s := &Storage{
				keygroups: tt.fields.keygroups,
				RWMutex:   tt.fields.RWMutex,
			}
			if err := s.Delete(tt.args.kgname, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStorage_DeleteKeygroup(t *testing.T) {
	type fields struct {
		keygroups map[string]Keygroup
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
			s := &Storage{
				keygroups: tt.fields.keygroups,
				RWMutex:   tt.fields.RWMutex,
			}
			if err := s.DeleteKeygroup(tt.args.kgname); (err != nil) != tt.wantErr {
				t.Errorf("DeleteKeygroup() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStorage_Read(t *testing.T) {
	type fields struct {
		keygroups map[string]Keygroup
		RWMutex   sync.RWMutex
	}
	type args struct {
		kgname string
		id     uint64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				keygroups: tt.fields.keygroups,
				RWMutex:   tt.fields.RWMutex,
			}
			got, err := s.Read(tt.args.kgname, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Read() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_Update(t *testing.T) {
	type fields struct {
		keygroups map[string]Keygroup
		RWMutex   sync.RWMutex
	}
	type args struct {
		kgname string
		id     uint64
		data   string
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
			s := &Storage{
				keygroups: tt.fields.keygroups,
				RWMutex:   tt.fields.RWMutex,
			}
			if err := s.Update(tt.args.kgname, tt.args.id, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}