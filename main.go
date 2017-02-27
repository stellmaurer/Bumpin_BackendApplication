package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/tables", tables)
	http.HandleFunc("/person", person)
	http.HandleFunc("/party", party)
	http.HandleFunc("/bar", bar)
	http.HandleFunc("/myParties", myParties)
	http.HandleFunc("/barsCloseToMe", barsCloseToMe)
	http.ListenAndServe(":8080", nil)
}
