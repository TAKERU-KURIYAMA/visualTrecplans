-- Create muscle groups master table
CREATE TABLE muscle_groups (
    id SERIAL PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL,
    name_ja VARCHAR(100) NOT NULL,
    name_en VARCHAR(100) NOT NULL,
    category VARCHAR(50) NOT NULL, -- upper, lower, core, full_body
    color_code VARCHAR(7),
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create exercises master table
CREATE TABLE exercises (
    id SERIAL PRIMARY KEY,
    muscle_group_code VARCHAR(50) REFERENCES muscle_groups(code),
    name_ja VARCHAR(100) NOT NULL,
    name_en VARCHAR(100) NOT NULL,
    icon_name VARCHAR(50),
    is_custom BOOLEAN DEFAULT FALSE,
    created_by UUID REFERENCES users(id),
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create exercise icons master table
CREATE TABLE exercise_icons (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    svg_path TEXT,
    category VARCHAR(50),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX idx_muscle_groups_code ON muscle_groups(code);
CREATE INDEX idx_muscle_groups_category ON muscle_groups(category);
CREATE INDEX idx_exercises_muscle_group_code ON exercises(muscle_group_code);
CREATE INDEX idx_exercises_is_custom ON exercises(is_custom);
CREATE INDEX idx_exercises_created_by ON exercises(created_by);
CREATE INDEX idx_exercise_icons_name ON exercise_icons(name);

-- Update timestamp triggers
CREATE TRIGGER update_muscle_groups_updated_at BEFORE UPDATE ON muscle_groups
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_exercises_updated_at BEFORE UPDATE ON exercises
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_exercise_icons_updated_at BEFORE UPDATE ON exercise_icons
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();