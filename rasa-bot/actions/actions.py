from typing import Dict, Text, Any, List
from rasa_sdk import Action, Tracker
from rasa_sdk.executor import CollectingDispatcher
from rasa_sdk.events import SlotSet, AllSlotsReset, UserUtteranceReverted
from rasa_sdk.forms import FormValidationAction
import requests
import re

class ValidateManifiestoForm(FormValidationAction):
    def name(self) -> Text:
        return "validate_manifiesto_form"

    async def validate_flete(
        self,
        slot_value: Any,
        dispatcher: CollectingDispatcher,
        tracker: Tracker,
        domain: Dict[Text, Any],
    ) -> Dict[Text, Any]:
        """Valida y normaliza el slot 'flete'"""
        if slot_value is None:
            dispatcher.utter_message(template="utter_pedir_flete")
            return {"flete": None}

        # limpiar el valor
        try:
            value = float(str(slot_value).replace("$", "").replace(".", "").replace(",", "").strip())
            return {"flete": value}
        except ValueError:
            dispatcher.utter_message(template="utter_pedir_flete")
            return {"flete": None}


class ActionSubmitManifiesto(Action):
    def name(self) -> Text:
        return "action_submit_manifiesto"

    def run(
        self,
        dispatcher: CollectingDispatcher,
        tracker: Tracker,
        domain: Dict[Text, Any],
    ) -> List[Dict[Text, Any]]:
        
        # Recopilar datos de los slots
        data = {
            "flete": tracker.get_slot("flete"),
            "descripcion": tracker.get_slot("descripcion"),
            "peso": tracker.get_slot("peso"),
            "fecha_cargue": tracker.get_slot("fecha_cargue"),
            "fecha_descargue": tracker.get_slot("fecha_descargue"),
            "tarjeta": tracker.get_slot("tarjeta"),
            "licencia": tracker.get_slot("licencia"),
            "origen": tracker.get_slot("origen"),
            "destino": tracker.get_slot("destino"),
        }

        # Llamada al backend
        try:
            response = requests.post("http://localhost:8000/api/manifiestos", json=data)
            if response.status_code == 200:
                file_url = response.json().get("file_url", "URL no disponible")
                dispatcher.utter_message(text=f"✅ Tu manifiesto ha sido generado. Puedes descargarlo aquí: {file_url}")
            else:
                dispatcher.utter_message(text="⚠️ Hubo un problema generando el manifiesto. Inténtalo más tarde.")
        except Exception as e:
            dispatcher.utter_message(text=f"❌ Error al contactar el servicio de manifiestos: {str(e)}")

        # Limpiar los slots para que no quede la info en memoria
        return [AllSlotsReset()]


class ActionDefaultFallback(Action):
    def name(self) -> Text:
        return "action_default_fallback"

    def run(
        self,
        dispatcher: CollectingDispatcher,
        tracker: Tracker,
        domain: Dict[Text, Any],
    ) -> List[Dict[Text, Any]]:
        dispatcher.utter_message(template="utter_fallback")
        return [UserUtteranceReverted()]
