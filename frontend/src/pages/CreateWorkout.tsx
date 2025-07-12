import React from 'react';
import { useNavigate } from 'react-router-dom';
import { ArrowLeftIcon } from '@heroicons/react/24/outline';
import { WorkoutForm } from '@/components/forms/WorkoutForm';
import { useWorkoutForm } from '@/hooks/useWorkoutForm';
import { toast } from 'react-hot-toast';

export const CreateWorkout: React.FC = () => {
  const navigate = useNavigate();
  
  const { handleSubmit, isLoading, error } = useWorkoutForm({
    onSuccess: () => {
      toast.success('ワークアウトを保存しました');
      navigate('/dashboard');
    },
    onError: (error) => {
      toast.error(error.message);
    }
  });

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <div className="bg-white shadow-sm border-b">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex items-center justify-between h-16">
            <div className="flex items-center">
              <button
                onClick={() => navigate(-1)}
                className="flex items-center text-gray-600 hover:text-gray-900 transition-colors"
              >
                <ArrowLeftIcon className="w-5 h-5 mr-2" />
                戻る
              </button>
            </div>
            <h1 className="text-xl font-semibold text-gray-900">
              ワークアウト記録
            </h1>
            <div className="w-16" />
          </div>
        </div>
      </div>

      {/* Main Content */}
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="bg-white rounded-lg shadow p-6">
          <div className="mb-6">
            <h2 className="text-2xl font-bold text-gray-900 mb-2">
              新しいワークアウトを記録
            </h2>
            <p className="text-gray-600">
              トレーニングの詳細を入力してください。
            </p>
          </div>

          {/* Error Display */}
          {error && (
            <div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-lg">
              <div className="flex">
                <div className="flex-shrink-0">
                  <svg className="h-5 w-5 text-red-400" viewBox="0 0 20 20" fill="currentColor">
                    <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clipRule="evenodd" />
                  </svg>
                </div>
                <div className="ml-3">
                  <h3 className="text-sm font-medium text-red-800">
                    エラーが発生しました
                  </h3>
                  <div className="mt-2 text-sm text-red-700">
                    {error}
                  </div>
                </div>
              </div>
            </div>
          )}

          {/* Workout Form */}
          <WorkoutForm
            onSubmit={handleSubmit}
            isLoading={isLoading}
            submitText="ワークアウトを保存"
            showCancel
            onCancel={() => navigate(-1)}
          />
        </div>
      </div>
    </div>
  );
};