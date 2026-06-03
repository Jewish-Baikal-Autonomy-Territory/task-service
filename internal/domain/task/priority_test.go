package task

import "testing"

func TestParsePriority(t *testing.T) {
	type args struct {
		priority string
	}
	tests := []struct {
		name    string
		args    args
		want    Priority
		wantErr bool
	}{
		{
			name: "unknown priority",
			args: args{
				priority: "superrandomstring",
			},
			want:    PriorityUnknown,
			wantErr: true,
		},
		{
			name: "low priority",
			args: args{
				priority: "low",
			},
			want:    PriorityLow,
			wantErr: false,
		},
		{
			name: "medium priority",
			args: args{
				priority: "medium",
			},
			want:    PriorityMedium,
			wantErr: false,
		},
		{
			name: "high priority",
			args: args{
				priority: "high",
			},
			want:    PriorityHigh,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePriority(tt.args.priority)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParsePriority() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParsePriority() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPriority_String(t *testing.T) {
	tests := []struct {
		name string
		p    Priority
		want string
	}{
		{
			name: "unknown priority",
			p:    PriorityUnknown,
			want: "unknown",
		},
		{
			name: "low priority",
			p:    PriorityLow,
			want: "low",
		},
		{
			name: "medium priority",
			p:    PriorityMedium,
			want: "medium",
		},
		{
			name: "high priority",
			p:    PriorityHigh,
			want: "high",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPriority_Valid(t *testing.T) {
	tests := []struct {
		name string
		p    Priority
		want bool
	}{
		{
			name: "unknown priority",
			p:    PriorityUnknown,
			want: false,
		},
		{
			name: "low priority",
			p:    PriorityLow,
			want: true,
		},
		{
			name: "medium priority",
			p:    PriorityMedium,
			want: true,
		},
		{
			name: "high priority",
			p:    PriorityHigh,
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.Valid(); got != tt.want {
				t.Errorf("Valid() = %v, want %v", got, tt.want)
			}
		})
	}
}
