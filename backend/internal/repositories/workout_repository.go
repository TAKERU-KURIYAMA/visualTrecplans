package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/visualtrecplans/backend/internal/models"
)

// WorkoutRepository defines the interface for workout data operations
type WorkoutRepository interface {
	Create(ctx context.Context, workout *models.Workout) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Workout, error)
	FindByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Workout, error)
	FindByUserIDAndDateRange(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time, limit, offset int) ([]*models.Workout, error)
	FindByUserIDAndMuscleGroup(ctx context.Context, userID uuid.UUID, muscleGroup string, limit, offset int) ([]*models.Workout, error)
	Update(ctx context.Context, workout *models.Workout) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetStats(ctx context.Context, userID uuid.UUID, period string) (*models.WorkoutStats, error)
	Count(ctx context.Context, userID uuid.UUID) (int64, error)
	CountByMuscleGroup(ctx context.Context, userID uuid.UUID, muscleGroup string) (int64, error)
}

// workoutRepository implements WorkoutRepository interface
type workoutRepository struct {
	db *gorm.DB
}

// NewWorkoutRepository creates a new workout repository
func NewWorkoutRepository(db *gorm.DB) WorkoutRepository {
	return &workoutRepository{db: db}
}

// Create creates a new workout record
func (r *workoutRepository) Create(ctx context.Context, workout *models.Workout) error {
	if err := r.db.WithContext(ctx).Create(workout).Error; err != nil {
		return fmt.Errorf("failed to create workout: %w", err)
	}
	return nil
}

// FindByID finds a workout by ID
func (r *workoutRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Workout, error) {
	var workout models.Workout
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&workout).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("workout not found")
		}
		return nil, fmt.Errorf("failed to find workout: %w", err)
	}
	return &workout, nil
}

// FindByUserID finds workouts by user ID with pagination
func (r *workoutRepository) FindByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Workout, error) {
	var workouts []*models.Workout
	
	query := r.db.WithContext(ctx).Where("user_id = ?", userID)
	
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	
	if err := query.Order("performed_at DESC").Find(&workouts).Error; err != nil {
		return nil, fmt.Errorf("failed to find workouts: %w", err)
	}
	
	return workouts, nil
}

// FindByUserIDAndDateRange finds workouts by user ID within a date range
func (r *workoutRepository) FindByUserIDAndDateRange(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time, limit, offset int) ([]*models.Workout, error) {
	var workouts []*models.Workout
	
	query := r.db.WithContext(ctx).
		Where("user_id = ? AND performed_at >= ? AND performed_at <= ?", userID, startDate, endDate)
	
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	
	if err := query.Order("performed_at DESC").Find(&workouts).Error; err != nil {
		return nil, fmt.Errorf("failed to find workouts: %w", err)
	}
	
	return workouts, nil
}

// FindByUserIDAndMuscleGroup finds workouts by user ID and muscle group
func (r *workoutRepository) FindByUserIDAndMuscleGroup(ctx context.Context, userID uuid.UUID, muscleGroup string, limit, offset int) ([]*models.Workout, error) {
	var workouts []*models.Workout
	
	query := r.db.WithContext(ctx).Where("user_id = ? AND muscle_group = ?", userID, muscleGroup)
	
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	
	if err := query.Order("performed_at DESC").Find(&workouts).Error; err != nil {
		return nil, fmt.Errorf("failed to find workouts: %w", err)
	}
	
	return workouts, nil
}

// Update updates a workout record
func (r *workoutRepository) Update(ctx context.Context, workout *models.Workout) error {
	if err := r.db.WithContext(ctx).Save(workout).Error; err != nil {
		return fmt.Errorf("failed to update workout: %w", err)
	}
	return nil
}

// Delete soft deletes a workout record
func (r *workoutRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&models.Workout{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("failed to delete workout: %w", err)
	}
	return nil
}

// GetStats calculates workout statistics for a user
func (r *workoutRepository) GetStats(ctx context.Context, userID uuid.UUID, period string) (*models.WorkoutStats, error) {
	stats := &models.WorkoutStats{
		MuscleGroupStats: make(map[string]int),
		WeeklyStats:      []models.WeeklyWorkoutStats{},
		MonthlyStats:     []models.MonthlyWorkoutStats{},
	}
	
	// Total workouts count
	var totalCount int64
	if err := r.db.WithContext(ctx).Model(&models.Workout{}).
		Where("user_id = ?", userID).Count(&totalCount).Error; err != nil {
		return nil, fmt.Errorf("failed to count workouts: %w", err)
	}
	stats.TotalWorkouts = int(totalCount)
	
	// Calculate total volume and muscle group stats
	var workouts []models.Workout
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&workouts).Error; err != nil {
		return nil, fmt.Errorf("failed to find workouts for stats: %w", err)
	}
	
	for _, workout := range workouts {
		// Calculate total volume
		stats.TotalVolume += workout.CalculateVolume()
		
		// Count muscle group stats
		stats.MuscleGroupStats[workout.MuscleGroup]++
	}
	
	// Calculate period-specific stats based on the period parameter
	switch period {
	case "week":
		if err := r.calculateWeeklyStats(ctx, userID, stats); err != nil {
			return nil, err
		}
	case "month":
		if err := r.calculateMonthlyStats(ctx, userID, stats); err != nil {
			return nil, err
		}
	case "all":
		// Calculate both weekly and monthly stats
		if err := r.calculateWeeklyStats(ctx, userID, stats); err != nil {
			return nil, err
		}
		if err := r.calculateMonthlyStats(ctx, userID, stats); err != nil {
			return nil, err
		}
	}
	
	return stats, nil
}

