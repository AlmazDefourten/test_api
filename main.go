package main

import (
	"Pears/avl"
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"net/http"
)

var users []User

func main() {

	r := mux.NewRouter()

	// подключение к бд
	connStr := "user=postgres password=mypas dbname=pears sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	//получение данных с бд users
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		panic(err)
	}

	tree := avl.New()

	defer rows.Close()
	for rows.Next() {
		row := user{}
		err := rows.Scan(&row.id, &row.name, &row.mail)
		if err != nil {
			panic(err)
		}
		fmt.Println(row.id, row.name, row.mail)
		tree.Insert(row.name)
	}

	r.HandleFunc("/users", getUsers).Methods("GET")
	r.HandleFunc("/users/{id}", getUser).Methods("GET")
	r.HandleFunc("/users", createUser).Methods("POST")
	r.HandleFunc("/users/{id}", updateUser).Methods("PUT")
	r.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8000", r))

}

func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

}
