package httpx

import (
	"log/slog"
	"net"
	"net/http"
	"strings"
	"time"
)

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (r *statusRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func RequestLoggingMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			recorder := &statusRecorder{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			// 요청 처리 시간을 측정해서 접근 로그에 함께 남김
			startedAt := time.Now()
			next.ServeHTTP(recorder, r)

			logger.Info("rcv",
				"ip", clientIP(r),
				"method", r.Method,
				"path", r.URL.Path,
				"status", recorder.statusCode,
				"duration", time.Since(startedAt),
			)
		})
	}
}

func clientIP(r *http.Request) string {
	// 프록시 체인에서는 첫 번째 IP 를 클라이언트 주소로 사용
	if forwardedFor := strings.TrimSpace(r.Header.Get("X-Forwarded-For")); forwardedFor != "" {
		if clientIP, _, _ := strings.Cut(forwardedFor, ","); strings.TrimSpace(clientIP) != "" {
			return strings.TrimSpace(clientIP)
		}
	}

	if realIP := strings.TrimSpace(r.Header.Get("X-Real-IP")); realIP != "" {
		return realIP
	}

	host, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	if err == nil && host != "" {
		return host
	}

	return strings.TrimSpace(r.RemoteAddr)
}

func RecoverMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if recovered := recover(); recovered != nil {
					logger.Error("panic recovered", "panic", recovered)
					WriteError(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "서버 내부 오류입니다", nil)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
