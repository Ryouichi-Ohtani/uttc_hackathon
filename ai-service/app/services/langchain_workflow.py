from langchain_google_genai import ChatGoogleGenerativeAI
from langchain.prompts import ChatPromptTemplate
from langchain.chains import LLMChain
from langgraph.graph import Graph, END
from typing import Dict, Any, List
import os

class ProductAnalysisWorkflow:
    """
    LangChain/LangGraph workflow for comprehensive product analysis
    """

    def __init__(self):
        self.llm = ChatGoogleGenerativeAI(
            model="gemini-flash-latest",
            google_api_key=os.getenv("GOOGLE_API_KEY")
        )

    def create_workflow(self) -> Graph:
        """
        Create a stateful workflow for product analysis:
        1. Image Analysis
        2. Description Generation
        3. Price Estimation
        4. Safety Check
        """
        workflow = Graph()

        # Define nodes
        workflow.add_node("analyze_images", self.analyze_images_node)
        workflow.add_node("generate_description", self.generate_description_node)
        workflow.add_node("estimate_price", self.estimate_price_node)
        workflow.add_node("safety_check", self.safety_check_node)

        # Define edges
        workflow.add_edge("analyze_images", "generate_description")
        workflow.add_edge("generate_description", "estimate_price")
        workflow.add_edge("estimate_price", "safety_check")
        workflow.add_edge("safety_check", END)

        # Set entry point
        workflow.set_entry_point("analyze_images")

        return workflow.compile()

    def analyze_images_node(self, state: Dict[str, Any]) -> Dict[str, Any]:
        """Node 1: Analyze images and extract features"""
        prompt = ChatPromptTemplate.from_template("""
Based on the product images and title "{title}", identify:
1. Key visual features
2. Condition indicators
3. Brand/manufacturer if visible
4. Material composition

Product category: {category}

Provide a structured analysis.
""")

        chain = LLMChain(llm=self.llm, prompt=prompt)
        result = chain.run(title=state["title"], category=state["category"])

        state["image_analysis"] = result
        return state

    def generate_description_node(self, state: Dict[str, Any]) -> Dict[str, Any]:
        """Node 2: Generate compelling product description"""
        prompt = ChatPromptTemplate.from_template("""
Create an engaging product description for a flea market listing.

Product: {title}
Category: {category}
Image Analysis: {image_analysis}

Write a 100-150 word description that:
- Highlights key features and benefits
- Mentions condition and quality
- Appeals to eco-conscious buyers
- Is SEO-friendly

Description:
""")

        chain = LLMChain(llm=self.llm, prompt=prompt)
        result = chain.run(
            title=state["title"],
            category=state["category"],
            image_analysis=state.get("image_analysis", "")
        )

        state["generated_description"] = result.strip()
        return state

    def estimate_price_node(self, state: Dict[str, Any]) -> Dict[str, Any]:
        """Node 3: Estimate fair market price"""
        prompt = ChatPromptTemplate.from_template("""
Estimate a fair resale price in Japanese Yen for:

Product: {title}
Category: {category}
Description: {description}

Consider:
- Current market trends
- Condition and age
- Brand value
- Supply and demand

Provide ONLY a numeric value in Yen.
""")

        chain = LLMChain(llm=self.llm, prompt=prompt)
        result = chain.run(
            title=state["title"],
            category=state["category"],
            description=state.get("generated_description", "")
        )

        # Extract numeric value
        try:
            price = int(''.join(filter(str.isdigit, result)))
            state["suggested_price"] = price
        except:
            state["suggested_price"] = 1000  # Default fallback

        return state

    def safety_check_node(self, state: Dict[str, Any]) -> Dict[str, Any]:
        """Node 4: Check for prohibited content"""
        prompt = ChatPromptTemplate.from_template("""
Review this product listing for prohibited content:

Title: {title}
Category: {category}
Description: {description}

Check for:
- Illegal items (weapons, drugs)
- Counterfeit goods
- Hazardous materials
- Adult content

Respond with: SAFE or UNSAFE: <reason>
""")

        chain = LLMChain(llm=self.llm, prompt=prompt)
        result = chain.run(
            title=state["title"],
            category=state["category"],
            description=state.get("generated_description", "")
        )

        is_safe = "SAFE" in result.upper()
        state["is_safe"] = is_safe
        state["safety_reason"] = "" if is_safe else result.replace("UNSAFE:", "").strip()

        return state

    def run_analysis(
        self,
        title: str,
        category: str,
        images: List[bytes]
    ) -> Dict[str, Any]:
        """
        Execute the full workflow
        """
        workflow = self.create_workflow()

        initial_state = {
            "title": title,
            "category": category,
            "images": images,
        }

        # Run workflow
        final_state = workflow.invoke(initial_state)

        return {
            "generated_description": final_state.get("generated_description", ""),
            "suggested_price": final_state.get("suggested_price", 1000),
            "is_inappropriate": not final_state.get("is_safe", True),
            "inappropriate_reason": final_state.get("safety_reason", "")
        }
