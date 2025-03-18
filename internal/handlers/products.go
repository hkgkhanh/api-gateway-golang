package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"api-gateway-golang/internal/model"
	"api-gateway-golang/pkg/logger"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// GetProducts godoc
//
//	@Summary		Get all products
//	@Tags			Products
//	@Produce		json
//	@Success		200	{object} []model.GetProductResponse
//
// @Router			/api/products [get]
func (h *Handlers) GetProductsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	products, err := h.Storage.GetProducts(ctx)
	if err != nil {
		h.Sender.JSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	productResponse := []model.GetProductResponse{}
	for _, product := range products {
		productResponse = append(productResponse, model.GetProductResponse{
			ID:           product.ID,
			Name:         product.Name,
			Manufacturer: product.Manufacturer,
			Quantity:     product.Quantity,
		})

	}

	err = h.Sender.JSON(w, http.StatusOK, productResponse)
	if err != nil {
		logger.OutputLog.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Fatal("Error when requesting /products")

		panic(err)
	}
}

// GetProduct godoc
//
//	@Summary		Get a specific product
//	@Tags			Products
//	@Produce		json
//	@Param			id path int	true "Product ID"
//	@Success		200	{object} model.GetProductResponse
//
// @Router			/api/product/{id} [get]
func (h *Handlers) GetProductHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.Sender.JSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	var product model.Product
	product, err = h.Storage.GetProduct(ctx, id)
	if err != nil {
		h.Sender.JSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	if (model.Product{}) == product {
		h.Sender.JSON(w, http.StatusBadRequest, "Product with id="+fmt.Sprint(product.ID)+" not found")
		if err != nil {
			panic(err)
		}
		return
	}

	productResponse := model.GetProductResponse{
		ID:           product.ID,
		Name:         product.Name,
		Manufacturer: product.Manufacturer,
		Quantity:     product.Quantity,
	}

	err = h.Sender.JSON(w, http.StatusOK, productResponse)
	if err != nil {
		logger.OutputLog.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Fatal(fmt.Sprint("Error when requesting /product/", product.ID))

		panic(err)
	}
}

// AddProduct godoc
//
//		@Summary		Add a specific product
//		@Tags			Products
//		@Produce		json
//	 	@Accept			json
//		@Param			name body string true "Product name"
//		@Param			manufacturer body string true "Product manufacturer"
//		@Param			quantity body string true "Product quantity"
//		@Success		200	{object} model.IDResponse
//
// @Router			/api/product/add [post]
func (h *Handlers) AddProductHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var product model.AddProductRequest

	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		h.Sender.JSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = Validate.Struct(product)
	if err != nil {
		var errs []string
		for _, err := range err.(validator.ValidationErrors) {
			errs = append(errs, err.Field()+" "+err.Tag())
		}
		h.Sender.JSON(w, http.StatusBadRequest, strings.Join(errs, ", "))
		return
	}

	id, err := h.Storage.AddProduct(ctx, product)
	if err != nil {
		h.Sender.JSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := model.IDResponse{ID: id}
	err = h.Sender.JSON(w, http.StatusOK, response)
	if err != nil {
		panic(err)
	}
}

// UpdateProduct godoc
//
//		@Summary		Update a specific product
//		@Tags			Products
//		@Produce		json
//	 	@Accept			json
//		@Param			name body string false "Product name"
//		@Param			manufacturer body string false "Product manufacturer"
//		@Param			quantity body string false "Product quantity"
//		@Success		200	{object} model.IDResponse
//
// @Router			/api/product/update [patch]
func (h *Handlers) UpdateProductHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var product model.UpdateProductRequest

	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		h.Sender.JSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = Validate.Struct(product)
	if err != nil {
		var errs []string
		for _, err := range err.(validator.ValidationErrors) {
			errs = append(errs, err.Field()+" "+err.Tag())
		}
		h.Sender.JSON(w, http.StatusBadRequest, strings.Join(errs, ", "))
		return
	}

	exists, err := h.Storage.VerifyProductExists(ctx, product.ID)
	if err != nil {
		h.Sender.JSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	if !exists {
		h.Sender.JSON(w, http.StatusBadRequest, "Product with id="+fmt.Sprint(product.ID)+" not found")
		return
	}

	id, err := h.Storage.UpdateProduct(ctx, product)
	if err != nil {
		h.Sender.JSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := model.IDResponse{ID: id}
	err = h.Sender.JSON(w, http.StatusOK, response)
	if err != nil {
		panic(err)
	}
}

// DeleteProduct godoc
//
//	@Summary		Delete a specific product
//	@Tags			Products
//	@Produce		json
//	@Param			id path int	true "Product ID"
//	@Success		200	{object} model.GetProductResponse
//
// @Router			/api/product/delete/{id} [delete]
func (h *Handlers) DeleteProductHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.Sender.JSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	exists, err := h.Storage.VerifyProductExists(ctx, id)
	if err != nil {
		h.Sender.JSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	if !exists {
		err = h.Sender.JSON(w, http.StatusBadRequest, "Product with id="+fmt.Sprint(id)+" not found")
		if err != nil {
			panic(err)
		}
		return
	}

	err = h.Storage.DeleteProduct(ctx, id)
	if err != nil {
		h.Sender.JSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = h.Sender.JSON(w, http.StatusOK, map[string]bool{"success": true})
	if err != nil {
		panic(err)
	}
}

// cSpell:ignore godoc logrus
