package memorysd

import (
	"reflect"
	"sync"
	"testing"

	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/data"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name  string
		wantS *Storage
	}{
		{"create new empty Storage",
			&Storage{
				keygroups: make(map[string]Keygroup),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotS := New(); !reflect.DeepEqual(gotS, tt.wantS) {
				t.Errorf("New() = %v, want %v", gotS, tt.wantS)
			}
		})
	}
}

func TestStorage_CreateKeygroup(t *testing.T) {
	type fields struct {
		keygroups map[string]Keygroup
		sync.RWMutex
	}
	type args struct {
		data.Item
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"create simple keygroup",
			fields{
				make(map[string]Keygroup),
				sync.RWMutex{},
			},
			args{data.Item{
				Keygroup: "keygroup",
			}},
			false,
		},
		{"create keygroup with empty name",
			fields{
				make(map[string]Keygroup),
				sync.RWMutex{},
			},
			args{data.Item{
				Keygroup: "",
			}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				keygroups: tt.fields.keygroups,
				RWMutex:   tt.fields.RWMutex,
			}
			if err := s.CreateKeygroup(tt.args.Item); (err != nil) != tt.wantErr {
				t.Errorf("CreateKeygroup() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStorage_Delete(t *testing.T) {
	type fields struct {
		keygroups map[string]Keygroup
		sync.RWMutex
	}
	type args struct {
		data.Item
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"delete non-existent item from non-existent keygroup",
			fields{
				make(map[string]Keygroup),
				sync.RWMutex{},
			},
			args{
				data.Item{
					Keygroup: "keygroup",
					ID:       "0",
				},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				keygroups: tt.fields.keygroups,
				RWMutex:   tt.fields.RWMutex,
			}
			if err := s.Delete(tt.args.Item); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStorage_DeleteKeygroup(t *testing.T) {
	type fields struct {
		keygroups map[string]Keygroup
		sync.RWMutex
	}
	type args struct {
		data.Item
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"delete non-existent keygroup",
			fields{
				make(map[string]Keygroup),
				sync.RWMutex{},
			},
			args{
				data.Item{
					Keygroup: "keygroup",
					ID:       "0",
				},
			},
			true,
		},
		{"delete keygroup with empty name",
			fields{
				make(map[string]Keygroup),
				sync.RWMutex{},
			},
			args{data.Item{
				Keygroup: "",
			}},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				keygroups: tt.fields.keygroups,
				RWMutex:   tt.fields.RWMutex,
			}
			if err := s.DeleteKeygroup(tt.args.Item); (err != nil) != tt.wantErr {
				t.Errorf("DeleteKeygroup() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStorage_Read(t *testing.T) {
	type fields struct {
		keygroups map[string]Keygroup
		sync.RWMutex
	}
	type args struct {
		data.Item
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    data.Item
		wantErr bool
	}{
		{"read non-existent item from non-existent keygroup",
			fields{
				make(map[string]Keygroup),
				sync.RWMutex{},
			},
			args{
				data.Item{
					Keygroup: "keygroup",
					ID:       "0",
				},
			},
			data.Item{
				Keygroup: "keygroup",
				ID:       "0",
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				keygroups: tt.fields.keygroups,
				RWMutex:   tt.fields.RWMutex,
			}
			got, err := s.Read(tt.args.Item)
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
		sync.RWMutex
	}
	type args struct {
		data.Item
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"delete non-existent item from non-existent keygroup",
			fields{
				make(map[string]Keygroup),
				sync.RWMutex{},
			},
			args{
				data.Item{
					"keygroup",
					"0",
					"a new hello",
				},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				keygroups: tt.fields.keygroups,
				RWMutex:   tt.fields.RWMutex,
			}
			if err := s.Update(tt.args.Item); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
