// Package classification of Product API
//
// Documentation for Prodcut API
//
// Schemes: http
// BasePath: /
// Version: 1.0.0
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
// swagger:meta
package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/DonnachaHeff/goMicro/data"
	"github.com/gorilla/mux"
)

// list of products returned in response
// swagger:response productsResponse
type productsResponseWrapper struct {
	// All products in the system
	// in: body
	Body []data.Product
}

// swagger:parameters deleteProduct updateProduct
// deleteProduct & updateProduct act as a reference to the DELETE/PUT route comment
type productIDParameterWrapper struct {
	// The id of the product to be deleted/updated from the db
	// in: path
	// required: true
	ID int `json:"id"`
}

type errorResponseWrapper struct {
	Body GenericError
}

type validationErrorWrapper struct {
	Body ValidationError
}

// swagger:reponse noContent
type productsNoContent struct {}

// http.Handler
type Products struct {
	l *log.Logger
}

type GenericError struct {
	Message string `json:"message"`
}

type ValidationError struct {
	Message string `json:"message"`
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

// swagger:route GET /products products listProducts
// Returns a list of products
// responses:
// 200: productsResponse
func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET Products")

	lp := data.GetProducts()
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

// swagger:route DELETE /products/{id} products deleteProducts
// responses:
// 201: noContent
// 404: errorResponse
// 500: errorResponse
func (p *Products) DeleteProduct(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	p.l.Println("Handle DELETE Product", id)

	err := data.DeleteProduct(id)

	if err == data.ErrProductNotFound {
		http.Error(rw, "Product Not Found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Product Not Found", http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}

// swagger:route POST /products products addProducts
// Adds a product to the list of products
// responses:
// 201: statusCreated
func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request){
	p.l.Println("Handle POST Products")

	prod := r.Context().Value(KeyProduct{}).(data.Product) // cast to product
	data.AddProduct(&prod)
	rw.WriteHeader(http.StatusCreated)
}

// swagger:route PUT /products/{id} products updateProducts
// Updates a product within the list of products
// responses:
// 202: statusAccepted
// 400: validationErrorResponse
// 404: errorResponse
// 500: errorResponse
func (p *Products) UpdateProducts(rw http.ResponseWriter, r*http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		http.Error(rw, "Unable to convert string", http.StatusBadRequest)
		return
	}

	p.l.Println("Handle PUT Products", id)

	prod := r.Context().Value(KeyProduct{}).(data.Product) // cast to product

	err = data.UpdateProduct(id, &prod)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product Not Found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Product Not Found", http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusAccepted)
}

type KeyProduct struct{}

func (p Products) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		prod := &data.Product{}

		err := prod.FromJSON(r.Body)
		if err != nil {
			p.l.Println("[Error] validating product", err)
			http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
			return
		}

		err = prod.Validate()
		if err != nil {
			p.l.Println("[Error] validating product", err)
			http.Error(rw, fmt.Sprintf("Error validating product: %s", err), http.StatusBadRequest)
			return
		}

		// add the product to the context
		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		req := r.WithContext(ctx)

		next.ServeHTTP(rw, req)
	})
}
