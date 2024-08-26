# Go RabbitMQ/SQLite

Installation:
```bash
go mod tidy
```

Running the app:

- Fill environment variables in a `.env` file as described in `.env.example`
- Execute RabbitMQ with docker
```bash
docker-compose up -d
```
- Run the app:
```bash
go run cmd/main.go
```

The application will create a sqlite file in the folder specified in the .env file and create a table (`messages`) to store the messages from RabbitMQ.

Build:

- Windows:
```bash
make build-windows
```

- Linux/mac
```bash
make build
```
