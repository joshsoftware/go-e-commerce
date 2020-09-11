package service

import (
	"encoding/json"
	"fmt"
	"net/http"

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
		userID, _, err := getDataFromToken(authToken[0])
		fmt.Println("User id :", userID)
		if err != nil {
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte("Unauthorized"))
			return
		}

		cart_products, err := deps.Store.GetCart(req.Context(), int(userID))
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
