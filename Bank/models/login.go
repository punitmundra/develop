package models

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	u "../utils"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// Token :
type Token struct {
	UserID uint
	jwt.StandardClaims
}

//User struct
type User struct {
	gorm.Model
	Email    string `json:"email"`
	Password string `json:"password"`
	JWTToken string `json:"token";sql:"-"`
}

// Create the entry in db : user table
func (account *User) Create() map[string]interface{} {

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	account.Password = string(hashedPassword)

	tokenid := &Token{UserID: account.ID}
	fmt.Println("tokenid : ", tokenid)
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tokenid)
	fmt.Println("jwt token :", token)

	// we can read this from .env file
	// this pwd will update in every 5 days
	// TODO : need to write a cron for the password updation
	// if user will not able to get the details.package models
	// user need to update the token or request for the new token
	// we need to write the gettoken/{email} api for the user

	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	account.JWTToken = tokenString

	GetDB().Create(account)

	if account.ID <= 0 {
		return u.Message(false, "Failed to create User, might be connection error.")
	}
	//removing the  password from the response
	account.Password = ""

	response := u.Message(true, "Account has been created")
	response["account"] = account
	return response
}

// Login api will return the token
func Login(email, password string) map[string]interface{} {

	account := &User{}
	err := GetDB().Table("users").Where("email = ?", email).First(account).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "Email address not found")
		}
		return u.Message(false, "Connection error. Please retry")
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		return u.Message(false, "Invalid login credentials. Please try again")
	}
	//Worked! Logged In
	account.Password = ""

	//Create JWT token
	tokenID := &Token{UserID: account.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tokenID)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	account.JWTToken = tokenString //Store the token in the response

	resp := u.Message(true, "Logged In")
	resp["account"] = account
	return resp
}

// getToken : get the token
func (account *User) getToken() map[string]interface{} {
	temp := &User{}
	err := GetDB().Table("users").Where("email = ?", account.Email).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "Email address not found")
	}
	start := time.Now()
	resp := u.Message(true, "Token expires by :"+start.AddDate(0, 0, 5).String())
	resp["token"] = temp.JWTToken
	return resp
}

func (account *User) updateToken() map[string]interface{} {
	temp := &User{}
	fmt.Println(" THis is the account.Email:", account.Email)

	err := GetDB().Table("users").Where("email = ?", account.Email).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "Email address not found")
	}

	tokenID := &Token{UserID: account.ID}
	fmt.Println(" THis is the tokenID:", tokenID)

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tokenID)
	fmt.Println(" THis is the token:", token)

	// token pwd is update by cron so
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	fmt.Println(" THis is the token string:", tokenString)

	temp.JWTToken = tokenString //Store the token in the response
	GetDB().Model(&temp).Where("email = ?", account.Email).Update(tokenString)
	GetDB().Save(&temp)
	resp := u.Message(true, "Token updated ")
	resp["token"] = tokenString
	return resp
}

// CreateUser : clouser for Create Account
var CreateUser = func(w http.ResponseWriter, r *http.Request) {

	account := &User{}
	err := json.NewDecoder(r.Body).Decode(account)
	if err != nil {
		fmt.Println("error:", err, " acc:", account)
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	resp := account.Create()
	u.Respond(w, resp)
}

// Authenticate : clouser for auth
var Authenticate = func(w http.ResponseWriter, r *http.Request) {

	account := &User{}
	err := json.NewDecoder(r.Body).Decode(account)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	resp := Login(account.Email, account.Password)
	u.Respond(w, resp)
}

// GetToken : clouser for Create Account
var GetToken = func(w http.ResponseWriter, r *http.Request) {

	account := &User{}
	err := json.NewDecoder(r.Body).Decode(account)
	if err != nil {
		fmt.Println("error:", err, " acc:", account)
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	resp := account.getToken()
	u.Respond(w, resp)
}

// UpdateToken : clouser for Create Account
var UpdateToken = func(w http.ResponseWriter, r *http.Request) {

	account := &User{}
	err := json.NewDecoder(r.Body).Decode(account)
	if err != nil {
		fmt.Println("error:", err, " acc:", account)
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	resp := account.updateToken()
	u.Respond(w, resp)
}
