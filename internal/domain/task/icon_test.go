package task

import "testing"

func TestIcon_String(t *testing.T) {
	tests := []struct {
		name string
		i    Icon
		want string
	}{
		{
			name: "unknown icon",
			i:    IconUnknown,
			want: "unknown",
		},
		{
			name: "mark icon",
			i:    IconMark,
			want: "mark",
		},
		{
			name: "home icon",
			i:    IconHome,
			want: "home",
		},
		{
			name: "job icon",
			i:    IconJob,
			want: "job",
		},
		{
			name: "supermarket icon",
			i:    IconSupermarket,
			want: "supermarket",
		},
		{
			name: "cafe icon",
			i:    IconCafe,
			want: "cafe",
		},
		{
			name: "activity icon",
			i:    IconActivity,
			want: "activity",
		},
		{
			name: "drive icon",
			i:    IconDrive,
			want: "drive",
		},
		{
			name: "flight icon",
			i:    IconFlight,
			want: "flight",
		},
		{
			name: "star icon",
			i:    IconStar,
			want: "star",
		},
		{
			name: "flag icon",
			i:    IconFlag,
			want: "flag",
		},
		{
			name: "hospital icon",
			i:    IconHospital,
			want: "hospital",
		},
		{
			name: "outdoor icon",
			i:    IconOutdoor,
			want: "outdoor",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIcon_Valid(t *testing.T) {
	tests := []struct {
		name string
		i    Icon
		want bool
	}{
		{
			name: "unknown icon",
			i:    IconUnknown,
			want: false,
		},
		{
			name: "mark icon",
			i:    IconMark,
			want: true,
		},
		{
			name: "home icon",
			i:    IconHome,
			want: true,
		},
		{
			name: "job icon",
			i:    IconJob,
			want: true,
		},
		{
			name: "supermarket icon",
			i:    IconSupermarket,
			want: true,
		},
		{
			name: "cafe icon",
			i:    IconCafe,
			want: true,
		},
		{
			name: "activity icon",
			i:    IconActivity,
			want: true,
		},
		{
			name: "drive icon",
			i:    IconDrive,
			want: true,
		},
		{
			name: "flight icon",
			i:    IconFlight,
			want: true,
		},
		{
			name: "star icon",
			i:    IconStar,
			want: true,
		},
		{
			name: "flag icon",
			i:    IconFlag,
			want: true,
		},
		{
			name: "hospital icon",
			i:    IconHospital,
			want: true,
		},
		{
			name: "outdoor icon",
			i:    IconOutdoor,
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.Valid(); got != tt.want {
				t.Errorf("Valid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseIcon(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    Icon
		wantErr bool
	}{
		{
			name: "unknown icon",
			args: args{
				s: "superrandomstring",
			},
			want:    IconUnknown,
			wantErr: true,
		},
		{
			name: "mark icon",
			args: args{
				s: "mark",
			},
			want:    IconMark,
			wantErr: false,
		},
		{
			name: "home icon",
			args: args{
				s: "home",
			},
			want:    IconHome,
			wantErr: false,
		},
		{
			name: "job icon",
			args: args{
				s: "job",
			},
			want:    IconJob,
			wantErr: false,
		},
		{
			name: "supermarket icon",
			args: args{
				s: "supermarket",
			},
			want:    IconSupermarket,
			wantErr: false,
		},
		{
			name: "cafe icon",
			args: args{
				s: "cafe",
			},
			want:    IconCafe,
			wantErr: false,
		},
		{
			name: "activity icon",
			args: args{
				s: "activity",
			},
			want:    IconActivity,
			wantErr: false,
		},
		{
			name: "drive icon",
			args: args{
				s: "drive",
			},
			want:    IconDrive,
			wantErr: false,
		},
		{
			name: "flight icon",
			args: args{
				s: "flight",
			},
			want:    IconFlight,
			wantErr: false,
		},
		{
			name: "star icon",
			args: args{
				s: "star",
			},
			want:    IconStar,
			wantErr: false,
		},
		{
			name: "flag icon",
			args: args{
				s: "flag",
			},
			want:    IconFlag,
			wantErr: false,
		},
		{
			name: "hospital icon",
			args: args{
				s: "hospital",
			},
			want:    IconHospital,
			wantErr: false,
		},
		{
			name: "outdoor icon",
			args: args{
				s: "outdoor",
			},
			want:    IconOutdoor,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseIcon(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseIcon() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseIcon() got = %v, want %v", got, tt.want)
			}
		})
	}
}
