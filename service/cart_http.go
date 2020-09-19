package service

import (
	"encoding/json"
	"net/http"
	"strings"

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

		authToken := req.Header.Get("Authorization")
		if strings.HasPrefix(strings.ToUpper(authToken), "BEARER") {
			authToken = authToken[len("BEARER "):]
		}

		userID, _, err := getDataFromToken(authToken)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error fetching data")
			responses(rw, http.StatusUnauthorized, errorResponse{
				Error: messageObject{
					Message: "Invalid Credentials",
				},
			})
			return
		}

		cart_products, err := deps.Store.GetCart(req.Context(), int(userID))
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error fetching data")
			responses(rw, http.StatusInternalServerError, errorResponse{
				Error: messageObject{
					Message: "Failure in getting data from database",
				},
			})
			return
		}

		respBytes, err := json.Marshal(cart_products)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error marshaling cart data")
			responses(rw, http.StatusInternalServerError, errorResponse{
				Error: messageObject{
					Message: "Failure in marshalling data",
				},
			})
			return
		}

		rw.Header().Add("Content-Type", "application/json")
		rw.Write(respBytes)
	})
}
