#!/bin/bash

# Stop all services
echo "Stopping AWSomeShop services..."

# Stop Docker containers
docker-compose --profile production down

echo "All services stopped!"
