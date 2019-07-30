package credential

import (
	"encoding/gob"
	"text/template"

	"adrdn/dit/user"

	"github.com/gorilla/sessions"
	"github.com/gorilla/securecookie"
)

var store *sessions.CookieStore
var tmpl *template.Template

func init() {
	authenticationKey 	:= securecookie.GenerateRandomKey(64)
	encryptionKey		:= securecookie.GenerateRandomKey(32)

	store = sessions.NewCookieStore(authenticationKey, encryptionKey)
	store.Options = &sessions.Options {
		Path: 		"/",
		MaxAge:		86400 * 7,
		HttpOnly:	true,	 
	}
	gob.Register(user.User{})

	tmpl = template.Must(template.ParseGlob("forms/credential/*"))
}

func getUser(s *sessions.Session) user.User {
	u 		:= user.User{}
	value 	:= s.Values["name"]
	u, ok 	:= value.(user.User)
	if !ok {
		return user.User{Authenticated: false}
	}
	return u
}