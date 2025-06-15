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

	r.Post("/users", userC.Create)
	r.Post("/authenticate", userC.Authenticate)
	r.Post("/signout", userC.SignOut)
	r.Get("/users/me", userC.CurrentUser)
	r.Post("/datapoint", userC.Datapoint)

	fmt.Printf("Server running on port %d\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}
