package task

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
)

func TestNewCompleteTaskCommand(t *testing.T) {
	var (
		testID      = uuid.New()
		testOwnerID = uuid.New()
	)
	type args struct {
		id      uuid.UUID
		ownerID uuid.UUID
	}
	tests := []struct {
		name    string
		args    args
		want    CompleteTaskCommand
		wantErr bool
	}{
		{
			name: "invalid id and owner id",
			args: args{
				id:      uuid.Nil,
				ownerID: uuid.Nil,
			},
			want:    CompleteTaskCommand{},
			wantErr: true,
		},
		{
			name: "valid id only",
			args: args{
				id:      testID,
				ownerID: uuid.Nil,
			},
			want:    CompleteTaskCommand{},
			wantErr: true,
		},
		{
			name: "valid owner id only",
			args: args{
				id:      uuid.Nil,
				ownerID: testOwnerID,
			},
			want:    CompleteTaskCommand{},
			wantErr: true,
		},
		{
			name: "valid id and owner id",
			args: args{
				id:      testID,
				ownerID: testOwnerID,
			},
			want: CompleteTaskCommand{
				ID:      testID,
				OwnerID: testOwnerID,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewCompleteTaskCommand(tt.args.id, tt.args.ownerID)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCompleteTaskCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCompleteTaskCommand() got = %v, want %v", got, tt.want)
			}
		})
	}
}
