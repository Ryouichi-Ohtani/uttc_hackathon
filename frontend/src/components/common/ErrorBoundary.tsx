import { Component, ErrorInfo, ReactNode } from 'react'
import { ExclamationTriangleIcon, ArrowPathIcon } from '@heroicons/react/24/outline'

interface Props {
  children: ReactNode
  fallback?: ReactNode
}

interface State {
  hasError: boolean
  error: Error | null
  errorInfo: ErrorInfo | null
}

export class ErrorBoundary extends Component<Props, State> {
  public state: State = {
    hasError: false,
    error: null,
    errorInfo: null,
  }

  public static getDerivedStateFromError(error: Error): Partial<State> {
    return { hasError: true, error }
  }

  public componentDidCatch(error: Error, errorInfo: ErrorInfo) {
    console.error('ErrorBoundary caught an error:', error, errorInfo)
    this.setState({
      error,
      errorInfo,
    })

    // Log to error reporting service (e.g., Sentry, LogRocket)
    // Example: Sentry.captureException(error, { extra: errorInfo })
  }

  private handleReset = () => {
    this.setState({
      hasError: false,
      error: null,
      errorInfo: null,
    })
    // Optionally reload the page or navigate to home
    window.location.href = '/'
  }

  private handleRetry = () => {
    this.setState({
      hasError: false,
      error: null,
      errorInfo: null,
    })
  }

  public render() {
    if (this.state.hasError) {
      // Custom fallback UI provided
      if (this.props.fallback) {
        return this.props.fallback
      }

      // Default error UI
      return (
        <div className="min-h-screen flex items-center justify-center bg-slate-50 dark:bg-dark p-4">
          <div className="max-w-2xl w-full">
            <div className="card p-8 text-center">
              {/* Error Icon */}
              <div className="mb-6 flex justify-center">
                <div className="relative">
                  <div className="absolute inset-0 bg-gradient-to-r from-red-500 to-orange-500 opacity-20 blur-3xl rounded-full animate-pulse" />
                  <div className="relative p-4 bg-gradient-to-br from-red-500 to-orange-600 rounded-full">
                    <ExclamationTriangleIcon className="w-16 h-16 text-white" />
                  </div>
                </div>
              </div>

              {/* Error Message */}
              <h1 className="text-3xl font-bold text-slate-900 dark:text-white mb-3">
                予期しないエラーが発生しました
              </h1>
              <p className="text-slate-600 dark:text-slate-400 mb-6">
                申し訳ございません。問題が発生しました。
                <br />
                ページを再読み込みするか、ホームに戻ってください。
              </p>

              {/* Error Details (Development only) */}
              {import.meta.env.DEV && this.state.error && (
                <details className="mb-6 text-left">
                  <summary className="cursor-pointer text-sm font-medium text-slate-700 dark:text-slate-300 mb-2">
                    エラーの詳細を表示
                  </summary>
                  <div className="p-4 bg-slate-100 dark:bg-slate-800 rounded-lg overflow-auto max-h-64">
                    <p className="text-xs font-mono text-red-600 dark:text-red-400 mb-2">
                      {this.state.error.toString()}
                    </p>
                    {this.state.errorInfo && (
                      <pre className="text-xs font-mono text-slate-700 dark:text-slate-300 whitespace-pre-wrap">
                        {this.state.errorInfo.componentStack}
                      </pre>
                    )}
                  </div>
                </details>
              )}

              {/* Action Buttons */}
              <div className="flex flex-col sm:flex-row gap-3 justify-center">
                <button
                  onClick={this.handleRetry}
                  className="btn-outline flex items-center justify-center gap-2 px-6 py-3"
                >
                  <ArrowPathIcon className="w-5 h-5" />
                  再試行
                </button>
                <button
                  onClick={this.handleReset}
                  className="btn-gradient flex items-center justify-center gap-2 px-6 py-3"
                >
                  ホームに戻る
                </button>
              </div>

              {/* Support Information */}
              <div className="mt-8 pt-6 border-t border-slate-200 dark:border-slate-700">
                <p className="text-sm text-slate-500 dark:text-slate-400">
                  問題が解決しない場合は、
                  <a
                    href="mailto:support@example.com"
                    className="text-primary-600 dark:text-primary-400 hover:underline ml-1"
                  >
                    サポートチーム
                  </a>
                  までお問い合わせください。
                </p>
              </div>
            </div>
          </div>
        </div>
      )
    }

    return this.props.children
  }
}
