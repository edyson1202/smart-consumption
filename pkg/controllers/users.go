package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"lenslocked/pkg/models"
	"log"
	"net/http"
)

type Users struct {
	Templates struct {
		New    Template
		SignIn Template
	}
	UserService    *models.UserService
	SessionService *models.SessionService
}

func (u Users) Create(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprint(w, "Email: ", r.FormValue("email"))
	//fmt.Fprint(w, "Password: ", r.FormValue("password"))

	newUser := models.NewUser{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	user, err := u.UserService.Create(newUser)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal Server Error", 500)
		return
	}

	_, token, err := u.SessionService.Create(user.ID)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}

	http.SetCookie(w, newCookie(CookieSessionName, token))
	//http.Redirect(w, r, "/users/me", http.StatusFound)
}

func (u Users) Authenticate(w http.ResponseWriter, r *http.Request) {
	newUser := models.NewUser{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	user, err := u.UserService.Authenticate(newUser)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal Server Error", 500)
		return
	}
	if errors.Is(err, models.IncorrectCredentialsError) {
		http.Error(w, "Incorrect Credentials", 400)
		return
	}

	_, token, err := u.SessionService.Create(user.ID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal Server Error", 500)
		return
	}

	http.SetCookie(w, newCookie(CookieSessionName, token))
	//http.Redirect(w, r, "/users/me", http.StatusFound)
}

func (u Users) CurrentUser(w http.ResponseWriter, r *http.Request) {

	sessionCookie, err := r.Cookie(CookieSessionName)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}

	user, err := u.SessionService.User(sessionCookie.Value)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}

	fmt.Fprintf(w, "CurrentUser: %s", user.Email)
}

func (u Users) SignOut(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie(CookieSessionName)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}

	err = u.SessionService.Delete(sessionCookie.Value)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal Server Error", 500)
		return
	}

	deleteCookie(w, CookieSessionName)
	//http.Redirect(w, r, "/signin", http.StatusFound)
}

func (u Users) Datapoint(w http.ResponseWriter, r *http.Request) {
	/*	sessionCookie, err := r.Cookie(CookieSessionName)

		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			return
		}

		user, err := u.SessionService.User(sessionCookie.Value)
		if err != nil {
			http.Error(w, "Invalid token!", http.StatusForbidden)
			return
		}*/

	var dp models.DataPoint

	// Decode JSON body into struct
	err := json.NewDecoder(r.Body).Decode(&dp)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err = u.UserService.CreateDatapoint(1, dp)

	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}

	// For demonstration: print the struct
	fmt.Printf("Received data point: %+v\n", dp)

	// Respond with status
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "received"})
}

func (u Users) GetDatapoints(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie(CookieSessionName)

	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}

	user, err := u.SessionService.User(sessionCookie.Value)
	if err != nil {
		http.Error(w, "Invalid token!", http.StatusForbidden)
		return
	}

	datapoints, err := u.UserService.ListDatapoints(user.ID)
	if err != nil {
		http.Error(w, "Failed to fetch datapoints", http.StatusInternalServerError)
		log.Printf("Fetch error: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(datapoints); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Printf("Encoding error: %v", err)
	}
}
