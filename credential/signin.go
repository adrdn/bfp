package credential

import (
	"fmt"
	"time"
	"net/http"

	"adrdn/dit/config"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

const listOneUser = "SELECT password FROM user where username = ?"

// Login revokes the login page
func Login(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "Login", nil)
}

// Authentication decides if the user can login or not
func Authentication(w http.ResponseWriter, r *http.Request) {
	//var creds Credentials

	db := config.DbConn()
	var hashedPassword string

	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")	

		checkData, err := db.Query(listOneUser, username)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		for checkData.Next() {
			err = checkData.Scan(&hashedPassword)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}

		if err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
			defer db.Close()
			w.WriteHeader(http.StatusUnauthorized)
			if err.Error() == "crypto/bcrypt: hashedPassword is not the hash of the given password" {
				fmt.Println("Invalid Password")
			} else {
				fmt.Println("Invalid Username")
			}
			return
		} 
		expirationTime := time.Now().Add(30 * time.Second)

		claims := &Claims{
			Username: username,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTime.Unix(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println(err)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name: "token",
			Value: tokenString,
			Expires: expirationTime,
		})

		fmt.Println("The cookie has been made")
		defer db.Close()
		http.Redirect(w, r, "/home", 301)
	}
}

// Home revokes the home page
func Home(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Println(err, "inja1")
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err, "inja2")
		return
	}

	tknStr := c.Value

	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Println(err)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Println(err)
		return
	}
	fmt.Println("Username: ")
	tmpl.ExecuteTemplate(w, "Home", nil)
}

// Refresh sth
func Refresh(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Println(err)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	tknStr := c.Value
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Println(err)
		return
	}
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Println(err)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30 *time.Second {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	expirationTime := time.Now().Add(30 * time.Second)
	claims.ExpiresAt = expirationTime.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	http.SetCookie(w, &http.Cookie {
		Name: "token",
		Value: tokenString,
		Expires: expirationTime,
	})
}