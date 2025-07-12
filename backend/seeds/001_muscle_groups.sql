-- Muscle groups master data
INSERT INTO muscle_groups (code, name_ja, name_en, category, color_code, sort_order) VALUES
-- Upper body
('chest', '胸', 'Chest', 'upper', '#ff6b6b', 1),
('back', '背中', 'Back', 'upper', '#4ecdc4', 2),
('shoulders', '肩', 'Shoulders', 'upper', '#dda0dd', 3),
('biceps', '二頭筋', 'Biceps', 'upper', '#96ceb4', 4),
('triceps', '三頭筋', 'Triceps', 'upper', '#74c69d', 5),
('forearms', '前腕', 'Forearms', 'upper', '#52b788', 6),

-- Core
('abs', '腹筋', 'Abs', 'core', '#ffd93d', 7),
('obliques', '腹斜筋', 'Obliques', 'core', '#ffd60a', 8),
('lower_back', '腰', 'Lower Back', 'core', '#fdc500', 9),

-- Lower body
('quadriceps', '大腿四頭筋', 'Quadriceps', 'lower', '#45b7d1', 10),
('hamstrings', 'ハムストリング', 'Hamstrings', 'lower', '#0077b6', 11),
('glutes', '臀筋', 'Glutes', 'lower', '#ff9999', 12),
('calves', 'ふくらはぎ', 'Calves', 'lower', '#00b4d8', 13),

-- Full body
('full_body', '全身', 'Full Body', 'full_body', '#b8b8b8', 14),
('cardio', '有酸素', 'Cardio', 'cardio', '#ff7b7b', 15);

-- Exercise icons master data
INSERT INTO exercise_icons (name, category) VALUES
-- Chest exercises
('bench_press', 'chest'),
('push_up', 'chest'),
('dumbbell_fly', 'chest'),
('chest_press', 'chest'),

-- Back exercises
('pull_up', 'back'),
('lat_pulldown', 'back'),
('rowing', 'back'),
('deadlift', 'back'),

-- Shoulder exercises
('shoulder_press', 'shoulders'),
('lateral_raise', 'shoulders'),
('front_raise', 'shoulders'),
('rear_delt_fly', 'shoulders'),

-- Arms exercises
('bicep_curl', 'biceps'),
('hammer_curl', 'biceps'),
('tricep_dip', 'triceps'),
('tricep_extension', 'triceps'),

-- Core exercises
('sit_up', 'abs'),
('plank', 'abs'),
('russian_twist', 'obliques'),
('leg_raise', 'abs'),

-- Leg exercises
('squat', 'quadriceps'),
('lunge', 'quadriceps'),
('leg_press', 'quadriceps'),
('calf_raise', 'calves'),

-- General
('dumbbell', 'general'),
('barbell', 'general'),
('machine', 'general'),
('bodyweight', 'general');