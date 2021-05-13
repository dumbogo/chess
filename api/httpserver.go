package api

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
)

// HTTPRouter Mux router for HTTP server
var HTTPRouter *mux.Router

// ConfigHTTPRouter HTTP router configuration
type ConfigHTTPRouter struct {
	URLLoc       url.URL
	GithubKey    string
	GithubSecret string
	// SessionKey Ensure your key is sufficiently random - i.e. use Go's
	// crypto/rand or securecookie.GenerateRandomKey(32) and persist the result.
	SessionKey string

	// Env environment
	Env string
}

// InitHTTPRouter initializes HTTPHandler
func InitHTTPRouter(c ConfigHTTPRouter) {
	// TODO: Refactor
	HTTPRouter = mux.NewRouter()
	HTTPRouter.HandleFunc("/auth/{provider}/callback", callbackHandler)
	HTTPRouter.HandleFunc("/auth/{provider}", gothic.BeginAuthHandler)
	HTTPRouter.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "<p><a href='/auth/github?provider=github'>Click to log in with github</a></p>")
	})

	provierCallbackURL := fmt.Sprintf("%s://%s/auth/github/callback?provider=github", c.URLLoc.Scheme, c.URLLoc.Host)
	goth.UseProviders(
		github.New(c.GithubKey, c.GithubSecret, provierCallbackURL),
	)
	gothic.GetState = func(r *http.Request) string {
		return r.URL.Query().Get("state")
	}
	initGothicStore(c.SessionKey, c.Env)
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

func initGothicStore(key string, env string) {
	// By default, gothic uses a CookieStore from the gorilla/sessions package to store session data.
	// As configured, this default store (gothic.Store) will generate cookies with Options:
	// &Options{
	// 	Path:   "/",
	// 	Domain: "",
	// 	MaxAge: 86400 * 30,
	// 	HttpOnly: true,
	// 	Secure: false,
	// }
	// To tailor these fields for your application, you can override the gothic.Store variable at startup.
	// The following snippet shows one way to do this:
	maxAge := 86400 * 30 // 30 days
	isProd := false      // Set to true when serving over https
	if env == EnvProduction {
		isProd = true
	}

	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(maxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true // HttpOnly should always be enabled
	store.Options.Secure = isProd

	gothic.Store = store

}
