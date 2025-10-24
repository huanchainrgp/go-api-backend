#!/bin/bash

# Go API Test1 - Development Script

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

# Function to setup development environment
setup_dev() {
    print_header "Setting up Development Environment"
    
    # Check if Go is installed
    if ! command_exists go; then
        print_error "Go is not installed. Please install Go 1.21 or higher."
        exit 1
    fi
    
    # Check Go version
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    print_status "Go version: $GO_VERSION"
    
    # Create .env file if it doesn't exist
    if [ ! -f .env ]; then
        cp env.example .env
        print_status "Created .env file from template"
    else
        print_warning ".env file already exists"
    fi
    
    # Install dependencies
    print_status "Installing Go dependencies..."
    go mod tidy
    
    # Install Swagger CLI if not present
    if ! command_exists swag; then
        print_status "Installing Swagger CLI..."
        go install github.com/swaggo/swag/cmd/swag@latest
    fi
    
    # Generate Swagger documentation
    print_status "Generating Swagger documentation..."
    swag init -g main.go
    
    print_status "Development environment setup complete!"
    print_status "Run './dev.sh run' to start the development server"
}

# Function to run the application
run_app() {
    print_header "Starting Go API Test1"
    
    if [ ! -f .env ]; then
        print_warning ".env file not found. Creating from template..."
        cp env.example .env
    fi
    
    print_status "Starting server on http://localhost:8080"
    print_status "Swagger UI available at: http://localhost:8080/swagger/index.html"
    
    go run main.go
}

# Function to run tests
run_tests() {
    print_header "Running Tests"
    go test ./...
}

# Function to build the application
build_app() {
    print_header "Building Application"
    go build -o bin/main main.go
    print_status "Binary built successfully: bin/main"
}

# Function to run with Docker
run_docker() {
    print_header "Running with Docker"
    
    if ! command_exists docker; then
        print_error "Docker is not installed. Please install Docker."
        exit 1
    fi
    
    print_status "Building Docker image..."
    docker build -t go-api-test1 .
    
    print_status "Running Docker container..."
    docker run -p 8080:8080 --env-file .env go-api-test1
}

# Function to run with Docker Compose
run_docker_compose() {
    print_header "Running with Docker Compose"
    
    if ! command_exists docker-compose; then
        print_error "Docker Compose is not installed. Please install Docker Compose."
        exit 1
    fi
    
    print_status "Starting services with Docker Compose..."
    docker-compose up -d
    
    print_status "Services started!"
    print_status "API: http://localhost:8080"
    print_status "PostgreSQL: localhost:5432"
    print_status "Swagger UI: http://localhost:8080/swagger/index.html"
}

# Function to stop Docker Compose
stop_docker_compose() {
    print_header "Stopping Docker Compose Services"
    docker-compose down
    print_status "Services stopped!"
}

# Function to clean up
cleanup() {
    print_header "Cleaning Up"
    
    # Remove build artifacts
    rm -rf bin/
    go clean
    
    # Remove generated docs
    rm -rf docs/
    
    print_status "Cleanup complete!"
}

# Function to show help
show_help() {
    echo "Go API Test1 - Development Script"
    echo ""
    echo "Usage: $0 [command]"
    echo ""
    echo "Commands:"
    echo "  setup     - Set up development environment"
    echo "  run       - Run the application"
    echo "  test      - Run tests"
    echo "  build     - Build the application"
    echo "  docker    - Run with Docker"
    echo "  compose   - Run with Docker Compose"
    echo "  stop      - Stop Docker Compose services"
    echo "  clean     - Clean up build artifacts"
    echo "  help      - Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 setup     # First time setup"
    echo "  $0 run       # Start development server"
    echo "  $0 compose   # Start with PostgreSQL"
}

# Main script logic
case "${1:-help}" in
    setup)
        setup_dev
        ;;
    run)
        run_app
        ;;
    test)
        run_tests
        ;;
    build)
        build_app
        ;;
    docker)
        run_docker
        ;;
    compose)
        run_docker_compose
        ;;
    stop)
        stop_docker_compose
        ;;
    clean)
        cleanup
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
