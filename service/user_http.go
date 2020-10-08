package service

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"joshsoftware/go-e-commerce/db"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	logger "github.com/sirupsen/logrus"

	"golang.org/x/crypto/bcrypt"
)

// @Title listUsers
// @Description list all User
// @Router /users [get]
// @Accept  json
// @Success 200 {object}
// @Failure 400 {object}
func listUsersHandler(deps Dependencies) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		users, err := deps.Store.ListUsers(req.Context())
		if err != nil {
			logger.WithField("err", err.Error()).Error("error fetching data")
			responses(rw, http.StatusInternalServerError, errorResponse{
				Error: messageObject{
					Message: "Internal Server Error",
				},
			})
			return
		}

		responses(rw, http.StatusOK, successResponse{
			Data: users,
		})
	})

}

// @Title registerUser
// @Description registers new user
// @Router /register [post]
// @Accept  mulipart/form-data
// @Success 201 {object}
// @Failure 400 {object}
func registerUserHandler(deps Dependencies) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		var user db.User
		var decoder = schema.NewDecoder()

		err := req.ParseMultipartForm(15 << 20) // 15 MB Max File Size
		if err != nil {
			logger.WithField("err", err.Error()).Error("error while parsing the user form")
			responses(rw, http.StatusBadRequest, errorResponse{
				Error: messageObject{
					Message: "Invalid Form Data",
				},
			})
			return
		}

		// Retrive file from posted data
		formdata := req.MultipartForm

		// grab the filename
		contents := formdata.Value
		images := formdata.File["profile_image"]

		//grab product
		err = decoder.Decode(&user, contents)
		if err != nil {
			logger.WithField("err", err.Error()).Error("error while decoding product")
			responses(rw, http.StatusBadRequest, errorResponse{
				Error: messageObject{
					Message: "Invalid Form Contents",
				},
			})
			return
		}

		err = user.Validate()
		if err != nil {
			logger.WithField("err", err.Error()).Error("some data missing in request")
			responses(rw, http.StatusBadRequest, errorResponse{
				Error: messageObject{
					Message: "Incomplete Form Data",
				},
			})
			return
		}

		// check if profile image is present or not
		if len(images) > 0 {
			image, err := images[0].Open()

			defer image.Close()
			if err != nil {
				logger.WithField("err", err.Error()).Error("error while decoding image Data")
				responses(rw, http.StatusBadRequest, errorResponse{
					Error: messageObject{
						Message: "Invalid Image",
					},
				})
				return
			}

			extension := filepath.Ext(images[0].Filename)

			if len(extension) < 2 || len(extension) > 5 {
				err = fmt.Errorf("couldn't get extension of file!")
				logger.WithField("err", err.Error()).Error("error while getting image extension")
				responses(rw, http.StatusBadRequest, errorResponse{
					Error: messageObject{
						Message: "Unexpected File Type",
					},
				})
				return
			}

			fileName := strings.ReplaceAll(user.FirstName, " ", "")
			tempFile, err := ioutil.TempFile("assets/users", fileName+"-*"+string(extension))
			if err != nil {
				logger.WithField("err", err.Error()).Error("error while creating a temporary file")
				responses(rw, http.StatusInternalServerError, errorResponse{
					Error: messageObject{
						Message: "Internal Server Error: Failed To Set Profile Pic",
					},
				})
				return
			}
			defer tempFile.Close()

			imageBytes, err := ioutil.ReadAll(image)
			if err != nil {
				logger.WithField("err", err.Error()).Error("error while reading image File")
				responses(rw, http.StatusInternalServerError, errorResponse{
					Error: messageObject{
						Message: "Internal Server Error: Failed To Set Profile Pic",
					},
				})
				return
			}
			tempFile.Write(imageBytes)

			user.ProfileImage = tempFile.Name()

		}

		// For checking if user already registered
		check, _, err := deps.Store.CheckUserByEmail(req.Context(), user.Email)

		// If check true then user is already registered
		if check {
			logger.WithField("err", "error while registering new user: already exist")
			responses(rw, http.StatusConflict, errorResponse{
				Error: messageObject{
					Message: "User Already Registered",
				},
			})
			return
		}

		// For checking error occured while looking already registered user
		if err != nil && err != sql.ErrNoRows {
			logger.WithField("err", err.Error()).Error("error while looking existing user")
			responses(rw, http.StatusInternalServerError, errorResponse{
				Error: messageObject{
					Message: "Internal Server Error: Invalid Email ID",
				},
			})
			return
		}
		// creating hash of the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
		if err != nil {
			logger.WithField("err", err.Error()).Error("error while creating hash of the password")
			responses(rw, http.StatusInternalServerError, errorResponse{
				Error: messageObject{
					Message: "Invalid Password: Please Choose Other Password",
				},
			})
			return
		}
		user.Password = string(hashedPassword)
		// TODO send email for verification of user account
		user.IsVerified = true

		// Storing new user's data in database
		_, err = deps.Store.CreateNewUser(req.Context(), user)
		if err != nil {
			logger.WithField("err", err.Error()).Error("error in inserting user in database")
			responses(rw, http.StatusInternalServerError, errorResponse{
				Error: messageObject{
					Message: "Internal Server Error",
				},
			})
			return
		}

		responses(rw, http.StatusCreated, successResponse{
			Data: messageObject{
				Message: "User Successfully Registered",
			},
		})
		return
	})
}

