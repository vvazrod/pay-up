package gmicro

import "github.com/gorilla/mux"

// NewRouter with the given handlers.
func NewRouter(h *Handlers) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", h.GetStatusHandler).Methods("GET")
	router.HandleFunc("/groups", h.PostGroupHandler).Methods("POST")
	router.HandleFunc("/groups/{groupid}", h.GetGroupHandler).Methods("GET")
	router.HandleFunc("/groups/{groupid}", h.PostMemberHandler).Methods("POST")
	router.HandleFunc("/groups/{groupid}", h.PutGroupHandler).Methods("PUT")
	router.HandleFunc("/groups/{groupid}", h.DeleteGroupHandler).Methods("DELETE")
	router.HandleFunc("/groups/{groupid}/members/{memberid}", h.GetMemberHandler).Methods("GET")
	router.HandleFunc("/groups/{groupid}/members/{memberid}", h.PutMemberHandler).Methods("PUT")
	router.HandleFunc("/groups/{groupid}/members/{memberid}", h.DeleteMemberHandler).Methods("DELETE")

	return router
}
