package service

import (
	"fmt"
	"joshsoftware/go-e-commerce/db"
	"math"
	"net/http"
	"strconv"

	logger "github.com/sirupsen/logrus"
)

// @Title getProductByFilters
// @Description list all Products with specified filters
// @Router /products/filters [GET]
// @Params /products/filters?categoryid=id&price=asc&brand=name&size=name&color=name
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
		filter.CategoryId = req.URL.Query().Get("category_id")
		filter.Price = req.URL.Query().Get("price")
		filter.Brand = req.URL.Query().Get("brand")
		filter.Size = req.URL.Query().Get("size")
		filter.Color = req.URL.Query().Get("color")

		pageStr := req.URL.Query().Get("page")
		limitStr := req.URL.Query().Get("limit")

		// Setting default limit as 5
		if limitStr == "" {
			limitStr = "5"
		}

		// Setting default page as 1
		if pageStr == "" {
			pageStr = "1"
		}

		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error while converting limitStr to int")
			Message := "Limits value invalid"
			responseMsg(rw, http.StatusBadRequest, Message)
			return
		}

		page, err := strconv.Atoi(pageStr)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error while converting pageStr to int")
			Message := "Page value invalid"
			responseMsg(rw, http.StatusBadRequest, Message)
			return
		}

		// Avoid divide by zero exception and -ve values for page and limit
		if (limit > 0 && page > 0) == false {
			err = fmt.Errorf("limit or page are non-positive")
			logger.WithField("err", err.Error()).Error("Error limit or page were invalid values")
			Message := "limits or page value are non-positive"
			responseMsg(rw, http.StatusBadRequest, Message)
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

		offset := (page - 1) * limit
		offsetStr := strconv.Itoa(offset)
		totalRecords, products, err := deps.Store.FilteredProducts(req.Context(), filter, limitStr, offsetStr)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error getting filtered records or Page not Found")
			Message := "Error getting filtered records or Page not Found"
			responseMsg(rw, http.StatusBadRequest, Message)
			return
		}

		var pagination db.Pagination
		pagination.TotalPages = int(math.Ceil(float64(totalRecords) / float64(limit)))
		pagination.Products = products
		response(rw, http.StatusOK, pagination)
		return
	})

}

// @Title getProductBySearch
// @Description list all Products with specified filters
// @Router /products/search [GET]
// @Params /products/search?text=apple+that+can+be+eaten
//  checking will take place in product name , brand and category name
//  brand, size, color will be also be checked case-insensitively string
// @Accept	json
// @Success 200 {object}
// @Failure 404 {object}

// TODO Optimize the queries

func getProductBySearchHandler(deps Dependencies) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		pageStr := req.URL.Query().Get("page")
		limitStr := req.URL.Query().Get("limit")
		text := req.URL.Query().Get("text")

		// Setting default limit as 5
		if limitStr == "" {
			limitStr = "5"
		}

		// Setting default page as 1
		if pageStr == "" {
			pageStr = "1"
		}

		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error while converting limit to int")
			Message := "limits value invalid"
			responseMsg(rw, http.StatusBadRequest, Message)
			return
		}

		page, err := strconv.Atoi(pageStr)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error while converting page to int")
			Message := "page value invalid"
			responseMsg(rw, http.StatusBadRequest, Message)
			return
		}

		// Avoid divide by zero exception and -ve values for page and limit
		if limit <= 0 || page <= 0 {
			err = fmt.Errorf("limit or page are non-positive")
			logger.WithField("err", err.Error()).Error("Error limit or page were invalid values")
			Message := "limits or page value are non-positive."
			responseMsg(rw, http.StatusBadRequest, Message)
			return
		}

		var totalRecords int
		var products []db.Product
		offset := (page - 1) * limit
		offsetStr := strconv.Itoa(offset)
		if text == "" {
			// Behave same as List All Products and return
			totalRecords, products, err = deps.Store.ListProducts(req.Context(), limit, offset)

			if err != nil {
				logger.WithField("err", err.Error()).Error("Error Couldn't find any Product records or Page out of range")
				Message := "Couldn't find any Products records or Page out of range"
				responseMsg(rw, http.StatusInternalServerError, Message)
				return
			}
			goto Skip
		}

		totalRecords, products, err = deps.Store.SearchProductsByText(req.Context(), text, limitStr, offsetStr)
		if err != nil || totalRecords == 0 {
			logger.WithField("err", err.Error()).Error("Error Couldn't find any matching search records or Page out of range")
			Message := "Couldn't find any matching search records or Page out of range"
			responseMsg(rw, http.StatusBadRequest, Message)
			return
		}

	Skip:
		var pagination db.Pagination
		pagination.TotalPages = int(math.Ceil(float64(totalRecords) / float64(limit)))
		pagination.Products = products
		response(rw, http.StatusOK, pagination)
		return
	})

}
