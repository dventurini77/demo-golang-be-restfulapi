package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// User model
type User struct {
	ID    int `sql:"primary_key;AUTO_INCREMENT"`
	Name  string
	Email string
}

var psqlInfo string

func getUsers(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open("postgres", psqlInfo)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	var users []User
	db.Find(&users)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func getUser(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open("postgres", psqlInfo)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	var params = mux.Vars(r) // Get params
	var user User

	db.Debug().Where("ID = ?", params["id"]).First(&user)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&user)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func createUser(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open("postgres", psqlInfo)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)

	db.Debug().Create(&user)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&user)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open("postgres", psqlInfo)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	var params = mux.Vars(r) // Get params
	var upduser, user User
	_ = json.NewDecoder(r.Body).Decode(&upduser)

	db.Debug().Where("ID = ?", params["id"]).First(&user).Update(&upduser)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&user)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open("postgres", psqlInfo)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	var params = mux.Vars(r) // Get params
	var user User

	db.Debug().Where("ID = ?", params["id"]).Delete(&user)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&user)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/api/users", getUsers).Methods("GET")
	myRouter.HandleFunc("/api/users/{id}", getUser).Methods("GET")
	myRouter.HandleFunc("/api/users", createUser).Methods("POST")
	myRouter.HandleFunc("/api/users/{id}", updateUser).Methods("PUT")
	myRouter.HandleFunc("/api/users/{id}", deleteUser).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func initialMigration() {
	db, err := gorm.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect database")
	}
	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&User{})

	db.LogMode(true)
	db.SetLogger(log.New(os.Stdout, "\r\n", 0))
}

func initialDbConnection() {
	var host = os.Getenv("PGDATABASE_HOST")
	var port = os.Getenv("PGDATABASE_PORT")
	var user = os.Getenv("PGDATABASE_USER")
	var password = os.Getenv("PGDATABASE_PASS")
	var dbname = os.Getenv("PGDATABASE_NAME")

	psqlInfo = fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	log.Println(psqlInfo)
}

func main() {
	log.Println("Service RESTful::User started")

	initialDbConnection()

	initialMigration()

	// Handle Subsequent requests
	handleRequests()
}
