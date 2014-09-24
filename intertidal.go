package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jh-bate/intertidal/api"
)

func main() {
	/*
	 * Service setup
	 */

	rtr := mux.NewRouter()
	api := api.InitApi()
	api.SetHandlers("", rtr)

	http.ListenAndServe("3000", rtr)

}
