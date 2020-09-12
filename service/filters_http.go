package service

import (
	"encoding/json"
	"joshsoftware/go-e-commerce/db"
	"net/http"
	"strconv"

	logger "github.com/sirupsen/logrus"
)

// @Title getProductByFilters
// @Description list all Products with specified filters
// @Router /products/filters [GET]
// @Params /products/filters?categoryid=id?price=asc?brand=name?size=name?color=name
// @Accept	json
// @Success 200 {object}
// @Failure 400 {object}
func getProductByFiltersHandler(deps Dependencies) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		var filter db.Filter

		// filter.Price will either be "asc" or "desc"
		filter.CategoryId = req.URL.Query().Get("categoryid")
		filter.Price = req.URL.Query().Get("price")
		filter.Brand = req.URL.Query().Get("brand")
		filter.Size = req.URL.Query().Get("size")
		filter.Color = req.URL.Query().Get("color")

		page := req.URL.Query().Get("page")
		limit := req.URL.Query().Get("limit")

		if limit == "" {
			limit = "5"
		}

		if page == "" {
			page = "1"
		}

		// Handle errors
		ls, err := strconv.Atoi(limit)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error while converting limit to int")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		ps, err := strconv.Atoi(page)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error while converting page to int")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		if filter.CategoryId != "" {
			filter.CategoryFlag = true
		}

		if filter.Price != "" {
			filter.PriceFlag = true
		}

		if filter.Brand != "" {
			filter.BrandFlag = true
		}

		if filter.Size != "" {
			filter.SizeFlag = true
		}

		if filter.Color != "" {
			filter.ColorFlag = true
		}

		count, err := deps.Store.FilteredRecordsCount(req.Context(), filter)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error getting count of filtered records")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		if count <= (ls * (int(ps) - 1)) {
			rw.WriteHeader(http.StatusOK)
			response(rw, http.StatusOK, successResponse{
				Data: messageObject{
					Message: "Page Not Found..!",
				},
			})
			return
		}

		products, err := deps.Store.FilteredRecords(req.Context(), filter, limit, page)
		//products, err := deps.Store.ListProducts(req.Context(), limit, page)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error fetching data")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		respBytes, err := json.Marshal(products)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error mashaling products data")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		rw.Header().Add("Content-Type", "application/json")
		rw.Write(respBytes)
	})

}
