package server

import (
	"context"
	"github.com/go-chi/chi/v5"
	"net/http"
	"os"
	"story-time-server/handler"
	"story-time-server/middlewares"
	"story-time-server/utils"
	"time"
)

type Server struct {
	chi.Router
	server *http.Server
}

func SetupRoutes() *Server {
	router := chi.NewRouter()
	router.Route("/story", func(story chi.Router) {
		story.Use(middlewares.CORSMiddlewares()...)
		story.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			utils.RespondJSON(w, http.StatusOK, struct {
				Status string `json:"status"`
			}{Status: "server is running"})
		})

		story.Route("/", func(public chi.Router) {
			public.Get("/avatar", handler.GetAvatars)
			public.Post("/register", handler.RegisterUser)
			public.Post("/login", handler.LoginUser)
		})

		story.Route("/kid", func(kid chi.Router) {
			kid.Use(middlewares.FirebaseAuthMiddleware)
			kid.Group(UserKidsRoutes)
		})

	})
	return &Server{
		Router: router,
	}
}

func (svc *Server) Run() error {
	port := ":" + os.Getenv("SERVER_PORT")
	svc.server = &http.Server{
		Addr:    port,
		Handler: svc.Router,
	}
	return svc.server.ListenAndServe()
}

func (svc *Server) Shutdown(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return svc.server.Shutdown(ctx)
}
