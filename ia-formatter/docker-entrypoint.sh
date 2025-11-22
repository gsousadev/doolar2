#!/bin/sh

ollama serve & sleep 5

ollama pull deepseek-r1

wait