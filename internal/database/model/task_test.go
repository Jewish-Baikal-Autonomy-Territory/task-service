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
			tp:   TaskPriorityUnknown,
			want: "unknown",
		},
		{
			name: "Low priority",
			tp:   TaskPriorityLow,
			want: "low",
		},
		{
			name: "Medium priority",
			tp:   TaskPriorityMedium,
			want: "medium",
		},
		{
			name: "High priority",
			tp:   TaskPriorityHigh,
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
		{
			name: "Unknown task status",
			ts:   TaskStatusUnknown,
			want: "unknown",
		},
		{
			name: "Active task status",
			ts:   TaskStatusActive,
			want: "active",
		},
		{
			name: "Completed task status",
			ts:   TaskStatusCompleted,
			want: "completed",
		},
		{
			name: "Deleting task status",
			ts:   TaskStatusPendingDeletion,
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
