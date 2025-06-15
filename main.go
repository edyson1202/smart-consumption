package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"lenslocked/migrations"
	"lenslocked/pkg/controllers"
	"lenslocked/pkg/models"
	"net/http"
)

const port = 8084

func main() {

	cfg := models.DefaultPostgresConfig()
	fmt.Println(cfg)
	db, err := models.Open(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Println("Error connecting to database")
		panic(err)
	}

	err = models.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	userC := controllers.Users{
		UserService: &models.UserService{
			DB: db,
		},
		SessionService: &models.SessionService{
			DB: db,
		},
	}

	r := chi.NewRouter()

	r.Post("/users", enableCORS(userC.Create))
	r.Post("/authenticate", enableCORS(userC.Authenticate))
	r.Post("/signout", enableCORS(userC.SignOut))
	r.Get("/users/me", enableCORS(userC.CurrentUser))
	r.Post("/datapoints", enableCORS(userC.Datapoint))
	r.Get("/datapoints", enableCORS(userC.GetDatapoints))

	fmt.Printf("Server running on port %d\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Allow all origins (wildcard) — adjust for production
		w.Header().Set("Access-Control-Allow-Origin", "localhost:4200")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight OPTIONS request
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next(w, r)
	}
}
