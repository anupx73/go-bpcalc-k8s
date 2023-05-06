[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=go-bp-calc-k8s&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=go-bp-calc-k8s)

# BP Calc Microservice

## Overview

This backend service is written in Go which provides the following endpoints and interacts with MongoDB.

| Service | Method | Endpoint       |
|---------|--------|----------------|
| List BP Readings | `GET` | `/api/bpcalc/` |
| Get BP Readings by Id | `GET` | `/api/bpcalc/{id}` |
| Insert BP Reading | `POST` | `/api/bpcalc/` |
| Delete BP Reading | `DELETE` | `/api/bpcalc/{id}` |

This Go project template is sourced from: [mmorejon](https://github.com/mmorejon/microservices-docker-go-mongodb) and the main purpose of this project remains to demonstrate microservice deployment to kubernetes.

## Build Commands

```
go build ./...
go run ./...
```

## Publishing Module

```
go mod tidy

git tag v0.1.0
git push origin v0.1.0

GOPROXY=proxy.golang.org go list -m github.com/anupx73/go-bpcalc-backend-k8s@v0.1.0
```

## Testing

```
curl  -X POST http://backend-service/api/bpcalc/ -H "Content-Type: application/json" -d '{"name":"Steven A","email":"steven.a@domain.com","systolic":"120","diastolic":"80"}'
```
