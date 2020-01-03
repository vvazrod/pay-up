package gmicro

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// HTTPHandlers for HTTP requests to the groups microservice.
type HTTPHandlers struct {
	GetStatusHandler    func(http.ResponseWriter, *http.Request)
	PostGroupHandler    func(http.ResponseWriter, *http.Request)
	GetGroupHandler     func(http.ResponseWriter, *http.Request)
	PutGroupHandler     func(http.ResponseWriter, *http.Request)
	DeleteGroupHandler  func(http.ResponseWriter, *http.Request)
	PostMemberHandler   func(http.ResponseWriter, *http.Request)
	GetMemberHandler    func(http.ResponseWriter, *http.Request)
	PutMemberHandler    func(http.ResponseWriter, *http.Request)
	DeleteMemberHandler func(http.ResponseWriter, *http.Request)
}

// NewHTTPHandlers that use a given data manager.
func NewHTTPHandlers(m Manager) *HTTPHandlers {
	return &HTTPHandlers{
		GetStatusHandler:    buildGetStatusHandler(),
		PostGroupHandler:    buildPostGroupHandler(m),
		GetGroupHandler:     buildGetGroupHandler(m),
		PutGroupHandler:     buildPutGroupHandler(m),
		DeleteGroupHandler:  buildDeleteGroupHandler(m),
		PostMemberHandler:   buildPostMemberHandler(m),
		GetMemberHandler:    buildGetMemberHandler(m),
		PutMemberHandler:    buildPutMemberHandler(m),
		DeleteMemberHandler: buildDeleteMemberHandler(m),
	}
}

func buildGetStatusHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		status := map[string]string{"status": "OK"}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&status)
	}
}

func buildPostGroupHandler(m Manager) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Body == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var data map[string]interface{}
		err = json.Unmarshal(body, &data)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		name, ok := data["name"].(string)
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		g, err := m.CreateGroup(name)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(&g)
	}
}

func buildGetGroupHandler(m Manager) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		strid := mux.Vars(r)["groupid"]

		id, err := uuid.Parse(strid)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		g, err := m.FetchGroup(id)
		if err != nil {
			if _, ok := err.(*NotFoundError); ok {
				w.WriteHeader(http.StatusNotFound)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}

			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&g)
	}
}

func buildPutGroupHandler(m Manager) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		strid := mux.Vars(r)["groupid"]

		id, err := uuid.Parse(strid)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if r.Body == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var data map[string]interface{}
		err = json.Unmarshal(body, &data)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		name, ok := data["name"].(string)
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = m.UpdateGroup(id, name)
		if err != nil {
			if _, ok := err.(*NotFoundError); ok {
				w.WriteHeader(http.StatusNotFound)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}

			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func buildDeleteGroupHandler(m Manager) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		strid := mux.Vars(r)["groupid"]

		id, err := uuid.Parse(strid)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = m.RemoveGroup(id)
		if err != nil {
			if _, ok := err.(*NotFoundError); ok {
				w.WriteHeader(http.StatusNotFound)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}

			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func buildPostMemberHandler(m Manager) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		strid := mux.Vars(r)["groupid"]

		id, err := uuid.Parse(strid)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if r.Body == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var data map[string]interface{}
		err = json.Unmarshal(body, &data)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		name, ok := data["name"].(string)
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		m, err := m.AddMember(id, name)
		if err != nil {
			if _, ok := err.(*NotFoundError); ok {
				w.WriteHeader(http.StatusNotFound)
			} else if _, ok := err.(*AlreadyPresentError); ok {
				w.WriteHeader(http.StatusConflict)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}

			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&m)
	}
}

func buildGetMemberHandler(m Manager) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		strgroupid, strmemberid := vars["groupid"], vars["memberid"]

		groupid, err := uuid.Parse(strgroupid)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		memberid, err := uuid.Parse(strmemberid)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		m, err := m.FetchMember(groupid, memberid)
		if err != nil {
			if _, ok := err.(*NotFoundError); ok {
				w.WriteHeader(http.StatusNotFound)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}

			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&m)
	}
}

func buildPutMemberHandler(m Manager) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		strgroupid, strmemberid := vars["groupid"], vars["memberid"]

		groupid, err := uuid.Parse(strgroupid)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		memberid, err := uuid.Parse(strmemberid)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if r.Body == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var data map[string]interface{}
		err = json.Unmarshal(body, &data)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		name, ok := data["name"].(string)
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = m.UpdateMember(groupid, memberid, name)
		if err != nil {
			if _, ok := err.(*NotFoundError); ok {
				w.WriteHeader(http.StatusNotFound)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}

			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func buildDeleteMemberHandler(m Manager) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		strgroupid, strmemberid := vars["groupid"], vars["memberid"]

		groupid, err := uuid.Parse(strgroupid)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		memberid, err := uuid.Parse(strmemberid)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = m.RemoveMember(groupid, memberid)
		if err != nil {
			if _, ok := err.(*NotFoundError); ok {
				w.WriteHeader(http.StatusNotFound)
			} else if _, ok := err.(*BalanceError); ok {
				w.WriteHeader(http.StatusConflict)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}

			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
