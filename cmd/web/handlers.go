package main

import (
	"encoding/json"
	"io/ioutil"
	"juntosajugar/pkg/models"
	"net/http"
	"net/url"
	"strconv"
)

func (app *application) healthCheck(w http.ResponseWriter, r *http.Request) {
	// Handler de ejemplo que devuele un Json indicando que el servidor esta ok

	// Creo un struct anonima con los valores que quiero mandar
	response := struct {
		Key   string
		Value string
	}{
		"servidor",
		"ok",
	}

	// convierto la string en un json
	js, err := json.Marshal(response)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(js)
	if err != nil {
		app.serverError(w, err)
		return
	}
}

func (app *application) userCreation(w http.ResponseWriter, r *http.Request) {
	// Handler que toma el body del request, trata de Unmarshalearlo en una struct de
	// tipo user, y si no hay duplicados responde con el mismo user en el cuerpo.
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		app.serverError(w, err)
		return
	}

	var nuevoUser models.User
	err = nuevoUser.FromJson(body)
	if err != nil {
		app.clientError(w, 400)
		return
	}

	err = app.db.Create(&nuevoUser).Error
	if duplicateError(err) {
		app.clientError(w, 409)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}
	app.responseJson(w, body)
	return
}

func (app *application) userDeletion(w http.ResponseWriter, r *http.Request) {
	// Manejador que dada una peticion con el id en la URI, elimina al usuario y devuelve un 200 vacío.
	userId, err := strconv.Atoi(r.URL.Query().Get(":id"))
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		app.serverError(w, err)
		return
	}

	var delUser models.User
	err = app.db.First(&delUser, userId).Error
	if err != nil && err.Error() == "record not found" {
		app.clientError(w, 404)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}
	err = app.db.Unscoped().Delete(&delUser).Error
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.responseJson(w, body)
	return
}

func (app *application) userRetrieval(w http.ResponseWriter, r *http.Request) {
	// Manejador que dada una peticion con el id de usuario en la URI, devuelve los datos del mismo en Json
	userId, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	var reqUser models.User
	err = app.db.First(&reqUser, userId).Error
	if err != nil && err.Error() == "record not found" {
		app.clientError(w, 404)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}
	body, err := json.Marshal(&reqUser)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.responseJson(w, body)
	return
}

func (app *application) userRetrievalByEmail(w http.ResponseWriter, r *http.Request) {
	// Manejador que dada una peticion con el id de usuario en la URI, devuelve los datos del mismo en Json
	encodedParameter := r.URL.Query().Get(":email")
	userEmail, err := url.QueryUnescape(encodedParameter)
	if err != nil {
		app.serverError(w, err)
		return
	}
	var reqUser models.User
	err = app.db.Where("email = ?", userEmail).First(&reqUser).Error
	if err != nil && err.Error() == "record not found" {
		app.clientError(w, 404)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}
	body, err := json.Marshal(&reqUser)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.responseJson(w, body)
	return
}
