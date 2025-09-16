#!/bin/bash

# PathwayDB IDE Startup Script
set -e

echo "🚀 Starting PathwayDB IDE..."

# Check if Node.js is installed
if ! command -v node &> /dev/null; then
    echo "❌ Node.js is not installed. Please install Node.js 16+ and try again."
    exit 1
fi

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21+ and try again."
    exit 1
fi

# Install frontend dependencies if needed
if [ ! -d "frontend/node_modules" ]; then
    echo "📦 Installing frontend dependencies..."
    cd frontend
    npm install
    cd ..
fi

# Install backend dependencies
echo "📦 Installing backend dependencies..."
cd backend
go mod tidy
cd ..

# Start Redis server in background
echo "🔴 Starting PathwayDB Redis server..."
cd ..
go run cmd/redis-server/main.go &
REDIS_PID=$!
cd ide

# Wait for Redis to start
sleep 2

# Start backend WebSocket bridge in background
echo "🌐 Starting WebSocket bridge server..."
cd backend
go run main.go &
BACKEND_PID=$!
cd ..

# Wait for backend to start
sleep 2

# Start frontend development server
echo "⚛️  Starting React development server..."
cd frontend
npm start &
FRONTEND_PID=$!
cd ..

echo "✅ PathwayDB IDE is starting up!"
echo ""
echo "🔗 Access the IDE at: http://localhost:3000"
echo "🔴 Redis server running on: localhost:6379"
echo "🌐 WebSocket bridge running on: localhost:8081"
echo ""
echo "Press Ctrl+C to stop all services"

# Function to cleanup background processes
cleanup() {
    echo ""
    echo "🛑 Stopping PathwayDB IDE services..."
    kill $REDIS_PID 2>/dev/null || true
    kill $BACKEND_PID 2>/dev/null || true
    kill $FRONTEND_PID 2>/dev/null || true
    exit 0
}

# Set trap to cleanup on script exit
trap cleanup SIGINT SIGTERM

# Wait for user to stop
wait
