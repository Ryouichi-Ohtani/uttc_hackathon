interface ChatbotButtonProps {
  onClick: () => void
  isOpen: boolean
}

export const ChatbotButton = ({ onClick, isOpen }: ChatbotButtonProps) => {
  if (isOpen) return null

  return (
    <button
      onClick={onClick}
      className="fixed bottom-8 right-8 bg-gradient-to-r from-green-500 to-emerald-600 text-white rounded-full shadow-2xl hover:shadow-3xl hover:scale-105 transition-all duration-300 flex items-center gap-3 px-6 py-4 z-40 group"
      aria-label="AIとchatで相談する"
    >
      <span className="text-2xl group-hover:scale-110 transition-transform">💬</span>
      <span className="font-bold text-sm whitespace-nowrap">AIとchatで相談する</span>
    </button>
  )
}
