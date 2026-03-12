package docs

import (
	"net/http"

	spec "go-api"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Register(mux *http.ServeMux) {
	// swagger 경로와 raw spec 경로를 함께 노출
	mux.HandleFunc("GET /swagger", h.redirectSwagger)
	mux.HandleFunc("GET /swagger/", h.getSwaggerUI)
	mux.HandleFunc("GET /openapi.yaml", h.getOpenAPI)
}

func (h *Handler) redirectSwagger(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/swagger/", http.StatusMovedPermanently)
}

func (h *Handler) getOpenAPI(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/yaml; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(spec.OpenAPI)
}

func (h *Handler) getSwaggerUI(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	// 외부 CDN 의 Swagger UI asset 을 사용해서 최소 구현 유지
	_, _ = w.Write([]byte(swaggerHTML))
}

const swaggerHTML = `<!doctype html>
<html lang="ko">
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>Go API Swagger UI</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css" />
    <style>
      html, body {
        margin: 0;
        padding: 0;
        background: #f5f7fb;
      }
      .topbar {
        display: none;
      }
    </style>
  </head>
  <body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
    <script>
      window.ui = SwaggerUIBundle({
        url: "/openapi.yaml",
        dom_id: "#swagger-ui",
        deepLinking: true,
        displayRequestDuration: true,
        persistAuthorization: true
      })
    </script>
  </body>
</html>
`
