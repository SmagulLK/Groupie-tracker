package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	Userregistration "groupie/User"
	CheckRequest "groupie/internal/checkRequest"
	"groupie/internal/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"text/template"
)

var Db *sql.DB
var username string
var email string

func main() {
	mux := http.NewServeMux()
	initDB()
	if err := initConfig(); err != nil {
		log.Fatal(err)
	}

	defer Db.Close()

	mux.HandleFunc("/", handlers.HomePage)
	mux.HandleFunc("/pageTwo/", handlers.PageTwo)
	mux.HandleFunc("/search", handlers.SearchHandler)
	mux.HandleFunc("/reg", RegisterHandler)
	mux.HandleFunc("/logout", LogoutHandler)
	mux.HandleFunc("/login", LoginHandler)
	mux.HandleFunc("/account", AccountHandler)
	mux.HandleFunc("/delete", DeleteHandler)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("ui/static"))))
	log.Print("Start server http://127.0.0.1:8000")
	go func() {
		log.Fatal(http.ListenAndServe(":8000", mux))
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	fmt.Println("Shutting down")

}

// initDB initializes the database connection
func initDB() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		"localhost", 5436, "postgres", "Admin", "postgres")
	fmt.Println(psqlInfo)
	var err error
	Db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	err = Db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected to the database!")
}

// createUser creates a new user in the database
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем, что запрос является POST-запросом
	if r.Method == http.MethodGet {
		tmpl, _ := template.ParseFiles("./ui/templates/register.html")
		_ = tmpl.Execute(w, nil)
		return
	}
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		fmt.Println(username)
		email := r.FormValue("email")
		password := r.FormValue("password")
		if !Userregistration.CheckUsername(username, Db) || !Userregistration.CheckEmail(email, Db) {
			fmt.Println("Username created")
			Userregistration.CreateUser(username, email, password, Db)

			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else {
			fmt.Println("Username or email already exists")
			http.Redirect(w, r, "/reg", http.StatusSeeOther)
		}

		//http.Redirect(w, r, "/login", http.StatusOK)
	}

	// Парсим форму
	//err := r.ParseForm()
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}

	// Создаем нового пользователя

}
func LoginHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		tmpl, _ := template.ParseFiles("./ui/templates/login.html")
		_ = tmpl.Execute(w, nil)
		return
	}
	if r.Method == http.MethodPost {
		username = r.FormValue("username")
		password := r.FormValue("password")
		if Userregistration.CheckUsername(username, Db) {
			if Userregistration.CheckPassword(username, password, Db) {
				CheckRequest.LoggedIn = true
				fmt.Println("Login successful")

				http.Redirect(w, r, "/", http.StatusSeeOther)
			} else {
				fmt.Println("Wrong password")
				http.Redirect(w, r, "/login", http.StatusSeeOther)
			}
		} else {
			fmt.Println("Wrong username")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}
	}
}
func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
func AccountHandler(w http.ResponseWriter, r *http.Request) {
	if CheckRequest.LoggedIn {
		tmpl, err1 := template.ParseFiles("./ui/templates/myacc.html")
		if err1 != nil {
			fmt.Println(err1)
			return
		}

		_ = Db.QueryRow("SELECT email FROM users WHERE username=$1", username).Scan(&email)
		UserR := Userregistration.User{username, email}
		err := tmpl.Execute(w, UserR)
		fmt.Println(UserR)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(UserR)

	} else {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
func ForgetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, _ := template.ParseFiles("./ui/templates/forget.html")
		_ = tmpl.Execute(w, nil)
		return
	}
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		if Userregistration.CheckEmail(email, Db) {
			fmt.Println("they is here")
			password := Userregistration.GeneratePassword()
			Userregistration.UpdatePassword(email, password, Db)

		}
	}

}
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if CheckRequest.LoggedIn {
		CheckRequest.LoggedIn = false
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	if CheckRequest.LoggedIn {
		Userregistration.DeleteUser(username, Db)
		CheckRequest.LoggedIn = false
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
