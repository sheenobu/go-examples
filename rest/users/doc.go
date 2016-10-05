// Package users implements the services for the user model
package users

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/sheenobu/go-examples/rest"
)

// This is the code for supporting HTML based documentation.
// We push our declarative API into our template and execute it.

var userDocTemplate *template.Template

func init() {
	tmpl, err := template.New("userDoc").Parse(t)
	if err != nil {
		panic(err)
	}

	userDocTemplate = tmpl
}

func userDocumentation(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	var m = make(map[string]interface{})
	m["users"] = UsersAPI

	var body []byte
	body, _ = json.Marshal(&rest.User{"example", time.Now(), time.Now()})

	var out bytes.Buffer
	json.Indent(&out, body, "", "\t")

	m["sampleUser"] = string(out.String())

	var tmp bytes.Buffer

	if err := userDocTemplate.Execute(&tmp, m); err != nil {
		sendError(w, r, err)
		return
	}

	io.Copy(w, &tmp)
}

var t = `
<html>
	<head>

	</head>
	<body>
		<h1>Users REST API</h1>

		<h2>User Model</h2>

		<pre><code>{{.sampleUser}}</code></pre>

		<h2>Endpoints</h2>

		<table>
			<thead>
				<tr>
					<th>Method</th>
					<th>URL</th>
					<th>Description</th>
				</tr>
			</thead>
			<tbody>
				{{ range .users }}
				<tr>
					<td>{{.Method}}</td>
					<td>{{.URL}}</td>
					<td>{{.Desc}}</td>
				</tr>
				{{ end  }}
			</tbody>
		</table>
	</body>
</html>
`
