package inbound

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"go.uber.org/fx"

	"github.com/kmyokoyama/go-template/internal/adapters"
	"github.com/kmyokoyama/go-template/internal/components"
	"github.com/kmyokoyama/go-template/internal/controllers"
	"github.com/kmyokoyama/go-template/internal/wire"
)

// Handlers.

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

func SignupHandler(w http.ResponseWriter, r *http.Request, c *components.Components) {
	c.Logger.Info("received request on POST /signup")

	var signupRequest wire.SignupRequest
	err := json.NewDecoder(r.Body).Decode(&signupRequest)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user, err := adapters.ToUserInternal(signupRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err = controllers.Signup(c, user, signupRequest.Password)
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

func LoginHandler(w http.ResponseWriter, r *http.Request, c *components.Components) {
	c.Logger.Info("received request on POST /login")

	var loginRequest wire.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginRequest)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token, err := controllers.Login(c, loginRequest.Username, loginRequest.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respJson, err := json.Marshal(wire.LoginResponse{Token: token})
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

func WorkHandler(w http.ResponseWriter, r *http.Request, c *components.Components) {
	c.Logger.Info("received request on POST /work")

	var workRequest wire.WorkRequest
	err := json.NewDecoder(r.Body).Decode(&workRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	controllers.ProcessWork(c, workRequest.Id, workRequest.Description)

	respJson, err := json.Marshal(wire.WorkResponse{Id: workRequest.Id, Status: "pending"})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write(respJson)
}

// Middlewares.

type MiddlewareWithComponents func(http.ResponseWriter, *http.Request, http.HandlerFunc, *components.Components)

func (f MiddlewareWithComponents) WithComponents(c *components.Components) negroni.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		f(w, r, next, c)
	}
}

func Logging(w http.ResponseWriter, r *http.Request, next http.HandlerFunc, c *components.Components) {
	fmt.Println("Logged, calling next")
	next(w, r)
}

func Authentication(w http.ResponseWriter, r *http.Request, next http.HandlerFunc, c *components.Components) {
	authHeader := r.Header.Get("Authorization")
	
	const BEARER_SCHEMA = "Bearer "
	var token string
	if strings.HasPrefix(authHeader, BEARER_SCHEMA) {
		token = authHeader[len(BEARER_SCHEMA):]
		valid := controllers.IsValidToken(token)

		if valid {
			next(w, r)
		} else {
			respJson, err := json.Marshal(map[string]string{"error": "forbidden"})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			w.Write(respJson)
		}
	}
}

func AsJson(w http.ResponseWriter, r *http.Request, next http.HandlerFunc, c *components.Components) {
	next(w, r)
	fmt.Println("AsJson, calling next")
}

// Routes.

type Route struct {
	Name        string
	Method      string
	Path        string
	HandlerFunc HandlerFuncWithComponents
	Middlewares []MiddlewareWithComponents
}

type Routes []Route

type HandlerFuncWithComponents func(http.ResponseWriter, *http.Request, *components.Components)

func (f HandlerFuncWithComponents) WithComponents(c *components.Components) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		f(w, r, c)
	}
}

var routes Routes = Routes{
	Route{Name: "version", Method: "GET", Path: "/version", HandlerFunc: VersionHandler},
	Route{Name: "signup", Method: "POST", Path: "/signup", HandlerFunc: SignupHandler},
	Route{Name: "login", Method: "POST", Path: "/login", HandlerFunc: LoginHandler},
	Route{Name: "find user", Method: "GET", Path: "/user/{id}", HandlerFunc: FindUserHandler},
	Route{
		Name:        "work",
		Method:      "POST",
		Path:        "/work",
		HandlerFunc: WorkHandler,
		Middlewares: []MiddlewareWithComponents{Logging, Authentication, AsJson},
	},
}

// Constructors.

func NewRouter(c *components.Components) http.Handler {
	router := mux.NewRouter().StrictSlash(true)

	for _, route := range routes {
		var mws []negroni.Handler
		for _, mw := range route.Middlewares {
			mws = append(mws, mw.WithComponents(c))
		}

		n := negroni.New(mws...)
		n.UseHandler(route.HandlerFunc.WithComponents(c))

		router.Name(route.Name).Methods(route.Method).Path(route.Path).Handler(n)
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
