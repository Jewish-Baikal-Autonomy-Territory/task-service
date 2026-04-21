package model

import "testing"

func TestTaskPriority_String(t *testing.T) {
	tests := []struct {
		name string
		tp   TaskPriority
		want string
	}{
		{
			name: "Unknown task priority",
			tp:   TASK_PRIORITY_UNKNOWN,
			want: "unknown",
		},
		{
			name: "Low priority",
			tp:   TASK_PRIORITY_LOW,
			want: "low",
		},
		{
			name: "Medium priority",
			tp:   TASK_PRIORITY_MEDIUM,
			want: "medium",
		},
		{
			name: "High priority",
			tp:   TASK_PRIORITY_HIGH,
			want: "high",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tp.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTaskStatus_String(t *testing.T) {
	tests := []struct {
		name string
		ts   TaskStatus
		want string
	}{
		// TODO: Add test cases.
		{
			name: "Unknown task status",
			ts:   TASK_STATUS_UNKNOWN,
			want: "unknown",
		},
		{
			name: "Active task status",
			ts:   TASK_STATUS_ACTIVE,
			want: "active",
		},
		{
			name: "Completed task status",
			ts:   TASK_STATUS_COMPLETED,
			want: "completed",
		},
		{
			name: "Deleting task status",
			ts:   TASK_STATUS_PENDING_DELETION,
			want: "deleting",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ts.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
