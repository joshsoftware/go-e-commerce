package service

import (
	"strconv"
	"encoding/json"
	"net/http"
	logger "github.com/sirupsen/logrus"
)

type successResponse struct {
	Message string `json: "message"`
}

type errorResponse struct {
	Error string `json: "error"`
}

func response(rw http.ResponseWriter, status int, responseData interface{}){
	respBody, err := json.Marshal(responseData)
	if err != nil {
		logger.WithField("err", err.Error()).Error("error while marshling")
		rw.WriteHeader(http.StatusInternalServerError)
		return 
	}
	rw.Header().Add("Content-Type","application/json")
	rw.WriteHeader(status)
	rw.Write(respBody)
}

func addToCartHandler(deps Dependencies) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		authToken := req.Header["Token"]
		cartID, _, err := getDataFromToken(authToken[0])
		if err != nil {
			logger.WithField("err", err.Error()).Error("Unauthorized user")
			error := errorResponse {
				Error : "Unauthorized user",
			}
			response(rw, http.StatusUnauthorized, error)
			return
		}

		productID, err := strconv.Atoi(req.URL.Query()["productID"][0])
		if err != nil {
			logger.WithField("err", err.Error()).Error("product_id is missing")
			error := errorResponse {
				Error : "product_id missing",
			}
			response(rw, http.StatusBadRequest, error)
			return
		}

		rowsAffected, err := deps.Store.AddToCart(req.Context(), int(cartID), productID)
		if err != nil {
			logger.WithField("err", err.Error()).Error("error while adding to cart")
			error := errorResponse {
				Error : "could not add item",
			}
			response(rw, http.StatusInternalServerError, error)
			return
		}

		if rowsAffected != 1 {
			success := successResponse {
				Message : "zero rows affected",
			}
			response(rw, http.StatusOK, success)
			return
		}		

		success := successResponse{
			Message: "Item added successfully",
		}
		response(rw, http.StatusOK, success)
	})
}

func deleteFromCartHandler(deps Dependencies) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		authToken := req.Header["Token"]
		cartID, _, err := getDataFromToken(authToken[0])
		if err != nil {
			logger.WithField("err", err.Error()).Error("Unauthorized user")
			error := errorResponse {
				Error : "Unauthorized user",
			}
			response(rw, http.StatusUnauthorized, error)
			return
		}

		productID, err := strconv.Atoi(req.URL.Query()["productID"][0])
		if err != nil {
			logger.WithField("err", err.Error()).Error("product_id is missing")
			error := errorResponse {
				Error : "product_id missing",
			}
			response(rw, http.StatusBadRequest, error)
			return
		}

		rowsAffected, err := deps.Store.DeleteFromCart(req.Context(), int(cartID), productID)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error while removing from cart")
			error := errorResponse {
				Error : "could not remove item",
			}
			response(rw, http.StatusInternalServerError, error)
			return 
		}

		if rowsAffected != 1 {
			success := successResponse {
				Message : "zero rows affected",
			}
			response(rw, http.StatusOK, success)
			return
		}

		success := successResponse{
			Message: "Item removed successfully",
		}
		response(rw, http.StatusOK, success)
	})
}

func updateIntoCartHandler(deps Dependencies) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		authToken := req.Header["Token"]
		cartID, _, err := getDataFromToken(authToken[0])
		if err != nil {
			logger.WithField("err", err.Error()).Error("Unauthorized user")
			error := errorResponse {
				Error : "Unauthorized user",
			}
			response(rw, http.StatusUnauthorized, error)
			return
		}
		
		productID, err := strconv.Atoi(req.URL.Query()["productID"][0])
		if err != nil {
			logger.WithField("err", err.Error()).Error("product_id is missing")
			error := errorResponse {
				Error : "product_id missing",
			}
			response(rw, http.StatusBadRequest, error)
			return
		}

		quantity, err := strconv.Atoi(req.URL.Query()["quantity"][0])
		if err != nil {
			logger.WithField("err", err.Error()).Error("quantity is missing")
			error := errorResponse {
				Error : "quantity missing",
			}
			response(rw, http.StatusBadRequest, error)
			return 
		}
		
		rowsAffected, err := deps.Store.UpdateIntoCart(req.Context(), quantity, int(cartID), productID)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error while updating to cart")
			error := errorResponse {
				Error : "could not update quantity",
			}
			response(rw, http.StatusInternalServerError, error)
			return
		}

		if rowsAffected != 1 {
			success := successResponse {
				Message : "zero rows affected",
			}
			response(rw, http.StatusOK, success)
			return
		}

		success := successResponse{
			Message : "Quantity updated successfully", 
		}
		response(rw, http.StatusOK, success)
	})
}