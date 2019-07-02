package main

import (
	"fmt"
	"net/http"

	app "../Bank/app"
	"../Bank/models"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.Use(app.JwtAuthentication)
	router.HandleFunc("/api/bank/{ifsc}", models.GetBranchDetail).Methods("GET")
	router.HandleFunc("/api/banks/{branch}/{city}", models.GetBranchDetails).Methods("GET")
	router.HandleFunc("/v1/user/login", models.Authenticate).Methods("POST")
	router.HandleFunc("/v1/user/new", models.CreateUser).Methods("POST")
	router.HandleFunc("/v1/user/gettoken", models.GetToken).Methods("GET")
	router.HandleFunc("/v1/user/updatetoken", models.UpdateToken).Methods("PUT")
	err := http.ListenAndServe(":5000", router)
	if err != nil {
		fmt.Print(err)
	}
}
