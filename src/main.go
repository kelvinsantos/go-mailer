package main

import (
	"kelvin.com/mailer/src/controllers"
	"kelvin.com/mailer/src/env"
	"kelvin.com/mailer/src/services"
	"kelvin.com/mailer/src/types"
	"kelvin.com/mailer/src/utils"
	"encoding/json"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func Hello() string {
	return "Hello, world."
}

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		token := r.Header.Get("Authorization")

		isValidated := utils.ValidateAuthSecret(token)
		if isValidated != "" {
			json.NewEncoder(w).Encode(types.ApiResponse{
				Message: isValidated,
				Success: false,
			})
			return
		}

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func main() {
	fmt.Println("Starting the application...")

	// initialize store and env
	utils.Init()

	// initialize mongo client
	services.SetClient(utils.GetClient())

	govalidator.SetFieldsRequiredByDefault(true)

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})

	// Endpoints to expose
	router.HandleFunc("/api/inbox/{email_address}", controllers.GetInboxByEmail).Methods("GET")
	router.HandleFunc("/api/inbox/{email_address}/{message_id}", controllers.GetMessageById).Methods("GET")
	router.HandleFunc("/api/send-email", controllers.SendRawEmail).Methods("POST")

	// Add custom middleware
	router.Use(middleware)

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	port := env.GO_MAILER_PORT
	fmt.Println("Server listening on port " + port)
	log.Fatal(http.ListenAndServe(port, handlers.CORS(originsOk, headersOk, methodsOk)(router)))
}