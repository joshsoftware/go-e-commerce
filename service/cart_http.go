package service

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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
		params := mux.Vars(req)
		user_id, err := strconv.Atoi(params["user_id"])
		cart, err := deps.Store.GetCart(req.Context(), user_id)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error fetching data")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		respBytes, err := json.Marshal(cart)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error marshaling cart data")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		rw.Header().Add("Content-Type", "application/json")
		rw.Write(respBytes)
	})
}
