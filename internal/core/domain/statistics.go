package domain

import "time"

type Statistics struct {
	TaskCreated               int
	TaskCompleted             int
	TasksCompletedRate        *float64
	TaskAverageCompletionTime *time.Duration
}

func NewStatistics(
	tasksCreated int,
	tasksCompleted int,
	tasksCompletedRate *float64,
	tasksAverageCompletionTime *time.Duration,
) Statistics {
	return Statistics{
		TaskCreated:               tasksCreated,
		TaskCompleted:             tasksCompleted,
		TasksCompletedRate:        tasksCompletedRate,
		TaskAverageCompletionTime: tasksAverageCompletionTime,
	}
}
