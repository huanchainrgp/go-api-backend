#!/bin/bash

# Go API Test1 - Quick Start Script
# This script handles Docker PATH issues and provides easy startup options

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_header() {
    echo -e "${BLUE}================================${NC}"
    echo -e "${BLUE} $1${NC}"
    echo -e "${BLUE}================================${NC}"
}

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to setup Docker PATH
setup_docker_path() {
    if [ -d "/Applications/Docker.app/Contents/Resources/bin" ]; then
        export PATH="/Applications/Docker.app/Contents/Resources/bin:$PATH"
        print_status "Docker PATH configured"
        return 0
    else
        print_warning "Docker.app not found in /Applications/"
        return 1
    fi
}

# Function to get docker compose command
get_docker_compose_cmd() {
    # Try modern docker compose first
    if command_exists docker && docker compose version >/dev/null 2>&1; then
        echo "docker compose"
        return 0
    fi
    
    # Try legacy docker-compose
    if command_exists docker-compose; then
        echo "docker-compose"
        return 0
    fi
    
    # Try with Docker Desktop path
    if [ -f "/Applications/Docker.app/Contents/Resources/bin/docker" ]; then
        export PATH="/Applications/Docker.app/Contents/Resources/bin:$PATH"
        if docker compose version >/dev/null 2>&1; then
            echo "docker compose"
            return 0
        elif command_exists docker-compose; then
            echo "docker-compose"
            return 0
        fi
    fi
    
    return 1
}

# Function to check Docker status
check_docker() {
    if ! command_exists docker; then
        print_warning "Docker CLI not found, attempting to configure..."
        if setup_docker_path; then
            if ! command_exists docker; then
                print_error "Docker CLI still not available. Please ensure Docker Desktop is running."
                return 1
            fi
        else
            print_error "Could not configure Docker PATH. Please install Docker Desktop."
            return 1
        fi
    fi
    
    # Check if Docker daemon is running
    if ! docker info >/dev/null 2>&1; then
        print_error "Docker daemon is not running. Please start Docker Desktop."
        return 1
    fi
    
    print_status "Docker is ready"
    return 0
}

# Function to start development environment
start_dev() {
    print_header "Starting Development Environment"
    
    # Check if .env exists
    if [ ! -f .env ]; then
        print_warning ".env file not found. Creating from template..."
        if [ -f env.example ]; then
            cp env.example .env
            print_status "Created .env file from template"
        else
            print_error "env.example file not found"
            exit 1
        fi
    fi
    
    # Check Docker
    if ! check_docker; then
        print_error "Docker is not available. Please install and start Docker Desktop."
        exit 1
    fi
    
    # Get docker compose command
    DOCKER_COMPOSE_CMD=$(get_docker_compose_cmd)
    if [ $? -ne 0 ]; then
        print_error "Docker Compose not available. Please ensure Docker Desktop is installed and running."
        exit 1
    fi
    
    # Start development services
    print_status "Starting development services with Docker Compose..."
    $DOCKER_COMPOSE_CMD -f docker-compose.dev.yml up -d
    
    print_status "Development environment started!"
    print_status "PostgreSQL: localhost:5433"
    print_status "pgAdmin: http://localhost:5050 (admin@example.com / admin)"
    print_status ""
    print_status "To start the Go API server, run:"
    print_status "  go run main.go"
    print_status ""
    print_status "Or use the full dev script:"
    print_status "  ./dev.sh run"
}

# Function to start production environment
start_prod() {
    print_header "Starting Production Environment"
    
    # Check if .env exists
    if [ ! -f .env ]; then
        print_warning ".env file not found. Creating from template..."
        if [ -f env.example ]; then
            cp env.example .env
            print_status "Created .env file from template"
        else
            print_error "env.example file not found"
            exit 1
        fi
    fi
    
    # Check Docker
    if ! check_docker; then
        print_error "Docker is not available. Please install and start Docker Desktop."
        exit 1
    fi
    
    # Get docker compose command
    DOCKER_COMPOSE_CMD=$(get_docker_compose_cmd)
    if [ $? -ne 0 ]; then
        print_error "Docker Compose not available. Please ensure Docker Desktop is installed and running."
        exit 1
    fi
    
    # Start production services
    print_status "Starting production services with Docker Compose..."
    $DOCKER_COMPOSE_CMD up -d
    
    print_status "Production environment started!"
    print_status "API: http://localhost:8080"
    print_status "PostgreSQL: localhost:5432"
    print_status "Swagger UI: http://localhost:8080/swagger/index.html"
}

# Function to stop all services
stop_all() {
    print_header "Stopping All Services"
    
    # Check Docker
    if ! check_docker; then
        print_error "Docker is not available"
        exit 1
    fi
    
    # Get docker compose command
    DOCKER_COMPOSE_CMD=$(get_docker_compose_cmd)
    if [ $? -ne 0 ]; then
        print_error "Docker Compose not available"
        exit 1
    fi
    
    print_status "Stopping development services..."
    $DOCKER_COMPOSE_CMD -f docker-compose.dev.yml down 2>/dev/null || true
    
    print_status "Stopping production services..."
    $DOCKER_COMPOSE_CMD down 2>/dev/null || true
    
    print_status "All services stopped!"
}

# Function to show status
show_status() {
    print_header "Service Status"
    
    if ! check_docker; then
        print_error "Docker is not available"
        exit 1
    fi
    
    # Get docker compose command
    DOCKER_COMPOSE_CMD=$(get_docker_compose_cmd)
    if [ $? -ne 0 ]; then
        print_error "Docker Compose not available"
        exit 1
    fi
    
    print_status "Development services:"
    $DOCKER_COMPOSE_CMD -f docker-compose.dev.yml ps 2>/dev/null || print_warning "No development services running"
    
    print_status ""
    print_status "Production services:"
    $DOCKER_COMPOSE_CMD ps 2>/dev/null || print_warning "No production services running"
}

# Function to show help
show_help() {
    echo "Go API Test1 - Quick Start Script"
    echo ""
    echo "Usage: $0 [command]"
    echo ""
    echo "Commands:"
    echo "  dev       - Start development environment (PostgreSQL + pgAdmin)"
    echo "  prod      - Start production environment (API + PostgreSQL)"
    echo "  stop      - Stop all services"
    echo "  status    - Show service status"
    echo "  help      - Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 dev     # Start development database"
    echo "  $0 prod    # Start full production stack"
    echo "  $0 stop    # Stop everything"
    echo ""
    echo "After starting services, run the API with:"
    echo "  go run main.go"
    echo ""
    echo "Or use the full dev script:"
    echo "  ./dev.sh run"
}

# Main script logic
case "${1:-help}" in
    dev)
        start_dev
        ;;
    prod)
        start_prod
        ;;
    stop)
        stop_all
        ;;
    status)
        show_status
        ;;
    help|--help|-h)
        show_help
        ;;
    *)
        print_error "Unknown command: $1"
        show_help
        exit 1
        ;;
esac
