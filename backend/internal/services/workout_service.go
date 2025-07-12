package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/visualtrecplans/backend/internal/handlers/workout"
	"github.com/visualtrecplans/backend/internal/models"
	"github.com/visualtrecplans/backend/internal/repositories"
	"github.com/visualtrecplans/backend/pkg/logger"
)

// WorkoutService interface defines workout-related business logic
type WorkoutService interface {
	Create(ctx context.Context, userID uuid.UUID, req *workout.CreateWorkoutRequest) (*models.Workout, error)
	GetByID(ctx context.Context, workoutID, userID uuid.UUID) (*models.Workout, error)
	GetUserWorkouts(ctx context.Context, userID uuid.UUID, filter *workout.WorkoutFilter) ([]*models.Workout, int64, error)
	Update(ctx context.Context, workoutID, userID uuid.UUID, req *workout.UpdateWorkoutRequest) (*models.Workout, error)
	Delete(ctx context.Context, workoutID, userID uuid.UUID) error
	GetStats(ctx context.Context, userID uuid.UUID, period string) (*models.WorkoutStats, error)
}

// workoutService implements WorkoutService
type workoutService struct {
	workoutRepo repositories.WorkoutRepository
	masterRepo  repositories.MasterRepository
	logger      *logger.Logger
}

// NewWorkoutService creates a new workout service
func NewWorkoutService(workoutRepo repositories.WorkoutRepository, masterRepo repositories.MasterRepository, logger *logger.Logger) WorkoutService {
	return &workoutService{
		workoutRepo: workoutRepo,
		masterRepo:  masterRepo,
		logger:      logger,
	}
}

// Create creates a new workout record
func (s *workoutService) Create(ctx context.Context, userID uuid.UUID, req *workout.CreateWorkoutRequest) (*models.Workout, error) {
	// Validate muscle group
	muscleGroups, err := s.masterRepo.GetAllMuscleGroups(ctx)
	if err != nil {
		s.logger.Error("Failed to get muscle groups", "error", err)
		return nil, fmt.Errorf("failed to validate muscle group: %w", err)
	}

	validMuscleGroup := false
	for _, mg := range muscleGroups {
		if mg.Code == req.MuscleGroup {
			validMuscleGroup = true
			break
		}
	}

	if !validMuscleGroup {
		return nil, ErrInvalidMuscleGroup
	}

	// Set default performed_at if not provided
	performedAt := req.PerformedAt
	if performedAt.IsZero() {
		performedAt = time.Now()
	}

	// Create workout model
	w := &models.Workout{
		UserID:       userID,
		MuscleGroup:  req.MuscleGroup,
		ExerciseName: req.ExerciseName,
		ExerciseIcon: req.ExerciseIcon,
		WeightKg:     req.WeightKg,
		Reps:         req.Reps,
		Sets:         req.Sets,
		Notes:        req.Notes,
		PerformedAt:  performedAt,
	}

	// Save to database
	if err := s.workoutRepo.Create(ctx, w); err != nil {
		s.logger.Error("Failed to create workout", "error", err, "userID", userID)
		return nil, fmt.Errorf("failed to create workout: %w", err)
	}

	s.logger.Info("Workout created successfully", "workoutID", w.ID, "userID", userID)
	return w, nil
}

// GetByID retrieves a workout by ID
func (s *workoutService) GetByID(ctx context.Context, workoutID, userID uuid.UUID) (*models.Workout, error) {
	w, err := s.workoutRepo.FindByID(ctx, workoutID)
	if err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			return nil, ErrWorkoutNotFound
		}
		s.logger.Error("Failed to get workout", "error", err, "workoutID", workoutID)
		return nil, fmt.Errorf("failed to get workout: %w", err)
	}

	// Check ownership
	if w.UserID != userID {
		return nil, ErrUnauthorized
	}

	return w, nil
}

