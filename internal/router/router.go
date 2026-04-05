package router

import (
	"codecommerceapi/internal/handler"
	"codecommerceapi/internal/service"

	"github.com/gorilla/mux"
)

func New(catSvc *service.CategoryService, prodSvc *service.ProductService) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/api/health", handler.Health).Methods("GET")

	catHandler := handler.NewCategoryHandler(catSvc)
	r.HandleFunc("/api/categories", catHandler.CreateCategory).Methods("POST")
	r.HandleFunc("/api/categories", catHandler.ListCategories).Methods("GET")
	r.HandleFunc("/api/categories/{id}", catHandler.GetCategory).Methods("GET")
	r.HandleFunc("/api/categories/{id}", catHandler.UpdateCategory).Methods("PUT")
	r.HandleFunc("/api/categories/{id}", catHandler.DeleteCategory).Methods("DELETE")

	prodHandler := handler.NewProductHandler(prodSvc)
	r.HandleFunc("/api/products", prodHandler.CreateProduct).Methods("POST")
	r.HandleFunc("/api/products", prodHandler.ListProducts).Methods("GET")
	r.HandleFunc("/api/products/{id}", prodHandler.GetProduct).Methods("GET")
	r.HandleFunc("/api/products/{id}", prodHandler.UpdateProduct).Methods("PUT")
	r.HandleFunc("/api/products/{id}", prodHandler.DeleteProduct).Methods("DELETE")
	r.HandleFunc("/api/categories/{id}/products", prodHandler.ListProductsByCategory).Methods("GET")

	return r
}
