package users

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

// this interface allows use to send an error and gets
// stored in the context. Error handling code calls
// this to send an error in the appropriate format.
type errorSender interface {
	Send(w http.ResponseWriter, r *http.Request, err error)
}

// The error sender gets injected into the context via middleware.

type errorSenderKeyType int

var errorSenderKey errorSenderKeyType

// sendError is our context-aware error sending function and is what
// the end users call.

func sendError(w http.ResponseWriter, r *http.Request, err error) {

	sender, ok := r.Context().Value(errorSenderKey).(errorSender)
	if !ok {
		// fallback to send the error as plain text
		code := getCode(err)
		w.WriteHeader(code)
		w.Write([]byte(err.Error()))
		return
	}

	sender.Send(w, r, err)
}

// The JSON middleware injects an error sender that sends the error as JSON.

type jsonErrorSender int

func (s jsonErrorSender) Send(w http.ResponseWriter, r *http.Request, err error) {
	code := getCode(err)

	m := make(map[string]interface{})
	m["error"] = errors.Cause(err).Error()
	m["code"] = code

	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(&m); err != nil {
		panic(err)
	}
}

// The HTML middleware injects an error sender that sends the error as HTML.

type htmlErrorSender int

func (s htmlErrorSender) Send(w http.ResponseWriter, r *http.Request, err error) {
	code := getCode(err)

	w.WriteHeader(code)

	w.Write([]byte(fmt.Sprintf(`
	<html>
		<head>
		</head>
		<body>
		<h1>Error '%d'</h1>
		<p>%v</p>
		</body>
	</html>
	`, code, errors.Cause(err).Error())))
}
