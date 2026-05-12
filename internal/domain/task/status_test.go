package task

import "testing"

func TestParseStatus(t *testing.T) {
	type args struct {
		status string
	}
	tests := []struct {
		name    string
		args    args
		want    Status
		wantErr bool
	}{
		{
			name: "invalid status",
			args: args{
				status: "undefined",
			},
			want:    StatusUnknown,
			wantErr: true,
		},
		{
			name: "pending status",
			args: args{
				status: "pending",
			},
			want:    StatusPending,
			wantErr: false,
		},
		{
			name: "complete status",
			args: args{
				status: "completed",
			},
			want:    StatusCompleted,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseStatus(tt.args.status)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseStatus() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStatus_String(t *testing.T) {
	tests := []struct {
		name string
		s    Status
		want string
	}{
		{
			name: "invalid status",
			s:    StatusUnknown,
			want: "unknown",
		},
		{
			name: "pending status",
			s:    StatusPending,
			want: "pending",
		},
		{
			name: "complete status",
			s:    StatusCompleted,
			want: "completed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStatus_Valid(t *testing.T) {
	tests := []struct {
		name string
		s    Status
		want bool
	}{
		{
			name: "invalid status",
			s:    StatusUnknown,
			want: false,
		},
		{
			name: "pending status",
			s:    StatusPending,
			want: true,
		},
		{
			name: "complete status",
			s:    StatusCompleted,
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Valid(); got != tt.want {
				t.Errorf("Valid() = %v, want %v", got, tt.want)
			}
		})
	}
}
