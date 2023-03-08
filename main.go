package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"github.com/gorilla/mux"
	_"github.com/lib/pq"
)

type User struct{
	ID  int `json:"id"`
	Name  string `json:"name"`
	Email  string `json:"email"`
}

func main(){
	//Connect to database
	db , err := sql.Open("Postrgres", os.Getenv("DATABASE_URL"))
	if err != nil{
		log.Fatal(err)
	}
	defer db.Close()

	//Create routes
	router := mux.NewRouter()
	router.HandleFunc("/users", getUsers(db)).Methods("GET")
	router.HandleFunc("/users/{id}", getUser(db)).Methods("GET")
	router.HandleFunc("/users", createUSer(db)).Methods("POST")
	router.HandleFunc("/users/{id}", updateUser(db)).Methods("PUT")
	router.HandleFunc("/users/{id}", deleteUser(db)).Methods("DELETE")

	//Start server
	log.Fatal(http.ListenAndServe(":8000", jsonContentTypeMiddleWare(router)))

}

func jsonContentTypeMiddleWare(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}


func getUsers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT * FROM users")
		if err != nil{
			log.Fatal(err)
		}
		defer rows.Close()

		users := []User{}
		for rows.Next() {
			var u User
			if err:= rows.Scan(&u.ID, &u.Name, &u.Email); err != nil{
				log.Fatal(err)
			}
			users = append(users, u)
		}
		if err := rows.Err(); err != nil{
			log.Fatal(err)
		}

		json.NewEncoder(w).Encode(users)
	}
}


//Get user by ID
func getUser(db *sql.DB) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var u User
		err := db.QueryRow("SELECT * FROM users WHERE id = $1", id).Scan(&u.ID, &u.Name, &u.Email)
		if err != nil{
			log.Fatal(err)
		}
		
		json.NewEncoder(w).Encode(u)
	}
}