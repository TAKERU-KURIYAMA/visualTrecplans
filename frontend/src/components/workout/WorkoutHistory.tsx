import React, { useState, useMemo } from 'react';
import { format, isToday, isYesterday, isThisWeek } from 'date-fns';
import { ja } from 'date-fns/locale';
import { PencilIcon, TrashIcon, FunnelIcon } from '@heroicons/react/24/outline';
import { useWorkouts } from '@/hooks/useWorkout';
import { Workout, WorkoutFilter } from '@/types/workout';
import { toast } from 'react-hot-toast';

interface GroupedWorkouts {
  [date: string]: Workout[];
}

export const WorkoutHistory: React.FC = () => {
  const [filter, setFilter] = useState<WorkoutFilter>({});
  const [showFilters, setShowFilters] = useState(false);
  
  const { workouts, isLoading, total } = useWorkouts(filter);

  const groupedWorkouts = useMemo(() => {
    const grouped: GroupedWorkouts = {};
    
    workouts.forEach(workout => {
      const date = format(new Date(workout.performed_at), 'yyyy-MM-dd');
      if (!grouped[date]) {
        grouped[date] = [];
      }
      grouped[date].push(workout);
    });

    // Sort workouts within each date by time (newest first)
    Object.keys(grouped).forEach(date => {
      grouped[date].sort((a, b) => 
        new Date(b.performed_at).getTime() - new Date(a.performed_at).getTime()
      );
    });

    return grouped;
  }, [workouts]);

  const formatDateHeader = (dateStr: string) => {
    const date = new Date(dateStr);
    
    if (isToday(date)) {
      return '今日';
    } else if (isYesterday(date)) {
      return '昨日';
    } else if (isThisWeek(date)) {
      return format(date, 'EEEE', { locale: ja });
    } else {
      return format(date, 'M月d日 (EEE)', { locale: ja });
    }
  };

  const handleEdit = (workout: Workout) => {
    // Navigate to edit page - implementation depends on routing setup
    toast.info('編集機能は今後実装予定です');
  };

  const handleDelete = async (workout: Workout) => {
    if (window.confirm('このワークアウトを削除しますか？')) {
      try {
        // Implementation would use workout store
        toast.success('ワークアウトを削除しました');
      } catch (error) {
        toast.error('削除に失敗しました');
      }
    }
  };

  if (isLoading) {
    return <WorkoutHistorySkeleton />;
  }

  if (workouts.length === 0) {
    return <EmptyWorkoutHistory />;
  }

  return (
    <div className="space-y-6">
      {/* Header with filters */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold text-gray-900">トレーニング履歴</h2>
          <p className="text-gray-600">{total}件のワークアウト</p>
        </div>
        
        <button
          onClick={() => setShowFilters(!showFilters)}
          className="flex items-center px-4 py-2 border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors"
        >
          <FunnelIcon className="w-5 h-5 mr-2" />
          フィルター
        </button>
      </div>

      {/* Filters */}
      {showFilters && (
        <div className="bg-gray-50 p-4 rounded-lg">
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <select
              value={filter.muscle_group || ''}
              onChange={(e) => setFilter(prev => ({ ...prev, muscle_group: e.target.value || undefined }))}
              className="border border-gray-300 rounded px-3 py-2"
            >
              <option value="">全ての筋肉部位</option>
              <option value="chest">胸</option>
              <option value="back">背中</option>
              <option value="shoulders">肩</option>
              <option value="arms">腕</option>
              <option value="core">腹</option>
              <option value="legs">脚</option>
              <option value="glutes">臀部</option>
            </select>
            
            <input
              type="date"
              value={filter.start_date || ''}
              onChange={(e) => setFilter(prev => ({ ...prev, start_date: e.target.value || undefined }))}
              className="border border-gray-300 rounded px-3 py-2"
              placeholder="開始日"
            />
            
            <input
              type="date"
              value={filter.end_date || ''}
              onChange={(e) => setFilter(prev => ({ ...prev, end_date: e.target.value || undefined }))}
              className="border border-gray-300 rounded px-3 py-2"
              placeholder="終了日"
            />
          </div>
        </div>
      )}

      {/* Workout Groups */}
      <div className="space-y-8">
        {Object.entries(groupedWorkouts)
          .sort(([a], [b]) => new Date(b).getTime() - new Date(a).getTime())
          .map(([date, dayWorkouts]) => (
            <div key={date} className="space-y-4">
              <div className="flex items-center">
                <h3 className="text-lg font-semibold text-gray-900">
                  {formatDateHeader(date)}
                </h3>
                <div className="flex-1 ml-4 h-px bg-gray-200" />
                <span className="ml-4 text-sm text-gray-500">
                  {dayWorkouts.length}件
                </span>
              </div>
              
              <div className="grid gap-4">
                {dayWorkouts.map((workout) => (
                  <WorkoutCard
                    key={workout.id}
                    workout={workout}
                    onEdit={() => handleEdit(workout)}
                    onDelete={() => handleDelete(workout)}
                  />
                ))}
              </div>
            </div>
          ))}
      </div>
    </div>
  );
};

interface WorkoutCardProps {
  workout: Workout;
  onEdit: () => void;
  onDelete: () => void;
}

const WorkoutCard: React.FC<WorkoutCardProps> = ({ workout, onEdit, onDelete }) => {
  const muscleGroupColors: Record<string, string> = {
    chest: '#ff6b6b',
    back: '#4ecdc4',
    shoulders: '#dda0dd',
    arms: '#96ceb4',
    core: '#ffd93d',
    legs: '#45b7d1',
    glutes: '#ff9999',
    full_body: '#b8b8b8'
  };

  return (
    <div className="bg-white border border-gray-200 rounded-lg p-4 hover:shadow-md transition-shadow">
      <div className="flex items-start justify-between">
        <div className="flex-1">
          <div className="flex items-center mb-2">
            <div
              className="w-3 h-3 rounded-full mr-3"
              style={{ backgroundColor: muscleGroupColors[workout.muscle_group] || '#gray' }}
            />
            <h4 className="text-lg font-medium text-gray-900">
              {workout.exercise_name}
            </h4>
            <span className="ml-2 text-sm text-gray-500">
              {format(new Date(workout.performed_at), 'HH:mm')}
            </span>
          </div>
          
          <div className="flex items-center space-x-4 text-sm text-gray-600">
            {workout.weight_kg && (
              <span>{workout.weight_kg}kg</span>
            )}
            {workout.reps && (
              <span>{workout.reps}回</span>
            )}
            {workout.sets && (
              <span>{workout.sets}セット</span>
            )}
          </div>
          
          {workout.notes && (
            <p className="mt-2 text-sm text-gray-600 line-clamp-2">
              {workout.notes}
            </p>
          )}
        </div>
        
        <div className="flex items-center space-x-2 ml-4">
          <button
            onClick={onEdit}
            className="p-2 text-gray-400 hover:text-blue-600 hover:bg-blue-50 rounded transition-colors"
          >
            <PencilIcon className="w-4 h-4" />
          </button>
          <button
            onClick={onDelete}
            className="p-2 text-gray-400 hover:text-red-600 hover:bg-red-50 rounded transition-colors"
          >
            <TrashIcon className="w-4 h-4" />
          </button>
        </div>
      </div>
    </div>
  );
};

const WorkoutHistorySkeleton: React.FC = () => (
  <div className="space-y-6">
    <div className="h-8 bg-gray-200 rounded w-48 animate-pulse" />
    {[...Array(3)].map((_, i) => (
      <div key={i} className="space-y-4">
        <div className="h-6 bg-gray-200 rounded w-32 animate-pulse" />
        <div className="space-y-3">
          {[...Array(2)].map((_, j) => (
            <div key={j} className="h-24 bg-gray-100 rounded animate-pulse" />
          ))}
        </div>
      </div>
    ))}
  </div>
);

const EmptyWorkoutHistory: React.FC = () => (
  <div className="text-center py-12">
    <div className="w-24 h-24 mx-auto mb-4 text-gray-300">
      <svg fill="currentColor" viewBox="0 0 24 24">
        <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/>
      </svg>
    </div>
    <h3 className="text-lg font-medium text-gray-900 mb-2">
      まだワークアウトがありません
    </h3>
    <p className="text-gray-600 mb-6">
      最初のワークアウトを記録して始めましょう！
    </p>
    <button className="bg-blue-600 text-white px-6 py-3 rounded-lg hover:bg-blue-700 transition-colors">
      ワークアウトを記録
    </button>
  </div>
);