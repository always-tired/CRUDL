# Subscription Aggregator Service

## Requirements
- Go 1.24
- Docker + Docker Compose

## Quick start (Docker)
```bash
docker compose --env-file .env up --build
```

- API: http://localhost:8080
- Swagger UI: http://localhost:8080/swagger-ui.html
- Swagger JSON: http://localhost:8080/swagger/doc.json

## Configuration (ENV only)
The service reads configuration from environment variables.

- `ENV`
- `HTTP_PORT`
- `HTTP_READ_TIMEOUT`
- `HTTP_WRITE_TIMEOUT`
- `HTTP_IDLE_TIMEOUT`
- `DB_URL`
- `LOG_LEVEL`

Environment template: `.env.example`

## Migrations (goose)
The `migrate` service runs goose against the database on startup.

## Endpoints
- `POST /subscriptions`
- `GET /subscriptions/{id}`
- `PUT /subscriptions/{id}`
- `DELETE /subscriptions/{id}`
- `GET /subscriptions`
- `GET /subscriptions/summary?start=MM-YYYY&end=MM-YYYY&user_id=&service_name=`

## Sample request
```bash
curl -X POST http://localhost:8080/subscriptions \
  -H 'Content-Type: application/json' \
  -d '{
    "service_name": "Yandex Plus",
    "price": 400,
    "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
    "start_date": "07-2025"
  }'
```
