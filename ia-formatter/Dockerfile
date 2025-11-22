FROM ollama/ollama:latest

# Copie o Modelfile para dentro da imagem
COPY Modelfile /root/.ollama/Modelfile

COPY docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh
RUN chmod +x /usr/local/bin/docker-entrypoint.sh

ENTRYPOINT ["docker-entrypoint.sh"]