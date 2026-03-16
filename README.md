# Go API

Go API 템플릿

## 주요 기술 스택

- Go 1.26
- `net/http`
- MySQL
- `sqlc`
- `caarlos0/env`
- `go-playground/validator/v10`
- `log/slog`
- `oapi-codegen`

## 실행 포트

- `8080`

## 빠른 시작

### 1. MySQL 실행

```sh
docker compose up -d
```

초기 실행 시 `sql/schema/*.sql` 이 자동 적용돼서 테이블이 생성된다

### 2. 환경 변수 확인

기본 개발용 설정은 [/.env](/Users/jeongju/projects/study/github_public/go-api/.env) 에 들어있다

```env
APP_ENV=local
HTTP_PORT=8080
LOG_LEVEL=debug
MYSQL_DSN=root:password@tcp(localhost:3306)/go_api?parseTime=true&loc=Local
```

### 3. 서버 실행

```sh
go run ./cmd/server
```

또는

```sh
make run
```

## Swagger

- UI: [http://localhost:8080/swagger/](http://localhost:8080/swagger/)
- Raw spec: [http://localhost:8080/openapi.yaml](http://localhost:8080/openapi.yaml)

## 주요 API

- `GET /health`
- `GET /users?page=1&size=20`
- `POST /users`
- `GET /users/{id}`
- `PUT /users/{id}`
- `DELETE /users/{id}`

### 목록 조회 응답 예시

```json
{
  "items": [
    {
      "id": 1,
      "name": "홍길동",
      "email": "hong@example.com",
      "createdAt": "2026-03-12T10:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "size": 20,
    "total_elements": 1,
    "total_pages": 1,
    "has_next": false,
    "has_previous": false
  }
}
```

### 수정 요청 예시

```json
{
  "name": "홍길동",
  "email": "hong@example.com"
}
```

삭제 성공 시에는 `204 No Content` 를 반환한다

## 코드 생성

### OpenAPI codegen

```sh
make generate-openapi
```

생성 위치

- [internal/gen/openapi/api.gen.go](/Users/jeongju/projects/study/github_public/go-api/internal/gen/openapi/api.gen.go)

### sqlc

```sh
make generate-sqlc
```

생성 위치

- [internal/gen/sqlc](/Users/jeongju/projects/study/github_public/go-api/internal/gen/sqlc)

### 전체 생성

```sh
make generate
```

## 개발 명령

```sh
make fmt
make vet
make test
make build
```

빌드 결과물

- [bin/go-api](/Users/jeongju/projects/study/github_public/go-api/bin/go-api)

## 디렉터리 구조

```text
cmd/server           서버 엔트리포인트
api                  OpenAPI 스펙
internal/api         OpenAPI 인터페이스 조립
internal/config      환경 변수 로드
internal/docs        Swagger UI, spec 노출
internal/httpx       공통 HTTP 유틸, 미들웨어
internal/user        user 도메인
platform/database    DB 연결 초기화
platform/logger      slog 초기화
sql/schema           DB schema
sql/queries          sqlc query
```

## 테스트

```sh
go test ./...
```

이 작업 환경에서는 캐시 충돌 방지를 위해 아래처럼 실행해도 된다

```sh
GOCACHE=$(pwd)/.gocache go test ./...
```

## 종료

```sh
docker compose down
```

볼륨까지 제거하려면

```sh
docker compose down -v
```
