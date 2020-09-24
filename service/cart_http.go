package service

import (
	"strings"
	"strconv"
	// "encoding/json"
	"net/http"
	logger "github.com/sirupsen/logrus"
)

// type successResponse struct {
// 	Message string `json: "message"`
// }

// type errorResponse struct {
// 	Error string `json: "error"`
// }

// func response(rw http.ResponseWriter, status int, responseData interface{}){
// 	respBody, err := json.Marshal(responseData)
// 	if err != nil {
// 		logger.WithField("err", err.Error()).Error("error while marshling")
// 		rw.WriteHeader(http.StatusInternalServerError)
// 		return 
// 	}
// 	rw.Header().Add("Content-Type","application/json")
// 	rw.WriteHeader(status)
// 	rw.Write(respBody)
// }

func addToCartHandler(deps Dependencies) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		authToken := req.Header.Get("Authorization")
    if strings.HasPrefix(strings.ToUpper(authToken), "BEARER") {
        authToken = authToken[len("BEARER "):]
    }

		cartID, _, err := getDataFromToken(authToken)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Unauthorized user")
			error := errorResponse {
				Error : "Unauthorized user",
			}
			responses(rw, http.StatusUnauthorized, error)
			return
		}

		productID, err := strconv.Atoi(req.URL.Query()["productID"][0])
		if err != nil {
			logger.WithField("err", err.Error()).Error("product_id is missing")
			error := errorResponse {
				Error : "product_id missing",
			}
			responses(rw, http.StatusBadRequest, error)
			return
		}

		rowsAffected, err := deps.Store.AddToCart(req.Context(), int(cartID), productID)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error while adding to cart")
			error := errorResponse {
				Error : "could not add item",
			}
			responses(rw, http.StatusInternalServerError, error)
			return
		}

		if rowsAffected != 1 {
			success := successResponse {
				Data : "zero rows affected",
			}
			responses(rw, http.StatusOK, success)
			return
		}		

		success := successResponse{
			Data: "Item added successfully",
		}
		responses(rw, http.StatusOK, success)
	})
}

func deleteFromCartHandler(deps Dependencies) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		authToken := req.Header.Get("Authorization")
    if strings.HasPrefix(strings.ToUpper(authToken), "BEARER") {
        authToken = authToken[len("BEARER "):]
    }

    cartID, _, err := getDataFromToken(authToken)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Unauthorized user")
			error := errorResponse {
				Error : "Unauthorized user",
			}
			responses(rw, http.StatusUnauthorized, error)
			return
		}

		productID, err := strconv.Atoi(req.URL.Query()["productID"][0])
		if err != nil {
			logger.WithField("err", err.Error()).Error("product_id is missing")
			error := errorResponse {
				Error : "product_id missing",
			}
			responses(rw, http.StatusBadRequest, error)
			return
		}

		rowsAffected, err := deps.Store.DeleteFromCart(req.Context(), int(cartID), productID)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error while removing from cart")
			error := errorResponse {
				Error : "could not remove item",
			}
			responses(rw, http.StatusInternalServerError, error)
			return 
		}

		if rowsAffected != 1 {
			success := successResponse {
				Data : "zero rows affected",
			}
			responses(rw, http.StatusOK, success)
			return
		}

		success := successResponse{
			Data: "Item removed successfully",
		}
		responses(rw, http.StatusOK, success)
	})
}

func updateIntoCartHandler(deps Dependencies) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		authToken := req.Header.Get("Authorization")
    if strings.HasPrefix(strings.ToUpper(authToken), "BEARER") {
        authToken = authToken[len("BEARER "):]
    }
		
		cartID, _, err := getDataFromToken(authToken)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Unauthorized user")
			error := errorResponse {
				Error : "Unauthorized user",
			}
			responses(rw, http.StatusUnauthorized, error)
			return
		}
		
		productID, err := strconv.Atoi(req.URL.Query()["productID"][0])
		if err != nil {
			logger.WithField("err", err.Error()).Error("product_id is missing")
			error := errorResponse {
				Error : "product_id missing",
			}
			responses(rw, http.StatusBadRequest, error)
			return
		}

		quantity, err := strconv.Atoi(req.URL.Query()["quantity"][0])
		if err != nil {
			logger.WithField("err", err.Error()).Error("quantity is missing")
			error := errorResponse {
				Error : "quantity missing",
			}
			responses(rw, http.StatusBadRequest, error)
			return 
		}
		
		rowsAffected, err := deps.Store.UpdateIntoCart(req.Context(), quantity, int(cartID), productID)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error while updating to cart")
			error := errorResponse {
				Error : "could not update quantity",
			}
			responses(rw, http.StatusInternalServerError, error)
			return
		}

		if rowsAffected != 1 {
			success := successResponse {
				Data : "zero rows affected",
			}
			responses(rw, http.StatusOK, success)
			return
		}

		success := successResponse{
			Data : "Quantity updated successfully", 
		}
		responses(rw, http.StatusOK, success)
	})
}