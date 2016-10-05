package users

import (
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type storageKeyType int

var storageKey storageKeyType

// injectStorageMiddleware wraps with a handler that injects
// the storage into the context.
func injectStorageMiddleware(handler httprouter.Handle, st Storage) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, storageKey, st)
		handler(w, r.WithContext(ctx), ps)
	}
}

// jsonMiddleware wraps with a handler that both injects
// a json appropriate error handler and sets the content type to application/json
func jsonMiddleware(handler httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// set content type
		w.Header().Set("Content-Type", "application/json")

		// inject error sender
		ctx := r.Context()
		ctx = context.WithValue(ctx, errorSenderKey, jsonErrorSender(0))

		// call handler
		handler(w, r.WithContext(ctx), ps)
	}
}

// htmlMiddleware wraps with a handler that both injects
// an html appropriate error handler and sets the content type to text/html
func htmlMiddleware(handler httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// set content type
		w.Header().Set("Content-Type", "text/html")

		// inject error sender
		ctx := r.Context()
		ctx = context.WithValue(ctx, errorSenderKey, htmlErrorSender(0))

		// call handler
		handler(w, r.WithContext(ctx), ps)
	}
}