// GetUserWorkouts retrieves workouts for a user with filtering
func (s *workoutService) GetUserWorkouts(ctx context.Context, userID uuid.UUID, filter *workout.WorkoutFilter) ([]*models.Workout, int64, error) {
	// Convert filter to repository filter
	repoFilter := &repositories.WorkoutFilter{
		UserID:       userID,
		MuscleGroup:  filter.MuscleGroup,
		StartDate:    filter.StartDate,
		EndDate:      filter.EndDate,
		ExerciseName: filter.ExerciseName,
		OrderBy:      filter.OrderBy,
		Order:        filter.Order,
		Limit:        filter.PerPage,
		Offset:       (filter.Page - 1) * filter.PerPage,
	}

	workouts, total, err := s.workoutRepo.FindByUserID(ctx, repoFilter)
	if err != nil {
		s.logger.Error("Failed to get user workouts", "error", err, "userID", userID)
		return nil, 0, fmt.Errorf("failed to get user workouts: %w", err)
	}

	return workouts, total, nil
}

// Update updates a workout record
func (s *workoutService) Update(ctx context.Context, workoutID, userID uuid.UUID, req *workout.UpdateWorkoutRequest) (*models.Workout, error) {
	// Get existing workout
	existing, err := s.GetByID(ctx, workoutID, userID)
	if err != nil {
		return nil, err
	}

	// Validate muscle group if provided
	if req.MuscleGroup != "" {
		muscleGroups, err := s.masterRepo.GetAllMuscleGroups(ctx)
		if err != nil {
			s.logger.Error("Failed to get muscle groups", "error", err)
			return nil, fmt.Errorf("failed to validate muscle group: %w", err)
		}

		validMuscleGroup := false
		for _, mg := range muscleGroups {
			if mg.Code == req.MuscleGroup {
				validMuscleGroup = true
				break
			}
		}

		if !validMuscleGroup {
			return nil, ErrInvalidMuscleGroup
		}
		existing.MuscleGroup = req.MuscleGroup
	}

	// Update fields if provided
	if req.ExerciseName != "" {
		existing.ExerciseName = req.ExerciseName
	}
	if req.ExerciseIcon != "" {
		existing.ExerciseIcon = req.ExerciseIcon
	}
	if req.WeightKg != nil {
		existing.WeightKg = req.WeightKg
	}
	if req.Reps != nil {
		existing.Reps = req.Reps
	}
	if req.Sets != nil {
		existing.Sets = req.Sets
	}
	if req.Notes != "" {
		existing.Notes = req.Notes
	}
	if !req.PerformedAt.IsZero() {
		existing.PerformedAt = req.PerformedAt
	}

	// Update in database
	if err := s.workoutRepo.Update(ctx, existing); err != nil {
		s.logger.Error("Failed to update workout", "error", err, "workoutID", workoutID)
		return nil, fmt.Errorf("failed to update workout: %w", err)
	}

	s.logger.Info("Workout updated successfully", "workoutID", workoutID, "userID", userID)
	return existing, nil
}

// Delete deletes a workout record
func (s *workoutService) Delete(ctx context.Context, workoutID, userID uuid.UUID) error {
	// Check ownership
	w, err := s.GetByID(ctx, workoutID, userID)
	if err != nil {
		return err
	}

	// Delete from database
	if err := s.workoutRepo.Delete(ctx, w.ID); err != nil {
		s.logger.Error("Failed to delete workout", "error", err, "workoutID", workoutID)
		return fmt.Errorf("failed to delete workout: %w", err)
	}

	s.logger.Info("Workout deleted successfully", "workoutID", workoutID, "userID", userID)
	return nil
}

// GetStats retrieves workout statistics for a user
func (s *workoutService) GetStats(ctx context.Context, userID uuid.UUID, period string) (*models.WorkoutStats, error) {
	stats, err := s.workoutRepo.GetStats(ctx, userID, period)
	if err != nil {
		s.logger.Error("Failed to get workout stats", "error", err, "userID", userID, "period", period)
		return nil, fmt.Errorf("failed to get workout stats: %w", err)
	}

	return stats, nil
}

// Service errors
var (
	ErrInvalidMuscleGroup = errors.New("invalid muscle group")
	ErrWorkoutNotFound    = errors.New("workout not found")
	ErrUnauthorized       = errors.New("unauthorized to access this workout")
)