//get user by id
func getUserHandler(deps Dependencies) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		//fetch usedId from request
		authToken := req.Header.Get("Token")
		payload, err := getDataFromToken(authToken)
		if err != nil {
			logger.WithField("err", err.Error()).Error("invalid token")
			responses(rw, http.StatusBadRequest, errorResponse{
				Error: messageObject{
					Message: "Bad Request",
				},
			})
			return
		}
		user, err := deps.Store.GetUser(req.Context(), int(payload.UserID))
		if err != nil {
			logger.WithField("err", err.Error()).Error("error in fetching userid")
			responses(rw, http.StatusNotFound, errorResponse{
				Error: messageObject{
					Message: "No User Found",
				},
			})
			return
		}
		responses(rw, http.StatusOK, successResponse{Data: user})
	})

}

//update user by id
func updateUserHandler(deps Dependencies) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		authToken := req.Header["Token"]
		payload, err := getDataFromToken(authToken[0])
		if err != nil {
			logger.WithField("err", err.Error()).Error("invalid token")
			responses(rw, http.StatusBadRequest, errorResponse{
				Error: messageObject{
					Message: "Invalid User Token",
				},
			})
			return
		}
		var user db.User
		var decoder = schema.NewDecoder()

		err = req.ParseMultipartForm(15 << 20) // 15 MB Max File Size
		if err != nil {
			logger.WithField("err", err.Error()).Error("error while parsing the product form")
			responses(rw, http.StatusBadRequest, errorResponse{
				Error: messageObject{
					Message: "Invalid Form Data",
				},
			})
			return
		}

		// Retrive file from posted data
		formdata := req.MultipartForm

		// grab the filename
		contents := formdata.Value
		images := formdata.File["profile_image"]

		//grab product
		err = decoder.Decode(&user, contents)
		if err != nil {
			logger.WithField("err", err.Error()).Error("error while decoding product")
			responses(rw, http.StatusBadRequest, errorResponse{
				Error: messageObject{
					Message: "Invalid form contents",
				},
			})
			return
		}

		dbUser, err := deps.Store.GetUser(req.Context(), int(payload.UserID))
		if err != nil {
			logger.WithField("err", err.Error()).Error("error in fetching userid")
			responses(rw, http.StatusNotFound, errorResponse{
				Error: messageObject{
					Message: "No User Found",
				},
			})
			return
		}
		previousFile := dbUser.ProfileImage

		// check if profile image is present or not
		if len(images) > 0 {
			image, err := images[0].Open()

			defer image.Close()
			if err != nil {
				logger.WithField("err", err.Error()).Error("error while decoding image data")
				responses(rw, http.StatusBadRequest, errorResponse{
					Error: messageObject{
						Message: "Invalid Image",
					},
				})
				return
			}

			extension := filepath.Ext(images[0].Filename)

			if len(extension) < 2 || len(extension) > 5 {
				logger.WithField("err", err.Error()).Error("error while getting image extension.")
				responses(rw, http.StatusBadRequest, errorResponse{
					Error: messageObject{
						Message: "Unexpected File Type",
					},
				})
				return
			}

			if user.FirstName == "" {
				user.FirstName = dbUser.FirstName
			}
			fileName := strings.ReplaceAll(user.FirstName, " ", "")
			tempFile, err := ioutil.TempFile("assets/users", fileName+"-*"+string(extension))
			if err != nil {
				logger.WithField("err", err.Error()).Error("error while creating a temporary file")
				responses(rw, http.StatusInternalServerError, errorResponse{
					Error: messageObject{
						Message: "Internal Server Error",
					},
				})
				return
			}
			defer tempFile.Close()

			imageBytes, err := ioutil.ReadAll(image)
			if err != nil {
				logger.WithField("err", err.Error()).Error("error while reading image File")
				responses(rw, http.StatusInternalServerError, errorResponse{
					Error: messageObject{
						Message: "Internal Server Error",
					},
				})
				return
			}
			tempFile.Write(imageBytes)

			dbUser.ProfileImage = tempFile.Name()

		}

		if user.Email != "" {
			logger.WithField("err", " cannot update email")
			responses(rw, http.StatusBadRequest, errorResponse{
				Error: messageObject{
					Message: "Cannot Update Email ID",
				},
			})
			return
		}

		err = dbUser.ValidatePatchParams(user)
		if err != nil {
			logger.WithField("err", err.Error()).Error("error in validate patch params")
			responses(rw, http.StatusBadRequest, errorResponse{
				Error: messageObject{
					Message: "Invalid User Data",
				},
			})
			return
		}

		err = deps.Store.UpdateUserByID(req.Context(), dbUser, int(payload.UserID))
		if err != nil {
			logger.WithField("err", err.Error()).Error("error while updating user's profile")
			responses(rw, http.StatusInternalServerError, errorResponse{
				Error: messageObject{
					Message: "Internal Server Error",
				},
			})
			return
		}

		err = os.Remove(previousFile)
		if err != nil {
			logger.WithField("err", err.Error()).Error("error while deleting the previous file of users profile image")
		}

		responses(rw, http.StatusOK, successResponse{Data: dbUser})
		return

	})

}

