package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"joshsoftware/go-e-commerce/db"
	"math"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
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
			logger.WithField("err", err.Error()).Error("Error while converting limit to int")
			rw.Header().Add("Content-Type", "application/json")
			rw.WriteHeader(http.StatusInternalServerError)
			responses(rw, http.StatusInternalServerError, errorResponse{
				Error: messageObject{
					Message: "Error while converting limitStr to int",
				},
			})
			return
		}

		page, err := strconv.Atoi(pageStr)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error while converting page to int")
			rw.Header().Add("Content-Type", "application/json")
			rw.WriteHeader(http.StatusInternalServerError)
			responses(rw, http.StatusInternalServerError, errorResponse{
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
			rw.Header().Add("Content-Type", "application/json")
			rw.WriteHeader(http.StatusBadRequest)
			responses(rw, http.StatusBadRequest, errorResponse{
				Error: messageObject{
					Message: "limits or page value invalid",
				},
			})
			return
		}

		count, products, err := deps.Store.ListProducts(req.Context(), limit, page)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error Couldn't find any Product records or Page out of range")
			rw.Header().Add("Content-Type", "application/json")
			rw.WriteHeader(http.StatusInternalServerError)
			responses(rw, http.StatusInternalServerError, errorResponse{
				Error: messageObject{
					Message: "Couldn't find any Products records or Page out of range",
				},
			})
			return
		}

		var pagination db.Pagination
		pagination.TotalPages = int(math.Ceil(float64(count) / float64(limit)))

		pagination.Products = products

		respBytes, err := json.Marshal(pagination)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error mashaling pagination data")
			rw.Header().Add("Content-Type", "application/json")
			rw.WriteHeader(http.StatusInternalServerError)
			responses(rw, http.StatusInternalServerError, errorResponse{
				Error: messageObject{
					Message: "Couldn't mashal pagination data",
				},
			})
			return
		}

		rw.Header().Add("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		rw.Write(respBytes)
	})
}

