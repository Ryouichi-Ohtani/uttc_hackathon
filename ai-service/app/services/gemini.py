import os
import google.generativeai as genai
from typing import List, Dict, Any
import base64
from PIL import Image
import io

class GeminiService:
    def __init__(self):
        api_key = os.getenv("GOOGLE_API_KEY")
        if not api_key:
            raise ValueError("GOOGLE_API_KEY environment variable is required")

        genai.configure(api_key=api_key)
        self.model = genai.GenerativeModel('gemini-flash-latest')

    def analyze_product_images(
        self,
        images: List[bytes],
        title: str,
        category: str
    ) -> Dict[str, Any]:
        """
        Analyze product images using Gemini's multimodal capabilities
        """
        # Convert bytes to PIL Images
        pil_images = []
        for img_bytes in images:
            pil_images.append(Image.open(io.BytesIO(img_bytes)))

        prompt = f"""
You are an expert product analyst for a sustainable flea market app. Analyze these images of a "{title}" in the "{category}" category.

Provide a JSON response with the following structure:
{{
    "description": "A compelling, SEO-friendly product description (100-150 words) highlighting key features and condition",
    "suggested_price_jpy": <estimated fair market price in Japanese Yen>,
    "estimated_weight_kg": <estimated weight in kg>,
    "manufacturer_country": "Most likely country of manufacture",
    "estimated_manufacturing_year": <approximate year of manufacture>,
    "detected_objects": ["object1", "object2", ...],
    "is_inappropriate": <true if content is inappropriate/illegal, false otherwise>,
    "inappropriate_reason": "Reason if inappropriate, empty string otherwise"
}}

Focus on accuracy and helpfulness. Consider the item's condition, brand, and current market trends.
"""

        try:
            # Include first image (or multiple if available)
            content = [prompt]
            content.extend(pil_images[:3])  # Limit to first 3 images

            response = self.model.generate_content(content)

            # Parse JSON response
            import json
            result = json.loads(response.text.strip().replace("```json", "").replace("```", ""))

            return result
        except Exception as e:
            print(f"Gemini API error: {e}")
            # Return fallback response
            return {
                "description": f"Quality {title} in {category} category. Great condition.",
                "suggested_price_jpy": 1000,
                "estimated_weight_kg": 0.5,
                "manufacturer_country": "Unknown",
                "estimated_manufacturing_year": 2020,
                "detected_objects": [category],
                "is_inappropriate": False,
                "inappropriate_reason": ""
            }

    def detect_inappropriate_content(self, images: List[bytes]) -> tuple[bool, str]:
        """
        Detect inappropriate or prohibited content
        """
        pil_images = [Image.open(io.BytesIO(img)) for img in images]

        prompt = """
Analyze these images for inappropriate or prohibited content including:
- Weapons, drugs, or illegal items
- Counterfeit or fake branded goods
- Adult/NSFW content
- Live animals
- Hazardous materials

Respond with JSON:
{
    "is_inappropriate": <true/false>,
    "reason": "Brief explanation if inappropriate, empty string otherwise"
}
"""

        try:
            content = [prompt] + pil_images[:2]
            response = self.model.generate_content(content)

            import json
            result = json.loads(response.text.strip().replace("```json", "").replace("```", ""))

            return result.get("is_inappropriate", False), result.get("reason", "")
        except:
            return False, ""

    def translate_search_query(self, query: str, target_language: str = "ja") -> Dict[str, str]:
        """
        Translate search query to multiple languages for multilingual search
        Returns translations in Japanese, English, and other relevant languages
        """
        prompt = f"""
Translate the following search query to help with multilingual product search.
Original query: "{query}"

Provide translations and search-relevant keywords in JSON format:
{{
    "japanese": "Japanese translation/keywords",
    "english": "English translation/keywords",
    "romanized": "Romanized version if applicable",
    "keywords": ["keyword1", "keyword2", "keyword3"],
    "detected_language": "original language code (ja/en/etc)",
    "search_intent": "brief description of what user is looking for"
}}

For example:
- If query is "スマホ", return english "smartphone", keywords ["phone", "mobile", "iPhone", "Android"]
- If query is "laptop", return japanese "ノートパソコン", keywords ["PC", "MacBook", "computer"]
- Be creative with synonyms and related terms for better search results
"""

        try:
            response = self.model.generate_content(prompt)
            import json
            result = json.loads(response.text.strip().replace("```json", "").replace("```", ""))
            return result
        except Exception as e:
            print(f"Translation error: {e}")
            # Fallback: return original query
            return {
                "japanese": query,
                "english": query,
                "romanized": query,
                "keywords": [query],
                "detected_language": "unknown",
                "search_intent": query
            }
