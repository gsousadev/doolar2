#!/bin/bash
set -e

echo "Starting official Mongo entrypoint..."
/usr/local/bin/docker-entrypoint.sh "$@" &

MONGO_PID=$!

echo "Waiting for Mongo to accept connections..."
until mongosh --eval "db.adminCommand('ping')" >/dev/null 2>&1; do
  sleep 1
done

echo "Waiting for root user to be created..."
until mongosh -u "$MONGO_INITDB_ROOT_USERNAME" \
              -p "$MONGO_INITDB_ROOT_PASSWORD" \
              --authenticationDatabase admin \
              --eval "db.runCommand({connectionStatus:1})" \
              >/dev/null 2>&1; do
  sleep 1
done

echo "Root user is ready. Checking replica set..."

RS_OK=$(mongosh -u "$MONGO_INITDB_ROOT_USERNAME" \
                -p "$MONGO_INITDB_ROOT_PASSWORD" \
                --authenticationDatabase admin \
                --quiet \
                --eval "rs.status().ok" || echo "0")

if [ "$RS_OK" != "1" ]; then
  echo "Initializing replica set..."
  mongosh -u "$MONGO_INITDB_ROOT_USERNAME" \
          -p "$MONGO_INITDB_ROOT_PASSWORD" \
          --authenticationDatabase admin \
          --eval "rs.initiate({_id:'rs0', members:[{_id:0, host:'db:27017'}]})"
else
  echo "Replica set already initialized"
fi

wait $MONGO_PID
