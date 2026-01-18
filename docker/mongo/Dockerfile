FROM mongo:6

RUN openssl rand -base64 756 > /etc/mongo-keyfile \
    && chown mongodb:mongodb /etc/mongo-keyfile \
    && chmod 600 /etc/mongo-keyfile

COPY docker/mongo/entrypoint.sh /custom-entrypoint.sh
RUN chmod +x /custom-entrypoint.sh

ENTRYPOINT ["/custom-entrypoint.sh"]

CMD ["mongod", "--replSet", "rs0", "--keyFile", "/etc/mongo-keyfile", "--bind_ip_all"]