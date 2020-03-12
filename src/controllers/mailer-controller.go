package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"kelvin.com/mailer/src/decorator"
	"kelvin.com/mailer/src/services"
	"kelvin.com/mailer/src/types"
	"kelvin.com/mailer/src/utils"
	"log"
	"net/http"
	"strconv"
)

func GetInboxByEmail(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	skip := r.FormValue("skip")
	var skipValue int64
	if len(skip) > 0 {
		value, err := strconv.ParseInt(skip, 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		skipValue = value
	} else {
		skipValue = 0
	}

	limit := r.FormValue("limit")
	var limitValue int64
	if len(limit) > 0 {
		value, err := strconv.ParseInt(limit, 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		limitValue = value
	} else {
		limitValue = 20
	}

	sort := r.FormValue("sort")
	var sortValue int64
	if len(sort) > 0 {
		value, err := strconv.ParseInt(sort, 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		sortValue = value
	} else {
		sortValue = -1
	}

	getInboxRequestJson := types.GetInboxRequestJson{
		Email: params["email_address"],
		Skip:  skipValue,
		Limit: limitValue,
		Sort:  sortValue,
	}

	err, response := services.GetInboxService(getInboxRequestJson)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetMessageById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	getMessageRequestJson := types.GetMessageRequestJson{
		Email: params["email_address"],
		MessageId: params["message_id"],
	}

	err, response := services.GetMessageService(getMessageRequestJson)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func SendRawEmail(w http.ResponseWriter, r *http.Request) {
	var requestJson types.SendEmailRequestJson

	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &requestJson)

	buildRawEmailInput := decorator.BuildEmailInputDecorate(services.BuildRawEmailInput)
	input := buildRawEmailInput(requestJson)

	if input.Error != nil {
		utils.HandleError(w, input.Error)
	}

	svc, err := services.CreateSesSession(input)
	if err != nil {
		utils.HandleError(w, err)
	}

	// Attempt to send the email.
	result, err := svc.SendRawEmail(input.SendRawEmailInput)

	// Display error messages if they occur.
	if err != nil {
		utils.HandleSesError(w, err)
	}

	fmt.Println("Email has been successfully sent!")
	fmt.Println(result)

	// Log sent emails
	services.SaveEmail(requestJson)

	//auditLog := buildAuditLog(requestJson)
	//saveAuditLog(auditLog)

	return
}

func SendEmail(w http.ResponseWriter, r *http.Request) {
	var requestJson types.SendEmailRequestJson

	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &requestJson)

	buildEmailInput := decorator.BuildEmailInputDecorate(services.BuildRawEmailInput)
	input := buildEmailInput(requestJson)

	if input.Error != nil {
		utils.HandleError(w, input.Error)
	}

	svc, err := services.CreateSesSession(input)
	if err != nil {
		utils.HandleError(w, err)
	}

	// Attempt to send the email.
	result, err := svc.SendEmail(input.SendEmailInput)

	// Display error messages if they occur.
	if err != nil {
		utils.HandleSesError(w, err)
	}

	fmt.Println("Email has been successfully sent!")
	fmt.Println(result)

	return
}
