package main

import (
	"fmt"
	"log"
	"net/http"
	"encoding/json"
	"os"
	"html/template"
	"strings"
)

type potluckItem struct {
	Name string
	Dish string
}

var potluck []potluckItem
var dataFilename string = "potluck.json"

func loadData() error {
	content, err := os.ReadFile(dataFilename)
	if err != nil {
		return err
	}

	err = json.Unmarshal(content, &potluck)
	if err != nil {
		return err
	}

	return nil
}

func saveData() error {
	jsonString, err := json.Marshal(potluck)
	if err != nil {
		return err
	}

	err = os.WriteFile(dataFilename, jsonString, 0600)
	if err != nil {
		return err
	}

	return nil
}

func handleIt(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		switch r.URL.Path {
		case "/favicon.ico":
			http.ServeFile(w, r, "site" + r.URL.Path)
		case "/js/snowstorm.js":
			http.ServeFile(w, r, "site" + r.URL.Path)
		case "/img/hero_opt.png":
			http.ServeFile(w, r, "site" + r.URL.Path)
		case "/img/bg_opt.png":
			http.ServeFile(w, r, "site" + r.URL.Path)
		default:
			message := "404 not found (" + r.URL.Path + ")"
			http.Error(w, message, http.StatusNotFound)
		}
		return
	}

	switch r.Method {
	case "GET":
		err := loadData()
		if err != nil {
			message := "500 internal server error (" + r.URL.Path + ")"
			http.Error(w, message, http.StatusInternalServerError)
			return
		}
		t, _ := template.ParseFiles("index.html")
		err = t.Execute(w, potluck)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	case "POST":
		if err := r.ParseForm(); err != nil {
			log.Fatalf("%v", err)
			return
		}

		name := r.FormValue("name")
		dish := r.FormValue("dish")
		code := r.FormValue("code")

		if strings.ToLower(code) != "" {
			http.ServeFile(w, r, "site/error.html")
			return
		}

		pItem := potluckItem{Name: name, Dish: dish}
		potluck = append(potluck, pItem)
		saveData()
		http.Redirect(w, r, "/", http.StatusFound)
		return

	default:
		message := "500 internal server error (" + r.URL.Path + ")"
		http.Error(w, message, http.StatusInternalServerError)
		return
	}
}

func main() {
	http.HandleFunc("/", handleIt)
	fmt.Printf("Starting server for HTTP requests...\n")
	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatal(err)
	}
}
