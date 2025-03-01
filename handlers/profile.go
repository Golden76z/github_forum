package handlers

import (
	"forum/cmd/lib"
	"net/http"
)

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	// Storing Db data into a variable
	db := lib.GetDB()

	// Checking the cookie values
	cookie, _ := r.Cookie("session_id")

	// Getting the User Data
	data := lib.GetData(db, cookie.Value, "logged", "profile", w, r)

	// Redirect User to the profile html page and sending the data to it
	lib.RenderTemplate(w, "layout/index", "page/profile", data)
}
