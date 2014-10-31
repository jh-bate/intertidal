package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	/*
	 * Service setup
	 */

	rtr := mux.NewRouter()
	api := InitApi()
	api.SetHandlers("", rtr)

	http.ListenAndServe("3000", rtr)

}
