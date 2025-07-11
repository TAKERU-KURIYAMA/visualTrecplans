import { Link } from 'react-router-dom'
import {
  HeartIcon,
  ChartBarIcon,
  UserGroupIcon,
  GlobeAltIcon,
} from '@heroicons/react/24/outline'

export const Home = () => {
  const features = [
    {
      icon: HeartIcon,
      title: '人体図UI',
      description:
        'インタラクティブな人体図で直感的にトレーニング部位を選択・記録できます。',
    },
    {
      icon: ChartBarIcon,
      title: '進捗可視化',
      description:
        'トレーニング記録とサプリメント摂取状況を分かりやすいグラフで可視化します。',
    },
    {
      icon: UserGroupIcon,
      title: 'ソーシャル機能',
      description:
        'フィットネス仲間と記録を共有し、モチベーションを維持できます。',
    },
    {
      icon: GlobeAltIcon,
      title: '多言語対応',
      description:
        '日本語・英語・スペイン語・フランス語に対応した国際的なアプリです。',
    },
  ]

  return (
    <div className="bg-white">
      {/* ヒーローセクション */}
      <div className="relative isolate px-6 pt-14 lg:px-8">
        <div className="absolute inset-x-0 -top-40 -z-10 transform-gpu overflow-hidden blur-3xl sm:-top-80">
          <div className="relative left-[calc(50%-11rem)] aspect-[1155/678] w-[36.125rem] -translate-x-1/2 rotate-[30deg] bg-gradient-to-tr from-[#ff80b5] to-[#9089fc] opacity-30 sm:left-[calc(50%-30rem)] sm:w-[72.1875rem]"></div>
        </div>

        <div className="mx-auto max-w-2xl py-32 sm:py-48 lg:py-56">
          <div className="text-center">
            <h1 className="text-4xl font-bold tracking-tight text-gray-900 sm:text-6xl">
              ビジュアル重視の
              <span className="text-blue-600">トレーニング管理</span>
            </h1>
            <p className="mt-6 text-lg leading-8 text-gray-600">
              人体図UIと直感的なデザインで、言語に依存しないフィットネス記録・管理アプリケーション。
              トレーニングとサプリメントの両方を効率的に管理できます。
            </p>
            <div className="mt-10 flex items-center justify-center gap-x-6">
              <Link
                to="/register"
                className="rounded-md bg-blue-600 px-3.5 py-2.5 text-sm font-semibold text-white shadow-sm hover:bg-blue-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600"
              >
                今すぐ始める
              </Link>
              <Link
                to="/login"
                className="text-sm font-semibold leading-6 text-gray-900"
              >
                ログイン <span aria-hidden="true">→</span>
              </Link>
            </div>
          </div>
        </div>
      </div>

      {/* 機能セクション */}
      <div className="py-24 sm:py-32">
        <div className="mx-auto max-w-7xl px-6 lg:px-8">
          <div className="mx-auto max-w-2xl lg:text-center">
            <h2 className="text-base font-semibold leading-7 text-blue-600">
              主要機能
            </h2>
            <p className="mt-2 text-3xl font-bold tracking-tight text-gray-900 sm:text-4xl">
              フィットネスライフを革新する機能
            </p>
            <p className="mt-6 text-lg leading-8 text-gray-600">
              VisualTrecplansは、従来のフィットネスアプリとは違う、
              ビジュアルとユーザビリティに特化した新しいアプローチを提供します。
            </p>
          </div>

          <div className="mx-auto mt-16 max-w-2xl sm:mt-20 lg:mt-24 lg:max-w-4xl">
            <dl className="grid max-w-xl grid-cols-1 gap-x-8 gap-y-10 lg:max-w-none lg:grid-cols-2 lg:gap-y-16">
              {features.map((feature) => (
                <div key={feature.title} className="relative pl-16">
                  <dt className="text-base font-semibold leading-7 text-gray-900">
                    <div className="absolute left-0 top-0 flex h-10 w-10 items-center justify-center rounded-lg bg-blue-600">
                      <feature.icon
                        className="h-6 w-6 text-white"
                        aria-hidden="true"
                      />
                    </div>
                    {feature.title}
                  </dt>
                  <dd className="mt-2 text-base leading-7 text-gray-600">
                    {feature.description}
                  </dd>
                </div>
              ))}
            </dl>
          </div>
        </div>
      </div>
    </div>
  )
}
