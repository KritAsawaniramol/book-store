<h1 align="center">
  Book-Store Microservice
  <br>
</h1>

<p align="center">
  <a href="#built-with">Built With</a> •
  <a href="#architecture">Architecture</a> •
  <a href="#installation">Installation</a> •
  <a href="#endpoints">Endpoints</a> •
  <a href="#license">License</a>
</p>

## Built With
<hr />

- [![Go][GO.dev]][GO-url]
- [![Gorm][GORM.io]][GORM-url]
- [![Gin-gonic][Gin-badge]][Gin-url]
- [![Postgresql][Postgresql-badge]][Postgresql-url]
- [![Kafka][Kafka-badge]][Kafka-url]
- [![Grpc][Grpc-badge]][Grpc-url]
- [![Stripe][Stripe-badge]][Stripe-url]
- [![Supabase][Supabase-badge]][Supabase-url]
- [![Confluent][Confluent-badge]][Confluent-url]
- [![Google Cloud][GoogleCloud-badge]][GoogleCloud-url]
- [![Docker][Docker-badge]][Docker-url]
- [![Kubernetes][Kubernetes-badge]][Kubernetes-url]

## Architecture
<hr />

<img src="./screenshots/book-store-arch.svg" alt="arch_v3.png">

## Saga pattern for distributed transaction
<img src="./screenshots/sagaPattern2.svg" alt="sagaPattern2.svg">

<img src="./screenshots/sagaPattern.svg" alt="sagaPattern.svg">


## Endpoints
<hr />

### Authentication Service Endpoints

#### Public Endpoints
- `GET /auth_v1/auth` - Health check

#### Authentication Endpoints
- `POST /auth_v1/auth/login`
- `POST /auth_v1/auth/logout`
- `POST /auth_v1/auth/refresh-token`


### Book Service Endpoints

#### Public Endpoints
- `GET /book_v1/book` - Health check
- `GET /book_v1/book/cover/:fileName`
- `GET /book_v1/book`
- `GET /book_v1/book/:id`
- `GET /book_v1/book/tags`

#### Authorized Endpoints (User who bought book or Admin Access)
- `GET /book_v1/book/read/:bookID`

#### Admin-only CRUD Endpoints
- `GET /book_v1/admin/book` 
- `GET /book_v1/admin/book/:id`
- `POST /book_v1/book`
- `PATCH /book_v1/book/:id`
- `PATCH /book_v1/book/cover/:id`
- `PATCH /book_v1/book/file/:id`

### Order Service Endpoints

#### Public Endpoints
- `GET /order_v1/order` - Health check

#### Authorized Endpoints (User Access Only)
- `POST /order_v1/order/buy`
- `GET /order_v1/order/myorder`
- `GET /order_v1/order`

### Shelf Service Endpoints

#### Public Endpoints
- `GET /shelf_v1/shelf` - Health check

#### Authorized Endpoints (User Access Only)
- `GET /shelf_v1/shelf`


### User Service Endpoints

#### Public Endpoints
- `GET /user_v1/user` - Health check
- `POST /user_v1/user/register`
- `POST /user_v1/webhook` - Stripe webhook

#### Authorized Endpoints (User Access Only)
- `POST /user_v1/user/top-up`
- `GET /user_v1/user/top-up/:id`
- `GET /user_v1/user/balance`
- `GET /user_v1/user/profile`

#### Admin-only Endpoints
- `POST /user_v1/user/transaction`
- `GET /user_v1/user/transaction`

