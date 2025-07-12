package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Workout represents a training record
type Workout struct {
	ID           uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	UserID       uuid.UUID      `json:"user_id" gorm:"type:uuid;not null"`
	User         User           `json:"user,omitempty" gorm:"constraint:OnDelete:CASCADE"`
	MuscleGroup  string         `json:"muscle_group" gorm:"type:varchar(50);not null"`
	ExerciseName string         `json:"exercise_name" gorm:"type:varchar(100);not null"`
	ExerciseIcon string         `json:"exercise_icon" gorm:"type:varchar(50)"`
	WeightKg     *float64       `json:"weight_kg" gorm:"type:decimal(5,2)"`
	Reps         *int           `json:"reps" gorm:"check:reps > 0"`
	Sets         *int           `json:"sets" gorm:"check:sets > 0"`
	Notes        string         `json:"notes" gorm:"type:text"`
	PerformedAt  time.Time      `json:"performed_at" gorm:"not null"`
	CreatedAt    time.Time      `json:"created_at" gorm:"not null"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"not null"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

// WorkoutResponse represents the response structure for API
type WorkoutResponse struct {
	ID           string     `json:"id"`
	UserID       string     `json:"user_id"`
	MuscleGroup  string     `json:"muscle_group"`
	ExerciseName string     `json:"exercise_name"`
	ExerciseIcon string     `json:"exercise_icon"`
	WeightKg     *float64   `json:"weight_kg"`
	Reps         *int       `json:"reps"`
	Sets         *int       `json:"sets"`
	Notes        string     `json:"notes"`
	PerformedAt  time.Time  `json:"performed_at"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// ToResponse converts Workout to WorkoutResponse
func (w *Workout) ToResponse() *WorkoutResponse {
	return &WorkoutResponse{
		ID:           w.ID.String(),
		UserID:       w.UserID.String(),
		MuscleGroup:  w.MuscleGroup,
		ExerciseName: w.ExerciseName,
		ExerciseIcon: w.ExerciseIcon,
		WeightKg:     w.WeightKg,
		Reps:         w.Reps,
		Sets:         w.Sets,
		Notes:        w.Notes,
		PerformedAt:  w.PerformedAt,
		CreatedAt:    w.CreatedAt,
		UpdatedAt:    w.UpdatedAt,
	}
}

// WorkoutStats represents workout statistics
type WorkoutStats struct {
	TotalWorkouts   int                    `json:"total_workouts"`
	TotalVolume     float64                `json:"total_volume"` // sets * reps * weight
	MuscleGroupStats map[string]int        `json:"muscle_group_stats"`
	WeeklyStats     []WeeklyWorkoutStats   `json:"weekly_stats"`
	MonthlyStats    []MonthlyWorkoutStats  `json:"monthly_stats"`
}

// WeeklyWorkoutStats represents weekly statistics
type WeeklyWorkoutStats struct {
	Week        string  `json:"week"`        // YYYY-WW format
	Workouts    int     `json:"workouts"`
	TotalVolume float64 `json:"total_volume"`
}

// MonthlyWorkoutStats represents monthly statistics
type MonthlyWorkoutStats struct {
	Month       string  `json:"month"`       // YYYY-MM format
	Workouts    int     `json:"workouts"`
	TotalVolume float64 `json:"total_volume"`
}

// MuscleGroup represents the muscle group master data
type MuscleGroup struct {
	ID        int       `json:"id" gorm:"primary_key"`
	Code      string    `json:"code" gorm:"type:varchar(50);unique;not null"`
	NameJa    string    `json:"name_ja" gorm:"type:varchar(100);not null"`
	NameEn    string    `json:"name_en" gorm:"type:varchar(100);not null"`
	Category  string    `json:"category" gorm:"type:varchar(50);not null"`
	ColorCode string    `json:"color_code" gorm:"type:varchar(7)"`
	SortOrder int       `json:"sort_order" gorm:"default:0"`
	CreatedAt time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null"`
}

// Exercise represents the exercise master data
type Exercise struct {
	ID              int       `json:"id" gorm:"primary_key"`
	MuscleGroupCode string    `json:"muscle_group_code" gorm:"type:varchar(50)"`
	MuscleGroup     MuscleGroup `json:"muscle_group,omitempty" gorm:"foreignKey:MuscleGroupCode;references:Code"`
	NameJa          string    `json:"name_ja" gorm:"type:varchar(100);not null"`
	NameEn          string    `json:"name_en" gorm:"type:varchar(100);not null"`
	IconName        string    `json:"icon_name" gorm:"type:varchar(50)"`
	IsCustom        bool      `json:"is_custom" gorm:"default:false"`
	CreatedBy       *uuid.UUID `json:"created_by" gorm:"type:uuid"`
	Creator         *User     `json:"creator,omitempty" gorm:"constraint:OnDelete:SET NULL"`
	SortOrder       int       `json:"sort_order" gorm:"default:0"`
	CreatedAt       time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"not null"`
}

// ExerciseIcon represents the exercise icon master data
type ExerciseIcon struct {
	ID        int       `json:"id" gorm:"primary_key"`
	Name      string    `json:"name" gorm:"type:varchar(50);unique;not null"`
	SvgPath   string    `json:"svg_path" gorm:"type:text"`
	Category  string    `json:"category" gorm:"type:varchar(50)"`
	CreatedAt time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null"`
}

// TableName specifies the table name for GORM
func (Workout) TableName() string {
	return "workouts"
}

func (MuscleGroup) TableName() string {
	return "muscle_groups"
}

func (Exercise) TableName() string {
	return "exercises"
}

func (ExerciseIcon) TableName() string {
	return "exercise_icons"
}

// BeforeCreate sets the ID if not already set
func (w *Workout) BeforeCreate(tx *gorm.DB) error {
	if w.ID == uuid.Nil {
		w.ID = uuid.New()
	}
	return nil
}

// CalculateVolume calculates the volume (sets * reps * weight) for the workout
func (w *Workout) CalculateVolume() float64 {
	if w.Sets == nil || w.Reps == nil || w.WeightKg == nil {
		return 0
	}
	return float64(*w.Sets) * float64(*w.Reps) * *w.WeightKg
}

// IsCardio checks if the workout is cardio based
func (w *Workout) IsCardio() bool {
	return w.MuscleGroup == "cardio"
}

// HasWeight checks if the workout has weight data
func (w *Workout) HasWeight() bool {
	return w.WeightKg != nil && *w.WeightKg > 0
}