package service

import (
	"encoding/json"
	"fmt"
	"joshsoftware/go-e-commerce/db"
	"math"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	logger "github.com/sirupsen/logrus"
)

// @Title listProducts
// @Description list all Products
// @Router /products [GET]
// @Accept	json
// @Success 200 {object}
// @Failure 400 {object}
func listProductsHandler(deps Dependencies) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		limitStr := req.URL.Query().Get("limit")
		pageStr := req.URL.Query().Get("page")

		if limitStr == "" {
			limitStr = "5"
		}

		if pageStr == "" {
			pageStr = "1"
		}

		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error while converting limitStr to int")
			response(rw, http.StatusInternalServerError, errorResponse{
				Error: messageObject{
					Message: "Error while converting limitStr to int",
				},
			})
			return
		}

		page, err := strconv.Atoi(pageStr)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error while converting pageStr to int")
			response(rw, http.StatusInternalServerError, errorResponse{
				Error: messageObject{
					Message: "Error while converting pageStr to int",
				},
			})
			return
		}

		// Avoid divide by zero exception and -ve values for page and limit
		if limit <= 0 || page <= 0 {
			err = fmt.Errorf("limit or page are non-positive")
			logger.WithField("err", err.Error()).Error("Error limit or page contained invalid value")
			response(rw, http.StatusBadRequest, errorResponse{
				Error: messageObject{
					Message: "limits or page value invalid",
				},
			})
			return
		}
		offset := (page - 1) * limit
		totalRecords, products, err := deps.Store.ListProducts(req.Context(), limit, offset)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error Couldn't find any Product records or Page out of range")
			response(rw, http.StatusInternalServerError, errorResponse{
				Error: messageObject{
					Message: "Couldn't find any Products records or Page out of range",
				},
			})
			return
		}

		var pagination db.Pagination
		pagination.TotalPages = int(math.Ceil(float64(totalRecords) / float64(limit)))

		pagination.Products = products

		response(rw, http.StatusOK, pagination)
		return
	})
}

// @ Title getProductById
// @ Description get single product by its id
// @ Router /product/product_id [get]
// @ Accept json
// @ Success 200 {object}
// @ Failure 400 {object}

func getProductByIdHandler(deps Dependencies) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		vars := mux.Vars(req)
		id, err := strconv.Atoi(vars["product_id"])
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error id key is missing")
			response(rw, http.StatusBadRequest, errorResponse{
				Error: messageObject{
					Message: "Error product_id is invalid",
				},
			})
			return
		}

		product, err := deps.Store.GetProductByID(req.Context(), id)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error fetching Product data, no Product found")
			response(rw, http.StatusBadRequest, errorResponse{
				Error: messageObject{
					Message: "Error feching data, no row found.",
				},
			})
			return
		}

		response(rw, http.StatusOK, product)
		return
	})
}

// @Title createProduct
// @Description create a Product, insert into DB
// @Router /createProduct [POST]
// @Accept	json
// @Success 200 {object}
// @Failure 400 {object}
func createProductHandler(deps Dependencies) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		var product db.Product
		err := json.NewDecoder(req.Body).Decode(&product)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error while decoding product")
			response(rw, http.StatusBadRequest, errorResponse{
				Error: messageObject{
					Message: "Invalid json body",
				},
			})
			return
		}

		errRes, valid := product.Validate()
		if !valid {
			response(rw, http.StatusBadRequest, errRes)
			return
		}

		createdProductID, err := deps.Store.CreateProduct(req.Context(), product)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error while inserting product")
			response(rw, http.StatusBadRequest, errorResponse{
				Error: messageObject{
					Message: "Error inserting the product, product already exists",
				},
			})
			return
		}
		response(rw, http.StatusOK, successResponse{Data: createdProductID})
		return
	})
}

// @ Title deleteProductById
// @ Description delete product by its id
// @ Router /product/product_id [delete]
// @ Accept json
// @ Success 200 {object}
// @ Failure 400 {object}

