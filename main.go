package main

import (
	"log"
	"net/http"
	"time"
	"html/template"

	"github.com/google/uuid"
)

func cssHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "styles/output.css")
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()
		username := r.Form.Get("username")
		password := r.Form.Get("password")

		expectedPassword, ok := users[username]
		if !ok || expectedPassword != password {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println("Login failed for user:", username)
			return
		}

		sessionToken := uuid.NewString()
		expiresAt := time.Now().Add(120 * time.Second)

		sessions[sessionToken] = session{
			username: username,
			expiry:   expiresAt,
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   sessionToken,
			Expires: expiresAt,
		})

		log.Println("Username:", username)
		log.Println("Password:", password)

		w.Header().Set("HX-Redirect", "/")
	} else {
		tmpl, _ := template.New("").ParseFiles("./templates/login.html", "./templates/base.html")
		tmpl.ExecuteTemplate(w, "base", nil)
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
	if err != nil {
		log.Println("No session token found")
		if err == http.ErrNoCookie {
			w.Header().Set("HX-Redirect", "/login")
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sessionToken := c.Value
	for k := range sessions {
		if k == sessionToken {
			delete(sessions, k)
		}
	}

	w.Header().Set("HX-Redirect", "/login")
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("session_token")
		if err != nil {
			log.Println("No session token found")
			if err == http.ErrNoCookie {
				http.Redirect(w, r, "/login", http.StatusFound)

				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		sessionToken := c.Value
		sess, ok := sessions[sessionToken]
		if !ok || sess.isExpired() {
			log.Println("Session expired")
			http.Redirect(w, r, "/login", http.StatusFound)

			for k := range sessions {
				if k == sessionToken {
					delete(sessions, k)
				}
			}
			return
		}

		tmpl, _ := template.New("").ParseFiles("./templates/index.html", "./templates/base.html")
		tmpl.ExecuteTemplate(w, "base", nil)
	})
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/styles/output.css", cssHandler)

	log.Println("Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