## Installation
<hr />

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/book-store.git
   cd book-store
   ```

2. Install dependencies:

   ```bash
   go mod tidy
   ```

3. Set up environment variables (fill in spaces in every .env file in env/dev)
4. Run the project:

   ```bash
   go run main.go ./env/dev/.env.auth
   ```

   ```bash
   go run main.go ./env/dev/.env.item
   ```

   ```bash
   go run main.go ./env/dev/.env.player
   ```

   ```bash
   go run main.go ./env/dev/.env.inventory
   ```

   ```bash
   go run main.go ./env/dev/.env.payment
   ```

## Usage

### Running Tests
<hr />

To run tests, use:

```bash
go test ./...
```

### Running Docker Compose for PostgreSQL, pgAdmin, and Kafka
<hr />

start docker compose for start pgAdmin, Kafka and create PostgreSQL database for each service

Run Docker Compose to start all services:

```bash
docker-compose up -d
```

Accessing a Container's Shell:

```bash
docker exec -it <container name> bash
```

Stopping Services:

```bash
docker-compose down
```

### Deploy to Kubernetes Engine
<hr />
Make sure you have the following installed:

- Docker
- kubectl (configured to access your Kubernetes cluster)
- Kubernetes cluster (e.g., GKE, EKS, or Minikube)
- NGINX Ingress Controller

#### Docker Build and Push

To build and push the Docker image to Docker Hub:

1. **Build the Docker Image**:

   ```bash
   docker build -f ./build/Dockerfile -t <docker_hub_username>/book-store:latest .
   ```

2. **Push the Image to Docker Hub**:

   ```bash
   docker image push <docker_hub_username>/book-store:latest
   ```

   Replace `<docker_hub_username>` with your Docker Hub username. This will upload your image to Docker Hub for deployment.

#### Kubernetes Configurations

##### Creating a ConfigMap

To create a ConfigMap from an environment file:

```bash
kubectl create configmap book-store-env --from-file=./env/prod/.env
```

#### Setting Up NGINX Ingress
1. Deploy the NGINX Ingress Controller:
```bash
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.12.0-beta.0/deploy/static/provider/cloud/deploy.yaml
```

2. Apply the Ingress Resource:
```bash
kubectl apply -f ./build/book-store-ingress.yml
```

#### Applying Services and Deployments
1. Apply Services:
```bash
kubectl apply -f ./build/auth/auth-service.yml
kubectl apply -f ./build/user/user-service.yml
kubectl apply -f ./build/book/book-service.yml
kubectl apply -f ./build/shelf/shelf-service.yml
kubectl apply -f ./build/order/order-service.yml
```
2. Apply Deployments:
```bash
kubectl apply -f ./build/auth/auth-deployment.yml
kubectl apply -f ./build/user/user-deployment.yml
kubectl apply -f ./build/book/book-deployment.yml
kubectl apply -f ./build/shelf/shelf-deployment.yml
kubectl apply -f ./build/order/order-deployment.yml
```



### Migration Database
<hr />
dev

```bash
go run pkg/database/migration/migration.go ./env/dev/.env.auth
go run pkg/database/migration/migration.go ./env/dev/.env.book
go run pkg/database/migration/migration.go ./env/dev/.env.order
go run pkg/database/migration/migration.go ./env/dev/.env.shelf
go run pkg/database/migration/migration.go ./env/dev/.env.user
```

prod

```bash
go run pkg/database/migration/migration.go ./env/prod/.env.auth
go run pkg/database/migration/migration.go ./env/prod/.env.book
go run pkg/database/migration/migration.go ./env/prod/.env.order
go run pkg/database/migration/migration.go ./env/prod/.env.shelf
go run pkg/database/migration/migration.go ./env/prod/.env.user
```

### Kafka 
<hr />
Create topic
dev:

```bash
go run pkg/queue/topic/topic.go ./env/prod/.env.book
```

prod:

```bash
go run pkg/queue/topic/topic.go ./env/prod/.env.book
```

### Generate a Proto File Command
<hr />
User

```bash
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    ./module/user/userPb/userPb.proto
```

Auth

```bash
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    ./module/auth/authPb/authPb.proto
```

Book

```bash
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    ./module/book/bookPb/bookPb.proto
```

Shelf

```bash
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    ./module/shelf/shelfPb/shelfPb.proto
```




## License

Distributed under the MIT License. See LICENSE for more information.

---

> GitHub [@kritAsawaniramol](https://github.com/kritAsawaniramol) &nbsp;&middot;&nbsp;

<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->

[Kubernetes-url]: https://kubernetes.io/
[Kubernetes-badge]: https://img.shields.io/badge/kubernetes-326ce5.svg?&style=for-the-badge&logo=kubernetes&logoColor=white
[Docker-url]: https://www.docker.com/
[Docker-badge]: https://img.shields.io/badge/Docker-2CA5E0?style=for-the-badge&logo=docker&logoColor=white
[GO-url]: https://go.dev/
[GO.dev]: https://img.shields.io/badge/golang-00ADD8?&style=for-the-badge&logo=go&logoColor=white
[GORM-url]: https://gorm.io/
[GORM.io]: https://img.shields.io/badge/gorm-ORM-00ADD8?style=for-the-badge&logo=go&logoColor=white
[Gin-url]: https://gin-gonic.com/
[Gin-badge]: https://img.shields.io/badge/gin-008ECF?style=for-the-badge&logo=gin&logoColor=white
[Postgresql-badge]: https://img.shields.io/badge/PostgreSQL-316192?style=for-the-badge&logo=postgresql&logoColor=white
[Postgresql-url]: https://www.postgresql.org/
[Kafka-badge]: https://img.shields.io/badge/Apache_Kafka-231F20?style=for-the-badge&logo=apache-kafka&logoColor=white
[Kafka-url]: https://kafka.apache.org/
[GoogleCloud-badge]: https://img.shields.io/badge/Google_Cloud-4285F4?style=for-the-badge&logo=google-cloud&logoColor=white
[GoogleCloud-url]: https://cloud.google.com/cloud-console/?utm_source=google&utm_medium=cpc&utm_campaign=japac-TH-all-en-dr-BKWS-all-super-trial-PHR-dr-1605216&utm_content=text-ad-none-none-DEV_c-CRE_667174148378-ADGP_Hybrid+%7C+BKWS+-+BRO+%7C+Txt+-Management+Tools-Cloud+Console-google+cloud+console-main-KWID_43700077643446620-kwd-296393718382&userloc_9198024-network_g&utm_term=KW_google%20cloud%20console&gad_source=1&gclid=Cj0KCQjw4Oe4BhCcARIsADQ0csntshHEjPLB2CbT3X2VhKImmhD1HPdaONsmf7DNd4QnfHX3EanwatsaAselEALw_wcB&gclsrc=aw.ds
[Grpc-badge]: https://img.shields.io/badge/gRPC-244c5a?style=for-the-badge&logoColor=white
[Grpc-url]: https://grpc.io/
[Stripe-badge]: https://img.shields.io/badge/Stripe-626CD9?style=for-the-badge&logo=Stripe&logoColor=white
[Stripe-url]: https://stripe.com/
[Supabase-badge]: https://img.shields.io/badge/Supabase-181818?style=for-the-badge&logo=supabase&logoColor=white
[Supabase-url]: https://supabase.com/
[Confluent-url]: https://www.confluent.io/
[Confluent-badge]: https://img.shields.io/badge/Confluent-0077B5?style=for-the-badge&logoColor=white&logo=data:image/svg%2bxml;base64,PHN2ZyB2ZXJzaW9uPSIxLjEiIGlkPSJMYXllcl8xIiB4bWxuczp4PSJuc19leHRlbmQ7IiB4bWxuczppPSJuc19haTsiIHhtbG5zOmdyYXBoPSJuc19ncmFwaHM7IiB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHhtbG5zOnhsaW5rPSJodHRwOi8vd3d3LnczLm9yZy8xOTk5L3hsaW5rIiB4PSIwcHgiIHk9IjBweCIgdmlld0JveD0iMCAwIDQzIDQzIiBzdHlsZT0iZW5hYmxlLWJhY2tncm91bmQ6bmV3IDAgMCA0MyA0MzsiIHhtbDpzcGFjZT0icHJlc2VydmUiPgogPHN0eWxlIHR5cGU9InRleHQvY3NzIj4KICAuc3Qwe2ZpbGwtcnVsZTpldmVub2RkO2NsaXAtcnVsZTpldmVub2RkO2ZpbGw6I0ZGRkZGRjt9CiA8L3N0eWxlPgogPG1ldGFkYXRhPgogIDxzZncgeG1sbnM9Im5zX3NmdzsiPgogICA8c2xpY2VzPgogICA8L3NsaWNlcz4KICAgPHNsaWNlU291cmNlQm91bmRzIGJvdHRvbUxlZnRPcmlnaW49InRydWUiIGhlaWdodD0iNDMiIHdpZHRoPSI0MyIgeD0iNjguNyIgeT0iLTIyMS43Ij4KICAgPC9zbGljZVNvdXJjZUJvdW5kcz4KICA8L3Nmdz4KIDwvbWV0YWRhdGE+CiA8Zz4KICA8cGF0aCBjbGFzcz0ic3QwIiBkPSJNMzAuOCwyMi43YzEuMS0wLjEsMi4yLTAuMSwzLjMtMC4ydi0wLjhjLTEuMS0wLjEtMi4yLTAuMS0zLjMtMC4ybC0zLjMtMC4xYy0xLjktMC4xLTMuOS0wLjEtNS44LTAuMQoJCWMwLTEuOSwwLTMuOS0wLjEtNS44bC0wLjEtMy4zYy0wLjEtMS4xLTAuMS0yLjItMC4yLTMuM2gtMC45Yy0wLjEsMS4xLTAuMSwyLjItMC4yLDMuM2wtMC4xLDMuM2MwLDAuOSwwLDEuOCwwLDIuNwoJCWMtMC40LTAuOC0wLjctMS43LTEuMS0yLjVsLTEuNC0zYy0wLjUtMS0wLjktMi0xLjQtM0wxNS40LDEwYzAuNCwxLjEsMC43LDIuMSwxLjEsMy4xbDEuMiwzLjFjMC4zLDAuOSwwLjcsMS43LDEsMi42CgkJYy0wLjctMC42LTEuMy0xLjMtMi0xLjlsLTIuNC0yLjNjLTAuOC0wLjctMS42LTEuNS0yLjUtMi4ybC0wLjcsMC43YzAuNywwLjgsMS41LDEuNiwyLjIsMi41bDIuMywyLjRjMC42LDAuNywxLjMsMS4zLDEuOSwyCgkJYy0wLjgtMC4zLTEuNy0wLjctMi42LTFsLTMuMS0xLjJjLTEtMC40LTIuMS0wLjgtMy4xLTEuMWwtMC40LDAuOWMxLDAuNSwyLDAuOSwzLDEuNGwzLDEuNGMwLjgsMC40LDEuNywwLjcsMi41LDEuMQoJCWMtMC45LDAtMS44LDAtMi43LDBsLTMuMywwLjFjLTEuMSwwLjEtMi4yLDAuMS0zLjMsMC4ydjAuOWMxLjEsMC4xLDIuMiwwLjEsMy4zLDAuMmwzLjMsMC4xYzIsMC4xLDMuOSwwLjEsNS44LDAuMQoJCWMwLDEuOSwwLDMuOSwwLjEsNS44bDAuMSwzLjNjMC4xLDEuMSwwLjEsMi4yLDAuMiwzLjNoMC44YzAuMS0xLjEsMC4xLTIuMiwwLjItMy4zbDAuMS0zLjNjMC0wLjksMC0xLjksMC4xLTIuOAoJCWMwLjQsMC45LDAuNywxLjcsMS4xLDIuNmwxLjQsM2MwLjUsMSwwLjksMiwxLjQsM2wwLjgtMC4zYy0wLjMtMS4xLTAuNy0yLjEtMS4xLTMuMUwyNC4xLDI4Yy0wLjMtMC45LTAuNy0xLjctMS0yLjYKCQljMC43LDAuNywxLjMsMS4zLDIsMS45bDIuNCwyLjNjMC44LDAuNywxLjYsMS41LDIuNSwyLjJsMC42LTAuNmMtMC43LTAuOC0xLjUtMS42LTIuMi0yLjVsLTIuMy0yLjRjLTAuNi0wLjctMS4zLTEuNC0xLjktMgoJCWMwLjksMC4zLDEuNywwLjcsMi42LDFsMy4xLDEuMmMxLDAuNCwyLjEsMC44LDMuMSwxLjFsMC4zLTAuOGMtMS0wLjUtMi0xLTMtMS40bC0zLTEuNGMtMC45LTAuNC0xLjctMC44LTIuNi0xLjEKCQljMC45LDAsMS45LDAsMi44LTAuMUMyNy41LDIyLjgsMzAuOCwyMi43LDMwLjgsMjIuN3oiPgogIDwvcGF0aD4KICA8cGF0aCBjbGFzcz0ic3QwIiBkPSJNMjEuNSw0M0M5LjYsNDMsMCwzMy40LDAsMjEuNVM5LjYsMCwyMS41LDBTNDMsOS42LDQzLDIxLjVTMzMuNCw0MywyMS41LDQzIE0yMS41LDJDMTAuOCwyLDIsMTAuOCwyLDIxLjUKCQlTMTAuOCw0MSwyMS41LDQxUzQxLDMyLjIsNDEsMjEuNVMzMi4yLDIsMjEuNSwyIj4KICA8L3BhdGg+CiA8L2c+Cjwvc3ZnPg==
