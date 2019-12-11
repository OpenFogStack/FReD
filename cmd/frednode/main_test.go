package main

import (
	"github.com/gin-gonic/gin"
	"reflect"
	"testing"
)

func Test_setupRouter(t *testing.T) {
	type args struct {
		sd Storage
		kg Keygroups
	}
	tests := []struct {
		name  string
		args  args
		wantR *gin.Engine
	}{
		// TODO: Add Test Cases
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotR := setupRouter(tt.args.sd, tt.args.kg); !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("setupRouter() = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}