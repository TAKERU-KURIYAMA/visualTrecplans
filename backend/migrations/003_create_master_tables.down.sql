-- Drop triggers
DROP TRIGGER IF EXISTS update_muscle_groups_updated_at ON muscle_groups;
DROP TRIGGER IF EXISTS update_exercises_updated_at ON exercises;
DROP TRIGGER IF EXISTS update_exercise_icons_updated_at ON exercise_icons;

-- Drop indexes
DROP INDEX IF EXISTS idx_muscle_groups_code;
DROP INDEX IF EXISTS idx_muscle_groups_category;
DROP INDEX IF EXISTS idx_exercises_muscle_group_code;
DROP INDEX IF EXISTS idx_exercises_is_custom;
DROP INDEX IF EXISTS idx_exercises_created_by;
DROP INDEX IF EXISTS idx_exercise_icons_name;

-- Drop tables (order matters due to foreign keys)
DROP TABLE IF EXISTS exercises;
DROP TABLE IF EXISTS exercise_icons;
DROP TABLE IF EXISTS muscle_groups;