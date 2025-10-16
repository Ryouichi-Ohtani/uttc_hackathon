import grpc
from concurrent import futures
import sys
import os

# Add root directory to path
sys.path.append(os.path.join(os.path.dirname(__file__), '../../'))

from protos import product_analysis_pb2, product_analysis_pb2_grpc
from app.services.gemini import GeminiService
from app.services.co2_calculator import CO2Calculator
from app.services.langchain_workflow import ProductAnalysisWorkflow

class ProductAnalysisServicer(product_analysis_pb2_grpc.ProductAnalysisServiceServicer):
    def __init__(self):
        self.gemini = GeminiService()
        self.co2_calc = CO2Calculator()
        self.workflow = ProductAnalysisWorkflow()

    def AnalyzeProduct(self, request, context):
        """
        Comprehensive product analysis using LangChain workflow
        """
        try:
            # Use LangChain workflow for advanced analysis
            workflow_result = self.workflow.run_analysis(
                title=request.title,
                category=request.category,
                images=list(request.images)
            )

            # Gemini multimodal analysis for additional details
            gemini_result = self.gemini.analyze_product_images(
                images=list(request.images),
                title=request.title,
                category=request.category
            )

            # Calculate CO2 impact
            co2_result = self.co2_calc.calculate_co2_impact(
                category=request.category,
                weight_kg=gemini_result.get("estimated_weight_kg", 0.5),
                manufacturer_country=gemini_result.get("manufacturer_country", "Unknown"),
                manufacturing_year=gemini_result.get("estimated_manufacturing_year", 2020)
            )

            # Combine results
            description = workflow_result.get("generated_description") or gemini_result.get("description", "")
            if request.user_provided_description:
                description = request.user_provided_description

            return product_analysis_pb2.AnalyzeProductResponse(
                generated_description=description,
                suggested_price=workflow_result.get("suggested_price", gemini_result.get("suggested_price_jpy", 1000)),
                estimated_weight_kg=gemini_result.get("estimated_weight_kg", 0.5),
                manufacturer_country=gemini_result.get("manufacturer_country", "Unknown"),
                estimated_manufacturing_year=gemini_result.get("estimated_manufacturing_year", 2020),
                co2_impact_kg=co2_result["saved_kg"],
                is_inappropriate=workflow_result.get("is_inappropriate", False),
                inappropriate_reason=workflow_result.get("inappropriate_reason", ""),
                detected_objects=gemini_result.get("detected_objects", [])
            )

        except Exception as e:
            print(f"Error in AnalyzeProduct: {e}")
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(str(e))
            return product_analysis_pb2.AnalyzeProductResponse()

    def CalculateCO2Impact(self, request, context):
        """
        Calculate CO2 savings
        """
        try:
            result = self.co2_calc.calculate_co2_impact(
                category=request.category,
                weight_kg=request.weight_kg,
                manufacturer_country=request.manufacturer_country,
                manufacturing_year=request.manufacturing_year
            )

            return product_analysis_pb2.CalculateCO2Response(
                buying_new_kg=result["buying_new_kg"],
                buying_used_kg=result["buying_used_kg"],
                saved_kg=result["saved_kg"]
            )

        except Exception as e:
            print(f"Error in CalculateCO2Impact: {e}")
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(str(e))
            return product_analysis_pb2.CalculateCO2Response()

def serve():
    port = os.getenv("GRPC_PORT", "50051")
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))

    product_analysis_pb2_grpc.add_ProductAnalysisServiceServicer_to_server(
        ProductAnalysisServicer(), server
    )

    server.add_insecure_port(f'[::]:{port}')
    server.start()
    print(f"gRPC server started on port {port}")
    server.wait_for_termination()

if __name__ == '__main__':
    serve()
