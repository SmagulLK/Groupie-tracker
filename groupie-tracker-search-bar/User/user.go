package User

import (
	sql "database/sql"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"math/rand"
	"net/mail"
	"net/smtp"
)

type User struct {
	Username string
	Email    string
}

func CreateUser(username, email, password string, db *sql.DB) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {

		return err
	}
	fmt.Println(len(hash))
	_, err = db.Exec("INSERT INTO Users (username, user_password, email) VALUES ($1, $2, $3)", username, hash, email)
	if err != nil {
		return err
	}

	return nil
}

// checkUsername checks if the given username is already taken
func CheckUsername(username string, db *sql.DB) bool {
	var id int
	err := db.QueryRow("SELECT userid FROM users WHERE username = $1", username).Scan(&id)
	if err == sql.ErrNoRows {
		return false
	}
	if err != nil {
		log.Fatal(err)
	}
	return true
}

// checkEmail checks if the given email is already taken
func CheckEmail(email string, db *sql.DB) bool {
	var id int
	err := db.QueryRow("SELECT userid FROM users WHERE email = $1", email).Scan(&id)
	if err == sql.ErrNoRows {
		return false
	}
	if err != nil {
		log.Fatal(err)
	}
	return true
}
func CheckPassword(username string, password string, Db *sql.DB) bool {
	var pass string
	err := Db.QueryRow("SELECT user_password FROM users WHERE username = $1", username).Scan(&pass)
	if err != nil {
		log.Fatal(err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(pass), []byte(password))
	fmt.Println(err)
	return err == nil
}
func GeneratePassword() string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789@#$%^&*()_+")
	s := make([]rune, 10)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}

	return string(s)
}
func UpdatePassword(email, password string, Db *sql.DB) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = Db.Exec("UPDATE users SET user_password = $1 WHERE email =$2", hash, email)
	if err != nil {
		return err
	}
	return nil
}

// func SendEmail(email, password string) error {
//
//		return nil
//	}
func DeleteUser(username string, Db *sql.DB) error {
	_, err := Db.Exec("DELETE FROM users WHERE username = $1", username)
	if err != nil {
		return err
	}
	return nil
}
func SendEmail(to string, password string) {
	// Set up the email
	from := mail.Address{"Reset Password", "smagul.alkey@gmail.com"}
	toAddress := mail.Address{"", to}
	subject := "Reset Password"
	body := "Your new password is: " + password

	// Set up the message
	header := make(map[string]string)
	header["From"] = from.String()
	header["To"] = toAddress.String()
	header["Subject"] = subject
	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// Connect to the SMTP server
	smtpServer := "smtp.elasticemail.com:2525"
	auth := smtp.PlainAuth("", "smagul.alkey@gmail.com", "2C7B18DE61F94A629AD7A1E04126C86310AE", "smtp.elasticemail.com")
	err := smtp.SendMail(smtpServer, auth, from.Address, []string{to}, []byte(message))
	if err != nil {
		log.Fatal(err)
	}
}
