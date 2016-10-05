package users

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/sheenobu/go-examples/rest"
)

var exampleUserRequest = rest.User{Username: "example"}
var exampleUserResponse = rest.User{"example", time.Now(), time.Now()}

// UsersAPI defines the REST API for the user resource
// Adding new resources that are also documentated is just
// a matter of adding a new entry into this list with the
// required fields.
var UsersAPI = []struct {
	Method         string
	URL            string
	Handler        httprouter.Handle
	Desc           string
	SampleRequest  interface{}
	SampleResponse interface{}
}{
	{"GET", "/users", listUsers, "List users", nil, &[]rest.User{exampleUserRequest, exampleUserRequest}},
	{"POST", "/users", createUser, "Create a user", &exampleUserRequest, &exampleUserResponse},
	{"GET", "/users/:id", getUser, "Get a user", nil, &exampleUserResponse},
	{"DELETE", "/users/:id", deleteUser, "Delete a user", nil, nil},
}

// RegisterHTTP registers the HTTP endpoints to the given router
func RegisterHTTP(router *httprouter.Router, storage Storage) {

	for _, api := range UsersAPI {
		handler := api.Handler

		// json content type
		handler = jsonMiddleware(handler) // set content-type to JSON

		// inject the user storage into the context
		handler = injectStorageMiddleware(handler, storage)

		// register method and URL to middleware wrapped handler
		router.Handle(api.Method, api.URL, handler)
	}

	router.Handle("GET", "/users.html", htmlMiddleware(userDocumentation))
}

// everything below here is just standard go 1.7 HTTP code, except for the
// usage of httprouter.Params and the magic behind sendError

func listUsers(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	st := r.Context().Value(storageKey).(Storage)

	users, err := st.List()
	if err != nil {
		sendError(w, r, err)
		return
	}

	if err := json.NewEncoder(w).Encode(&users); err != nil {
		sendError(w, r, err)
		return
	}
}

func createUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	st := r.Context().Value(storageKey).(Storage)

	var user rest.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		sendError(w, r, err)
		return
	}

	if err := st.Save(&user); err != nil {
		sendError(w, r, err)
		return
	}

	if err := json.NewEncoder(w).Encode(&user); err != nil {
		sendError(w, r, err)
		return
	}
}

func getUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	st := r.Context().Value(storageKey).(Storage)
	id := ps.ByName("id")

	user, err := st.Get(id)
	if err != nil {
		sendError(w, r, err)
		return
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		sendError(w, r, err)
		return
	}
}

func deleteUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	st := r.Context().Value(storageKey).(Storage)
	id := ps.ByName("id")

	user, err := st.Get(id)
	if err != nil {
		sendError(w, r, err)
		return
	}

	if err := st.Delete(user); err != nil {
		sendError(w, r, err)
		return
	}

}