// calculateWeeklyStats calculates weekly workout statistics
func (r *workoutRepository) calculateWeeklyStats(ctx context.Context, userID uuid.UUID, stats *models.WorkoutStats) error {
	// Get workouts for the last 12 weeks
	twelveWeeksAgo := time.Now().AddDate(0, 0, -84) // 12 * 7 days
	
	var weeklyData []struct {
		Week        string  `json:"week"`
		Workouts    int     `json:"workouts"`
		TotalVolume float64 `json:"total_volume"`
	}
	
	err := r.db.WithContext(ctx).Raw(`
		SELECT 
			TO_CHAR(performed_at, 'YYYY-"W"WW') as week,
			COUNT(*) as workouts,
			COALESCE(SUM(
				CASE 
					WHEN sets IS NOT NULL AND reps IS NOT NULL AND weight_kg IS NOT NULL 
					THEN sets * reps * weight_kg 
					ELSE 0 
				END
			), 0) as total_volume
		FROM workouts 
		WHERE user_id = ? AND performed_at >= ? AND deleted_at IS NULL
		GROUP BY TO_CHAR(performed_at, 'YYYY-"W"WW')
		ORDER BY week DESC
		LIMIT 12
	`, userID, twelveWeeksAgo).Scan(&weeklyData).Error
	
	if err != nil {
		return fmt.Errorf("failed to calculate weekly stats: %w", err)
	}
	
	for _, data := range weeklyData {
		stats.WeeklyStats = append(stats.WeeklyStats, models.WeeklyWorkoutStats{
			Week:        data.Week,
			Workouts:    data.Workouts,
			TotalVolume: data.TotalVolume,
		})
	}
	
	return nil
}

// calculateMonthlyStats calculates monthly workout statistics
func (r *workoutRepository) calculateMonthlyStats(ctx context.Context, userID uuid.UUID, stats *models.WorkoutStats) error {
	// Get workouts for the last 12 months
	twelveMonthsAgo := time.Now().AddDate(0, -12, 0)
	
	var monthlyData []struct {
		Month       string  `json:"month"`
		Workouts    int     `json:"workouts"`
		TotalVolume float64 `json:"total_volume"`
	}
	
	err := r.db.WithContext(ctx).Raw(`
		SELECT 
			TO_CHAR(performed_at, 'YYYY-MM') as month,
			COUNT(*) as workouts,
			COALESCE(SUM(
				CASE 
					WHEN sets IS NOT NULL AND reps IS NOT NULL AND weight_kg IS NOT NULL 
					THEN sets * reps * weight_kg 
					ELSE 0 
				END
			), 0) as total_volume
		FROM workouts 
		WHERE user_id = ? AND performed_at >= ? AND deleted_at IS NULL
		GROUP BY TO_CHAR(performed_at, 'YYYY-MM')
		ORDER BY month DESC
		LIMIT 12
	`, userID, twelveMonthsAgo).Scan(&monthlyData).Error
	
	if err != nil {
		return fmt.Errorf("failed to calculate monthly stats: %w", err)
	}
	
	for _, data := range monthlyData {
		stats.MonthlyStats = append(stats.MonthlyStats, models.MonthlyWorkoutStats{
			Month:       data.Month,
			Workouts:    data.Workouts,
			TotalVolume: data.TotalVolume,
		})
	}
	
	return nil
}

// Count returns the total number of workouts for a user
func (r *workoutRepository) Count(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.Workout{}).
		Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count workouts: %w", err)
	}
	return count, nil
}

// CountByMuscleGroup returns the number of workouts for a specific muscle group
func (r *workoutRepository) CountByMuscleGroup(ctx context.Context, userID uuid.UUID, muscleGroup string) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.Workout{}).
		Where("user_id = ? AND muscle_group = ?", userID, muscleGroup).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count workouts by muscle group: %w", err)
	}
	return count, nil
}