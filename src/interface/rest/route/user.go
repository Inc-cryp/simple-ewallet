package route

import (
	"net/http"

	handlers "github.com/mrdolles/wallet-app/src/interface/rest/handlers/user"

	"github.com/go-chi/chi/v5"
)

// HealthRouter a completely separate router for health check routes
func UserRouter(h handlers.UserHandlerInterface) http.Handler {
	r := chi.NewRouter()

	r.Post("/create_user", h.Register)
	r.Post("/login", h.Login)

	return r
}
