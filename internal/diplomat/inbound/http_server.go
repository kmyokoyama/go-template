package inbound

import (
	"context"
	"encoding/json"
	"net/http"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/fx"

	"github.com/kmyokoyama/go-template/internal/adapters"
	"github.com/kmyokoyama/go-template/internal/components"
	"github.com/kmyokoyama/go-template/internal/controllers"
	"github.com/kmyokoyama/go-template/internal/models"
	"github.com/kmyokoyama/go-template/internal/wire"
)

type Route struct {
	Name        string
	Method      string
	Path        string
	HandlerFunc HandlerFuncWithComponents
}

type Routes []Route

type HandlerFuncWithComponents func(http.ResponseWriter, *http.Request, *components.Components)

func (f HandlerFuncWithComponents) WithComponents(c *components.Components) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		f(w, r, c)
	}
}

func VersionHandler(w http.ResponseWriter, r *http.Request, c *components.Components) {
	c.Logger.Info("received request on GET /version")

	version, err := controllers.GetVersion(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := adapters.ToVersionResponse(version)

	respJson, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.Write(respJson)
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request, c *components.Components) {
	c.Logger.Info("received request on POST /user")

	var createUserRequest wire.CreateUserRequest
	err := json.NewDecoder(r.Body).Decode(&createUserRequest)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user, err := adapters.ToUserInternal(createUserRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err = controllers.Signup(c, user, createUserRequest.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := adapters.ToUserResponse(user)

	respJson, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.Write(respJson)
}

func FindUserHandler(w http.ResponseWriter, r *http.Request, c *components.Components) {
	c.Logger.Info("received request on GET /user")

	vars := mux.Vars(r)

	id, _ := uuid.Parse(vars["id"])

	user, err := controllers.FindUser(c, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := adapters.ToUserResponse(user)

	respJson, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.Write(respJson)
}

var routes Routes = Routes{
	Route{"version", "GET", "/version", VersionHandler},
	Route{"create user", "POST", "/user", CreateUserHandler},
	Route{"find user", "GET", "/user/{id}", FindUserHandler},
}

func NewRouter(c *components.Components) http.Handler {
	router := mux.NewRouter().StrictSlash(true)

	for _, route := range routes {
		router.Name(route.Name).Methods(route.Method).Path(route.Path).HandlerFunc(route.HandlerFunc.WithComponents(c))
	}

	return router
}

func NewHttpServer(lc fx.Lifecycle, h http.Handler, c *components.Components) *http.Server {
	srv := &http.Server{Addr: ":8080"}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			c.Logger.Info("Started on", "port", ":8080")
			go http.ListenAndServe(":8080", h)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		}})

	return srv
}

// JWT.

func NewToken(id uuid.UUID, role models.Role, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"user-id": id,
		"role":    role.String(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(secret) // TODO: Handle this error.

	return tokenString, err
}
