import React from 'react';
import { useNavigate } from 'react-router-dom';
import { 
  User,
  Settings,
  Activity,
  BarChart3,
  LogOut,
  Calendar,
  Target,
  TrendingUp
} from 'lucide-react';
import { useAuth, useAuthActions } from '../stores/authStore';

const Dashboard: React.FC = () => {
  const { user } = useAuth();
  const { logout } = useAuthActions();
  const navigate = useNavigate();

  const handleLogout = async () => {
    try {
      await logout();
      navigate('/login');
    } catch (error) {
      console.error('ログアウトエラー:', error);
    }
  };

  const stats = [
    {
      icon: Activity,
      name: 'トレーニング回数',
      value: '24',
      unit: '回',
      change: '+12%',
      changeType: 'positive' as const,
    },
    {
      icon: Target,
      name: '今月の目標達成',
      value: '85',
      unit: '%',
      change: '+5%',
      changeType: 'positive' as const,
    },
    {
      icon: Calendar,
      name: '平均トレーニング時間',
      value: '72',
      unit: '分',
      change: '+8%',
      changeType: 'positive' as const,
    },
    {
      icon: TrendingUp,
      name: '連続記録日数',
      value: '15',
      unit: '日',
      change: '+3日',
      changeType: 'positive' as const,
    },
  ];

  const recentWorkouts = [
    {
      date: '2025-01-10',
      type: 'チェスト',
      duration: '65分',
      exercises: 4,
    },
    {
      date: '2025-01-08',
      type: 'バック',
      duration: '70分',
      exercises: 5,
    },
    {
      date: '2025-01-06',
      type: 'レッグ',
      duration: '80分',
      exercises: 6,
    },
  ];

  const getUserDisplayName = () => {
    if (user?.first_name && user?.last_name) {
      return `${user.last_name} ${user.first_name}`;
    }
    if (user?.first_name) {
      return user.first_name;
    }
    if (user?.last_name) {
      return user.last_name;
    }
    return user?.email || 'ユーザー';
  };

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <div className="bg-white shadow">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center py-4">
            <div className="flex items-center">
              <div className="w-10 h-10 bg-primary-600 rounded-lg flex items-center justify-center">
                <User className="w-6 h-6 text-white" />
              </div>
              <div className="ml-4">
                <h1 className="text-xl font-semibold text-gray-900">
                  ようこそ、{getUserDisplayName()}さん
                </h1>
                <p className="text-sm text-gray-600">{user?.email}</p>
              </div>
            </div>
            <div className="flex items-center space-x-4">
              <button className="p-2 text-gray-400 hover:text-gray-600 rounded-lg hover:bg-gray-100">
                <Settings className="w-5 h-5" />
              </button>
              <button
                onClick={handleLogout}
                className="flex items-center px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-lg hover:bg-gray-50"
              >
                <LogOut className="w-4 h-4 mr-2" />
                ログアウト
              </button>
            </div>
          </div>
        </div>
      </div>

      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="mb-8">
          <h2 className="text-2xl font-bold text-gray-900">ダッシュボード</h2>
          <p className="mt-2 text-gray-600">
            あなたのフィットネス記録を確認しましょう
          </p>
        </div>

        {/* Stats Cards */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
          {stats.map((stat) => (
            <div key={stat.name} className="bg-white overflow-hidden shadow rounded-lg">
              <div className="p-6">
                <div className="flex items-center">
                  <div className="flex-shrink-0">
                    <stat.icon className="h-8 w-8 text-primary-600" />
                  </div>
                  <div className="ml-4">
                    <p className="text-sm font-medium text-gray-600">
                      {stat.name}
                    </p>
                    <p className="text-2xl font-bold text-gray-900">
                      {stat.value}
                      <span className="text-sm font-normal text-gray-500 ml-1">
                        {stat.unit}
                      </span>
                    </p>
                    <p
                      className={`text-sm ${
                        stat.changeType === 'positive'
                          ? 'text-green-600'
                          : 'text-red-600'
                      }`}
                    >
                      {stat.change} 先月比
                    </p>
                  </div>
                </div>
              </div>
            </div>
          ))}
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
          {/* Recent Workouts */}
          <div className="bg-white overflow-hidden shadow rounded-lg">
            <div className="p-6">
              <h3 className="text-lg font-semibold text-gray-900 mb-4">
                最近のトレーニング
              </h3>
              <div className="space-y-4">
                {recentWorkouts.map((workout, index) => (
                  <div
                    key={index}
                    className="flex items-center justify-between p-4 bg-gray-50 rounded-lg"
                  >
                    <div>
                      <p className="font-medium text-gray-900">{workout.type}</p>
                      <p className="text-sm text-gray-600">{workout.date}</p>
                    </div>
                    <div className="text-right">
                      <p className="text-sm font-medium text-gray-900">
                        {workout.duration}
                      </p>
                      <p className="text-sm text-gray-600">
                        {workout.exercises} 種目
                      </p>
                    </div>
                  </div>
                ))}
              </div>
              <div className="mt-4">
                <button className="w-full px-4 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-700 bg-white hover:bg-gray-50">
                  すべてのトレーニングを見る
                </button>
              </div>
            </div>
          </div>

          {/* Progress Chart Placeholder */}
          <div className="bg-white overflow-hidden shadow rounded-lg">
            <div className="p-6">
              <h3 className="text-lg font-semibold text-gray-900 mb-4">
                今月の進捗
              </h3>
              <div className="h-64 flex items-center justify-center bg-gray-50 rounded-lg">
                <div className="text-center">
                  <BarChart3 className="h-12 w-12 text-gray-400 mx-auto mb-2" />
                  <p className="text-gray-600">チャート機能は実装予定です</p>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Quick Actions */}
        <div className="mt-8 grid grid-cols-1 md:grid-cols-3 gap-6">
          <button className="bg-white overflow-hidden shadow rounded-lg hover:shadow-md transition-shadow">
            <div className="p-6 text-center">
              <div className="w-12 h-12 bg-primary-100 rounded-lg flex items-center justify-center mx-auto mb-3">
                <Activity className="h-6 w-6 text-primary-600" />
              </div>
              <h3 className="text-lg font-medium text-gray-900">
                新しいトレーニング
              </h3>
              <p className="text-sm text-gray-600 mt-1">
                トレーニングを記録する
              </p>
            </div>
          </button>

          <button className="bg-white overflow-hidden shadow rounded-lg hover:shadow-md transition-shadow">
            <div className="p-6 text-center">
              <div className="w-12 h-12 bg-green-100 rounded-lg flex items-center justify-center mx-auto mb-3">
                <BarChart3 className="h-6 w-6 text-green-600" />
              </div>
              <h3 className="text-lg font-medium text-gray-900">進捗を確認</h3>
              <p className="text-sm text-gray-600 mt-1">詳細な分析を表示</p>
            </div>
          </button>

          <button className="bg-white overflow-hidden shadow rounded-lg hover:shadow-md transition-shadow">
            <div className="p-6 text-center">
              <div className="w-12 h-12 bg-purple-100 rounded-lg flex items-center justify-center mx-auto mb-3">
                <Target className="h-6 w-6 text-purple-600" />
              </div>
              <h3 className="text-lg font-medium text-gray-900">目標設定</h3>
              <p className="text-sm text-gray-600 mt-1">新しい目標を設定</p>
            </div>
          </button>
        </div>
      </div>
    </div>
  );
};

export default Dashboard;
