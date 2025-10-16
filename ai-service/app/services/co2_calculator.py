from typing import Dict
import math

class CO2Calculator:
    """
    Calculate CO2 emissions saved by buying used instead of new products.

    Based on research data:
    - Manufacturing emissions vary by product category
    - Shipping distance affects total emissions
    - Product age affects degradation/efficiency
    """

    # Base CO2 emissions for manufacturing (kg CO2 per kg of product)
    CATEGORY_EMISSIONS = {
        "electronics": 50.0,      # High energy manufacturing
        "clothing": 15.0,          # Textile production
        "furniture": 8.0,          # Wood/metal processing
        "books": 2.5,              # Paper production
        "toys": 12.0,              # Plastic manufacturing
        "appliances": 45.0,        # Complex manufacturing
        "sports": 10.0,            # Mixed materials
        "accessories": 8.0,        # Fashion items
        "default": 10.0            # Generic fallback
    }

    # Shipping emissions (kg CO2 per km per kg)
    SHIPPING_EMISSIONS_PER_KM = 0.00014

    # Country distances from Japan (approximate km)
    COUNTRY_DISTANCES = {
        "Japan": 0,
        "China": 3000,
        "USA": 10000,
        "Germany": 9000,
        "Vietnam": 4000,
        "Thailand": 4500,
        "Korea": 1200,
        "Taiwan": 2200,
        "Unknown": 5000  # Default
    }

    def calculate_co2_impact(
        self,
        category: str,
        weight_kg: float,
        manufacturer_country: str,
        manufacturing_year: int
    ) -> Dict[str, float]:
        """
        Calculate CO2 saved by buying used vs new
        """
        # Get base manufacturing emissions
        category_lower = category.lower()
        base_emission_rate = self.CATEGORY_EMISSIONS.get(
            category_lower,
            self.CATEGORY_EMISSIONS["default"]
        )

        # Calculate new product emissions
        manufacturing_co2 = weight_kg * base_emission_rate

        # Calculate shipping emissions
        distance = self.COUNTRY_DISTANCES.get(manufacturer_country, 5000)
        shipping_co2_new = weight_kg * distance * self.SHIPPING_EMISSIONS_PER_KM

        # Total for new product
        total_new_co2 = manufacturing_co2 + shipping_co2_new

        # Used product only has local shipping (assume 50km average)
        shipping_co2_used = weight_kg * 50 * self.SHIPPING_EMISSIONS_PER_KM

        # Account for degradation (older items may have higher maintenance emissions)
        current_year = 2024
        age = current_year - manufacturing_year
        degradation_factor = 1.0 + (age * 0.02)  # 2% per year

        total_used_co2 = shipping_co2_used * degradation_factor

        # CO2 saved
        co2_saved = total_new_co2 - total_used_co2

        # Ensure non-negative
        co2_saved = max(co2_saved, 0)

        return {
            "buying_new_kg": round(total_new_co2, 2),
            "buying_used_kg": round(total_used_co2, 2),
            "saved_kg": round(co2_saved, 2)
        }

    def calculate_environmental_equivalents(self, co2_saved_kg: float) -> Dict[str, float]:
        """
        Convert CO2 savings to relatable equivalents
        """
        return {
            "trees_planted": round(co2_saved_kg / 20, 2),  # 1 tree absorbs ~20kg/year
            "car_km_avoided": round(co2_saved_kg / 0.12, 2),  # Car emits ~0.12kg/km
            "plastic_bottles_recycled": round(co2_saved_kg / 0.082, 0),  # 1 bottle = 82g CO2
            "light_bulb_hours": round(co2_saved_kg / 0.0006, 0)  # 1 hour = 0.6g CO2
        }
