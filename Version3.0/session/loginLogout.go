//Package session contains the relevant structures and functions to handle:
//1. 		User Login/Logout
//2. 		User session management
//3. 		New User signup
package session

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"text/template"
	"time"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
)

var tpl = template.Must(template.ParseGlob("templates/*"))

//User struct contains the information of User. Password will be hashed by bcrypt function and stored as a hashed value.
type User struct {
	Username string
	Password []byte
	First    string
	Last     string
}

//MapUsers is a map of key:username and value: User struct.
var MapUsers = map[string]User{}

//MapSessions is a map of key:uuid string (to store cookie session IDs) and value: username
var MapSessions = map[string]string{}
var validate = validator.New()

//func Index is the handler for main Index page.
func Index(res http.ResponseWriter, req *http.Request) {
	myUser := GetUser(res, req)
	err := tpl.ExecuteTemplate(res, "index.gohtml", myUser)
	if err != nil {
		log.Fatalln("Error with Index template: ", err)
	}
}

//func Signup is the handler for user sign up.
func Signup(res http.ResponseWriter, req *http.Request) {
	if AlreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	var myUser User
	// process form submission
	if req.Method == http.MethodPost {
		// get form values
		username := req.FormValue("username")
		err1 := validate.Var(username, "required,min=3,max=30,alphanum")
		if err1 != nil {
			http.Error(res, "Invalid/missing username, please try again.", http.StatusForbidden)
			log.Warning("Attempt to signup with invalid username - ", err1)
		}

		password := req.FormValue("password")
		err2 := validate.Var(password, "required,min=6,max=20,alphanum")
		if err2 != nil {
			http.Error(res, "Attempt to signup with invalid password. Password should be alphanumberical and consist between 6 to 20 characters.", http.StatusForbidden)
			log.Warning("Attempt to signup with invalid password - ", err2)
		}

		firstname := req.FormValue("firstname")
		err3 := validate.Var(firstname, "required,min=2,max=30,alphanum")
		if err3 != nil {
			http.Error(res, "Invalid/empty first name, please try again.", http.StatusForbidden)
			log.Warning("Attempt to signup with invalid first name - ", err3)
		}

		lastname := (req.FormValue("lastname"))
		err4 := validate.Var(lastname, "required,min=2,max=30,alphanum")
		if err4 != nil {
			http.Error(res, "Invalid/empty last name, please try again.", http.StatusForbidden)
			log.Warning("Attempt to signup with invalid last name - ", err4)
		}

		if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
			return
		}

		if username != "" {
			// check if username exist/ taken
			if _, ok := MapUsers[username]; ok {
				http.Error(res, "Username already taken", http.StatusForbidden)
				log.Warning("Unsuccessful user signup attempt")
				return
			}

			// create session
			id := uuid.NewV4()
			expirytime := time.Now().Add(30 * time.Minute)

			myCookie := &http.Cookie{
				Name:     "myCookie",
				Value:    id.String(),
				Expires:  expirytime,
				HttpOnly: true,
				Path:     "/",
				Domain:   "127.0.0.1",
				Secure:   true,
			}

			http.SetCookie(res, myCookie)
			MapSessions[myCookie.Value] = username

			bPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
			if err != nil {
				http.Error(res, "Internal server error", http.StatusInternalServerError)
				log.Error("Error with bcrypt password generator.", err)
				return
			}

			myUser = User{username, bPassword, firstname, lastname}
			MapUsers[username] = myUser
			users, err := json.MarshalIndent(MapUsers, "", " ")
			if err != nil {
				log.Error(err)
			}

			err = ioutil.WriteFile("users.json", users, 0644)
			if err != nil {
				log.Error(err)
			}
		}
		// redirect to main index
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return

	}
	tpl.ExecuteTemplate(res, "signup.gohtml", myUser)
}

//func Login is the handler for user login.
func Login(res http.ResponseWriter, req *http.Request) {
	if AlreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	// process form submission
	if req.Method == http.MethodPost {
		username := req.FormValue("username")
		password := req.FormValue("password")
		// check if user exist with username
		myUser, ok := MapUsers[username]
		if !ok {
			http.Error(res, "Username and/or password do not match", http.StatusUnauthorized)
			log.Warning("Invalid login attempt - user does not exist.")
			return
		}

		// Matching of password entered
		err := bcrypt.CompareHashAndPassword(myUser.Password, []byte(password))
		if err != nil { //passwords do not match, mismatch will be logged
			http.Error(res, "Username and/or password do not match", http.StatusForbidden)
			errorString := "Failed authentication attempt (password mismatch) by " + myUser.Username
			log.Warning(errorString)
			return
		}

		if multiLogin(username) == true {
			http.Error(res, "Multiple login attempt detected, please log out from other device before proceeding.", http.StatusUnauthorized)
			log.Warning("Multiple session login attempted by username:", username)
			return
		}

		// create session
		id := uuid.NewV4()
		expirytime := time.Now().Add(30 * time.Minute)

		myCookie := &http.Cookie{
			Name:     "myCookie",
			Value:    id.String(),
			Expires:  expirytime,
			HttpOnly: true,
			Path:     "/",
			Domain:   "127.0.0.1",
			Secure:   true,
		}

		http.SetCookie(res, myCookie)
		MapSessions[myCookie.Value] = username
		http.Redirect(res, req, "/", http.StatusTemporaryRedirect)
		return
	}

	tpl.ExecuteTemplate(res, "login.gohtml", nil)
}

// Logout is the handler for user logout.
func Logout(res http.ResponseWriter, req *http.Request) {
	if !AlreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	myCookie, _ := req.Cookie("myCookie")
	// delete the session
	delete(MapSessions, myCookie.Value)
	// remove the cookie
	myCookie = &http.Cookie{
		Name:   "myCookie",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(res, myCookie)

	http.Redirect(res, req, "/", http.StatusSeeOther)
}

//func GetUser checks user cookie session ID and returns user details
func GetUser(res http.ResponseWriter, req *http.Request) User {
	var myUser User
	// get current session cookie
	myCookie, err := req.Cookie("myCookie")

	if err != nil {
		id := uuid.NewV4()

		expirytime := time.Now().Add(30 * time.Minute)

		myCookie = &http.Cookie{
			Name:     "myCookie",
			Value:    id.String(),
			Expires:  expirytime,
			HttpOnly: true,
			Path:     "/",
			Domain:   "127.0.0.1",
			Secure:   true,
		}

		http.SetCookie(res, myCookie)
	}

	// if the user exists already, get user
	if username, ok := MapSessions[myCookie.Value]; ok {
		myUser = MapUsers[username]
	}

	return myUser
}

//Function AlreadyLoggedIn takes in the HTTP request cookie and returns a boolean value of whether a logged user is currently logged.
func AlreadyLoggedIn(req *http.Request) bool {
	myCookie, err := req.Cookie("myCookie")
	if err != nil {
		return false
	}
	username := MapSessions[myCookie.Value]
	_, ok := MapUsers[username]
	return ok
}

//Function multiLogin takes in a username and returns a boolean value of whether user is currently logged in session.
func multiLogin(username string) bool {
	for _, v := range MapSessions {
		if v == username {
			return true
		}
	}
	return false
}
