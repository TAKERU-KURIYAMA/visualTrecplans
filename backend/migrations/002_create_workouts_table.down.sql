-- Drop triggers first
DROP TRIGGER IF EXISTS update_workouts_updated_at ON workouts;

-- Drop indexes
DROP INDEX IF EXISTS idx_workouts_user_id;
DROP INDEX IF EXISTS idx_workouts_performed_at;
DROP INDEX IF EXISTS idx_workouts_muscle_group;
DROP INDEX IF EXISTS idx_workouts_exercise_name;

-- Drop table
DROP TABLE IF EXISTS workouts;