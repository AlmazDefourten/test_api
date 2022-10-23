package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	id    int
	login string
	mail  string
}

var DB *gorm.DB

func main() {
	// подключение к бд
	dsn := "host=localhost user=postgres password=mypas dbname=pears port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	DB = db

	// создание роутера
	r := mux.NewRouter()

	r.HandleFunc("/users", getUsers).Methods("GET")
	r.HandleFunc("/users/{id}", getUser).Methods("GET")
	r.HandleFunc("/users", createUser).Methods("POST")
	r.HandleFunc("/users/{id}", updateUser).Methods("PUT")
	r.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")
	http.ListenAndServe(":4000", r)
	log.Println("API is running!")
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	// получение данных из бд из таблицы users
	var users []User
	DB.Find(&users)

	// вывод данных в консоль
	log.Println(users)

	// возвращение ответа
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	// чтение параметров запроса
	params := mux.Vars(r)

	// получение пользователя из бд
	var user User
	if err := DB.First(&user, params["id"]).Error; err != nil {
		log.Fatalln(err)
	}

	// возвращение ответа
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	// чтение тела запроса
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Fatalln(err)
	}

	// преобразование тела запроса в структуру
	var user *User
	json.Unmarshal(body, &user)

	// добавление пользователя в бд если нет ошибок
	if err := DB.Create(user).Error; err != nil {
		log.Fatalln(err)
	}

	// возвращение ответа
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("User created successfully")
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	// чтение параметров запроса
	params := mux.Vars(r)

	// чтение тела запроса
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Fatalln(err)
	}

	// преобразование тела запроса в структуру
	var user *User
	json.Unmarshal(body, &user)

	// обновление пользователя в бд если нет ошибок
	if err := DB.Model(&User{}).Where("id = ?", params["id"]).Updates(user).Error; err != nil {
		log.Fatalln(err)
	}

	// возвращение ответа
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("User updated successfully")
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	// чтение параметров запроса
	params := mux.Vars(r)

	// удаление пользователя из бд
	if err := DB.Delete(&User{}, params["id"]).Error; err != nil {
		log.Fatalln(err)
	}

	// возвращение ответа
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("User deleted successfully")
}
