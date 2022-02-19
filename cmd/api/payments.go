package main

import (
	"errors"
	"fmt"
	"github.com/alma-amirseitov/finance/internal/data"
	"github.com/alma-amirseitov/finance/internal/validator"
	"net/http"
)

func (app *application) addPaymentHandler(w http.ResponseWriter, r *http.Request){
	var input struct {
		Name string `json:"name"`
		PaymentType string `json:"payment_type"`
		Comment string `json:"comment"`
		CategoriesName string `json:"category_name"`
		Price int32 `json:"price"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	category,err := app.models.Categories.GetByName(input.CategoriesName)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			category = &data.Categories{
			Name: input.CategoriesName,
			}
			err = app.models.Categories.Insert(category)
			if err != nil {
				app.serverErrorResponse(w, r, err)
				return
			}
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	payment := &data.Payments{
		Name:   input.Name,
		PaymentType:    input.PaymentType,
		Comment: input.Comment,
		Categories:  *category,
		Price: input.Price,
	}

	v := validator.New()

	if data.ValidatePayment(v, payment); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	err = app.models.Payments.Insert(payment)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("api/payments/%d", payment.Id))

	err = app.writeJSON(w, http.StatusCreated, envelope{"payments": payment}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showPaymentHandler(w http.ResponseWriter, r *http.Request){
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	payment, err := app.models.Payments.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"payment": payment}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listPaymentsHandler(w http.ResponseWriter, r *http.Request){
	var input struct {
		Name string
		CategoryName string

		data.Filters
	}

	v := validator.New()

	qs := r.URL.Query()

	input.Name = app.readString(qs, "name", "")
	input.CategoryName = app.readString(qs, "category_name", "")

	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "name", "payment_type", "date","category_name","price", "-id", "-name", "-payment_type", "-date","-category_name","-price"}


	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	payments, err := app.models.Payments.GetAll(input.Name,input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"payments": payments}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updatePaymentHandler(w http.ResponseWriter, r *http.Request){
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	payment, err := app.models.Payments.Get(id)

	category,err := app.models.Categories.GetByName(payment.CategoriesName)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	payment.Categories.Id = category.Id
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Name *string `json:"name"`
		PaymentType *string `json:"payment_type"`
		Comment *string `json:"comment"`
		CategoryId *int64 `json:"category_id"`
		Price *int32 `json:"price"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Name != nil {
		payment.Name = *input.Name
	}
	if input.PaymentType != nil {
		payment.PaymentType = *input.PaymentType
	}
	if input.Comment != nil {
		payment.Comment = *input.Comment
	}
	if input.CategoryId != nil {
		payment.Categories.Id = *input.CategoryId
	}
	if input.Price != nil {
		payment.Price = *input.Price
	}

	v := validator.New()

	if data.ValidatePayment(v, payment); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Payments.Update(payment)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"payments": payment}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deletePaymentHandler(w http.ResponseWriter, r *http.Request){
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Payments.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "payment successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
