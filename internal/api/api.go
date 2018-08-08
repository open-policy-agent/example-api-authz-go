// Copyright 2018 The OPA Authors.  All rights reserved.
// Use of this source code is governed by an Apache2
// license that can be found in the LICENSE file.

package api

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"

	"github.com/gorilla/mux"
	"github.com/open-policy-agent/example-api-authz-go/internal/opa"
)

type apiError struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message,omitempty"`
	} `json:"error"`
}

type apiRouteNotFoundError struct {
	Error struct {
		Code   string   `json:"code"`
		Routes []string `json:"routes"`
	} `json:"error"`
}

type apiWrapper struct {
	Result interface{} `json:"result"`
}

const (
	apiCodeNotFound      = "not_found"
	apiCodeParseError    = "parse_error"
	apiCodeInternalError = "internal_error"
	apiCodeNotAuthorized = "not_authorized"
)

// API implements a simple HTTP API server to expose Car data.
type API struct {
	engine *opa.OPA
	router *mux.Router
	db     DB
}

// New returns a instance of the API.
func New(engine *opa.OPA) *API {

	api := &API{
		engine: engine,
		db:     mockDB(),
	}

	api.router = mux.NewRouter()
	api.router.HandleFunc("/", api.handleIndex).Methods(http.MethodGet)
	api.router.HandleFunc("/cars", api.handlGetCars).Methods(http.MethodGet)
	api.router.HandleFunc("/cars/{id}", api.handlePutCar).Methods(http.MethodPut)
	api.router.HandleFunc("/cars/{id}", api.handleGetCar).Methods(http.MethodGet)
	api.router.HandleFunc("/cars/{id}", api.handleDeleteCar).Methods(http.MethodDelete)
	api.router.HandleFunc("/cars/{id}/status", api.handlePutCarStatus).Methods(http.MethodPut)
	api.router.HandleFunc("/cars/{id}/status", api.handleGetCarStatus).Methods(http.MethodGet)
	api.router.NotFoundHandler = http.HandlerFunc(api.handleNotFound)
	api.router.Use(api.authorize)

	return api
}

// Run starts the HTTP server.
func (api *API) Run(ctx context.Context) error {
	return http.ListenAndServe(":8080", api.router)
}

func (api *API) getRoutes() []string {
	routes := make([]string, 0)
	api.router.Walk(func(route *mux.Route, _ *mux.Router, _ []*mux.Route) error {
		tmpl, err := route.GetPathTemplate()
		if err != nil {
			return err
		}
		found := false
		for i := range routes {
			if routes[i] == tmpl {
				found = true
				break
			}
		}
		if !found {
			routes = append(routes, tmpl)
		}
		return nil
	})
	return routes
}

func (api *API) handleNotFound(w http.ResponseWriter, r *http.Request) {
	var resp apiRouteNotFoundError
	resp.Error.Routes = api.getRoutes()
	resp.Error.Code = apiCodeNotFound
	writeJSON(w, http.StatusNotFound, resp)
}

func (api *API) handleIndex(w http.ResponseWriter, r *http.Request) {
	routes := api.getRoutes()
	writeJSON(w, http.StatusOK, routes)
}

func (api *API) handlGetCars(w http.ResponseWriter, r *http.Request) {

	cars := make([]Car, 0, len(api.db.Cars))

	for _, car := range api.db.Cars {
		cars = append(cars, car)
	}

	sort.Slice(cars, func(i, j int) bool {
		return cars[i].ID < cars[j].ID
	})

	writeJSON(w, http.StatusOK, apiWrapper{
		Result: cars,
	})
}

func (api *API) handlePutCar(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusInternalServerError, apiCodeInternalError, err)
		return
	}

	var car Car
	if err := json.Unmarshal(bs, &car); err != nil {
		writeError(w, http.StatusBadRequest, apiCodeParseError, err)
		return
	}

	id := vars["id"]

	api.db.Cars[id] = car

	writeJSON(w, http.StatusOK, car)
}

func (api *API) handleDeleteCar(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	id := vars["id"]
	if car, ok := api.db.Cars[id]; !ok {
		writeError(w, http.StatusNotFound, apiCodeNotFound, nil)
	} else {
		delete(api.db.Cars, id)
		delete(api.db.Statuses, id)
		writeJSON(w, http.StatusOK, car)
	}
}

func (api *API) handleGetCar(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	id := vars["id"]
	if car, ok := api.db.Cars[id]; !ok {
		writeError(w, http.StatusNotFound, apiCodeNotFound, nil)
	} else {
		writeJSON(w, http.StatusOK, car)
	}
}

func (api *API) handleGetCarStatus(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	id := vars["id"]
	if status, ok := api.db.Statuses[id]; !ok {
		writeError(w, http.StatusNotFound, apiCodeNotFound, nil)
	} else {
		writeJSON(w, http.StatusOK, status)
	}
}

func (api *API) handlePutCarStatus(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusInternalServerError, apiCodeInternalError, err)
		return
	}

	var status CarStatus
	if err := json.Unmarshal(bs, &status); err != nil {
		writeError(w, http.StatusBadRequest, apiCodeParseError, err)
		return
	}

	id := vars["id"]

	api.db.Statuses[id] = status
	writeJSON(w, http.StatusOK, status)
}

func (api *API) authorize(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		user := r.Header.Get("Authorization")
		path := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

		input := map[string]interface{}{
			"method": r.Method,
			"path":   path,
			"user":   user,
		}

		allowed, err := api.engine.Bool(r.Context(), input)

		if err != nil {
			writeError(w, http.StatusInternalServerError, apiCodeInternalError, err)
		} else if !allowed {
			writeError(w, http.StatusForbidden, apiCodeNotAuthorized, nil)
		} else {
			next.ServeHTTP(w, r)
		}
	})

}

func writeError(w http.ResponseWriter, status int, code string, err error) {
	var resp apiError
	resp.Error.Code = code
	if err != nil {
		resp.Error.Message = err.Error()
	}
	writeJSON(w, status, resp)
}

func writeJSON(w http.ResponseWriter, status int, x interface{}) {
	bs, _ := json.Marshal(x)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(bs)
}
