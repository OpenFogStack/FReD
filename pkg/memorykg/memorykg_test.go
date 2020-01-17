package memorykg

import (
	"reflect"
	"sync"
	"testing"

	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/commons"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/keygroup"
)

func TestKeygroupStorage_Create(t *testing.T) {
	type fields struct {
		keygroups map[commons.KeygroupName]struct{}
		sync.RWMutex
	}
	type args struct {
		keygroup.Keygroup
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"create simple keygroup",
			fields{
				make(map[commons.KeygroupName]struct{}),
				sync.RWMutex{},
			},
			args{keygroup.Keygroup{
				Name: "keygroup",
			}},
			false,
		},
		{"create keygroup with empty name",
			fields{
				make(map[commons.KeygroupName]struct{}),
				sync.RWMutex{},
			},
			args{keygroup.Keygroup{
				Name: "",
			}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := &KeygroupStorage{
				keygroups: tt.fields.keygroups,
				RWMutex:   tt.fields.RWMutex,
			}
			if err := k.Create(tt.args.Keygroup); (err != nil) != tt.wantErr {
				t.Errorf("Create() errors = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestKeygroupStorage_Delete(t *testing.T) {
	type fields struct {
		keygroups map[commons.KeygroupName]struct{}
		RWMutex   sync.RWMutex
	}
	type args struct {
		keygroup.Keygroup
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"delete non-existent keygroup",
			fields{
				make(map[commons.KeygroupName]struct{}),
				sync.RWMutex{},
			},
			args{keygroup.Keygroup{
				Name: "keygroup",
			}},
			true,
		},
		{"delete keygroup with empty name",
			fields{
				make(map[commons.KeygroupName]struct{}),
				sync.RWMutex{},
			},
			args{keygroup.Keygroup{
				Name: "",
			}},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := &KeygroupStorage{
				keygroups: tt.fields.keygroups,
				RWMutex:   tt.fields.RWMutex,
			}
			if err := k.Delete(tt.args.Keygroup); (err != nil) != tt.wantErr {
				t.Errorf("Delete() errors = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestKeygroupStorage_Exists(t *testing.T) {
	type fields struct {
		keygroups map[commons.KeygroupName]struct{}
		RWMutex   sync.RWMutex
	}
	type args struct {
		keygroup.Keygroup
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{"check non-existent keygroup",
			fields{
				make(map[commons.KeygroupName]struct{}),
				sync.RWMutex{},
			},
			args{keygroup.Keygroup{
				Name: "",
			}},
			false,
		},
		{"check keygroup with empty name",
			fields{
				make(map[commons.KeygroupName]struct{}),
				sync.RWMutex{},
			},
			args{keygroup.Keygroup{
				Name: "",
			}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := &KeygroupStorage{
				keygroups: tt.fields.keygroups,
				RWMutex:   tt.fields.RWMutex,
			}
			if got := k.Exists(tt.args.Keygroup); got != tt.want {
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
		{"create new empty KeygroupStorage",
			&KeygroupStorage{
				keygroups: make(map[commons.KeygroupName]struct{}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotK := New(); !reflect.DeepEqual(gotK, tt.wantK) {
				t.Errorf("New() = %v, want %v", gotK, tt.wantK)
			}
		})
	}
}