func deleteProductByIdHandler(deps Dependencies) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		vars := mux.Vars(req)
		id, err := strconv.Atoi(vars["product_id"])
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error id key is missing")
			response(rw, http.StatusBadRequest, errorResponse{
				Error: messageObject{
					Message: "Error id is missing/invalid",
				},
			})
			return
		}

		err = deps.Store.DeleteProductById(req.Context(), id)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error fetching data no row found")
			response(rw, http.StatusInternalServerError, errorResponse{
				Error: messageObject{
					Message: "Internal server error  (Error feching data)",
				},
			})
			return
		}

		response(rw, http.StatusOK, successResponse{
			Data: "Product Deleted Successfully!",
		})
		return
	})
}

// @ Title updateProductStockById
// @ Description update product by its id
// @ Router /product/product_id [put]
// @ Accept json
// @ Success 200 {object}
// @ Failure 400 {object}
func updateProductStockByIdHandler(deps Dependencies) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		Id := req.URL.Query().Get("product_id")
		Count := req.URL.Query().Get("stock")
		var err error

		// Handle errors
		productId, err := strconv.Atoi(Id)

		if Id == "" || err != nil {
			logger.WithField("err", err.Error()).Error("Error product_id parameter is missing or corrupt")
			response(rw, http.StatusBadRequest, errorResponse{
				Error: messageObject{
					Message: "Error id is missing/invalid",
				},
			})
			return
		}

		count, err := strconv.Atoi(Count)

		if Count == "" || err != nil {
			logger.WithField("err", err.Error()).Error("Error stock parameter is missing or corrupt")
			response(rw, http.StatusBadRequest, errorResponse{
				Error: messageObject{
					Message: "Error id is missing/invalid",
				},
			})
			return
		}

		var product db.Product
		product, err = deps.Store.GetProductByID(req.Context(), productId)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error while fetching product with stated id ")
			response(rw, http.StatusBadRequest, errorResponse{
				Error: messageObject{
					Message: "Error Product is wasn't found in database.",
				},
			})
			return
		}

		// Decrement available Quantity
		product.Quantity -= count

		// Validate if Quantity is less than 0
		errRes, valid := product.Validate()
		if !valid {
			response(rw, http.StatusBadRequest, errRes)
			return
		}

		var updatedProduct db.Product
		updatedProduct, err = deps.Store.UpdateProductStockById(req.Context(), product, productId)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error while updating product attribute")
			response(rw, http.StatusInternalServerError, errorResponse{
				Error: messageObject{
					Message: "Internal server error",
				},
			})
			return
		}

		response(rw, http.StatusOK, successResponse{Data: updatedProduct})
		return
	})
}

// @ Title updateProductById
// @ Description update product by its id
// @ Router /product/product_id [put]
// @ Accept json
// @ Success 200 {object}
// @ Failure 400 {object}

func updateProductByIdHandler(deps Dependencies) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		vars := mux.Vars(req)
		id, err := strconv.Atoi(vars["product_id"])
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error id key is missing")
			rw.WriteHeader(http.StatusBadRequest)
			response(rw, http.StatusBadRequest, errorResponse{
				Error: messageObject{
					Message: "Error id is missing/invalid",
				},
			})
			return
		}

		var product db.Product
		err = json.NewDecoder(req.Body).Decode(&product)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			logger.WithField("err", err.Error()).Error("Error while decoding user")
			response(rw, http.StatusBadRequest, errorResponse{
				Error: messageObject{
					Message: "Internal server error",
				},
			})
			return
		}

		var updatedProduct db.Product
		updatedProduct, err = deps.Store.UpdateProductById(req.Context(), product, id)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error while updating product attribute")
			response(rw, http.StatusInternalServerError, errorResponse{
				Error: messageObject{
					Message: "Internal server error",
				},
			})
			return
		}

		response(rw, http.StatusOK, successResponse{Data: updatedProduct})

		return
	})
}

// TODO test case to delete category
