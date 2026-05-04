#!/bin/bash

echo "🚀 Initiating build..."
docker build -t tech-case-api .

echo "📦 Running server..."
docker run --rm -p 8080:8080 tech-case-api