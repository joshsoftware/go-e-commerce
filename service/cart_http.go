package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	logger "github.com/sirupsen/logrus"
)

// @Title listCart
// @Description list all Product inside cart
// @Router /user/id/cart [get]
// @Accept  json
// @Success 200 {object}
// @Failure 400 {object}
func getCartHandler(deps Dependencies) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// request_params := mux.Vars(req)
		// user_id, err := strconv.Atoi(request_params["user_id"])

		authToken := req.Header["Token"]
		fmt.Println("auth Token : ", authToken[0])
		payload, err := getDataFromToken(authToken[0])
		fmt.Println("User id :", payload.UserID)
		if err != nil {
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte("Unauthorized"))
			return
		}

		cart_products, err := deps.Store.GetCart(req.Context(), int(payload.UserID))
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error fetching data")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		respBytes, err := json.Marshal(cart_products)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error marshaling cart data")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		rw.Header().Add("Content-Type", "application/json")
		rw.Write(respBytes)
	})
}

func addToCartHandler(deps Dependencies) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		authToken := req.Header["Token"]
		payload, err := getDataFromToken(authToken[0])
		if err != nil {
			logger.WithField("err", err.Error()).Error("Unauthorized user")
			error := errorResponse{
				Error: "Unauthorized user",
			}
			responses(rw, http.StatusUnauthorized, error)
			return
		}

		productID, err := strconv.Atoi(req.URL.Query()["productID"][0])
		if err != nil {
			logger.WithField("err", err.Error()).Error("product_id is missing")
			error := errorResponse{
				Error: "product_id missing",
			}
			responses(rw, http.StatusBadRequest, error)
			return
		}

		rowsAffected, err := deps.Store.AddToCart(req.Context(), int(payload.UserID), productID)
		if err != nil {
			logger.WithField("err", err.Error()).Error("error while adding to cart")
			error := errorResponse{
				Error: "could not add item",
			}
			responses(rw, http.StatusInternalServerError, error)
			return
		}

		if rowsAffected != 1 {
			success := successResponse{
				Data: "zero rows affected",
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
		authToken := req.Header["Token"]
		payload, err := getDataFromToken(authToken[0])
		if err != nil {
			logger.WithField("err", err.Error()).Error("Unauthorized user")
			error := errorResponse{
				Error: "Unauthorized user",
			}
			responses(rw, http.StatusUnauthorized, error)
			return
		}

		productID, err := strconv.Atoi(req.URL.Query()["productID"][0])
		if err != nil {
			logger.WithField("err", err.Error()).Error("product_id is missing")
			error := errorResponse{
				Error: "product_id missing",
			}
			responses(rw, http.StatusBadRequest, error)
			return
		}

		rowsAffected, err := deps.Store.DeleteFromCart(req.Context(), int(payload.UserID), productID)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error while removing from cart")
			error := errorResponse{
				Error: "could not remove item",
			}
			responses(rw, http.StatusInternalServerError, error)
			return
		}

		if rowsAffected != 1 {
			success := successResponse{
				Data: "zero rows affected",
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
		authToken := req.Header["Token"]
		payload, err := getDataFromToken(authToken[0])
		if err != nil {
			logger.WithField("err", err.Error()).Error("Unauthorized user")
			error := errorResponse{
				Error: "Unauthorized user",
			}
			responses(rw, http.StatusUnauthorized, error)
			return
		}

		productID, err := strconv.Atoi(req.URL.Query()["productID"][0])
		if err != nil {
			logger.WithField("err", err.Error()).Error("product_id is missing")
			error := errorResponse{
				Error: "product_id missing",
			}
			responses(rw, http.StatusBadRequest, error)
			return
		}

		quantity, err := strconv.Atoi(req.URL.Query()["quantity"][0])
		if err != nil {
			logger.WithField("err", err.Error()).Error("quantity is missing")
			error := errorResponse{
				Error: "quantity missing",
			}
			responses(rw, http.StatusBadRequest, error)
			return
		}

		rowsAffected, err := deps.Store.UpdateIntoCart(req.Context(), quantity, int(payload.UserID), productID)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error while updating to cart")
			error := errorResponse{
				Error: "could not update quantity",
			}
			responses(rw, http.StatusInternalServerError, error)
			return
		}

		if rowsAffected != 1 {
			success := successResponse{
				Data: "zero rows affected",
			}
			responses(rw, http.StatusOK, success)
			return
		}

		success := successResponse{
			Data: "Quantity updated successfully",
		}
		responses(rw, http.StatusOK, success)
	})
}
