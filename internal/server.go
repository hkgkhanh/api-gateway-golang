package api

import (
	"fmt"
	"net/http"
	"time"

	_ "api-gateway-golang/cmd/server/docs"
	"api-gateway-golang/config"
	"api-gateway-golang/internal/handlers"
	"api-gateway-golang/internal/middlewares"
	"api-gateway-golang/internal/storage"
	"api-gateway-golang/pkg/httputils"
	"api-gateway-golang/pkg/logger"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"github.com/unrolled/render"
	"github.com/unrolled/secure"

	"github.com/urfave/negroni"
)

type AppServer struct {
	Env     string
	Port    string
	Version string
	handlers.Handlers
}

func (app *AppServer) Run(appConfig config.ApiEnvConfig) {
	app.Env = appConfig.Env
	app.Port = appConfig.Port
	app.Version = appConfig.Version
	app.Sender = &httputils.Sender{
		Render: render.New(render.Options{
			IndentJSON: true,
		}),
	}

	// can change DB implementation from here
	storage, err := storage.NewPostgresDB()
	if err != nil {
		logger.Log.Error(err)
		panic(err.Error())
	}
	// Migrations which will update the DB or create it if it doesn't exist.
	if err := storage.MigratePostgres("file://migrations"); err != nil {
		logger.Log.Fatal(err)
	}
	app.Storage = storage

	router := mux.NewRouter().StrictSlash(true)
	router.MethodNotAllowedHandler = http.HandlerFunc(app.NotAllowedHandler)
	router.NotFoundHandler = http.HandlerFunc(app.NotFoundHandler)
	// router.Methods("GET").Path("/api/books").HandlerFunc(app.GetBooksHandler)
	// router.Methods("GET").Path("/api/book/{id:[0-9]+}").HandlerFunc(app.GetBookHandler)
	// router.Methods("POST").Path("/api/book/add").HandlerFunc(app.AddBookHandler)
	// router.Methods("PATCH").Path("/api/book/update").HandlerFunc(app.UpdateBookHandler)
	// router.Methods("DELETE").Path("/api/book/delete/{id:[0-9]+}").HandlerFunc(app.DeleteBookHandler)
	// other handlers
	router.Methods("GET").Path("/api/users").HandlerFunc(app.GetUsersHandler)
	router.Methods("GET").Path("/api/user/{id:[0-9]+}").HandlerFunc(app.GetUserHandler)
	router.Methods("POST").Path("/api/user/add").HandlerFunc(app.AddUserHandler)
	router.Methods("PATCH").Path("/api/user/update").HandlerFunc(app.UpdateUserHandler)
	router.Methods("DELETE").Path("/api/user/delete/{id:[0-9]+}").HandlerFunc(app.DeleteUserHandler)

	router.Methods("GET").Path("/api/products").HandlerFunc(app.GetProductsHandler)
	router.Methods("GET").Path("/api/product/{id:[0-9]+}").HandlerFunc(app.GetProductHandler)
	router.Methods("POST").Path("/api/product/add").HandlerFunc(app.AddProductHandler)
	router.Methods("PATCH").Path("/api/product/update").HandlerFunc(app.UpdateProductHandler)
	router.Methods("DELETE").Path("/api/product/delete/{id:[0-9]+}").HandlerFunc(app.DeleteProductHandler)

	if app.Env != config.PROD_ENV {
		router.Methods("GET").PathPrefix("/api/docs/").Handler(httpSwagger.Handler(
			httpSwagger.URL(fmt.Sprint("http://localhost:", app.Port, "/api/docs/doc.json")),
			httpSwagger.DeepLinking(true),
			httpSwagger.DocExpansion("none"),
			httpSwagger.DomID("swagger-ui"),
		))
	}

	// Security Middlewares
	secureMiddleware := secure.New(secure.Options{
		IsDevelopment:      app.Env == "DEV",
		ContentTypeNosniff: true,
		SSLRedirect:        true,
		// If the app is behind a proxy
		// SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
	})

	// Usual Middlewares
	n := negroni.New()
	n.Use(negroni.NewLogger())
	n.Use(negroni.NewRecovery())
	n.Use(negroni.HandlerFunc(secureMiddleware.HandlerFuncWithNext))
	n.Use(negroni.HandlerFunc(middlewares.TrackRequestMiddleware))
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Allows all origins
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           86400,
	})
	// router with cors middleware
	wrappedRouter := corsMiddleware.Handler(router)
	n.UseHandler(wrappedRouter)

	startupMessage := "Starting API server (v" + app.Version + ")"
	startupMessage = startupMessage + " on port " + app.Port
	startupMessage = startupMessage + " in " + app.Env + " mode."
	logger.Log.Info(startupMessage)

	addr := ":" + app.Port
	if app.Env == "DEV" {
		addr = "localhost:" + app.Port
	}

	server := http.Server{
		Addr:         addr,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 90 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      n,
	}

	logger.Log.Info("Listening...")

	server.ListenAndServe()
}

// OnShutdown is called when the server has a panic.
func (app *AppServer) OnShutdown() {
	// Do cleanup or logging
	logger.OutputLog.Error("Executed OnShutdown")
}

// Special server handlers, outside of specific routes we have
func (app *AppServer) NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	err := app.Sender.JSON(w, http.StatusNotFound, fmt.Sprint("Not Found ", r.URL))
	if err != nil {
		panic(err)
	}
}

func (app *AppServer) NotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	err := app.Sender.JSON(w, http.StatusMethodNotAllowed, fmt.Sprint(r.Method, " method not allowed"))
	if err != nil {
		panic(err)
	}
}
