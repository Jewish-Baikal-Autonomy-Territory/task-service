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
			name: "missing environment variable",
			args: args{
				key:      "MISSING_ENV_VAR",
				fallback: "default",
			},
			want: "default",
		},
		{
			name: "existing environment variable",
			args: args{
				key:      "ENV_TEST_VAR",
				fallback: "default",
			},
			want: "env_test_var",
		},
	}

	t.Setenv("ENV_TEST_VAR", "env_test_var")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetEnv(tt.args.key, tt.args.fallback); got != tt.want {
				t.Errorf("GetEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetEnvBool1(t *testing.T) {
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
			name: "missing environment variable",
			args: args{
				key:      "MISSING_ENV_VAR",
				fallback: true,
			},
			want: true,
		},
		{
			name: "incorrect environment variable",
			args: args{
				key:      "ENV_INCORRECT_VAR",
				fallback: true,
			},
			want: true,
		},
		{
			name: "existing environment variable",
			args: args{
				key:      "ENV_TEST_VAR",
				fallback: false,
			},
			want: false,
		},
	}

	t.Setenv("ENV_INCORRECT_VAR", "some-random-string")
	t.Setenv("ENV_TEST_VAR", "false")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetEnvBool(tt.args.key, tt.args.fallback); got != tt.want {
				t.Errorf("GetEnvBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetEnvDuration1(t *testing.T) {
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
			name: "missing environment variable",
			args: args{
				key:      "MISSING_ENV_VAR",
				fallback: time.Second,
			},
			want: time.Second,
		},
		{
			name: "incorrect environment variable",
			args: args{
				key:      "ENV_INCORRECT_VAR",
				fallback: time.Second,
			},
			want: time.Second,
		},
		{
			name: "existing environment variable",
			args: args{
				key:      "ENV_TEST_VAR",
				fallback: 67*time.Hour + 52*time.Minute,
			},
			want: 67*time.Hour + 52*time.Minute,
		},
	}

	t.Setenv("ENV_INCORRECT_VAR", "some-random-string")
	t.Setenv("ENV_TEST_VAR", "67h52m")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetEnvDuration(tt.args.key, tt.args.fallback); got != tt.want {
				t.Errorf("GetEnvDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetEnvFloat(t *testing.T) {
	type args struct {
		key      string
		fallback float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "missing environment variable",
			args: args{
				key:      "MISSING_ENV_VAR",
				fallback: 320.0,
			},
			want: 320.0,
		},
		{
			name: "incorrect environment variable",
			args: args{
				key:      "ENV_INCORRECT_VAR",
				fallback: 120.0,
			},
			want: 120.0,
		},
		{
			name: "existing environment variable",
			args: args{
				key:      "ENV_TEST_VAR",
				fallback: 52.0,
			},
			want: 67.0,
		},
	}

	t.Setenv("ENV_INCORRECT_VAR", "some-random-string")
	t.Setenv("ENV_TEST_VAR", "67.0")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetEnvFloat(tt.args.key, tt.args.fallback); got != tt.want {
				t.Errorf("GetEnvFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetEnvInt1(t *testing.T) {
	type args struct {
		key      string
		fallback int32
	}
	tests := []struct {
		name string
		args args
		want int32
	}{
		{
			name: "missing environment variable",
			args: args{
				key:      "MISSING_ENV_VAR",
				fallback: 32,
			},
			want: 32,
		},
		{
			name: "incorrect environment variable",
			args: args{
				key:      "ENV_INCORRECT_VAR",
				fallback: 67,
			},
			want: 67,
		},
		{
			name: "existing environment variable",
			args: args{
				key:      "ENV_TEST_VAR",
				fallback: 10,
			},
			want: 52,
		},
	}

	t.Setenv("ENV_INCORRECT_VAR", "some-random-string")
	t.Setenv("ENV_TEST_VAR", "52")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetEnvInt(tt.args.key, tt.args.fallback); got != tt.want {
				t.Errorf("GetEnvInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetEnvUInt(t *testing.T) {
	type args struct {
		key      string
		fallback uint32
	}
	tests := []struct {
		name string
		args args
		want uint32
	}{
		{
			name: "missing environment variable",
			args: args{
				key:      "MISSING_ENV_VAR",
				fallback: 120,
			},
			want: 120,
		},
		{
			name: "incorrect environment variable",
			args: args{
				key:      "ENV_INCORRECT_VAR",
				fallback: 52,
			},
			want: 52,
		},
		{
			name: "existing environment variable",
			args: args{
				key:      "ENV_TEST_VAR",
				fallback: 100,
			},
			want: 67,
		},
	}

	t.Setenv("ENV_INCORRECT_VAR", "some-random-string")
	t.Setenv("ENV_TEST_VAR", "67")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetEnvUInt(tt.args.key, tt.args.fallback); got != tt.want {
				t.Errorf("GetEnvUInt() = %v, want %v", got, tt.want)
			}
		})
	}
}
