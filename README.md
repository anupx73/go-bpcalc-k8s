# BP Calc Microservice

## Overview

This project demonstrates the use of microservice and its deployment to kubernetes.  
This backend service is powered by Go which provides some endpoints and interacts with MongoDB.

 * BP Calc Service: Provides blood pressure category calculation based on inputs.

The Go project template is sourced from: [mmorejon](https://github.com/mmorejon/microservices-docker-go-mongodb)

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
