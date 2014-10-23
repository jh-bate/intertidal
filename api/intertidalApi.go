package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

type (
	Api struct{}

	varsHandler func(http.ResponseWriter, *http.Request, map[string]string)

	//Enum type's
	Source       string
	Notification string
)

const (
	//Available Sources
	SourceTrackThis Source = "trackthis"

	NotificationEmail Notification = "email"
	NotificationSms   Notification = "sms"

	TP_SESSION_TOKEN      = "x-tidepool-session-token"
	STATUS_NO_USR_DETAILS = "No user id was given"
)

func InitApi() *Api {
	return &Api{}
}

func (a *Api) SetHandlers(prefix string, rtr *mux.Router) {
	rtr.Handle("/sync/{userid}", varsHandler(a.SyncData)).Methods("GET")
	rtr.Handle("/register/source/{userid}/{type}", varsHandler(a.RegisterSource)).Methods("POST")
	rtr.Handle("/register/notification/{userid}/{type}", varsHandler(a.RegisterNotification)).Methods("POST")
}

func (h varsHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	h(res, req, vars)
}

func (a *Api) RegisterSource(res http.ResponseWriter, req *http.Request, vars map[string]string) {

	id := vars["userid"]
	theType := vars["type"]

	if theType == string(SourceSms) || theType == string(SourceTrackThis) {
		key := req.URL.Query().Get("key")

		if id == "" || key == "" {
			res.WriteHeader(http.StatusBadRequest)
			res.Write([]byte(STATUS_NO_USR_DETAILS))
			return
		}
	}

	res.WriteHeader(http.StatusNotImplemented)
	return
}

func (a *Api) SyncData(res http.ResponseWriter, req *http.Request, vars map[string]string) {

	id := vars["userid"]

	if id == "" {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte(STATUS_NO_USR_DETAILS))
		return
	}
	res.WriteHeader(http.StatusNotImplemented)
	return
}

func (a *Api) RegisterNotification(res http.ResponseWriter, req *http.Request, vars map[string]string) {

	id := vars["userid"]
	theType := vars["type"]

	if theType == string(NotificationSms) || theType == string(NotificationEmail) {

		if id == "" {
			res.WriteHeader(http.StatusBadRequest)
			res.Write([]byte(STATUS_NO_USR_DETAILS))
			return
		}
	}
	res.WriteHeader(http.StatusNotImplemented)
	return
}
