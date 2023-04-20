# BP Calc Microservice

## Overview

This project demonstrates the use of microservice and its deployment to kubernetes.  
This backend service is powered by Go which provides some endpoints and interacts with MongoDB.

 * BP Calc Service: Provides blood pressure category calculation based on inputs.

The Go project template is sourced from: [mmorejon](https://github.com/mmorejon/microservices-docker-go-mongodb)

## Go Module Publish

To publish this module to be used by the frontend service use the followings:

```
go mod tidy

git tag v0.1.0
git push origin v0.1.0

GOPROXY=proxy.golang.org go list -m github.com/anupx73/go-bpcalc-backend-k8s@v0.1.0
```

## Endpoints

| Service | Method | Endpoint       |
|---------|--------|----------------|
| List BP Readings | `GET` | `/api/bpcalc/` |
| Get BP Readings by Id | `GET` | `/api/bpcalc/{id}` |
| Insert BP Reading | `POST` | `/api/bpcalc/` |
| Delete BP Reading | `DELETE` | `/api/bpcalc/{id}` |

### POST 

```
curl  -X POST http://localhost/api/bpcalc/ \
      -H "Content-Type: application/json" \
      -d '{"name":"Steven A","email":"steven.a@domain.com","systolic":120,"diastolic":80}'
```
