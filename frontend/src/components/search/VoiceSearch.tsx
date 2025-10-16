import React, { useState, useEffect } from 'react';
import { aiService } from '@/services/ai';

interface VoiceSearchProps {
  onSearch?: (query: string) => void;
  onTranscript?: (query: string) => void;
  onIntentDetected?: (intent: string, keywords: string[]) => void;
  placeholder?: string;
}

export const VoiceSearch: React.FC<VoiceSearchProps> = ({
  onSearch,
  onTranscript,
  onIntentDetected,
  placeholder = '音声で検索...',
}) => {
  const [isListening, setIsListening] = useState(false);
  const [transcript, setTranscript] = useState('');
  const [supported, setSupported] = useState(false);
  const [recognition, setRecognition] = useState<any>(null);
  const [isProcessing, setIsProcessing] = useState(false);
  const [detectedIntent, setDetectedIntent] = useState<string>('');

  useEffect(() => {
    // Check if Speech Recognition API is supported
    const SpeechRecognition =
      (window as any).SpeechRecognition || (window as any).webkitSpeechRecognition;

    if (SpeechRecognition) {
      const recognitionInstance = new SpeechRecognition();
      recognitionInstance.continuous = false;
      recognitionInstance.interimResults = true;
      recognitionInstance.lang = 'ja-JP'; // Japanese

      recognitionInstance.onresult = async (event: any) => {
        const current = event.resultIndex;
        const transcriptText = event.results[current][0].transcript;
        setTranscript(transcriptText);

        if (event.results[current].isFinal) {
          setIsListening(false);
          setIsProcessing(true);

          try {
            // Use AI to understand intent and enhance search
            const translation = await aiService.translateSearch(transcriptText);
            setDetectedIntent(translation.search_intent);

            // Notify parent with intent and keywords
            if (onIntentDetected) {
              onIntentDetected(translation.search_intent, translation.keywords);
            }

            // Build enhanced query
            const enhancedQuery = [
              transcriptText,
              translation.japanese,
              translation.english,
              ...translation.keywords
            ].filter(Boolean).join(' ');

            if (onSearch) onSearch(enhancedQuery);
            if (onTranscript) onTranscript(enhancedQuery);
          } catch (error) {
            console.error('Intent detection error:', error);
            // Fallback to original transcript
            if (onSearch) onSearch(transcriptText);
            if (onTranscript) onTranscript(transcriptText);
          } finally {
            setIsProcessing(false);
          }
        }
      };

      recognitionInstance.onerror = (event: any) => {
        console.error('Speech recognition error:', event.error);
        setIsListening(false);
      };

      recognitionInstance.onend = () => {
        setIsListening(false);
      };

      setRecognition(recognitionInstance);
      setSupported(true);
    }
  }, [onSearch, onTranscript]);

  const startListening = () => {
    if (recognition) {
      setTranscript('');
      setIsListening(true);
      recognition.start();
    }
  };

  const stopListening = () => {
    if (recognition) {
      recognition.stop();
      setIsListening(false);
    }
  };

  const handleManualSearch = () => {
    if (transcript) {
      if (onSearch) onSearch(transcript);
      if (onTranscript) onTranscript(transcript);
    }
  };

  if (!supported) {
    return (
      <div className="text-sm text-gray-500">
        音声検索はこのブラウザではサポートされていません
      </div>
    );
  }

  return (
    <div className="flex flex-col gap-2">
      <div className="flex items-center gap-2">
        <button
          onClick={isListening ? stopListening : startListening}
          disabled={isProcessing}
          className={`p-3 rounded-lg transition flex items-center gap-2 ${
            isListening
              ? 'bg-red-500 text-white hover:bg-red-600'
              : isProcessing
              ? 'bg-gray-400 text-white cursor-not-allowed'
              : 'bg-green-600 text-white hover:bg-green-700'
          }`}
          title={isListening ? '停止' : isProcessing ? 'Processing...' : '音声検索'}
        >
          {isListening ? '⏹️ Listening...' : isProcessing ? '🔄 Processing...' : '🎤 Voice'}
        </button>
      </div>

      {transcript && (
        <div className="text-sm p-2 bg-gray-100 rounded-lg">
          <span className="font-medium">Transcript:</span> {transcript}
        </div>
      )}

      {detectedIntent && (
        <div className="text-sm p-2 bg-primary-50 rounded-lg text-primary-700">
          <span className="font-medium">🎯 Intent:</span> {detectedIntent}
        </div>
      )}
    </div>
  );
};
