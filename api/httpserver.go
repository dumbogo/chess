package api

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
)

// HTTPRouter Mux router for HTTP server
var HTTPRouter *mux.Router

// InitHTTPRouter initializes HTTPHandler
func InitHTTPRouter(urlLoc url.URL, githubKey, githubSecret string) {
	// TODO: Refactor
	HTTPRouter = mux.NewRouter()
	HTTPRouter.HandleFunc("/auth/{provider}/callback", callbackHandler)
	HTTPRouter.HandleFunc("/auth/{provider}", gothic.BeginAuthHandler)
	HTTPRouter.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "<p><a href='/auth/github?provider=github'>Click to log in with github</a></p>")
	})
	provierCallbackURL := fmt.Sprintf("%s://%s/auth/github/callback?provider=github", urlLoc.Scheme, urlLoc.Host)
	goth.UseProviders(
		github.New(githubKey, githubSecret, provierCallbackURL),
	)
	gothic.GetState = func(r *http.Request) string {
		return r.URL.Query().Get("state")
	}
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Received request: %+v\n", r)
	fmt.Println(gothic.GetState(r))
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		panic(err)
	}
	fmt.Fprintln(w, "logged in!", user)
}
