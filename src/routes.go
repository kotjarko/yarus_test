package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"math"
	"net/http"
	"strconv"
	"strings"
)

const templateError = "{\"error\":\"%s\"}"
const templateCurrency = "{\"%s\":\"%s\"}"
const validXApiKey = "123321"

func AuthCheck(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-API-KEY") != validXApiKey {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(makeErrorResponse("unauthorized"))
		} else {
			next.ServeHTTP(w, r)
		}
	}

	return http.HandlerFunc(fn)
}

func (a *Application) createHandler() http.Handler {
	mux := chi.NewMux()
	mux.Route("/exchange", func(r chi.Router) {
		r.Use(middleware.RequestID)
		r.Use(middleware.RealIP)
		if a.Config.Debug {
			r.Use(middleware.Logger)
		}
		r.Use(middleware.Recoverer)
		
		// use middleware for authorization
		r.Use(AuthCheck)

		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			item, err := a.CurrencyBase.GetRandomCurrency()
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write(makeErrorResponse(err.Error()))
			} else {
				w.WriteHeader(http.StatusOK)
				w.Write(makeResponse(item.CharCode, MakeCurrencyMessage(item.Nominal, item.Value, item.Name)))
			}
		})

		r.Get("/{currency_name}", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			currencyName := chi.URLParam(r, "currency_name")
			item, err := a.CurrencyBase.GetCurrency(currencyName)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write(makeErrorResponse(err.Error()))
			} else {
				w.WriteHeader(http.StatusOK)
				w.Write(makeResponse(item.CharCode, MakeCurrencyMessage(item.Nominal, item.Value, item.Name)))
			}
		})

	})
	return mux
}

func makeErrorResponse(errorText string) []byte {
	str := fmt.Sprintf(templateError, errorText)
	return []byte(str)
}

func makeResponse(currency, message string) []byte {
	str := fmt.Sprintf(templateCurrency, currency, message)
	return []byte(str)
}

func MakeCurrencyMessage(nominal int, value float64, name string) string {
	resultStr := strconv.Itoa(nominal) + " " + strings.ToLower(name)
	if nominal == 1 {
		resultStr += " равен "
	} else {
		resultStr += " равны "
	}

	resultStr += strconv.FormatFloat(value, 'f', 4, 64)

	checkRub := int(math.Floor(value))
	if (checkRub % 100 != 11) && (checkRub % 10 == 1) {
		resultStr += " рублю"
	} else {
		resultStr += " рублям"
	}

	return resultStr
}
