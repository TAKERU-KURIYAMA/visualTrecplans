-- Exercises master data
INSERT INTO exercises (muscle_group_code, name_ja, name_en, icon_name, sort_order) VALUES

-- Chest exercises
('chest', 'ベンチプレス', 'Bench Press', 'bench_press', 1),
('chest', 'インクラインベンチプレス', 'Incline Bench Press', 'bench_press', 2),
('chest', 'ダンベルフライ', 'Dumbbell Fly', 'dumbbell_fly', 3),
('chest', 'プッシュアップ', 'Push Up', 'push_up', 4),
('chest', 'チェストプレス', 'Chest Press', 'chest_press', 5),
('chest', 'ディップス', 'Dips', 'tricep_dip', 6),

-- Back exercises
('back', 'プルアップ', 'Pull Up', 'pull_up', 1),
('back', 'ラットプルダウン', 'Lat Pulldown', 'lat_pulldown', 2),
('back', 'ローイング', 'Rowing', 'rowing', 3),
('back', 'デッドリフト', 'Deadlift', 'deadlift', 4),
('back', 'ワンハンドロー', 'One Hand Row', 'rowing', 5),
('back', 'シーテッドロー', 'Seated Row', 'rowing', 6),

-- Shoulder exercises  
('shoulders', 'ショルダープレス', 'Shoulder Press', 'shoulder_press', 1),
('shoulders', 'ラテラルレイズ', 'Lateral Raise', 'lateral_raise', 2),
('shoulders', 'フロントレイズ', 'Front Raise', 'front_raise', 3),
('shoulders', 'リアデルトフライ', 'Rear Delt Fly', 'rear_delt_fly', 4),
('shoulders', 'アップライトロー', 'Upright Row', 'rowing', 5),

-- Biceps exercises
('biceps', 'バイセップカール', 'Bicep Curl', 'bicep_curl', 1),
('biceps', 'ハンマーカール', 'Hammer Curl', 'hammer_curl', 2),
('biceps', 'プリーチャーカール', 'Preacher Curl', 'bicep_curl', 3),
('biceps', '21s', '21s', 'bicep_curl', 4),

-- Triceps exercises
('triceps', 'トライセップディップ', 'Tricep Dip', 'tricep_dip', 1),
('triceps', 'トライセップエクステンション', 'Tricep Extension', 'tricep_extension', 2),
('triceps', 'オーバーヘッドエクステンション', 'Overhead Extension', 'tricep_extension', 3),
('triceps', 'クローズグリップベンチプレス', 'Close Grip Bench Press', 'bench_press', 4),

-- Forearms exercises
('forearms', 'リストカール', 'Wrist Curl', 'bicep_curl', 1),
('forearms', 'リバースリストカール', 'Reverse Wrist Curl', 'bicep_curl', 2),
('forearms', 'ファーマーズウォーク', 'Farmers Walk', 'dumbbell', 3),

-- Abs exercises
('abs', 'シットアップ', 'Sit Up', 'sit_up', 1),
('abs', 'クランチ', 'Crunch', 'sit_up', 2),
('abs', 'プランク', 'Plank', 'plank', 3),
('abs', 'レッグレイズ', 'Leg Raise', 'leg_raise', 4),
('abs', 'マウンテンクライマー', 'Mountain Climber', 'plank', 5),

-- Obliques exercises
('obliques', 'ロシアンツイスト', 'Russian Twist', 'russian_twist', 1),
('obliques', 'サイドプランク', 'Side Plank', 'plank', 2),
('obliques', 'バイシクルクランチ', 'Bicycle Crunch', 'sit_up', 3),

-- Lower back exercises
('lower_back', 'ハイパーエクステンション', 'Hyperextension', 'deadlift', 1),
('lower_back', 'グッドモーニング', 'Good Morning', 'deadlift', 2),

-- Quadriceps exercises
('quadriceps', 'スクワット', 'Squat', 'squat', 1),
('quadriceps', 'レッグプレス', 'Leg Press', 'leg_press', 2),
('quadriceps', 'ランジ', 'Lunge', 'lunge', 3),
('quadriceps', 'レッグエクステンション', 'Leg Extension', 'machine', 4),
('quadriceps', 'ブルガリアンスクワット', 'Bulgarian Squat', 'lunge', 5),

-- Hamstrings exercises
('hamstrings', 'レッグカール', 'Leg Curl', 'machine', 1),
('hamstrings', 'ルーマニアンデッドリフト', 'Romanian Deadlift', 'deadlift', 2),
('hamstrings', 'スティッフレッグデッドリフト', 'Stiff Leg Deadlift', 'deadlift', 3),

-- Glutes exercises
('glutes', 'ヒップスラスト', 'Hip Thrust', 'machine', 1),
('glutes', 'グルートブリッジ', 'Glute Bridge', 'bodyweight', 2),
('glutes', 'サイドステップ', 'Side Step', 'bodyweight', 3),

-- Calves exercises
('calves', 'カーフレイズ', 'Calf Raise', 'calf_raise', 1),
('calves', 'シーテッドカーフレイズ', 'Seated Calf Raise', 'calf_raise', 2),

-- Full body exercises
('full_body', 'バーピー', 'Burpee', 'bodyweight', 1),
('full_body', 'スラスター', 'Thruster', 'shoulder_press', 2),
('full_body', 'クリーン&ジャーク', 'Clean & Jerk', 'deadlift', 3),

-- Cardio exercises
('cardio', 'ランニング', 'Running', 'bodyweight', 1),
('cardio', 'サイクリング', 'Cycling', 'machine', 2),
('cardio', 'ロウイング', 'Rowing', 'rowing', 3),
('cardio', '縄跳び', 'Jump Rope', 'bodyweight', 4);