package models

type ProjectAnalytics struct {
	TaskCompletionRate float64 `bson:"task_completion_rate" json:"taskCompletionRate"`
	OverdueTasks       int     `bson:"overdue_tasks" json:"overdueTasks"`
	ActiveMembers      int     `bson:"active_members" json:"activeMembers"`
}
