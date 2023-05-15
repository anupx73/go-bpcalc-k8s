[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=go-bp-calc-k8s&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=go-bp-calc-k8s)

# BP Calc Backend Microservice

## Overview

This backend service is written in Go providing the following endpoints and interacts with MongoDB. The main purpose of this project remains microservice deployment to kubernetes environment and establish ci/cd pipeline.

| Service | Method | Endpoint       |
|---------|--------|----------------|
| List BP Readings | `GET` | `/api/bpcalc/` |
| Get BP Readings by Id | `GET` | `/api/bpcalc/{id}` |
| Insert BP Reading | `POST` | `/api/bpcalc/` |
| Delete BP Reading | `DELETE` | `/api/bpcalc/{id}` |

Currently due to GCE ingress controller rule resolving [issue](https://www.googlecloudcommunity.com/gc/Google-Kubernetes-Engine-GKE/GCE-ingress-to-route-traffic-to-multiple-services/m-p/551562#M696) this app is limited to serve at root path, hence the above endpoints are for reference only. The routing code in `routes.go` is commented for future use and shall be experimented with Nginx ingress controller/rules to access multiple services using a single LoadBalancer IP.

## Useful Commands

```
# build
go build ./...
go run ./...

# go module publish (was used when experimental go frontend was created)
go mod tidy
git tag v0.1.0
git push origin v0.1.0
GOPROXY=proxy.golang.org go list -m github.com/anupx73/go-bpcalc-backend-k8s@v0.1.0

# local image build
docker build . --file Dockerfile --tag backend:v1.0.99-security;
docker tag backend:v1.0.99-security gcr.io/tudublin/backend:v1.0.99-security;
docker push gcr.io/tudublin/backend:v1.0.99-security

# deployment
helm upgrade backend helm/ --install --namespace ns-backend --create-namespace --wait

# api testing
curl -X POST http://backend-service/api/bpcalc/ -H "Content-Type: application/json" -d '{"name":"Steven A","email":"steven.a@domain.com","systolic":"120","diastolic":"80"}'
```

## Miscellaneous  

**1. Vault Helm Issue**  
The following annotation in deployment.yaml did not work and thrown. 
*Error: INSTALLATION FAILED: parse error at (backend-chart/templates/deployment.yaml:20): function "secret" not defined*

`vault.hashicorp.com/agent-inject-template-database-config.txt: |
  {{- with secret "internal/data/database/config" -}}
  mongodb+srv://{{ .Data.data.username }}:{{ .Data.data.password }}@{{ .Data.data.url }}/?retryWrites=true&w=majority
  {{- end -}}`

This seems to be an issue from vault-helm. Ref: [issue-853](https://github.com/hashicorp/vault-helm/issues/853). Currently, as a workaround `main.go` is parsing the raw Vault file in a funny way to extract database credentials.
