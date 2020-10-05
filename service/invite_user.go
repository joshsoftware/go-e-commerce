package service

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/kataras/go-mailer"
	logger "github.com/sirupsen/logrus"
	// "golang.org/x/crypto/bcrypt"
	"html/template"
	"io/ioutil"
	"joshsoftware/go-e-commerce/config"
	"joshsoftware/go-e-commerce/db"
	"log"
	// "math/rand"
	"net/http"
	// "strings"
	// "time"
)

//Email Struct
type Email struct {
	Email []string `json:"email"`
}

func inviteUsersHandler(deps Dependencies) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		var email Email

		reqBody, err := ioutil.ReadAll(req.Body)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error in reading request body")
			responses(rw, http.StatusInternalServerError, errorResponse{
				Error: messageObject{
					Message: "Internal Server Error",
				},
			})
			return
		}

		err = json.Unmarshal(reqBody, &email)

		if err != nil {
			logger.WithField("err", err.Error()).Error("JSON Decoding Failed")
			responses(rw, http.StatusBadRequest, errorResponse{
				Error: messageObject{
					Message: "JSON Decoding Failed",
				},
			})
			return
		}

		//Flag to check if user already exists
		existFlag := false
		var existingUsers []string

		for _, emailID := range email.Email {
			// For checking if user already registered
			var check bool
			check, _, err = deps.Store.CheckUserByEmail(req.Context(), emailID)

			// If check true then user is already registered
			if check {
				existFlag = true
				existingUsers = append(existingUsers, emailID)
				log.Printf("\nuser with e-mail id %d already exists", emailID)
				continue
			}
			// For checking error occured while looking already registered user
			if err != nil && err != sql.ErrNoRows {
				logger.WithField("err", err.Error()).Error("Error while looking existing user")
				continue
			}

			// randPassword := randomPassGenerator()
			// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(randPassword), 8)
			// if err != nil {
			// 	logger.WithField("err", err.Error()).Error("Error while creating hash of the password")
			// 	continue
			// }

			user := db.User{}
			dbUser := db.User{}

			user.Email = emailID

			dbUser, err = deps.Store.CreateNewUser(req.Context(), user)
			if err != nil {
				logger.WithField("err", err.Error()).Error("Error in inserting user in database")
				continue
			}
			token, err := generateJwt(dbUser.ID, dbUser.IsAdmin)

			temp, err := template.ParseFiles("assets/templates/mail_invite.html")
			if err != nil {
				fmt.Println("****Error has occured****")
			}

			var body bytes.Buffer
			var verificationURL = "https://joshreact-e-commerce.herokuapp.com/verify?Token=" + token

			temp.Execute(&body, struct {
				Email           string
				SetPasswordLink string
			}{
				Email:           emailID,
				SetPasswordLink: verificationURL,
			})

			subject := "Registration Successful"
			mail(subject, body.String(), emailID)

		}
		if existFlag {
			responses(rw, http.StatusConflict, Email{
				Email: existingUsers,
			})
			return
		}

		responses(rw, http.StatusOK, successResponse{
			Data: messageObject{
				Message: "All Users Registered Successfully",
			},
		})
		return

	})
}

func mail(subject string, content string, sendTo string) {

	host, port, username, from, pass := config.MailerConfig()

	config := mailer.Config{
		Host:       host,
		Username:   username,
		Password:   pass,
		FromAddr:   from,
		Port:       port,
		UseCommand: false,
	}

	sender := mailer.New(config)

	err := sender.Send(subject, content, sendTo)

	if err != nil {
		fmt.Println("Error While Sending Email : " + err.Error())
	}
}