// @Title deleteuser
// @Description delete the user record by given id
// @Router /user/{id} [Delete]
// @Accept  json
// @Success 200 {object}
// @Failure 400 {object}
func deleteUserHandler(deps Dependencies) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		params := mux.Vars(req)
		userID, err := strconv.Atoi(params["id"])
		if err != nil {
			logger.WithField("err", err.Error()).Error("error in fetching userid")
			responses(rw, http.StatusBadRequest, errorResponse{
				Error: messageObject{
					Message: "Invalid UserID",
				},
			})
			return
		}

		//check if the provided id is not of another Admin
		user, err := deps.Store.GetUser(req.Context(), userID)
		if err != nil {
			logger.WithField("err", err.Error()).Error("no user in database")
			responses(rw, http.StatusNotFound, errorResponse{
				Error: messageObject{
					Message: "No User Found",
				},
			})
			return
		}

		//check if the users is admin
		if user.IsAdmin == true {
			logger.WithField("err", err.Error()).Error("user trying to access admin resource")
			responses(rw, http.StatusForbidden, errorResponse{
				Error: messageObject{
					Message: "Only Admin Can Delete User",
				},
			})
			return
		}

		err = deps.Store.DeleteUserByID(req.Context(), userID)
		if err != nil {
			logger.WithField("err", err.Error()).Error("error in deleting user data")
			responses(rw, http.StatusInternalServerError, errorResponse{
				Error: messageObject{
					Message: "Internal Server Error",
				},
			})
			return
		}

		responses(rw, http.StatusOK, successResponse{
			Data: "User Deleted Successfully",
		})

		return

	})
}

// @Title disableuser
// @Description disable the user record by given id
// @Router /user/disable/{id} [Patch]
// @Accept  json
// @Success 200 {object}
// @Failure 400 {object}
func disableUserHandler(deps Dependencies) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		params := mux.Vars(req)
		userID, err := strconv.Atoi(params["id"])
		if err != nil {
			logger.WithField("err", err.Error()).Error("error in fetching userid")
			responses(rw, http.StatusBadRequest, errorResponse{
				Error: messageObject{
					Message: "Bad Request",
				},
			})
			return
		}

		//check if the provided id is not of another Admin
		user, err := deps.Store.GetUser(req.Context(), userID)
		if err != nil {
			logger.WithField("err", err.Error()).Error("no user in database")
			responses(rw, http.StatusNotFound, errorResponse{
				Error: messageObject{
					Message: "No User Found",
				},
			})
			return
		}

		//check if the user is admin
		if user.IsAdmin == true {
			logger.WithField("err", err.Error()).Error("user trying to disable admin")
			responses(rw, http.StatusForbidden, errorResponse{
				Error: messageObject{
					Message: "Admin Cannot be Disabled/Enabled",
				},
			})
			return
		}

		err = deps.Store.DisableUserByID(req.Context(), userID)
		if err != nil {
			logger.WithField("err", err.Error()).Error("error in disabling user data")
			responses(rw, http.StatusInternalServerError, errorResponse{
				Error: messageObject{
					Message: "Internal Server Error",
				},
			})
			return
		}

		responses(rw, http.StatusOK, successResponse{
			Data: "record disabled Successfully",
		})

		return

	})
}

