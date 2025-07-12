package workout

import (
	"time"
)

// CreateWorkoutRequest represents the request to create a new workout record
type CreateWorkoutRequest struct {
	MuscleGroup  string    `json:"muscle_group" binding:"required,muscle_group"`
	ExerciseName string    `json:"exercise_name" binding:"required,min=1,max=100"`
	ExerciseIcon string    `json:"exercise_icon" binding:"omitempty,max=50"`
	WeightKg     *float64  `json:"weight_kg" binding:"omitempty,min=0,max=999.99"`
	Reps         *int      `json:"reps" binding:"omitempty,min=1,max=999"`
	Sets         *int      `json:"sets" binding:"omitempty,min=1,max=99"`
	Notes        string    `json:"notes" binding:"omitempty,max=500"`
	PerformedAt  time.Time `json:"performed_at"`
}

// UpdateWorkoutRequest represents the request to update a workout record
type UpdateWorkoutRequest struct {
	MuscleGroup  string    `json:"muscle_group" binding:"omitempty,muscle_group"`
	ExerciseName string    `json:"exercise_name" binding:"omitempty,min=1,max=100"`
	ExerciseIcon string    `json:"exercise_icon" binding:"omitempty,max=50"`
	WeightKg     *float64  `json:"weight_kg" binding:"omitempty,min=0,max=999.99"`
	Reps         *int      `json:"reps" binding:"omitempty,min=1,max=999"`
	Sets         *int      `json:"sets" binding:"omitempty,min=1,max=99"`
	Notes        string    `json:"notes" binding:"omitempty,max=500"`
	PerformedAt  time.Time `json:"performed_at"`
}

// WorkoutResponse represents a workout in the API response
type WorkoutResponse struct {
	ID           string    `json:"id"`
	MuscleGroup  string    `json:"muscle_group"`
	ExerciseName string    `json:"exercise_name"`
	ExerciseIcon string    `json:"exercise_icon,omitempty"`
	WeightKg     *float64  `json:"weight_kg,omitempty"`
	Reps         *int      `json:"reps,omitempty"`
	Sets         *int      `json:"sets,omitempty"`
	Notes        string    `json:"notes,omitempty"`
	PerformedAt  time.Time `json:"performed_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// WorkoutListResponse represents a list of workouts with pagination info
type WorkoutListResponse struct {
	Workouts   []WorkoutResponse `json:"workouts"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	PerPage    int               `json:"per_page"`
	TotalPages int               `json:"total_pages"`
}

// WorkoutFilter represents filtering options for workout queries
type WorkoutFilter struct {
	MuscleGroup  string    `form:"muscle_group" binding:"omitempty,muscle_group"`
	StartDate    time.Time `form:"start_date" binding:"omitempty"`
	EndDate      time.Time `form:"end_date" binding:"omitempty"`
	ExerciseName string    `form:"exercise_name" binding:"omitempty"`
	Page         int       `form:"page,default=1" binding:"omitempty,min=1"`
	PerPage      int       `form:"per_page,default=20" binding:"omitempty,min=1,max=100"`
	OrderBy      string    `form:"order_by,default=performed_at" binding:"omitempty,oneof=performed_at created_at exercise_name"`
	Order        string    `form:"order,default=desc" binding:"omitempty,oneof=asc desc"`
}

// WorkoutStatsResponse represents workout statistics
type WorkoutStatsResponse struct {
	TotalWorkouts      int                           `json:"total_workouts"`
	TotalSets          int                           `json:"total_sets"`
	TotalReps          int                           `json:"total_reps"`
	TotalWeightLifted  float64                       `json:"total_weight_lifted"`
	WorkoutsByMuscle   map[string]int                `json:"workouts_by_muscle"`
	MostUsedExercises  []ExerciseCount               `json:"most_used_exercises"`
	RecentWorkouts     []WorkoutResponse             `json:"recent_workouts"`
	WeeklyProgress     []WeeklyProgress              `json:"weekly_progress"`
}

// ExerciseCount represents exercise usage count
type ExerciseCount struct {
	ExerciseName string `json:"exercise_name"`
	Count        int    `json:"count"`
}

// WeeklyProgress represents workout progress by week
type WeeklyProgress struct {
	Week          string  `json:"week"`
	WorkoutCount  int     `json:"workout_count"`
	TotalSets     int     `json:"total_sets"`
	TotalReps     int     `json:"total_reps"`
	TotalWeight   float64 `json:"total_weight"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string                 `json:"error"`
	Code    string                 `json:"code,omitempty"`
	Details map[string]interface{} `json:"details,omitempty"`
}