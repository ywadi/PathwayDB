#!/bin/bash

# PathwayDB IDE Startup Script
set -e

# Change to the directory where the script is located
cd "$(dirname "$0")"

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
if [ ! -d "ide/frontend/node_modules" ]; then
    echo "📦 Installing frontend dependencies..."
    (cd ide/frontend && npm install)
fi

# Install backend dependencies
echo "📦 Installing backend dependencies..."
(cd ide/backend && go mod tidy)

# Set default ports if not provided by environment variables
export REDIS_ADDR=${REDIS_ADDR:-":6379"}
export WEBSOCKET_ADDR=${WEBSOCKET_ADDR:-":8081"}
export PORT=${PORT:-"3000"}

# Construct the full WebSocket URL for the frontend
WEBSOCKET_HOST=$(echo $WEBSOCKET_ADDR | cut -d: -f1)
if [ -z "$WEBSOCKET_HOST" ]; then
    WEBSOCKET_HOST="localhost"
fi
WEBSOCKET_PORT_NUM=$(echo $WEBSOCKET_ADDR | cut -d: -f2)
export REACT_APP_WEBSOCKET_URL="ws://${WEBSOCKET_HOST}:${WEBSOCKET_PORT_NUM}/ws"
export REACT_APP_API_BASE_URL="http://${WEBSOCKET_HOST}:${WEBSOCKET_PORT_NUM}"

# Start Redis server in background
echo "🔴 Starting PathwayDB Redis server..."
if [ -n "$PROD" ]; then
    ./cmd/redis-server/redis-server &
else
    go run cmd/redis-server/main.go &
fi

# Wait for Redis to start
sleep 2

# Start backend WebSocket bridge in background
echo "🌐 Starting WebSocket bridge server..."
if [ -n "$PROD" ]; then
    (cd ide/backend && ./backend-server) &
else
    (cd ide/backend && go run main.go) &
fi

# Wait for backend to start
sleep 2

# Start frontend development server if not in production
if [ -z "$PROD" ]; then
    echo "⚛️  Starting React development server..."
    (cd ide/frontend && npm start) &
fi

echo "✅ PathwayDB IDE is starting up!"
echo ""
echo "🔗 Access the IDE at: http://localhost:$PORT"
echo "🔴 Redis server running on: $REDIS_ADDR"
echo "🌐 WebSocket bridge running on: $WEBSOCKET_ADDR"
echo ""
echo "Press Ctrl+C to stop all services"

# Function to cleanup background processes
cleanup() {
    echo ""
    echo "🛑 Stopping PathwayDB IDE services..."
    # Use pkill to reliably find and kill processes by their command line.
    pkill -f "cmd/redis-server/main.go" 2>/dev/null || true
    pkill -f "backend/main.go" 2>/dev/null || true
    pkill -f "react-scripts start" 2>/dev/null || true
    sleep 1 # Give processes a moment to shut down
    echo "✅ Services stopped."
    exit 0
}

# Set trap to cleanup on script exit
trap cleanup SIGINT SIGTERM

# Wait for user to stop
wait
