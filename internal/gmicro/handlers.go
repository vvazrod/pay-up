package gmicro

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// Handlers for HTTP requests to the groups microservice.
type Handlers struct {
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

// NewHandlers that use a given data manager.
func NewHandlers(gm *GroupsManager) *Handlers {
	return &Handlers{
		GetStatusHandler:    buildGetStatusHandler(),
		PostGroupHandler:    buildPostGroupHandler(gm),
		GetGroupHandler:     buildGetGroupHandler(gm),
		PutGroupHandler:     buildPutGroupHandler(gm),
		DeleteGroupHandler:  buildDeleteGroupHandler(gm),
		PostMemberHandler:   buildPostMemberHandler(gm),
		GetMemberHandler:    buildGetMemberHandler(gm),
		PutMemberHandler:    buildPutMemberHandler(gm),
		DeleteMemberHandler: buildDeleteMemberHandler(gm),
	}
}

func buildGetStatusHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		status := map[string]string{"status": "OK"}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&status)
	}
}

func buildPostGroupHandler(gm *GroupsManager) func(http.ResponseWriter, *http.Request) {
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

		g, err := gm.CreateGroup(name)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(&g)
	}
}

func buildGetGroupHandler(gm *GroupsManager) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		strid := mux.Vars(r)["groupid"]

		id, err := uuid.Parse(strid)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		g, err := gm.FetchGroup(id)
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

func buildPutGroupHandler(gm *GroupsManager) func(http.ResponseWriter, *http.Request) {
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

		err = gm.UpdateGroup(id, name)
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

func buildDeleteGroupHandler(gm *GroupsManager) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		strid := mux.Vars(r)["groupid"]

		id, err := uuid.Parse(strid)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = gm.RemoveGroup(id)
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

func buildPostMemberHandler(gm *GroupsManager) func(http.ResponseWriter, *http.Request) {
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

		m, err := gm.AddMember(id, name)
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

func buildGetMemberHandler(gm *GroupsManager) func(http.ResponseWriter, *http.Request) {
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

		m, err := gm.FetchMember(groupid, memberid)
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

func buildPutMemberHandler(gm *GroupsManager) func(http.ResponseWriter, *http.Request) {
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

		err = gm.UpdateMember(groupid, memberid, name)
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

func buildDeleteMemberHandler(gm *GroupsManager) func(http.ResponseWriter, *http.Request) {
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

		err = gm.RemoveMember(groupid, memberid)
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
