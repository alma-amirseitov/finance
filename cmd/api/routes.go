package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/api/payments",  app.listPaymentsHandler)
	router.HandlerFunc(http.MethodPost, "/api/payments", app.addPaymentHandler)
	router.HandlerFunc(http.MethodGet, "/api/payments/:id", app.showPaymentHandler)
	router.HandlerFunc(http.MethodPatch, "/api/payments/:id",  app.updatePaymentHandler)
	router.HandlerFunc(http.MethodDelete, "/api/payments/:id", app.deletePaymentHandler)

	router.HandlerFunc(http.MethodGet, "/api/categories",  app.listCategoriesHandler)
	router.HandlerFunc(http.MethodPost, "/api/categories", app.addCategoryHandler)
	router.HandlerFunc(http.MethodGet, "/api/categories/:id", app.showCategoryHandler)
	router.HandlerFunc(http.MethodPatch, "/api/categories/:id",  app.editCategoryHandler)
	router.HandlerFunc(http.MethodDelete, "/api/categories/:id", app.deleteCategoryHandler)


	return app.recoverPanic(app.enableCORS(app.rateLimit(router)))
}