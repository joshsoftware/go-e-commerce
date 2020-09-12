package service

import (
	"encoding/json"
	"joshsoftware/go-e-commerce/db"
	"math"
	"net/http"
	"strconv"

	logger "github.com/sirupsen/logrus"
)

// @Title getProductByFilters
// @Description list all Products with specified filters
// @Router /products/filters [GET]
// @Params /products/filters?categoryid=id?price=asc?brand=name?size=name?color=name
//  price can be asc or desc, it will stored as a string
//  categoryid will be an integer value, but for convinience it will be stored as string
//  brand, size, color will be case-sensitive string
// @Accept	json
// @Success 200 {object}
// @Failure 404 {object}
// @Features This API can replace ListProducts API, but time Complexity will be a bit high
func getProductByFiltersHandler(deps Dependencies) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		var filter db.Filter

		// filter.Price will either be "asc" or "desc"
		// All these object's are used as string itself
		// Using String really made some things easy in dynamic query writing
		filter.CategoryId = req.URL.Query().Get("categoryid")
		filter.Price = req.URL.Query().Get("price")
		filter.Brand = req.URL.Query().Get("brand")
		filter.Size = req.URL.Query().Get("size")
		filter.Color = req.URL.Query().Get("color")

		page := req.URL.Query().Get("page")
		limit := req.URL.Query().Get("limit")

		// Setting default limit as 5
		if limit == "" {
			limit = "5"
		}

		// Setting default page as 1
		if page == "" {
			page = "1"
		}

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

		// Checking for flags, true means we need to filter by that field
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

		// TODO Handle that last page getting as null
		if count < (ls * (int(ps) - 1)) {
			rw.WriteHeader(http.StatusBadRequest)
			response(rw, http.StatusBadRequest, errorResponse{
				Error: messageObject{
					Message: "Page Not Found..!",
				},
			})
			return
		}

		var pagination db.Pagination
		pagination.TotalPages = int(math.Ceil(float64(count) / float64(ls)))

		products, err := deps.Store.FilteredRecords(req.Context(), filter, limit, page)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error fetching data filtered products")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		pagination.Products = products

		respBytes, err := json.Marshal(pagination)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error mashaling pagination and product data")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		rw.Header().Add("Content-Type", "application/json")
		rw.Write(respBytes)
	})

}
