package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/visualtrecplans/backend/internal/handlers/master"
	"github.com/visualtrecplans/backend/internal/models"
	"github.com/visualtrecplans/backend/internal/repositories"
	"github.com/visualtrecplans/backend/pkg/logger"
)

// MasterService interface defines master data-related business logic
type MasterService interface {
	GetAllMuscleGroups(ctx context.Context, lang string) ([]*models.MuscleGroup, error)
	GetMuscleGroupsByCategory(ctx context.Context, category, lang string) ([]*models.MuscleGroup, error)
	GetAllExercises(ctx context.Context, muscleGroupCode, lang string) ([]*models.Exercise, error)
	GetUserCustomExercises(ctx context.Context, userID uuid.UUID) ([]*models.Exercise, error)
	CreateCustomExercise(ctx context.Context, userID uuid.UUID, req *master.CreateCustomExerciseRequest) (*models.Exercise, error)
	GetAllExerciseIcons(ctx context.Context, category string) ([]*models.ExerciseIcon, error)
}

// masterService implements MasterService
type masterService struct {
	masterRepo repositories.MasterRepository
	logger     *logger.Logger
}

// NewMasterService creates a new master service
func NewMasterService(masterRepo repositories.MasterRepository, logger *logger.Logger) MasterService {
	return &masterService{
		masterRepo: masterRepo,
		logger:     logger,
	}
}

// GetAllMuscleGroups retrieves all muscle groups
func (s *masterService) GetAllMuscleGroups(ctx context.Context, lang string) ([]*models.MuscleGroup, error) {
	muscleGroups, err := s.masterRepo.GetAllMuscleGroups(ctx)
	if err != nil {
		s.logger.Error("Failed to get muscle groups", "error", err)
		return nil, fmt.Errorf("failed to get muscle groups: %w", err)
	}

	return muscleGroups, nil
}

// GetMuscleGroupsByCategory retrieves muscle groups by category
func (s *masterService) GetMuscleGroupsByCategory(ctx context.Context, category, lang string) ([]*models.MuscleGroup, error) {
	muscleGroups, err := s.masterRepo.GetMuscleGroupsByCategory(ctx, category)
	if err != nil {
		s.logger.Error("Failed to get muscle groups by category", "error", err, "category", category)
		return nil, fmt.Errorf("failed to get muscle groups by category: %w", err)
	}

	return muscleGroups, nil
}

// GetAllExercises retrieves all exercises
func (s *masterService) GetAllExercises(ctx context.Context, muscleGroupCode, lang string) ([]*models.Exercise, error) {
	var exercises []*models.Exercise
	var err error

	if muscleGroupCode != "" {
		exercises, err = s.masterRepo.GetExercisesByMuscleGroup(ctx, muscleGroupCode)
	} else {
		exercises, err = s.masterRepo.GetAllExercises(ctx)
	}

	if err != nil {
		s.logger.Error("Failed to get exercises", "error", err, "muscleGroupCode", muscleGroupCode)
		return nil, fmt.Errorf("failed to get exercises: %w", err)
	}

	return exercises, nil
}

// GetUserCustomExercises retrieves custom exercises for a user
func (s *masterService) GetUserCustomExercises(ctx context.Context, userID uuid.UUID) ([]*models.Exercise, error) {
	exercises, err := s.masterRepo.GetCustomExercisesByUser(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get custom exercises", "error", err, "userID", userID)
		return nil, fmt.Errorf("failed to get custom exercises: %w", err)
	}

	return exercises, nil
}

// CreateCustomExercise creates a new custom exercise
func (s *masterService) CreateCustomExercise(ctx context.Context, userID uuid.UUID, req *master.CreateCustomExerciseRequest) (*models.Exercise, error) {
	// Validate muscle group
	muscleGroups, err := s.masterRepo.GetAllMuscleGroups(ctx)
	if err != nil {
		s.logger.Error("Failed to get muscle groups", "error", err)
		return nil, fmt.Errorf("failed to validate muscle group: %w", err)
	}

	validMuscleGroup := false
	for _, mg := range muscleGroups {
		if mg.Code == req.MuscleGroupCode {
			validMuscleGroup = true
			break
		}
	}

	if !validMuscleGroup {
		return nil, ErrInvalidMuscleGroup
	}

	// Create exercise model
	exercise := &models.Exercise{
		MuscleGroupCode: req.MuscleGroupCode,
		NameJa:          req.Name,
		NameEn:          req.Name, // Default to same name for custom exercises
		IconName:        req.IconName,
		IsCustom:        true,
		CreatedBy:       &userID,
		SortOrder:       999, // Custom exercises at the end
	}

	// Save to database
	if err := s.masterRepo.CreateCustomExercise(ctx, exercise); err != nil {
		s.logger.Error("Failed to create custom exercise", "error", err, "userID", userID)
		return nil, fmt.Errorf("failed to create custom exercise: %w", err)
	}

	s.logger.Info("Custom exercise created successfully", "exerciseID", exercise.ID, "userID", userID)
	return exercise, nil
}

// GetAllExerciseIcons retrieves all exercise icons
func (s *masterService) GetAllExerciseIcons(ctx context.Context, category string) ([]*models.ExerciseIcon, error) {
	var icons []*models.ExerciseIcon
	var err error

	if category != "" {
		icons, err = s.masterRepo.GetExerciseIconsByCategory(ctx, category)
	} else {
		icons, err = s.masterRepo.GetAllExerciseIcons(ctx)
	}

	if err != nil {
		s.logger.Error("Failed to get exercise icons", "error", err, "category", category)
		return nil, fmt.Errorf("failed to get exercise icons: %w", err)
	}

	return icons, nil
}