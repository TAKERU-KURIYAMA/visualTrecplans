package master

// MuscleGroupResponse represents a muscle group in the API response
type MuscleGroupResponse struct {
	Code      string `json:"code"`
	NameJa    string `json:"name_ja"`
	NameEn    string `json:"name_en"`
	Category  string `json:"category"`
	ColorCode string `json:"color_code"`
	SortOrder int    `json:"sort_order"`
}

// MuscleGroupListResponse represents a list of muscle groups
type MuscleGroupListResponse struct {
	Data []MuscleGroupResponse `json:"data"`
}

// ExerciseResponse represents an exercise in the API response
type ExerciseResponse struct {
	ID              int    `json:"id"`
	MuscleGroupCode string `json:"muscle_group_code"`
	NameJa          string `json:"name_ja"`
	NameEn          string `json:"name_en"`
	IconName        string `json:"icon_name"`
	IsCustom        bool   `json:"is_custom"`
	SortOrder       int    `json:"sort_order"`
}

// ExerciseListResponse represents a list of exercises
type ExerciseListResponse struct {
	Data []ExerciseResponse `json:"data"`
}

// ExerciseIconResponse represents an exercise icon in the API response
type ExerciseIconResponse struct {
	Name     string `json:"name"`
	SVGPath  string `json:"svg_path"`
	Category string `json:"category"`
}

// ExerciseIconListResponse represents a list of exercise icons
type ExerciseIconListResponse struct {
	Data []ExerciseIconResponse `json:"data"`
}

// CreateCustomExerciseRequest represents the request to create a custom exercise
type CreateCustomExerciseRequest struct {
	Name            string `json:"name" binding:"required,min=1,max=100"`
	MuscleGroupCode string `json:"muscle_group_code" binding:"required,muscle_group"`
	IconName        string `json:"icon_name" binding:"omitempty,max=50"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string                 `json:"error"`
	Code    string                 `json:"code,omitempty"`
	Details map[string]interface{} `json:"details,omitempty"`
}