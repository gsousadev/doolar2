FROM python:3.12-slim-trixie

ENV PYTHONDONTWRITEBYTECODE=1
ENV PYTHONUNBUFFERED=1

# Instala ffmpeg e dependências essenciais
RUN apt-get update && apt-get install -y --no-install-recommends \
        ffmpeg \
    && rm -rf /var/lib/apt/lists/*

# Instala dependências Python
WORKDIR /app

COPY ./app/requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copia aplicação
COPY ./app .

# Variáveis padrão
ENV WHISPER_MODEL=small
ENV WHISPER_LANGUAGE=pt
ENV WHISPER_DEVICE=cpu
ENV PORT=8000

EXPOSE 8000

CMD ["python", "main.py"]