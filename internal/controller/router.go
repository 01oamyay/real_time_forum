package controller

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"text/template"
	"time"

	"rlf/internal/entity"
	"rlf/internal/service"
	"rlf/pkg/config"
)

type Handler struct {
	service   *service.Service
	secret    string
	webSocket *WebSocket
	*RateLimiter
}

type Route struct {
	Path    string
	Handler http.HandlerFunc
	Role    uint
}

func NewHandler(service *service.Service, secret string) *Handler {
	return &Handler{
		service:     service,
		secret:      secret,
		webSocket:   newWS(),
		RateLimiter: NewRateLimiter(1000000, time.Minute),
	}
}

func (h *Handler) InitRoutes(conf *config.Conf) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/src/", func(w http.ResponseWriter, r *http.Request) {
		if _, err := os.Stat(path.Join("./web", r.URL.Path)); os.IsNotExist(err) {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		} else if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Printf("Error checking file: %v", err)
			return
		}
		http.ServeFile(w, r, path.Join("./web", r.URL.Path))
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("./web/index.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		if err = tmpl.Execute(w, fmt.Sprintf("%v:%v", conf.API.Host, conf.API.Port)); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
	})
	mux.HandleFunc("/api/is-valid", h.isValidToken)
	routes := h.createRoutes()
	for _, route := range routes {
		if route.Role == entity.Roles.Authorized {
			mux.Handle(route.Path, h.Limiter(h.corsMiddleWare(h.isAlreadyIdentified(route.Handler))))
		} else {
			mux.Handle(route.Path, h.Limiter(h.corsMiddleWare(h.identify(route.Role, route.Handler))))
		}
	}
	return mux
}

func (h *Handler) createRoutes() []Route {
	return []Route{
		{
			Path:    "/api/signup",
			Handler: h.signUp,
			Role:    entity.Roles.Authorized,
		},
		{
			Path:    "/api/signin",
			Handler: h.signIn,
			Role:    entity.Roles.Authorized,
		},
		{
			Path:    "/api/signout",
			Handler: h.signOut,
			Role:    entity.Roles.User,
		},
		{
			Path:    "/api/profile/posts/",
			Handler: h.getAllPostsByUserID,
			Role:    entity.Roles.User,
		},
		{
			Path:    "/api/profile/liked-posts/",
			Handler: h.getAllLikedPostsByUserID,
			Role:    entity.Roles.User,
		},
		{
			Path:    "/api/profile/disliked-posts/",
			Handler: h.getAllDisLikedPostsByUserID,
			Role:    entity.Roles.User,
		},
		{
			Path:    "/api/post/create",
			Handler: h.createPost,
			Role:    entity.Roles.User,
		},
		{
			Path:    "/api/posts/",
			Handler: h.getALLPosts,
			Role:    entity.Roles.User,
		},
		{
			Path:    "/api/post/",
			Handler: h.getPostbyID,
			Role:    entity.Roles.User,
		},
		{
			Path:    "/api/post/vote",
			Handler: h.votePost,
			Role:    entity.Roles.User,
		},
		{
			Path:    "/api/comment/create",
			Handler: h.createComment,
			Role:    entity.Roles.User,
		},
		{
			Path:    "/api/comment/vote",
			Handler: h.voteComment,
			Role:    entity.Roles.User,
		},
		{
			Path:    "/api/categories",
			Handler: h.getAllCategories,
			Role:    entity.Roles.User,
		},
		{
			Path:    "/api/contacts",
			Handler: h.GetContacts,
			Role:    entity.Roles.User,
		},
		{
			Path:    "/ws",
			Handler: h.WebSocketHandler,
			Role:    entity.Roles.User,
		},
		{
			Path:    "/api/chat/",
			Handler: h.GetMessages,
			Role:    entity.Roles.User,
		},
	}
}
