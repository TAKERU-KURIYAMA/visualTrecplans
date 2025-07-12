-- Create workouts table for training records
CREATE TABLE workouts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    muscle_group VARCHAR(50) NOT NULL,
    exercise_name VARCHAR(100) NOT NULL,
    exercise_icon VARCHAR(50),
    weight_kg DECIMAL(5,2),
    reps INTEGER CHECK (reps > 0),
    sets INTEGER CHECK (sets > 0),
    notes TEXT,
    performed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Indexes for performance optimization
CREATE INDEX idx_workouts_user_id ON workouts(user_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_workouts_performed_at ON workouts(performed_at DESC) WHERE deleted_at IS NULL;
CREATE INDEX idx_workouts_muscle_group ON workouts(muscle_group) WHERE deleted_at IS NULL;
CREATE INDEX idx_workouts_exercise_name ON workouts(exercise_name) WHERE deleted_at IS NULL;

-- Update timestamp trigger
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_workouts_updated_at BEFORE UPDATE ON workouts
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();