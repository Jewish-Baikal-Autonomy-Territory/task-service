package utility

import (
	"testing"
	"time"
)

func TestGetEnv(t *testing.T) {
	type args struct {
		key      string
		fallback string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Missing environment variable",
			args: args{
				key:      "TEST_ENV_MISSING_STR",
				fallback: "missing-test",
			},
			want: "missing-test",
		},
		{
			name: "Set environment variable",
			args: args{
				key:      "TEST_ENV_STR",
				fallback: "missing-test",
			},
			want: "set-test",
		},
	}
	t.Setenv("TEST_ENV_STR", "set-test")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetEnv(tt.args.key, tt.args.fallback); got != tt.want {
				t.Errorf("GetEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetEnvBool(t *testing.T) {
	type args struct {
		key      string
		fallback bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Missing environment variable",
			args: args{
				key:      "TEST_ENV_MISSING_BOOL",
				fallback: false,
			},
		},
		{
			name: "Set environment variable",
			args: args{
				key:      "TEST_ENV_BOOL",
				fallback: false,
			},
			want: true,
		},
	}
	t.Setenv("TEST_ENV_BOOL", "true")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetEnvBool(tt.args.key, tt.args.fallback); got != tt.want {
				t.Errorf("GetEnvBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetEnvDuration(t *testing.T) {
	type args struct {
		key      string
		fallback time.Duration
	}
	tests := []struct {
		name string
		args args
		want time.Duration
	}{
		{
			name: "Missing environment variable",
			args: args{
				key:      "TEST_ENV_MISSING_DURATION",
				fallback: time.Duration(0),
			},
			want: time.Duration(0),
		},
		{
			name: "Incorrect format in environment variable",
			args: args{
				key:      "TEST_ENV_INCORRECT_DURATION",
				fallback: time.Duration(0),
			},
			want: time.Duration(0),
		},
		{
			name: "Set environment variable",
			args: args{
				key:      "TEST_ENV_DURATION",
				fallback: time.Duration(0),
			},
			want: 5 * time.Hour,
		},
	}
	t.Setenv("TEST_ENV_INCORRECT_DURATION", "5p")
	t.Setenv("TEST_ENV_DURATION", "5h")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetEnvDuration(tt.args.key, tt.args.fallback); got != tt.want {
				t.Errorf("GetEnvDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetEnvInt(t *testing.T) {
	type args struct {
		key      string
		fallback int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Missing environment variable",
			args: args{
				key:      "TEST_ENV_MISSING_INT",
				fallback: 0,
			},
			want: 0,
		},
		{
			name: "Value out of range",
			args: args{
				key:      "TEST_ENV_INCORRECT_INT",
				fallback: 0,
			},
			want: 0,
		},
		{
			name: "Set environment variable",
			args: args{
				key:      "TEST_ENV_INT",
				fallback: 0,
			},
			want: 67,
		},
	}
	t.Setenv("TEST_ENV_INCORRECT_INT", "10000000000000000000000000000")
	t.Setenv("TEST_ENV_INT", "67")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetEnvInt(tt.args.key, tt.args.fallback); got != tt.want {
				t.Errorf("GetEnvInt() = %v, want %v", got, tt.want)
			}
		})
	}
}
