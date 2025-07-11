import {
  ChartBarIcon,
  HeartIcon,
  ClockIcon,
  TrophyIcon,
} from '@heroicons/react/24/outline'

export const Dashboard = () => {
  const stats = [
    {
      icon: ChartBarIcon,
      name: 'トレーニング回数',
      value: '24',
      unit: '回',
      change: '+12%',
      changeType: 'positive',
    },
    {
      icon: HeartIcon,
      name: '今月の目標達成',
      value: '85',
      unit: '%',
      change: '+5%',
      changeType: 'positive',
    },
    {
      icon: ClockIcon,
      name: '平均トレーニング時間',
      value: '72',
      unit: '分',
      change: '+8%',
      changeType: 'positive',
    },
    {
      icon: TrophyIcon,
      name: '連続記録日数',
      value: '15',
      unit: '日',
      change: '+3日',
      changeType: 'positive',
    },
  ]

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
  ]

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="container py-8">
        <div className="mb-8">
          <h1 className="text-2xl font-bold text-gray-900">ダッシュボード</h1>
          <p className="mt-2 text-gray-600">
            あなたのフィットネス記録を確認しましょう
          </p>
        </div>

        {/* 統計カード */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
          {stats.map((stat) => (
            <div key={stat.name} className="card p-6">
              <div className="flex items-center">
                <div className="flex-shrink-0">
                  <stat.icon className="h-8 w-8 text-blue-600" />
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
          ))}
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
          {/* 最近のトレーニング */}
          <div className="card p-6">
            <h2 className="text-lg font-semibold text-gray-900 mb-4">
              最近のトレーニング
            </h2>
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
              <button className="w-full btn btn-outline">
                すべてのトレーニングを見る
              </button>
            </div>
          </div>

          {/* 進捗チャート */}
          <div className="card p-6">
            <h2 className="text-lg font-semibold text-gray-900 mb-4">
              今月の進捗
            </h2>
            <div className="h-64 flex items-center justify-center bg-gray-50 rounded-lg">
              <div className="text-center">
                <ChartBarIcon className="h-12 w-12 text-gray-400 mx-auto mb-2" />
                <p className="text-gray-600">チャート機能は実装予定です</p>
              </div>
            </div>
          </div>
        </div>

        {/* クイックアクション */}
        <div className="mt-8 grid grid-cols-1 md:grid-cols-3 gap-6">
          <button className="card p-6 hover:shadow-md transition-shadow">
            <div className="text-center">
              <div className="w-12 h-12 bg-blue-100 rounded-lg flex items-center justify-center mx-auto mb-3">
                <HeartIcon className="h-6 w-6 text-blue-600" />
              </div>
              <h3 className="text-lg font-medium text-gray-900">
                新しいトレーニング
              </h3>
              <p className="text-sm text-gray-600 mt-1">
                トレーニングを記録する
              </p>
            </div>
          </button>

          <button className="card p-6 hover:shadow-md transition-shadow">
            <div className="text-center">
              <div className="w-12 h-12 bg-green-100 rounded-lg flex items-center justify-center mx-auto mb-3">
                <ChartBarIcon className="h-6 w-6 text-green-600" />
              </div>
              <h3 className="text-lg font-medium text-gray-900">進捗を確認</h3>
              <p className="text-sm text-gray-600 mt-1">詳細な分析を表示</p>
            </div>
          </button>

          <button className="card p-6 hover:shadow-md transition-shadow">
            <div className="text-center">
              <div className="w-12 h-12 bg-purple-100 rounded-lg flex items-center justify-center mx-auto mb-3">
                <TrophyIcon className="h-6 w-6 text-purple-600" />
              </div>
              <h3 className="text-lg font-medium text-gray-900">目標設定</h3>
              <p className="text-sm text-gray-600 mt-1">新しい目標を設定</p>
            </div>
          </button>
        </div>
      </div>
    </div>
  )
}
