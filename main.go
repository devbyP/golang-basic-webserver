package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"
)

var port = "3000"
var temp *template.Template

func main() {
	exe, err := os.Executable()
	if err != nil {
		log.Fatalln("cannot get exe path")
	}
	exePath := filepath.Join(filepath.Dir(exe), "..")

	temp = template.Must(template.ParseGlob(path.Join(exePath, "/templates/*.html")))

	mux := http.NewServeMux()

	mux.HandleFunc("/", handleHome)
	mux.HandleFunc("/hello", handleHello)
	mux.HandleFunc("/form", handleForm)
    mux.HandleFunc("/redirectToHello", handleRedirectHello)

	server := http.Server{
		Addr:         ":" + port,
		ReadTimeout:  time.Second * 15,
		WriteTimeout: time.Second * 15,
		IdleTimeout:  time.Second * 15,
		Handler:      mux,
	}

	server.ListenAndServe()
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		temp.ExecuteTemplate(w, "index.html", nil)
	default:
		temp.ExecuteTemplate(w, "methodnotallow.html", nil)
	}
}

func handleForm(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		temp.ExecuteTemplate(w, "form.html", nil)
	case http.MethodPost:
		formPost(w, r)
	default:
		temp.ExecuteTemplate(w, "methodnotallow.html", nil)
	}
}

func formPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	name := r.FormValue("name")
	http.Redirect(w, r, "/redirectToHello?name="+name, http.StatusTemporaryRedirect)
}

func handleRedirectHello(w http.ResponseWriter, r *http.Request) {
    name := r.URL.Query().Get("name")
	if name == "" {
		name = "unknown user"
	}
    t, err := template.New("redirectHello").Parse(`
        <head>
            <meta http-equiv="refresh" content="1; url='{{ . }}'" />
        </head>
        <body>
            redirecting you to result page
        </body>
    `)
    if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
        return
    }

    t.Execute(w, "/hello?name="+name)
}

func handleHello(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        temp.ExecuteTemplate(w, "methodnotallow.html", nil)
        return
    }
    name := r.URL.Query().Get("name")
    n := showName{Name: name}
    temp.ExecuteTemplate(w, "showResult.html", n)
}

type showName struct {
    Name string `json:"name"`
}
