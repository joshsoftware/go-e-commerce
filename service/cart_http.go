package service

import (
	"strconv"
	//"encoding/json"
	"net/http"
	logger "github.com/sirupsen/logrus"
)


func addToCartHandler(deps Dependencies) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		cartId, ok := req.URL.Query()["cartId"]
		productId, ok := req.URL.Query()["productId"]

		if !ok || len(cartId) < 1 || len(productId) < 0 {
			rw.WriteHeader(http.StatusInternalServerError)
			return 
		}
		cId,_ := strconv.Atoi(cartId[0])
		pId,_ := strconv.Atoi(productId[0]) 
		err := deps.Store.AddToCart(req.Context(), cId, pId)
		if err != nil{
			logger.WithField("err", err.Error()).Error("Error while adding to cart")
		}
	})
}

func removeFromCartHandler(deps Dependencies) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		cartId, ok := req.URL.Query()["cartId"]
		productId, ok := req.URL.Query()["productId"]

		if !ok || len(cartId) < 1 || len(productId) < 0 {
			rw.WriteHeader(http.StatusInternalServerError)
			return 
		}
		cId,_ := strconv.Atoi(cartId[0])
		pId,_ := strconv.Atoi(productId[0]) 
		err := deps.Store.RemoveFromCart(req.Context(), cId, pId)
		if err != nil{
			logger.WithField("err", err.Error()).Error("Error while adding to cart")
		}
	})
}

func updateIntoCartHandler(deps Dependencies) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		cartId, ok := req.URL.Query()["cartId"]
		productId, ok := req.URL.Query()["productId"]
		quantity, ok := req.URL.Query()["quantity"]
		if !ok || len(cartId) < 1 || len(productId) < 0 {
			rw.WriteHeader(http.StatusInternalServerError)
			return 
		}
		cId,_ := strconv.Atoi(cartId[0])
		pId,_ := strconv.Atoi(productId[0]) 
		qty,_ := strconv.Atoi(quantity[0])
		err := deps.Store.UpdateIntoCart(req.Context(), qty, cId, pId)
		if err != nil{
			logger.WithField("err", err.Error()).Error("Error while adding to cart")
		}
	})
}