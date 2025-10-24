# Go API Test1

A scalable Go backend API with Users, Assets, and Transactions management.

## Features

- **User Management**: Register, login, and manage user accounts
- **Asset Management**: Create and manage assets (cryptocurrencies, stocks, etc.)
- **Transaction Management**: Handle buy/sell/transfer transactions
- **JWT Authentication**: Secure API endpoints with JWT tokens
- **Swagger Documentation**: Interactive API documentation
- **Database Support**: PostgreSQL (production) and SQLite (development)
- **Docker Support**: Easy deployment with Docker

## Tech Stack

- **Framework**: Gin (HTTP web framework)
- **ORM**: GORM (Object-Relational Mapping)
- **Database**: PostgreSQL / SQLite
- **Authentication**: JWT (JSON Web Tokens)
- **Documentation**: Swagger/OpenAPI
- **Password Hashing**: bcrypt

## Quick Start

### Prerequisites

- Go 1.21 or higher
- PostgreSQL (optional, SQLite is used by default)

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd go-api-test1
```

2. Install dependencies:
```bash
go mod tidy
```

3. Set up environment variables:
```bash
cp env.example .env
# Edit .env with your configuration
```

4. Run the application:
```bash
go run main.go
```

The API will be available at `http://localhost:8080`

### Swagger Documentation

Once the server is running, visit `http://localhost:8080/swagger/index.html` to access the interactive API documentation.

## API Endpoints

### Authentication
- `POST /api/v1/auth/register` - Register a new user
- `POST /api/v1/auth/login` - Login user

### Users (Protected)
- `GET /api/v1/users` - Get all users
- `GET /api/v1/users/{id}` - Get user by ID
- `PUT /api/v1/users/{id}` - Update user
- `DELETE /api/v1/users/{id}` - Delete user

### Assets (Protected)
- `GET /api/v1/assets` - Get all assets
- `GET /api/v1/assets/{id}` - Get asset by ID
- `POST /api/v1/assets` - Create new asset
- `PUT /api/v1/assets/{id}` - Update asset
- `DELETE /api/v1/assets/{id}` - Delete asset

### Transactions (Protected)
- `GET /api/v1/transactions` - Get all transactions
- `GET /api/v1/transactions/{id}` - Get transaction by ID
- `POST /api/v1/transactions` - Create new transaction
- `PUT /api/v1/transactions/{id}` - Update transaction
- `DELETE /api/v1/transactions/{id}` - Delete transaction

## Database Models

### User
- ID, Email, Username, Password (hashed)
- FirstName, LastName, IsActive
- CreatedAt, UpdatedAt, DeletedAt

### Asset
- ID, Name, Symbol, Type, Description
- Price, IsActive
- CreatedAt, UpdatedAt, DeletedAt

### Transaction
- ID, UserID, AssetID, Type (buy/sell/transfer)
- Amount, Price, TotalValue, Status
- Description, CreatedAt, UpdatedAt, DeletedAt
- Relationships: User, Asset

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DATABASE_URL` | Database connection string | SQLite file |
| `JWT_SECRET` | Secret key for JWT tokens | Required |
| `PORT` | Server port | 8080 |
| `ENVIRONMENT` | Environment (development/production) | development |

## Docker Support

### Build and run with Docker:

```bash
# Build the image
docker build -t go-api-test1 .

# Run the container
docker run -p 8080:8080 --env-file .env go-api-test1
```

### Docker Compose:

```bash
docker-compose up -d
```

## Development

### Project Structure

```
go-api-test1/
├── main.go                 # Application entry point
├── go.mod                  # Go module file
├── internal/
│   ├── config/            # Configuration management
│   ├── database/          # Database connection and setup
│   ├── handlers/          # HTTP request handlers
│   ├── middleware/        # HTTP middleware
│   └── models/            # Data models and DTOs
├── docs/                  # Swagger documentation (generated)
├── Dockerfile             # Docker configuration
├── docker-compose.yml     # Docker Compose configuration
└── README.md              # This file
```

### Adding New Features

1. Define models in `internal/models/`
2. Create handlers in `internal/handlers/`
3. Add routes in `main.go`
4. Update Swagger documentation
5. Test with the Swagger UI

## Security Considerations

- Passwords are hashed using bcrypt
- JWT tokens are used for authentication
- CORS is configured for cross-origin requests
- Input validation is implemented
- SQL injection protection via GORM

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License.