// @Title enableuser
// @Description disable the user record by given id
// @Router /user/disable/{id} [Patch]
// @Accept  json
// @Success 200 {object}
// @Failure 400 {object}
func enableUserHandler(deps Dependencies) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		params := mux.Vars(req)
		userID, err := strconv.Atoi(params["id"])
		if err != nil {
			logger.WithField("err", err.Error()).Error("error in fetching userid")
			responses(rw, http.StatusBadRequest, errorResponse{
				Error: messageObject{
					Message: "Bad Request",
				},
			})
			return
		}

		//check if the provided id is not of another Admin
		user, err := deps.Store.GetUser(req.Context(), userID)
		if err != nil {
			logger.WithField("err", err.Error()).Error("no user in database")
			responses(rw, http.StatusNotFound, errorResponse{
				Error: messageObject{
					Message: "No User Found",
				},
			})
			return
		}

		//check if the user is admin
		if user.IsAdmin == true {
			logger.WithField("err", err.Error()).Error("user trying to enable admin")
			responses(rw, http.StatusForbidden, errorResponse{
				Error: messageObject{
					Message: "Admin Cannot be Enabled/Disabled",
				},
			})
			return
		}

		err = deps.Store.EnableUserByID(req.Context(), userID)
		if err != nil {
			logger.WithField("err", err.Error()).Error("error in enabling user data")
			responses(rw, http.StatusInternalServerError, errorResponse{
				Error: messageObject{
					Message: "Internal Server Error",
				},
			})
			return
		}

		responses(rw, http.StatusOK, successResponse{
			Data: "user enabled Successfully",
		})

		return

	})
}

func verifyUserHandler(deps Dependencies) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		//fetching the token from header
		authToken := req.Header.Get("Token")
		//fetching details from the token
		payload, err := getDataFromToken(authToken)
		if err != nil {
			logger.WithField("err", err.Error()).Error("invalid token")
			responses(rw, http.StatusBadRequest, errorResponse{
				Error: messageObject{
					Message: "Bad Request",
				},
			})
			return
		}

		// reading data from body
		reqBody, err := ioutil.ReadAll(req.Body)
		if err != nil {
			logger.WithField("err", err.Error()).Error("error in reading request body")
			responses(rw, http.StatusBadRequest, errorResponse{
				Error: messageObject{
					Message: "Invalid Request",
				},
			})
			return
		}

		user := db.User{}
		err = json.Unmarshal(reqBody, &user)
		if err != nil {
			logger.WithField("err", err.Error()).Error("error while unmarshalling request json")
			responses(rw, http.StatusBadRequest, errorResponse{
				Error: messageObject{
					Message: "Invalid Fields",
				},
			})
			return
		}

		// creating hash of the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error while creating hash of the password")
			responses(rw, http.StatusBadRequest, errorResponse{
				Error: messageObject{
					Message: "Invalid Password",
				},
			})
			return
		}
		user.Password = string(hashedPassword)

		dbuser := db.User{}
		//check if the provided id is not of another Admin
		dbuser, err = deps.Store.GetUser(req.Context(), int(payload.UserID))
		if dbuser.IsAdmin {
			logger.WithField("err", err.Error()).Error("admin cannot access this resource")
			responses(rw, http.StatusUnauthorized, errorResponse{
				Error: messageObject{
					Message: "Unauthorized User",
				},
			})
			return
		}

		if dbuser.IsDisabled {
			responses(rw, http.StatusForbidden, errorResponse{
				Error: messageObject{
					Message: "User Disabled",
				},
			})
			return
		}

		if err != nil {
			if err == sql.ErrNoRows {
				logger.WithField("err", err.Error()).Error("no user found")
				responses(rw, http.StatusNotFound, errorResponse{
					Error: messageObject{
						Message: "No user found",
					},
				})
				return
			}
			logger.WithField("err", err.Error()).Error("error in getting user from token")
			responses(rw, http.StatusBadRequest, errorResponse{
				Error: messageObject{
					Message: "Bad Request",
				},
			})
			return
		}

		if dbuser.IsVerified {
			logger.WithField("err", " email is already verified")
			responses(rw, http.StatusConflict, errorResponse{
				Error: messageObject{
					Message: "Already Verified",
				},
			})
			return
		}

		err = deps.Store.VerifyUserByID(req.Context(), int(payload.UserID))
		if err != nil {
			logger.WithField("err", err.Error()).Error("error while verifing user")
			responses(rw, http.StatusInternalServerError, errorResponse{
				Error: messageObject{
					Message: "Internal Server Error: Invalid UserID",
				},
			})
			return
		}

		err = deps.Store.SetUserPasswordByID(req.Context(), user.Password, int(payload.UserID))
		if err != nil {
			logger.WithField("err", err.Error()).Error("error while setting user password")
			responses(rw, http.StatusInternalServerError, errorResponse{
				Error: messageObject{
					Message: "Internal Server Error: Invalid Password",
				},
			})
			return
		}

		responses(rw, http.StatusOK, successResponse{
			Data: "Successfully verified, please login to continue",
		})
		return
	})
}
