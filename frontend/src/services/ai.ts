import axios from 'axios'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'

export interface TranslateSearchResponse {
  japanese: string
  english: string
  romanized: string
  keywords: string[]
  detected_language: string
  search_intent: string
}

class AIService {
  async translateSearch(query: string): Promise<TranslateSearchResponse> {
    const response = await axios.post(`${API_BASE_URL}/v1/ai/translate-search`, {
      query,
    })
    return response.data
  }
}

export const aiService = new AIService()
