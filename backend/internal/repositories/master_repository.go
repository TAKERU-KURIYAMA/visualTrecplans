package repositories

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/visualtrecplans/backend/internal/models"
)

// MasterRepository defines the interface for master data operations
type MasterRepository interface {
	// Muscle Group operations
	GetAllMuscleGroups(ctx context.Context) ([]*models.MuscleGroup, error)
	GetMuscleGroupByCode(ctx context.Context, code string) (*models.MuscleGroup, error)
	GetMuscleGroupsByCategory(ctx context.Context, category string) ([]*models.MuscleGroup, error)
	
	// Exercise operations
	GetAllExercises(ctx context.Context) ([]*models.Exercise, error)
	GetExercisesByMuscleGroup(ctx context.Context, muscleGroupCode string) ([]*models.Exercise, error)
	GetExerciseByID(ctx context.Context, id int) (*models.Exercise, error)
	CreateCustomExercise(ctx context.Context, exercise *models.Exercise) error
	GetCustomExercisesByUser(ctx context.Context, userID uuid.UUID) ([]*models.Exercise, error)
	
	// Exercise Icon operations
	GetAllExerciseIcons(ctx context.Context) ([]*models.ExerciseIcon, error)
	GetExerciseIconsByCategory(ctx context.Context, category string) ([]*models.ExerciseIcon, error)
	GetExerciseIconByName(ctx context.Context, name string) (*models.ExerciseIcon, error)
}

// masterRepository implements MasterRepository interface
type masterRepository struct {
	db *gorm.DB
}

// NewMasterRepository creates a new master repository
func NewMasterRepository(db *gorm.DB) MasterRepository {
	return &masterRepository{db: db}
}

// GetAllMuscleGroups retrieves all muscle groups ordered by sort_order
func (r *masterRepository) GetAllMuscleGroups(ctx context.Context) ([]*models.MuscleGroup, error) {
	var muscleGroups []*models.MuscleGroup
	
	if err := r.db.WithContext(ctx).Order("sort_order ASC, name_ja ASC").Find(&muscleGroups).Error; err != nil {
		return nil, fmt.Errorf("failed to get muscle groups: %w", err)
	}
	
	return muscleGroups, nil
}

// GetMuscleGroupByCode retrieves a muscle group by its code
func (r *masterRepository) GetMuscleGroupByCode(ctx context.Context, code string) (*models.MuscleGroup, error) {
	var muscleGroup models.MuscleGroup
	
	if err := r.db.WithContext(ctx).Where("code = ?", code).First(&muscleGroup).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("muscle group not found")
		}
		return nil, fmt.Errorf("failed to get muscle group: %w", err)
	}
	
	return &muscleGroup, nil
}

// GetMuscleGroupsByCategory retrieves muscle groups by category
func (r *masterRepository) GetMuscleGroupsByCategory(ctx context.Context, category string) ([]*models.MuscleGroup, error) {
	var muscleGroups []*models.MuscleGroup
	
	if err := r.db.WithContext(ctx).
		Where("category = ?", category).
		Order("sort_order ASC, name_ja ASC").
		Find(&muscleGroups).Error; err != nil {
		return nil, fmt.Errorf("failed to get muscle groups by category: %w", err)
	}
	
	return muscleGroups, nil
}

// GetAllExercises retrieves all exercises with muscle group information
func (r *masterRepository) GetAllExercises(ctx context.Context) ([]*models.Exercise, error) {
	var exercises []*models.Exercise
	
	if err := r.db.WithContext(ctx).
		Preload("MuscleGroup").
		Order("muscle_group_code ASC, sort_order ASC, name_ja ASC").
		Find(&exercises).Error; err != nil {
		return nil, fmt.Errorf("failed to get exercises: %w", err)
	}
	
	return exercises, nil
}

// GetExercisesByMuscleGroup retrieves exercises for a specific muscle group
func (r *masterRepository) GetExercisesByMuscleGroup(ctx context.Context, muscleGroupCode string) ([]*models.Exercise, error) {
	var exercises []*models.Exercise
	
	if err := r.db.WithContext(ctx).
		Preload("MuscleGroup").
		Where("muscle_group_code = ?", muscleGroupCode).
		Order("sort_order ASC, name_ja ASC").
		Find(&exercises).Error; err != nil {
		return nil, fmt.Errorf("failed to get exercises by muscle group: %w", err)
	}
	
	return exercises, nil
}

// GetExerciseByID retrieves an exercise by its ID
func (r *masterRepository) GetExerciseByID(ctx context.Context, id int) (*models.Exercise, error) {
	var exercise models.Exercise
	
	if err := r.db.WithContext(ctx).
		Preload("MuscleGroup").
		Preload("Creator").
		Where("id = ?", id).
		First(&exercise).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("exercise not found")
		}
		return nil, fmt.Errorf("failed to get exercise: %w", err)
	}
	
	return &exercise, nil
}

// CreateCustomExercise creates a new custom exercise
func (r *masterRepository) CreateCustomExercise(ctx context.Context, exercise *models.Exercise) error {
	exercise.IsCustom = true
	
	if err := r.db.WithContext(ctx).Create(exercise).Error; err != nil {
		return fmt.Errorf("failed to create custom exercise: %w", err)
	}
	
	return nil
}

// GetCustomExercisesByUser retrieves custom exercises created by a specific user
func (r *masterRepository) GetCustomExercisesByUser(ctx context.Context, userID uuid.UUID) ([]*models.Exercise, error) {
	var exercises []*models.Exercise
	
	if err := r.db.WithContext(ctx).
		Preload("MuscleGroup").
		Where("is_custom = ? AND created_by = ?", true, userID).
		Order("muscle_group_code ASC, sort_order ASC, name_ja ASC").
		Find(&exercises).Error; err != nil {
		return nil, fmt.Errorf("failed to get custom exercises: %w", err)
	}
	
	return exercises, nil
}

// GetAllExerciseIcons retrieves all exercise icons ordered by category
func (r *masterRepository) GetAllExerciseIcons(ctx context.Context) ([]*models.ExerciseIcon, error) {
	var icons []*models.ExerciseIcon
	
	if err := r.db.WithContext(ctx).Order("category ASC, name ASC").Find(&icons).Error; err != nil {
		return nil, fmt.Errorf("failed to get exercise icons: %w", err)
	}
	
	return icons, nil
}

// GetExerciseIconsByCategory retrieves exercise icons by category
func (r *masterRepository) GetExerciseIconsByCategory(ctx context.Context, category string) ([]*models.ExerciseIcon, error) {
	var icons []*models.ExerciseIcon
	
	if err := r.db.WithContext(ctx).
		Where("category = ?", category).
		Order("name ASC").
		Find(&icons).Error; err != nil {
		return nil, fmt.Errorf("failed to get exercise icons by category: %w", err)
	}
	
	return icons, nil
}

// GetExerciseIconByName retrieves an exercise icon by its name
func (r *masterRepository) GetExerciseIconByName(ctx context.Context, name string) (*models.ExerciseIcon, error) {
	var icon models.ExerciseIcon
	
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&icon).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("exercise icon not found")
		}
		return nil, fmt.Errorf("failed to get exercise icon: %w", err)
	}
	
	return &icon, nil
}