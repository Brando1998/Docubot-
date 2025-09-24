from typing import Dict, Text, Any, List
from rasa_sdk import Action, Tracker, FormValidationAction
from rasa_sdk.executor import CollectingDispatcher
from rasa_sdk.events import AllSlotsReset, EventType
from rasa_sdk.types import DomainDict

class ActionDefaultFallback(Action):
    def name(self) -> Text:
        return "action_default_fallback"

    def run(
        self,
        dispatcher: CollectingDispatcher,
        tracker: Tracker,
        domain: Dict[Text, Any],
    ) -> List[EventType]:
        dispatcher.utter_message(template="utter_fallback")
        return []

class ActionSubmitManifiesto(Action):
    """Acción que se ejecuta cuando se completa el formulario de manifiesto."""
    
    def name(self) -> Text:
        return "action_submit_manifiesto"

    def run(
        self,
        dispatcher: CollectingDispatcher,
        tracker: Tracker,
        domain: Dict[Text, Any],
    ) -> List[EventType]:
        
        # Obtener los datos del formulario
        flete = tracker.get_slot("flete")
        descripcion = tracker.get_slot("descripcion")
        peso = tracker.get_slot("peso")
        fecha_cargue = tracker.get_slot("fecha_cargue")
        fecha_descargue = tracker.get_slot("fecha_descargue")
        tarjeta = tracker.get_slot("tarjeta")
        licencia = tracker.get_slot("licencia")
        origen = tracker.get_slot("origen")
        destino = tracker.get_slot("destino")
        
        # Log para debug (opcional)
        print(f"📋 Manifiesto solicitado:")
        print(f"  💰 Flete: {flete}")
        print(f"  📦 Descripción: {descripcion}")
        print(f"  ⚖️ Peso: {peso}")
        print(f"  📅 Fecha cargue: {fecha_cargue}")
        print(f"  📅 Fecha descargue: {fecha_descargue}")
        print(f"  🚗 Tarjeta: {tarjeta}")
        print(f"  🪪 Licencia: {licencia}")
        print(f"  📍 Origen: {origen}")
        print(f"  🎯 Destino: {destino}")
        
        # Aquí podrías hacer la llamada a tu API o Playwright
        # Por ahora solo confirmamos que recibimos los datos
        
        dispatcher.utter_message(
            text=f"✅ Perfecto! He recibido todos los datos para el manifiesto:\n\n"
                 f"💰 Flete: {flete}\n"
                 f"📦 Carga: {descripcion}\n"
                 f"⚖️ Peso: {peso}\n"
                 f"📅 Cargue: {fecha_cargue}\n"
                 f"📅 Descarga: {fecha_descargue}\n"
                 f"🚗 Tarjeta: {tarjeta}\n"
                 f"🪪 Licencia: {licencia}\n"
                 f"📍 Origen: {origen}\n"
                 f"🎯 Destino: {destino}\n\n"
                 f"🔄 Procesando manifiesto..."
        )
        
        # Limpiar los slots después de procesar
        return [AllSlotsReset()]

class ValidateManifiestoForm(FormValidationAction):
    """Validador para el formulario de manifiesto."""
    
    def name(self) -> Text:
        return "validate_manifiesto_form"

    def validate_flete(
        self,
        slot_value: Any,
        dispatcher: CollectingDispatcher,
        tracker: Tracker,
        domain: DomainDict,
    ) -> Dict[Text, Any]:
        """Valida el slot de flete."""
        if slot_value is None:
            return {"flete": None}
        
        # Limpiar el valor (remover símbolos)
        clean_value = str(slot_value).replace("$", "").replace(",", "").replace(".", "").strip()
        
        # Verificar si es numérico
        if clean_value.isdigit():
            return {"flete": clean_value}
        else:
            dispatcher.utter_message(text="❌ El flete debe ser un valor numérico. Por ejemplo: 150000 o $150,000")
            return {"flete": None}

    def validate_peso(
        self,
        slot_value: Any,
        dispatcher: CollectingDispatcher,
        tracker: Tracker,
        domain: DomainDict,
    ) -> Dict[Text, Any]:
        """Valida el slot de peso."""
        if slot_value is None:
            return {"peso": None}
        
        # Aceptar formato con "kg" o solo números
        clean_value = str(slot_value).lower().replace("kg", "").replace("kilos", "").strip()
        
        try:
            # Verificar si es numérico (puede tener decimales)
            float(clean_value)
            return {"peso": slot_value}  # Mantener formato original
        except ValueError:
            dispatcher.utter_message(text="❌ El peso debe ser un valor numérico. Por ejemplo: 500 kg o 1.5 toneladas")
            return {"peso": None}