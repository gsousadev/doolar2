import os
import uvicorn
import tempfile
from contextlib import asynccontextmanager
from fastapi import FastAPI, UploadFile, File, HTTPException
from fastapi.responses import JSONResponse
from faster_whisper import WhisperModel

# ---------------------------------------------------------
# Configurações via env
# ---------------------------------------------------------
MODEL = os.getenv("WHISPER_MODEL", "small")
LANG = os.getenv("WHISPER_LANGUAGE", "pt")
DEVICE = os.getenv("WHISPER_DEVICE", "cpu")

ALLOWED_MODELS = ["tiny", "base", "small"]

if MODEL not in ALLOWED_MODELS:
    raise ValueError(f"Modelo '{MODEL}' não permitido. Use: {', '.join(ALLOWED_MODELS)}")

compute_type = "int8"

model = None  # será carregado no lifespan


# ---------------------------------------------------------
# Lifespan — substitui on_event("startup")
# ---------------------------------------------------------
@asynccontextmanager
async def lifespan(app: FastAPI):
    global model
    print(f"Carregando modelo: {MODEL} (device={DEVICE})")
    model = WhisperModel(MODEL, device=DEVICE, compute_type=compute_type)
    yield
    # Aqui você pode liberar recursos (shutdown)
    print("Encerrando servidor Whisper ASR...")


app = FastAPI(title="Whisper ASR Server", lifespan=lifespan)


# ---------------------------------------------------------
# Endpoint HTTP
# ---------------------------------------------------------

@app.get("/")
async def root():
    return JSONResponse({"message": "Servidor Whisper ASR ativo"})


@app.post("/transcribe")
async def transcribe_audio(audio: UploadFile = File(...)):
    if not audio.filename:
        raise HTTPException(status_code=400, detail="Nenhum arquivo enviado")

    allowed_ext = [
        ".wav", ".mp3", ".m4a", ".ogg", ".flac", ".webm"
    ]

    if not any(audio.filename.lower().endswith(ext) for ext in allowed_ext):
        raise HTTPException(status_code=400, detail="Tipo de arquivo não suportado")

    try:
        # Salva o arquivo temporário
        with tempfile.NamedTemporaryFile(delete=True, suffix=audio.filename) as tmp:
            tmp.write(await audio.read())
            tmp.flush()

            segments, info = model.transcribe(tmp.name, language=LANG)

            text = " ".join([seg.text for seg in segments]).strip()

        return JSONResponse({"text": text})

    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


# ---------------------------------------------------------
# Execução
# ---------------------------------------------------------
if __name__ == "__main__":
    uvicorn.run(
        "main:app",
        host="0.0.0.0",
        port=int(os.getenv("PORT", 8000)),
        reload=True,
        reload_dirs=["/app"],
    )
