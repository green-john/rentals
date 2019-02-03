package transport

import (
	"github.com/gorilla/handlers"
	"net/http"
	"os"
	"rentals/services"
	"strings"
)

// Middleware used to authenticate and authorize users.
// Uses the url to check which resource is being accessed
func (s *Server) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip if we are trying to login or creating a new client
		if r.URL.Path == "/login" || r.URL.Path == "/newClient" {
			next.ServeHTTP(w, r)
			return
		}

		authHeader := r.Header["Authorization"]
		if len(authHeader) == 0 {
			respond(w, http.StatusUnauthorized, "Not allowed")
			return
		}

		token := authHeader[0]
		user := s.authn.Verify(token)

		if user == nil {
			respond(w, http.StatusUnauthorized, "Not allowed")
			return
		}

		// Try to get the requested resource from the url.
		// If not users or apartments, then it should be login/createUser
		// and not authorization is needed.
		requestedResource := getResource(r.URL.Path)
		if requestedResource != "" {
			op := getOp(r.Method)

			if !s.authz.Allowed(user.Role, requestedResource, op) {
				respond(w, http.StatusForbidden, "Not allowed")
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) ContentTypeJsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func (s *Server) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next = handlers.LoggingHandler(os.Stderr, next)
		next.ServeHTTP(w, r)
	})
}

func (s *Server) CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Enable CORS for testing purposes. This should be
		// configured properly for production
		allOrigins := handlers.AllowedOrigins([]string{"ruizandr.es,localhost"})
		allMethods := handlers.AllowedMethods([]string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"})
		allHeaders := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})

		next = handlers.CORS(allOrigins, allMethods, allHeaders)(next)
		next.ServeHTTP(w, r)
	})
}

func getOp(method string) services.Permission {
	meth2Perm := make(map[string]services.Permission)
	meth2Perm["POST"] = services.Create
	meth2Perm["GET"] = services.Read
	meth2Perm["PATCH"] = services.Update
	meth2Perm["DELETE"] = services.Delete

	return meth2Perm[strings.ToUpper(method)]
}

// Returns the desired resource by inferring it from the url
func getResource(urlPath string) string {
	if strings.HasPrefix(urlPath, "/users") {
		return "users"
	} else if strings.HasPrefix(urlPath, "/apartments") {
		return "apartments"
	}
	return ""
}
