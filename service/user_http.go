package service

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"joshsoftware/go-e-commerce/db"
	"net/http"
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
			logger.WithField("err", err.Error()).Error("Error fetching data")
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
			logger.WithField("err", err.Error()).Error("Error while parsing the Product form")
			responses(rw, http.StatusBadRequest, errorResponse{
				Error: messageObject{
					Message: "Invalid Form Data!",
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
			logger.WithField("err", err.Error()).Error("Error while decoding product")
			responses(rw, http.StatusBadRequest, errorResponse{
				Error: messageObject{
					Message: "Invalid form contents",
				},
			})
			return
		}

		err = user.Validate()
		if err != nil {
			logger.WithField("err", err.Error()).Error("Some data missing in request")
			responses(rw, http.StatusBadRequest, errorResponse{
				Error: messageObject{
					Message: err.Error(),
				},
			})
			return
		}

		// check if profile image is present or not
		if len(images) > 0 {
			image, err := images[0].Open()

			defer image.Close()
			if err != nil {
				logger.WithField("err", err.Error()).Error("Error while decoding image Data")
				responses(rw, http.StatusBadRequest, errorResponse{
					Error: messageObject{
						Message: "Invalid Image !",
					},
				})
				return
			}

			extension := filepath.Ext(images[0].Filename)

			if len(extension) < 2 || len(extension) > 5 {
				err = fmt.Errorf("Couldn't get extension of file!")
				logger.WithField("err", err.Error()).Error("Error while getting image Extension.")
				responses(rw, http.StatusBadRequest, errorResponse{
					Error: messageObject{
						Message: "Re-check the image file extension!",
					},
				})
				return
			}

			fileName := strings.ReplaceAll(user.FirstName, " ", "")
			tempFile, err := ioutil.TempFile("assets/users", fileName+"-*"+string(extension))
			if err != nil {
				logger.WithField("err", err.Error()).Error("Error while Creating a Temporary File")
				responses(rw, http.StatusInternalServerError, errorResponse{
					Error: messageObject{
						Message: "Couldn't  create temporary storage!",
					},
				})
				return
			}
			defer tempFile.Close()

			imageBytes, err := ioutil.ReadAll(image)
			if err != nil {
				logger.WithField("err", err.Error()).Error("Error while reading image File")
				responses(rw, http.StatusInternalServerError, errorResponse{
					Error: messageObject{
						Message: "Couldn't read the image file!",
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
			e := errorResponse{
				Error: "user already registered",
			}
			respBytes, err := json.Marshal(e)
			if err != nil {
				logger.WithField("err", err.Error()).Error("Error while marshalling error msg ")
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
			rw.Header().Add("Content-Type", "application/json")
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write(respBytes)
			return
		}

		// For checking error occured while looking already registered user
		if err != nil && err != sql.ErrNoRows {
			logger.WithField("err", err.Error()).Error("Error while looking existing user")

			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		// creating hash of the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error while creating hash of the password")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		user.Password = string(hashedPassword)

		// Storing new user's data in database
		_, err = deps.Store.CreateNewUser(req.Context(), user)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error in inserting user in database")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		msg := successResponse{
			Data: "user successfully registered",
		}
		respBytes, err := json.Marshal(msg)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error while marshalling success msg ")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		rw.Header().Add("Content-Type", "application/json")
		rw.WriteHeader(http.StatusCreated)
		rw.Write(respBytes)
		return
	})
}

//get user by id
func getUserHandler(deps Dependencies) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		//fetch usedId from request
		authToken := req.Header.Get("Token")
		userID, _, _, err := getDataFromToken(authToken)
		if err != nil {
			responses(rw, http.StatusUnauthorized, errorResponse{
				Error: messageObject{
					Message: "Unauthorized User",
				},
			})
			return
		}
		user, err := deps.Store.GetUser(req.Context(), int(userID))
		if err != nil {
			logger.WithField("err", err.Error()).Error("error while fetching User")
			rw.WriteHeader(http.StatusNotFound)
			responses(rw, http.StatusNotFound, errorResponse{
				Error: messageObject{
					Message: "id Not Found",
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
		userID, _, _, err := getDataFromToken(authToken[0])
		if err != nil {
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte("Unauthorized"))
			return
		}
		var user db.User
		var decoder = schema.NewDecoder()

		err = req.ParseMultipartForm(15 << 20) // 15 MB Max File Size
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error while parsing the Product form")
			responses(rw, http.StatusBadRequest, errorResponse{
				Error: messageObject{
					Message: "Invalid Form Data!",
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
			logger.WithField("err", err.Error()).Error("Error while decoding product")
			responses(rw, http.StatusBadRequest, errorResponse{
				Error: messageObject{
					Message: "Invalid form contents",
				},
			})
			return
		}

		dbUser, err := deps.Store.GetUser(req.Context(), int(userID))
		if err != nil {
			logger.WithField("err", err.Error()).Error("error while fetching User")
			rw.WriteHeader(http.StatusNotFound)
			responses(rw, http.StatusNotFound, errorResponse{
				Error: messageObject{
					Message: "error while fetching users",
				},
			})
			return
		}

		// check if profile image is present or not
		if len(images) > 0 {
			image, err := images[0].Open()

			defer image.Close()
			if err != nil {
				logger.WithField("err", err.Error()).Error("Error while decoding image Data")
				responses(rw, http.StatusBadRequest, errorResponse{
					Error: messageObject{
						Message: "Invalid Image !",
					},
				})
				return
			}

			extension := filepath.Ext(images[0].Filename)

			if len(extension) < 2 || len(extension) > 5 {
				err = fmt.Errorf("Couldn't get extension of file!")
				logger.WithField("err", err.Error()).Error("Error while getting image Extension.")
				responses(rw, http.StatusBadRequest, errorResponse{
					Error: messageObject{
						Message: "Re-check the image file extension!",
					},
				})
				return
			}

			fileName := strings.ReplaceAll(user.FirstName, " ", "")
			tempFile, err := ioutil.TempFile("assets/users", fileName+"-*"+string(extension))
			if err != nil {
				logger.WithField("err", err.Error()).Error("Error while Creating a Temporary File")
				responses(rw, http.StatusInternalServerError, errorResponse{
					Error: messageObject{
						Message: "Couldn't  create temporary storage!",
					},
				})
				return
			}
			defer tempFile.Close()

			imageBytes, err := ioutil.ReadAll(image)
			if err != nil {
				logger.WithField("err", err.Error()).Error("Error while reading image File")
				responses(rw, http.StatusInternalServerError, errorResponse{
					Error: messageObject{
						Message: "Couldn't read the image file!",
					},
				})
				return
			}
			tempFile.Write(imageBytes)

			dbUser.ProfileImage = tempFile.Name()

		}

		if user.Email != "" {

			rw.WriteHeader(http.StatusBadRequest)
			logger.WithField("err", "cannot update email")
			responses(rw, http.StatusBadRequest, errorResponse{
				Error: messageObject{
					Message: "cannot update email id !!",
				},
			})
			return
		}

		err = dbUser.ValidatePatchParams(user)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			responses(rw, http.StatusInternalServerError, errorResponse{
				Error: messageObject{
					Message: "internal server error",
				},
			})
			logger.WithField("err", err.Error())
			return
		}

		err = deps.Store.UpdateUserByID(req.Context(), dbUser, int(userID))
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			responses(rw, http.StatusInternalServerError, errorResponse{
				Error: messageObject{
					Message: "internal server error",
				},
			})
			logger.WithField("err", err.Error()).Error("error while updating user's profile")
			return
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

		//TODO get the IsAdmin field of the user from the token
		isAdmin := true
		if !isAdmin {
			responses(rw, http.StatusForbidden, errorResponse{
				Error: messageObject{
					Message: "User Forbidden From Deleting Data",
				},
			})
			return
		}

		params := mux.Vars(req)
		userID, err := strconv.Atoi(params["id"])
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error in getting id")
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		//check if the provided id is not of another Admin
		user, err := deps.Store.GetUser(req.Context(), userID)
		if err != nil {
			logger.WithField("err", err.Error()).Error("error while fetching User")
			rw.WriteHeader(http.StatusNotFound)
			responses(rw, http.StatusNotFound, errorResponse{
				Error: messageObject{
					Message: "id Not Found",
				},
			})
			return
		}

		//check if the users is admin
		if user.IsAdmin == true {
			logger.WithField("err", "cannot delete an admin")
			rw.WriteHeader(http.StatusForbidden)
			responses(rw, http.StatusForbidden, errorResponse{
				Error: messageObject{
					Message: "access denied for delete",
				},
			})
			return
		}

		err = deps.Store.DeleteUserByID(req.Context(), userID)
		if err != nil {
			logger.WithField("err", err.Error()).Error("error in deleting user data")
			rw.WriteHeader(http.StatusBadRequest)
			responses(rw, http.StatusBadRequest, errorResponse{
				Error: messageObject{
					Message: "error in deleting user data",
				},
			})
			return
		}

		responses(rw, http.StatusOK, successResponse{
			Data: "record deleted Successfully",
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

		//TODO get the IsAdmin field of the user from the token
		isAdmin := true
		if !isAdmin {
			responses(rw, http.StatusForbidden, errorResponse{
				Error: messageObject{
					Message: "User Forbidden From disabling Data",
				},
			})
			return
		}

		params := mux.Vars(req)
		userID, err := strconv.Atoi(params["id"])
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error in getting id")
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		//check if the provided id is not of another Admin
		user, err := deps.Store.GetUser(req.Context(), userID)
		if err != nil {
			logger.WithField("err", err.Error()).Error("error while fetching User")
			rw.WriteHeader(http.StatusNotFound)
			responses(rw, http.StatusNotFound, errorResponse{
				Error: messageObject{
					Message: "id Not Found",
				},
			})
			return
		}

		//check if the user is admin
		if user.IsAdmin == true {
			logger.WithField("err", "cannot disable an admin")
			rw.WriteHeader(http.StatusForbidden)
			responses(rw, http.StatusForbidden, errorResponse{
				Error: messageObject{
					Message: "access denied for disable",
				},
			})
			return
		}

		err = deps.Store.DisableUserByID(req.Context(), userID)
		if err != nil {
			logger.WithField("err", err.Error()).Error("error in disabling user data")
			rw.WriteHeader(http.StatusBadRequest)
			responses(rw, http.StatusBadRequest, errorResponse{
				Error: messageObject{
					Message: "error in disabling user data",
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
			logger.WithField("err", err.Error()).Error("Error in getting id")
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		//check if the provided id is not of another Admin
		user, err := deps.Store.GetUser(req.Context(), userID)
		if err != nil {
			logger.WithField("err", err.Error()).Error("error while fetching User")
			rw.WriteHeader(http.StatusNotFound)
			responses(rw, http.StatusNotFound, errorResponse{
				Error: messageObject{
					Message: "id Not Found",
				},
			})
			return
		}

		//check if the user is admin
		if user.IsAdmin == true {
			logger.WithField("err", "cannot disable an admin")
			rw.WriteHeader(http.StatusForbidden)
			responses(rw, http.StatusForbidden, errorResponse{
				Error: messageObject{
					Message: "access denied for ensable",
				},
			})
			return
		}

		err = deps.Store.EnableUserByID(req.Context(), userID)
		if err != nil {
			logger.WithField("err", err.Error()).Error("error in enabling user data")
			rw.WriteHeader(http.StatusBadRequest)
			responses(rw, http.StatusBadRequest, errorResponse{
				Error: messageObject{
					Message: "error in enabling user data",
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
