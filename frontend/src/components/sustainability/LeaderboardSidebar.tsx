import React from 'react'
import { useTranslation } from '@/i18n/useTranslation'

export const LeaderboardSidebar: React.FC = () => {
  const { t } = useTranslation()

  return (
    <div className="w-80 bg-white/80 backdrop-blur-sm rounded-lg shadow-lg p-6 sticky top-4">
      <div className="flex items-center gap-2 mb-6">
        <h2 className="text-2xl font-bold text-gray-900">{t('aiSidebar.title')}</h2>
      </div>

      <div className="space-y-4">
        <div className="p-4 bg-gradient-to-r from-green-50 to-emerald-50 rounded-lg border border-green-200">
          <div className="flex items-center gap-2 mb-2">
            <p className="text-sm font-semibold text-gray-900">
              {t('aiSidebar.agentsTitle')}
            </p>
          </div>
          <p className="text-xs text-gray-600">
            {t('aiSidebar.agentsSubtitle')}
          </p>
        </div>

        <div className="grid grid-cols-1 gap-3">
          <div className="p-3 bg-blue-50 rounded-lg">
            <div className="text-xs text-gray-600">{t('aiSidebar.listingTitle')}</div>
            <div className="text-lg font-bold text-blue-600">{t('aiSidebar.listingDesc')}</div>
          </div>
          <div className="p-3 bg-purple-50 rounded-lg">
            <div className="text-xs text-gray-600">{t('aiSidebar.negotiationTitle')}</div>
            <div className="text-lg font-bold text-purple-600">{t('aiSidebar.negotiationDesc')}</div>
          </div>
          <div className="p-3 bg-orange-50 rounded-lg">
            <div className="text-xs text-gray-600">{t('aiSidebar.shippingTitle')}</div>
            <div className="text-lg font-bold text-orange-600">{t('aiSidebar.shippingDesc')}</div>
          </div>
        </div>

        <div className="mt-4 p-4 bg-green-50 rounded-lg border border-green-200">
          <div className="flex items-center gap-2 mb-2">
            <p className="text-sm font-semibold text-gray-900">
              {t('aiSidebar.timeSaveTitle')}
            </p>
          </div>
          <p className="text-xs text-gray-600">
            {t('aiSidebar.timeSaveDesc')}
          </p>
        </div>
      </div>
    </div>
  )
}
