package credential

import (
	"net/http"
	"encoding/gob"
	"text/template"

	"adrdn/dit/user"

	"github.com/gorilla/sessions"
	"github.com/gorilla/securecookie"
)

// Store will hold all session data
var Store *sessions.CookieStore
var tmpl *template.Template

func init() {
	authenticationKey 	:= securecookie.GenerateRandomKey(64)
	encryptionKey		:= securecookie.GenerateRandomKey(32)

	Store = sessions.NewCookieStore(authenticationKey, encryptionKey)
	Store.Options = &sessions.Options {
		Path: 		"/",
		MaxAge:		86400 * 7,
		//MaxAge:		60 * 1,
		HttpOnly:	true,	 
	}
	gob.Register(user.User{})

	tmpl = template.Must(template.ParseGlob("forms/credential/*"))
}

// CheckAuthentication checks if the user has access to the system
func CheckAuthentication(w http.ResponseWriter, r *http.Request) (bool, user.User) {
	session, err := Store.Get(r, "dit")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false, user.User{}
	}
	u := getUser(session)
	
	if auth := u.Authenticated; !auth {
		session.AddFlash("You don't have access:D")
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return false, u
		}
		http.Redirect(w, r, "/", http.StatusFound)
		return false, u
	} 
	return true, u
}

func getUser(s *sessions.Session) user.User {
	var u	 = user.User{}
	value 	:= s.Values["user"]
	u, ok 	:= value.(user.User)
	if !ok {
		return user.User{Authenticated: false}
	}
	return u
}