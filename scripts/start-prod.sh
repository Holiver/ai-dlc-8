#!/bin/bash

# Start production environment
echo "Starting AWSomeShop production environment..."

# Check if .env file exists
if [ ! -f .env ]; then
    echo "Error: .env file not found!"
    echo "Please create .env file from .env.example and configure it"
    exit 1
fi

# Build and start all services
echo "Building and starting services..."
docker-compose --profile production up -d --build

echo "Waiting for services to be ready..."
sleep 15

# Check service health
echo "Checking service health..."
docker-compose ps

echo ""
echo "Production environment started!"
echo "Application: https://localhost"
echo "API: https://localhost/api"
echo ""
echo "To view logs: docker-compose logs -f"
echo "To stop: docker-compose --profile production down"
