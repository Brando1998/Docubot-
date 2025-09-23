# 🤖 Rasa Bot - Docubot

Este módulo contiene el asistente conversacional inteligente desarrollado con **Rasa** para interactuar con usuarios a través de WhatsApp (usando Baileys) y ejecutar acciones automáticas como la expedición de manifiestos de carga mediante Playwright.

---

## 📁 Estructura del Proyecto

rasa-bot/
├── actions/ # Acciones personalizadas en Python
│ └── actions.py
├── data/ # Datos de entrenamiento
│ ├── nlu.yml # Intents y ejemplos de usuario
│ ├── rules.yml # Reglas de comportamiento
│ └── stories.yml # Historias de conversación
├── models/ # Modelos entrenados
├── config.yml # Configuración del pipeline y políticas
├── domain.yml # Dominio: intents, slots, entities, actions, responses
├── endpoints.yml # Conexión al servidor de acciones
└── credentials.yml # Canales de entrada (no usado aquí)


---

## 🚀 Cómo desplegar el bot

### Requisitos

- Python 3.8–3.10 recomendado
- pipenv o virtualenv
- Rasa 3.1+
- [Playwright](https://playwright.dev/) (si usas acciones que lo requieren)

### Paso a paso

```bash
# 1. Clona el repositorio
git clone https://github.com/tu-usuario/docubot.git
cd docubot/rasa-bot

# 2. Crea entorno virtual
python -m venv venv
source venv/bin/activate  # Windows: venv\Scripts\activate

# 3. Instala dependencias
pip install -U pip
pip install rasa
pip install -r requirements.txt  # (si existe)

🔧 Entrenar el modelo

Cada vez que modifiques los intents, stories, reglas o respuestas, debes reentrenar el modelo:

rasa train

El modelo entrenado se guardará en la carpeta /models.

▶️ Correr el servidor de Rasa y acciones
1. Levanta el servidor de acciones (acciones personalizadas):

rasa run actions

Esto ejecutará el archivo actions/actions.py y estará disponible en http://localhost:5055/webhook.
2. En otro terminal, levanta el bot:

rasa run --enable-api --cors "*" --port 5005

Este comando arranca el bot Rasa con su API en http://localhost:5005.


🔄 Cómo modificar el bot
Intents y ejemplos de usuario

Edita el archivo data/nlu.yml. Por ejemplo:

- intent: solicitar_manifiesto
  examples: |
    - necesito un manifiesto
    - quiero expedir un manifiesto

Historias

Edita data/stories.yml para definir flujos de conversación.
Slots y acciones personalizadas

Define slots y acciones en domain.yml y en actions/actions.py.
Respuestas automáticas

Modifica las respuestas en domain.yml dentro del bloque responses:.
🧪 Probar el bot

Puedes usar Rasa Shell para hacer pruebas rápidas:

rasa shell

🔌 Integración con WhatsApp

Este bot está diseñado para recibir mensajes desde WhatsApp a través del servicio baileys-ws. Los mensajes se reenvían a la API central, que a su vez reenvía a Rasa.

La API espera que Rasa devuelva una respuesta en el siguiente formato:

{
  "recipient_id": "usuario",
  "text": "Mensaje del bot"
}

🛠️ Personalización

    Agrega nuevos intents: edita nlu.yml y domain.yml

    Agrega nuevas acciones: edita actions/actions.py y domain.yml

    Ajusta las preguntas: modifica los forms en domain.yml

    Entrena: rasa train

🧹 Limpieza y mantenimiento

    Elimina modelos viejos: rm -rf models/*

    Reentrena con cambios: rasa train

    Verifica reglas con: rasa data validate

🧾 Licencia

Este proyecto es parte de Docubot y su uso está regulado por la licencia del repositorio principal.


---


D3v_S1C0P2025+*