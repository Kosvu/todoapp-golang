package statistic_service

import (
	"context"
	"fmt"
	"time"

	"github.com/Kosvu/todoapp-golang/internal/core/domain"
	core_errors "github.com/Kosvu/todoapp-golang/internal/core/errors"
)

func (s *StatisticsService) GetStatistics(
	ctx context.Context,
	userID *int,
	from *time.Time,
	to *time.Time,
) (domain.Statistics, error) {
	if from != nil && to != nil {
		//to.Before to до from (что невозможно)
		//to.Equal to равен from (что невозможно)
		if to.Before(*from) || to.Equal(*from) {
			return domain.Statistics{}, fmt.Errorf(
				"'to' must be after 'from': %w",
				core_errors.ErrInvalidArgument,
			)
		}
	}

	tasks, err := s.statisticsRepository.GetTasks(ctx, userID, from, to)
	if err != nil {
		return domain.Statistics{}, fmt.Errorf("get tasks from repository: %w", err)
	}

	statisics := calcStatistics(tasks)

	return statisics, nil

}

func calcStatistics(tasks []domain.Task) domain.Statistics {
	if len(tasks) == 0 {
		return domain.NewStatistics(0, 0, nil, nil)
	}

	tasksCreated := len(tasks)

	tasksCompleted := 0
	var totalCompletionDuration time.Duration
	for _, tasks := range tasks {
		if tasks.Completed {
			tasksCompleted++
		}

		completionDuration := tasks.CompletionDuration()
		if completionDuration != nil {
			totalCompletionDuration += *completionDuration
		}
	}

	tasksCompletedRate := float64(tasksCompleted) / float64(tasksCreated) * 100

	var tasksAverageCompletionTime *time.Duration
	if tasksCompleted > 0 && totalCompletionDuration != 0 {
		avg := totalCompletionDuration / time.Duration(tasksCompleted)

		tasksAverageCompletionTime = &avg
	}

	return domain.NewStatistics(tasksCreated, tasksCompleted, &tasksCompletedRate, tasksAverageCompletionTime)
}
