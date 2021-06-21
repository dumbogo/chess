package api

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"gorm.io/gorm/clause"
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

// InitHTTPRouter initializes HTTPHandler, create mux Handler and define HTTTP routes
func InitHTTPRouter(c ConfigHTTPRouter) {
	HTTPRouter = mux.NewRouter()
	HTTPRouter.HandleFunc("/auth/{provider}/callback", callbackHandler)
	HTTPRouter.HandleFunc("/auth/{provider}", gothic.BeginAuthHandler)
	HTTPRouter.HandleFunc("/", rootHandler)
	provierCallbackURL := fmt.Sprintf("%s://%s/auth/github/callback?provider=github", c.URLLoc.Scheme, c.URLLoc.Host)
	goth.UseProviders(
		github.New(c.GithubKey, c.GithubSecret, provierCallbackURL, "user:email"),
	)
	gothic.GetState = func(r *http.Request) string {
		return r.URL.Query().Get("state")
	}
	initGothicStore(c.SessionKey, c.Env)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<p><a href='/auth/github?provider=github'>Click to log in with github</a></p>")
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: needs to delete this, leaving it as project is on development
	log.Printf("Received request: %+v\n", r)
	log.Println(gothic.GetState(r))
	// Ends todo

	user, err := gothic.CompleteUserAuth(w, r)
	fmt.Fprintf(w, "<p>User token, copy to clipboard, be careful with spaces!!!: <b>%s</b></p>", user.AccessToken)
	if err != nil {
		panic(err)
	}

	userdb := User{
		Email:             user.Email,
		Name:              user.Name,
		FirstName:         user.FirstName,
		LastName:          user.LastName,
		NickName:          user.NickName,
		UserID:            user.UserID,
		AccessToken:       user.AccessToken,
		AccessTokenSecret: user.AccessTokenSecret,
		RefreshToken:      user.RefreshToken,
		IDToken:           user.IDToken,
	}
	if user.ExpiresAt.IsZero() {
		userdb.ExpiresAt = sql.NullTime{Valid: false}
	} else {
		userdb.ExpiresAt = sql.NullTime{Valid: true, Time: user.ExpiresAt}
	}
	// Update all columns, except primary keys, to new value on conflict
	DBConn.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "email"}},
		UpdateAll: true,
	}).Create(&userdb)
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
