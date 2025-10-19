import React from 'react';

export const LeaderboardSidebar: React.FC = () => {
  return (
    <div className="w-80 bg-white/80 backdrop-blur-sm rounded-lg shadow-lg p-6 sticky top-4">
      <div className="flex items-center gap-2 mb-6">
        <h2 className="text-2xl font-bold text-gray-900">AI Agent Stats</h2>
      </div>

      <div className="space-y-4">
        <div className="p-4 bg-gradient-to-r from-green-50 to-emerald-50 rounded-lg border border-green-200">
          <div className="flex items-center gap-2 mb-2">
            <p className="text-sm font-semibold text-gray-900">
              AI Autonomous Agents
            </p>
          </div>
          <p className="text-xs text-gray-600">
            Let AI handle your listing, negotiation, and shipping!
          </p>
        </div>

        <div className="grid grid-cols-1 gap-3">
          <div className="p-3 bg-blue-50 rounded-lg">
            <div className="text-xs text-gray-600">AI Listing Agent</div>
            <div className="text-lg font-bold text-blue-600">Auto-Create</div>
          </div>
          <div className="p-3 bg-purple-50 rounded-lg">
            <div className="text-xs text-gray-600">AI Negotiation</div>
            <div className="text-lg font-bold text-purple-600">Auto-Reply</div>
          </div>
          <div className="p-3 bg-orange-50 rounded-lg">
            <div className="text-xs text-gray-600">AI Shipping</div>
            <div className="text-lg font-bold text-orange-600">Auto-Prep</div>
          </div>
        </div>

        <div className="mt-4 p-4 bg-green-50 rounded-lg border border-green-200">
          <div className="flex items-center gap-2 mb-2">
            <p className="text-sm font-semibold text-gray-900">
              Save 90% of your time!
            </p>
          </div>
          <p className="text-xs text-gray-600">
            AI handles everything from listing to shipping preparation.
          </p>
        </div>
      </div>
    </div>
  );
};
