import { useState, useRef, useEffect } from 'react'
import { superClaudeService, SuperClaudeMessage, Attachment, CodeBlock, ConversationContext, SuperClaudeCapability } from '@/services/superClaude'
import {
  XMarkIcon,
  PaperAirplaneIcon,
  PaperClipIcon,
  SparklesIcon,
  CodeBracketIcon,
  DocumentDuplicateIcon,
  ShareIcon,
  TrashIcon,
  Cog6ToothIcon,
  FolderOpenIcon,
  ChevronLeftIcon,
  ChevronRightIcon,
  ArrowDownTrayIcon,
  PlayIcon,
  ChevronDownIcon,
  CheckIcon,
  ClockIcon,
  CpuChipIcon
} from '@heroicons/react/24/outline'
import { StarIcon } from '@heroicons/react/24/solid'
import toast from 'react-hot-toast'
import ReactMarkdown from 'react-markdown'
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter'
import { vscDarkPlus } from 'react-syntax-highlighter/dist/esm/styles/prism'

interface SuperClaudeProps {
  onClose: () => void
}

export const SuperClaude = ({ onClose }: SuperClaudeProps) => {
  const [input, setInput] = useState('')
  const [loading, setLoading] = useState(false)
  const [conversation, setConversation] = useState<ConversationContext | null>(null)
  const [conversations, setConversations] = useState<ConversationContext[]>([])
  const [attachments, setAttachments] = useState<Attachment[]>([])
  const [showSidebar, setShowSidebar] = useState(true)
  const [showCapabilities, setShowCapabilities] = useState(true)
  const [selectedModel, setSelectedModel] = useState<'claude-3-opus' | 'claude-3-sonnet' | 'claude-3-haiku'>('claude-3-opus')
  const [temperature, setTemperature] = useState(0.7)

  const fileInputRef = useRef<HTMLInputElement>(null)
  const messagesEndRef = useRef<HTMLDivElement>(null)
  const textareaRef = useRef<HTMLTextAreaElement>(null)

  useEffect(() => {
    // Load conversations and create/load current one
    const loadConversations = async () => {
      const history = superClaudeService.getConversationHistory()
      setConversations(history)

      let current = superClaudeService.getCurrentConversation()
      if (!current) {
        current = await superClaudeService.createConversation()
      }
      setConversation(current)
    }
    loadConversations()
  }, [])

  useEffect(() => {
    // Auto-scroll to bottom when new messages arrive
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [conversation?.messages])

  const handleSend = async () => {
    if (!input.trim() && attachments.length === 0) return

    setLoading(true)
    setShowCapabilities(false)

    try {
      const response = await superClaudeService.sendMessage(
        input,
        attachments,
        { model: selectedModel, temperature }
      )

      // Update conversation
      const updatedConversation = superClaudeService.getCurrentConversation()
      if (updatedConversation) {
        setConversation({ ...updatedConversation })
      }

      // Clear input
      setInput('')
      setAttachments([])

      // Auto-resize textarea
      if (textareaRef.current) {
        textareaRef.current.style.height = 'auto'
      }
    } catch (error) {
      toast.error('メッセージの送信に失敗しました')
    } finally {
      setLoading(false)
    }
  }

  const handleFileSelect = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const files = e.target.files
    if (!files) return

    for (const file of Array.from(files)) {
      try {
        const attachment = await superClaudeService.analyzeFile(file)
        setAttachments(prev => [...prev, attachment])
        toast.success(`${file.name} を追加しました`)
      } catch (error) {
        toast.error(`${file.name} の追加に失敗しました`)
      }
    }
  }

  const handleNewConversation = async () => {
    const newConv = await superClaudeService.createConversation()
    setConversation(newConv)
    setConversations([newConv, ...conversations])
    setShowCapabilities(true)
  }

  const handleSwitchConversation = (convId: string) => {
    const conv = superClaudeService.switchConversation(convId)
    if (conv) {
      setConversation(conv)
      setShowCapabilities(conv.messages.length === 0)
    }
  }

  const handleDeleteConversation = (convId: string) => {
    superClaudeService.deleteConversation(convId)
    setConversations(conversations.filter(c => c.id !== convId))
    if (conversation?.id === convId) {
      handleNewConversation()
    }
  }

  const handleExport = () => {
    const markdown = superClaudeService.exportConversation('markdown')
    const blob = new Blob([markdown], { type: 'text/markdown' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `conversation-${Date.now()}.md`
    a.click()
    URL.revokeObjectURL(url)
    toast.success('会話をエクスポートしました')
  }

  const handleShare = async () => {
    try {
      const shareUrl = await superClaudeService.shareConversation(['user@example.com'])
      navigator.clipboard.writeText(shareUrl)
      toast.success('共有リンクをコピーしました')
    } catch (error) {
      toast.error('共有に失敗しました')
    }
  }

  const handleExecuteCode = async (codeBlock: CodeBlock) => {
    try {
      const output = await superClaudeService.executeCode(codeBlock)
      toast.success('コードを実行しました')
      console.log(output)
    } catch (error) {
      toast.error('コード実行に失敗しました')
    }
  }

  const handleCapabilityClick = (capability: SuperClaudeCapability) => {
    const randomExample = capability.examples[Math.floor(Math.random() * capability.examples.length)]
    setInput(randomExample)
    textareaRef.current?.focus()
  }

  const modelInfo = {
    'claude-3-opus': { name: 'Claude 3 Opus', desc: '最高性能', icon: '🎯' },
    'claude-3-sonnet': { name: 'Claude 3 Sonnet', desc: 'バランス型', icon: '⚡' },
    'claude-3-haiku': { name: 'Claude 3 Haiku', desc: '高速応答', icon: '🚀' }
  }

  return (
    <div className="fixed inset-0 z-modal flex">
      {/* Backdrop */}
      <div className="absolute inset-0 bg-black/50 backdrop-blur-sm" onClick={onClose} />

      {/* Main Container */}
      <div className="relative flex w-full max-w-7xl mx-auto my-6 animate-scale">
        {/* Sidebar */}
        {showSidebar && (
          <div className="w-64 bg-white dark:bg-dark-card rounded-l-2xl border-r border-slate-200 dark:border-dark-border flex flex-col">
            {/* Sidebar Header */}
            <div className="p-4 border-b border-slate-200 dark:border-dark-border">
              <button
                onClick={handleNewConversation}
                className="w-full btn-gradient py-2 flex items-center justify-center gap-2"
              >
                <SparklesIcon className="w-5 h-5" />
                新しい会話
              </button>
            </div>

            {/* Conversation List */}
            <div className="flex-1 overflow-y-auto p-4 space-y-2">
              {conversations.map(conv => (
                <div
                  key={conv.id}
                  className={`group relative p-3 rounded-lg cursor-pointer transition-all ${
                    conversation?.id === conv.id
                      ? 'bg-primary-50 dark:bg-primary-900/20 border border-primary-200 dark:border-primary-800'
                      : 'hover:bg-slate-50 dark:hover:bg-slate-800'
                  }`}
                  onClick={() => handleSwitchConversation(conv.id)}
                >
                  <div className="flex items-start justify-between">
                    <div className="flex-1 min-w-0">
                      <h4 className="font-medium text-sm text-slate-900 dark:text-white truncate">
                        {conv.title}
                      </h4>
                      <p className="text-xs text-slate-500 dark:text-slate-400 mt-1">
                        {conv.messages.length} メッセージ
                      </p>
                      <p className="text-xs text-slate-400 dark:text-slate-500 mt-0.5">
                        {new Date(conv.updatedAt).toLocaleString('ja-JP', {
                          month: 'short',
                          day: 'numeric',
                          hour: '2-digit',
                          minute: '2-digit'
                        })}
                      </p>
                    </div>
                    {conversation?.id === conv.id && (
                      <button
                        onClick={(e) => {
                          e.stopPropagation()
                          handleDeleteConversation(conv.id)
                        }}
                        className="opacity-0 group-hover:opacity-100 p-1 hover:bg-red-100 dark:hover:bg-red-900/20 rounded transition-all"
                      >
                        <TrashIcon className="w-4 h-4 text-red-600 dark:text-red-400" />
                      </button>
                    )}
                  </div>
                </div>
              ))}
            </div>

            {/* Model Settings */}
            <div className="p-4 border-t border-slate-200 dark:border-dark-border space-y-3">
              {/* Model Selector */}
              <div>
                <label className="text-xs font-medium text-slate-600 dark:text-slate-400 mb-1 block">
                  AIモデル
                </label>
                <div className="relative">
                  <select
                    value={selectedModel}
                    onChange={(e) => setSelectedModel(e.target.value as any)}
                    className="w-full text-sm input pr-8"
                  >
                    {Object.entries(modelInfo).map(([key, info]) => (
                      <option key={key} value={key}>
                        {info.icon} {info.name}
                      </option>
                    ))}
                  </select>
                  <ChevronDownIcon className="absolute right-2 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400 pointer-events-none" />
                </div>
              </div>

              {/* Temperature Slider */}
              <div>
                <label className="text-xs font-medium text-slate-600 dark:text-slate-400 mb-1 block">
                  創造性: {temperature}
                </label>
                <input
                  type="range"
                  min="0"
                  max="1"
                  step="0.1"
                  value={temperature}
                  onChange={(e) => setTemperature(Number(e.target.value))}
                  className="w-full"
                />
              </div>
            </div>
          </div>
        )}

        {/* Main Chat Area */}
        <div className="flex-1 bg-white dark:bg-dark-card rounded-r-2xl flex flex-col">
          {/* Header */}
          <div className="flex items-center justify-between p-4 border-b border-slate-200 dark:border-dark-border">
            <div className="flex items-center gap-3">
              <button
                onClick={() => setShowSidebar(!showSidebar)}
                className="p-2 hover:bg-slate-100 dark:hover:bg-slate-800 rounded-lg transition-colors"
              >
                {showSidebar ? <ChevronLeftIcon className="w-5 h-5" /> : <ChevronRightIcon className="w-5 h-5" />}
              </button>
              <div className="flex items-center gap-2">
                <div className="w-10 h-10 bg-gradient-to-br from-primary-500 to-accent-500 rounded-xl flex items-center justify-center">
                  <SparklesIcon className="w-6 h-6 text-white" />
                </div>
                <div>
                  <h2 className="text-lg font-bold text-slate-900 dark:text-white">
                    Super Claude
                  </h2>
                  <p className="text-xs text-slate-500 dark:text-slate-400">
                    {modelInfo[selectedModel].icon} {modelInfo[selectedModel].desc}
                  </p>
                </div>
              </div>
            </div>

            <div className="flex items-center gap-2">
              <button
                onClick={handleExport}
                className="p-2 hover:bg-slate-100 dark:hover:bg-slate-800 rounded-lg transition-colors"
                title="エクスポート"
              >
                <ArrowDownTrayIcon className="w-5 h-5 text-slate-600 dark:text-slate-400" />
              </button>
              <button
                onClick={handleShare}
                className="p-2 hover:bg-slate-100 dark:hover:bg-slate-800 rounded-lg transition-colors"
                title="共有"
              >
                <ShareIcon className="w-5 h-5 text-slate-600 dark:text-slate-400" />
              </button>
              <button
                onClick={onClose}
                className="p-2 hover:bg-slate-100 dark:hover:bg-slate-800 rounded-lg transition-colors"
              >
                <XMarkIcon className="w-5 h-5 text-slate-600 dark:text-slate-400" />
              </button>
            </div>
          </div>

          {/* Messages Area */}
          <div className="flex-1 overflow-y-auto p-4 space-y-4">
            {/* Capabilities Grid (shown when no messages) */}
            {showCapabilities && conversation?.messages.length === 0 && (
              <div className="max-w-4xl mx-auto py-8">
                <h3 className="text-2xl font-bold text-center mb-2 gradient-text">
                  Super Claudeの機能
                </h3>
                <p className="text-center text-slate-600 dark:text-slate-400 mb-8">
                  以下の機能から選択するか、質問を入力してください
                </p>
                <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
                  {superClaudeService.capabilities.map(cap => (
                    <button
                      key={cap.id}
                      onClick={() => handleCapabilityClick(cap)}
                      className="card p-4 hover:shadow-lg hover:-translate-y-1 transition-all text-left"
                    >
                      <div className="text-3xl mb-2">{cap.icon}</div>
                      <h4 className="font-semibold text-sm text-slate-900 dark:text-white mb-1">
                        {cap.name}
                      </h4>
                      <p className="text-xs text-slate-600 dark:text-slate-400 line-clamp-2">
                        {cap.description}
                      </p>
                    </button>
                  ))}
                </div>
              </div>
            )}

            {/* Messages */}
            {conversation?.messages.map((message) => (
              <div
                key={message.id}
                className={`flex ${message.role === 'user' ? 'justify-end' : 'justify-start'} animate-fade-up`}
              >
                <div className={`max-w-3xl ${message.role === 'user' ? 'order-2' : 'order-1'}`}>
                  <div className="flex items-start gap-3 mb-1">
                    {message.role === 'assistant' && (
                      <div className="w-8 h-8 bg-gradient-to-br from-primary-500 to-accent-500 rounded-lg flex items-center justify-center flex-shrink-0">
                        <SparklesIcon className="w-5 h-5 text-white" />
                      </div>
                    )}
                    <div className="flex-1">
                      <div
                        className={`inline-block px-4 py-3 rounded-2xl ${
                          message.role === 'user'
                            ? 'bg-primary-600 text-white'
                            : 'bg-slate-100 dark:bg-slate-800 text-slate-900 dark:text-white'
                        }`}
                      >
                        {message.role === 'assistant' ? (
                          <ReactMarkdown
                            className="prose prose-sm dark:prose-invert max-w-none"
                            components={{
                              code({ node, inline, className, children, ...props }) {
                                const match = /language-(\w+)/.exec(className || '')
                                const codeString = String(children).replace(/\n$/, '')

                                if (!inline && match) {
                                  const codeBlock = message.codeBlocks?.find(b =>
                                    b.code.trim() === codeString.trim()
                                  )

                                  return (
                                    <div className="relative group my-3">
                                      <div className="flex items-center justify-between bg-slate-800 px-3 py-2 rounded-t-lg">
                                        <span className="text-xs text-slate-400">{match[1]}</span>
                                        <div className="flex items-center gap-2">
                                          {codeBlock?.executable && (
                                            <button
                                              onClick={() => handleExecuteCode(codeBlock)}
                                              className="text-xs text-green-400 hover:text-green-300 flex items-center gap-1"
                                            >
                                              <PlayIcon className="w-3 h-3" />
                                              実行
                                            </button>
                                          )}
                                          <button
                                            onClick={() => {
                                              navigator.clipboard.writeText(codeString)
                                              toast.success('コピーしました')
                                            }}
                                            className="text-xs text-slate-400 hover:text-white"
                                          >
                                            <DocumentDuplicateIcon className="w-4 h-4" />
                                          </button>
                                        </div>
                                      </div>
                                      <SyntaxHighlighter
                                        language={match[1]}
                                        style={vscDarkPlus}
                                        customStyle={{
                                          margin: 0,
                                          borderRadius: '0 0 0.5rem 0.5rem'
                                        }}
                                      >
                                        {codeString}
                                      </SyntaxHighlighter>
                                    </div>
                                  )
                                }

                                return inline ? (
                                  <code className="px-1 py-0.5 bg-slate-200 dark:bg-slate-700 rounded text-sm" {...props}>
                                    {children}
                                  </code>
                                ) : (
                                  <code {...props}>{children}</code>
                                )
                              }
                            }}
                          >
                            {message.content}
                          </ReactMarkdown>
                        ) : (
                          <div className="whitespace-pre-wrap">{message.content}</div>
                        )}

                        {/* Attachments */}
                        {message.attachments?.map(att => (
                          <div key={att.id} className="mt-2 p-2 bg-white/10 rounded-lg flex items-center gap-2">
                            <PaperClipIcon className="w-4 h-4" />
                            <span className="text-sm">{att.name}</span>
                            <span className="text-xs opacity-60">({att.type})</span>
                          </div>
                        ))}
                      </div>

                      {/* Metadata */}
                      {message.metadata && (
                        <div className="flex items-center gap-3 mt-1 px-1">
                          <span className="text-xs text-slate-400 flex items-center gap-1">
                            <ClockIcon className="w-3 h-3" />
                            {message.timestamp.toLocaleTimeString('ja-JP', {
                              hour: '2-digit',
                              minute: '2-digit'
                            })}
                          </span>
                          {message.metadata.processingTime && (
                            <span className="text-xs text-slate-400 flex items-center gap-1">
                              <CpuChipIcon className="w-3 h-3" />
                              {(message.metadata.processingTime / 1000).toFixed(1)}s
                            </span>
                          )}
                          {message.metadata.tokenCount && (
                            <span className="text-xs text-slate-400">
                              {message.metadata.tokenCount.input + message.metadata.tokenCount.output} tokens
                            </span>
                          )}
                        </div>
                      )}
                    </div>
                  </div>
                </div>
              </div>
            ))}

            {/* Loading */}
            {loading && (
              <div className="flex justify-start animate-fade-up">
                <div className="flex items-start gap-3">
                  <div className="w-8 h-8 bg-gradient-to-br from-primary-500 to-accent-500 rounded-lg flex items-center justify-center animate-pulse">
                    <SparklesIcon className="w-5 h-5 text-white" />
                  </div>
                  <div className="bg-slate-100 dark:bg-slate-800 px-4 py-3 rounded-2xl">
                    <div className="flex items-center gap-2">
                      <div className="w-2 h-2 bg-primary-600 rounded-full animate-bounce" />
                      <div className="w-2 h-2 bg-primary-600 rounded-full animate-bounce" style={{ animationDelay: '0.1s' }} />
                      <div className="w-2 h-2 bg-primary-600 rounded-full animate-bounce" style={{ animationDelay: '0.2s' }} />
                    </div>
                  </div>
                </div>
              </div>
            )}

            <div ref={messagesEndRef} />
          </div>

          {/* Input Area */}
          <div className="border-t border-slate-200 dark:border-dark-border p-4">
            {/* Attachments */}
            {attachments.length > 0 && (
              <div className="flex flex-wrap gap-2 mb-3">
                {attachments.map(att => (
                  <div
                    key={att.id}
                    className="inline-flex items-center gap-2 px-3 py-1.5 bg-slate-100 dark:bg-slate-800 rounded-lg"
                  >
                    <PaperClipIcon className="w-4 h-4 text-slate-500" />
                    <span className="text-sm text-slate-700 dark:text-slate-300">{att.name}</span>
                    <button
                      onClick={() => setAttachments(attachments.filter(a => a.id !== att.id))}
                      className="text-slate-400 hover:text-slate-600 dark:hover:text-slate-200"
                    >
                      <XMarkIcon className="w-4 h-4" />
                    </button>
                  </div>
                ))}
              </div>
            )}

            <div className="flex items-end gap-2">
              <button
                onClick={() => fileInputRef.current?.click()}
                className="p-2.5 hover:bg-slate-100 dark:hover:bg-slate-800 rounded-lg transition-colors"
                title="ファイルを添付"
              >
                <PaperClipIcon className="w-5 h-5 text-slate-500" />
              </button>

              <textarea
                ref={textareaRef}
                value={input}
                onChange={(e) => {
                  setInput(e.target.value)
                  // Auto-resize
                  e.target.style.height = 'auto'
                  e.target.style.height = `${e.target.scrollHeight}px`
                }}
                onKeyDown={(e) => {
                  if (e.key === 'Enter' && !e.shiftKey) {
                    e.preventDefault()
                    handleSend()
                  }
                }}
                placeholder="Super Claudeに質問する... (Shift+Enterで改行)"
                className="flex-1 resize-none input min-h-[44px] max-h-32"
                rows={1}
                disabled={loading}
              />

              <button
                onClick={handleSend}
                disabled={loading || (!input.trim() && attachments.length === 0)}
                className="btn-gradient p-2.5 rounded-lg disabled:opacity-50 disabled:cursor-not-allowed"
              >
                <PaperAirplaneIcon className="w-5 h-5" />
              </button>
            </div>

            <input
              ref={fileInputRef}
              type="file"
              multiple
              onChange={handleFileSelect}
              className="hidden"
              accept="image/*,.pdf,.txt,.md,.json,.csv,.xml,.yaml,.yml,.js,.ts,.py,.java,.cpp,.c,.go,.rs,.swift,.kt,.rb,.php,.sh,.sql"
            />
          </div>
        </div>
      </div>
    </div>
  )
}