var decoder = schema.NewDecoder()

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
			rw.Header().Add("Content-Type", "application/json")
			rw.WriteHeader(http.StatusBadRequest)
			responses(rw, http.StatusBadRequest, errorResponse{
				Error: messageObject{
					Message: "Error product_id is invalid",
				},
			})
			return
		}

		product, err := deps.Store.GetProductByID(req.Context(), id)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error fetching data no product found")
			rw.Header().Add("Content-Type", "application/json")
			rw.WriteHeader(http.StatusInternalServerError)
			responses(rw, http.StatusBadRequest, errorResponse{
				Error: messageObject{
					Message: "Error feching data, no row found.",
				},
			})
			return
		}

		_, err = json.Marshal(product)
		if err != nil {
			rw.Header().Add("Content-Type", "application/json")
			rw.WriteHeader(http.StatusInternalServerError)
			logger.WithField("err", err.Error()).Error("Error marshaling products data")
			responses(rw, http.StatusInternalServerError, errorResponse{
				Error: messageObject{
					Message: "Internal server error",
				},
			})
			return
		}
		responses(rw, http.StatusOK, product)
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

		// Parse input, multipart/form-data
		err := req.ParseMultipartForm(15 << 20) // 15 MB Max File Size
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error while parsing the Product form")
			responses(rw, http.StatusBadRequest, errorResponse{
				Error: messageObject{
					Message: "Invalid Form Data!",
				},
			})
			return
		}

		// Retrive file from posted data
		formdata := req.MultipartForm

		// grab the filename
		contents := formdata.Value
		images := formdata.File["images"]
		//err = req.ParseForm()

		//grab product
		fmt.Println(contents)
		err = decoder.Decode(&product, contents)
		fmt.Println(product)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error while decoding product")
			responses(rw, http.StatusBadRequest, errorResponse{
				Error: messageObject{
					Message: "Invalid form contents",
				},
			})
			return
		}

		errRes, valid := product.Validate()
		if !valid {
			_, err := json.Marshal(errRes)
			if err != nil {
				logger.WithField("err", err.Error()).Error("Error marshalling Product's data")
				responses(rw, http.StatusBadRequest, errorResponse{
					Error: messageObject{
						Message: "Invalid json body",
					},
				})
				return
			}
			responses(rw, http.StatusBadRequest, errRes)
			return
		}

		for i, _ := range images {
			image, err := images[i].Open()
			defer image.Close()
			if err != nil {
				logger.WithField("err", err.Error()).Error("Error while decoding image Data")
				responses(rw, http.StatusBadRequest, errorResponse{
					Error: messageObject{
						Message: "Invalid Image !",
					},
				})
				return
			}

			extensionRegex := regexp.MustCompile(`[.]+.*`)
			extension := extensionRegex.Find([]byte(images[i].Filename))
			if len(extension) < 2 || len(extension) > 5 {
				err = fmt.Errorf("Couldn't get extension of file!")
				logger.WithField("err", err.Error()).Error("Error while getting image Extension.")
				responses(rw, http.StatusInternalServerError, errorResponse{
					Error: messageObject{
						Message: "Re-check the image file extension!",
					},
				})
				return
			}

			tempFile, err := ioutil.TempFile("assets/products", product.Name+"-*"+string(extension))
			if err != nil {
				logger.WithField("err", err.Error()).Error("Error while Creating a Temporary File")
				responses(rw, http.StatusInternalServerError, errorResponse{
					Error: messageObject{
						Message: "Couldn't  create temporary storage!",
					},
				})
				return
			}
			defer tempFile.Close()

			imageBytes, err := ioutil.ReadAll(image)
			if err != nil {
				logger.WithField("err", err.Error()).Error("Error while reading image File")
				responses(rw, http.StatusInternalServerError, errorResponse{
					Error: messageObject{
						Message: "Couldn't read the image file!",
					},
				})
				return
			}
			tempFile.Write(imageBytes)
			product.URLs = append(product.URLs, tempFile.Name())
		}

		var createdProduct db.Product
		createdProduct, err = deps.Store.CreateProduct(req.Context(), product)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error while inserting product")
			rw.Header().Add("Content-Type", "application/json")
			rw.WriteHeader(http.StatusBadRequest)
			responses(rw, http.StatusBadRequest, errorResponse{
				Error: messageObject{
					Message: "Error inserting the product, product already exists",
				},
			})
			return
		}

		responses(rw, http.StatusOK, successResponse{Data: createdProduct})
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
			rw.Header().Add("Content-Type", "application/json")
			rw.WriteHeader(http.StatusBadRequest)
			responses(rw, http.StatusBadRequest, errorResponse{
				Error: messageObject{
					Message: "Error id is missing/invalid",
				},
			})
			return
		}

		err = deps.Store.DeleteProductById(req.Context(), id)
		if err != nil {
			rw.Header().Add("Content-Type", "application/json")
			rw.WriteHeader(http.StatusBadRequest)
			logger.WithField("err", err.Error()).Error("Error fetching data no row found")
			responses(rw, http.StatusInternalServerError, errorResponse{
				Error: messageObject{
					Message: "",
				},
			})
			return
		}

		// rw.WriteHeader(http.StatusOK)
		// rw.Header().Add("Content-Type", "application/json")
		responses(rw, http.StatusOK, successResponse{
			Data: messageObject{
				Message: "Product successfully deleted",
			},
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
			rw.Header().Add("Content-Type", "application/json")
			rw.WriteHeader(http.StatusBadRequest)
			responses(rw, http.StatusBadRequest, errorResponse{
				Error: messageObject{
					Message: "Error id is missing/invalid",
				},
			})
			return
		}

		count, err := strconv.Atoi(Count)

		if Count == "" || err != nil {
			logger.WithField("err", err.Error()).Error("Error stock parameter is missing or corrupt")
			rw.Header().Add("Content-Type", "application/json")
			rw.WriteHeader(http.StatusBadRequest)
			responses(rw, http.StatusBadRequest, errorResponse{
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
			rw.Header().Add("Content-Type", "application/json")
			rw.WriteHeader(http.StatusBadRequest)
			responses(rw, http.StatusBadRequest, errorResponse{
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
			_, err := json.Marshal(errRes)
			if err != nil {
				logger.WithField("err", err.Error()).Error("Error marshalling Product's data")
				responses(rw, http.StatusInternalServerError, errorResponse{
					Error: messageObject{
						Message: "Invalid json body",
					},
				})
				return
			}
			responses(rw, http.StatusBadRequest, errRes)
			return
		}

		var updatedProduct db.Product
		updatedProduct, err = deps.Store.UpdateProductStockById(req.Context(), product, productId)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error while updating product attribute")
			responses(rw, http.StatusInternalServerError, errorResponse{
				Error: messageObject{
					Message: "Internal server error",
				},
			})
			return
		}

		rw.Header().Add("Content-Type", "application/json")
		responses(rw, http.StatusOK, successResponse{Data: updatedProduct})
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
			responses(rw, http.StatusBadRequest, errorResponse{
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
			responses(rw, http.StatusBadRequest, errorResponse{
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
			responses(rw, http.StatusInternalServerError, errorResponse{
				Error: messageObject{
					Message: "Internal server error",
				},
			})
			return
		}

		responses(rw, http.StatusOK, successResponse{Data: updatedProduct})

		return
	})
}

// TODO test case to delete category
