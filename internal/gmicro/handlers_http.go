package gmicro

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/varrrro/pay-up/internal/gmicro/group"
	"github.com/varrrro/pay-up/internal/gmicro/member"
)

// StatusHandler returns a static message to know the server is working.
func StatusHandler(rw http.ResponseWriter, r *http.Request) {
	status := map[string]string{"status": "OK"}

	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(&status)
}

// GroupsHandler manages requests for creating new groups.
func GroupsHandler(m Manager) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		logger := log.WithFields(log.Fields{
			"uri":    r.URL,
			"method": r.Method,
		})

		// Parse JSON
		var g group.Group
		if err := json.NewDecoder(r.Body).Decode(&g); err != nil {
			logger.WithError(err).Error("Can't parse request body as group")
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		// Create group
		if err := m.CreateGroup(&g); err != nil {
			logger.WithError(err).Warn("Can't create group")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		rw.WriteHeader(http.StatusCreated)
	}
}

// GroupHandler manages requests for fetching, updating or deleting a group.
func GroupHandler(m Manager) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			getGroupHandler(m, rw, r)
			break
		case "PUT":
			putGroupHandler(m, rw, r)
			break
		case "DELETE":
			deleteGroupHandler(m, rw, r)
		}
	}
}

func getGroupHandler(m Manager, rw http.ResponseWriter, r *http.Request) {
	logger := log.WithFields(log.Fields{
		"uri":    r.URL,
		"method": r.Method,
	})

	// Get group ID from request path
	gid, err := uuid.Parse(mux.Vars(r)["groupid"])
	if err != nil {
		logger.WithError(err).Error("Can't parse group ID as UUID")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	// Fetch group
	g, err := m.FetchGroup(gid)
	if err != nil {
		logger.WithError(err).Warn("Can't fetch group")

		if _, ok := err.(*NotFoundError); ok {
			rw.WriteHeader(http.StatusNotFound)
		} else {
			rw.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(&g)
}

func putGroupHandler(m Manager, rw http.ResponseWriter, r *http.Request) {
	logger := log.WithFields(log.Fields{
		"uri":    r.URL,
		"method": r.Method,
	})

	// Get group ID from request path
	gid, err := uuid.Parse(mux.Vars(r)["groupid"])
	if err != nil {
		logger.WithError(err).Error("Can't parse group ID as UUID")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	// Parse JSON
	var g group.Group
	if err := json.NewDecoder(r.Body).Decode(&g); err != nil {
		logger.WithError(err).Error("Can't parse request body as group")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if path and body IDs match
	if g.ID != gid {
		logger.WithFields(log.Fields{
			"path": gid,
			"body": g.ID,
		}).Error("Group IDs in path and body don't match")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	// Update group
	if err := m.UpdateGroup(&g); err != nil {
		logger.WithError(err).Warn("Can't update group")

		if _, ok := err.(*NotFoundError); ok {
			rw.WriteHeader(http.StatusNotFound)
		} else {
			rw.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	rw.WriteHeader(http.StatusOK)
}

func deleteGroupHandler(m Manager, rw http.ResponseWriter, r *http.Request) {
	logger := log.WithFields(log.Fields{
		"uri":    r.URL,
		"method": r.Method,
	})

	// Get group ID from request path
	gid, err := uuid.Parse(mux.Vars(r)["groupid"])
	if err != nil {
		logger.WithError(err).Error("Can't parse group ID as UUID")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	// Remove group
	if err := m.RemoveGroup(gid); err != nil {
		logger.WithError(err).Warn("Can't remove group")

		if _, ok := err.(*NotFoundError); ok {
			rw.WriteHeader(http.StatusNotFound)
		} else {
			rw.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	rw.WriteHeader(http.StatusNoContent)
}

// MembersHandler manages requests for adding new members.
func MembersHandler(m Manager) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		logger := log.WithFields(log.Fields{
			"uri":    r.URL,
			"method": r.Method,
		})

		// Get group ID from request path
		gid, err := uuid.Parse(mux.Vars(r)["groupid"])
		if err != nil {
			logger.WithError(err).Error("Can't parse group ID as UUID")
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		// Parse JSON
		var mb member.Member
		if err := json.NewDecoder(r.Body).Decode(&mb); err != nil {
			logger.WithError(err).Error("Can't parse request body as member")
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		// Add member
		if err := m.AddMember(gid, &mb); err != nil {
			logger.WithError(err).Warn("Can't add member")

			if _, ok := err.(*NotFoundError); ok {
				rw.WriteHeader(http.StatusNotFound)
			} else if _, ok := err.(*AlreadyPresentError); ok {
				rw.WriteHeader(http.StatusConflict)
			} else {
				rw.WriteHeader(http.StatusInternalServerError)
			}

			return
		}

		rw.WriteHeader(http.StatusCreated)
	}
}

// MemberHandler manages requests for fetchinf, updating or deleting members.
func MemberHandler(m Manager) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			getMemberHandler(m, rw, r)
			break
		case "PUT":
			putMemberHandler(m, rw, r)
			break
		case "DELETE":
			deleteMemberHandler(m, rw, r)
		}
	}
}

func getMemberHandler(m Manager, rw http.ResponseWriter, r *http.Request) {
	logger := log.WithFields(log.Fields{
		"uri":    r.URL,
		"method": r.Method,
	})

	// Get group ID from request path
	gid, err := uuid.Parse(mux.Vars(r)["groupid"])
	if err != nil {
		logger.WithError(err).Error("Can't parse group ID as UUID")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get member ID from request path
	mid, err := uuid.Parse(mux.Vars(r)["memberid"])
	if err != nil {
		logger.WithError(err).Error("Can't parse member ID as UUID")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	// Fetch member
	mb, err := m.FetchMember(gid, mid)
	if err != nil {
		logger.WithError(err).Warn("Can't fetch member")

		if _, ok := err.(*NotFoundError); ok {
			rw.WriteHeader(http.StatusNotFound)
		} else {
			rw.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(&mb)
}

func putMemberHandler(m Manager, rw http.ResponseWriter, r *http.Request) {
	logger := log.WithFields(log.Fields{
		"uri":    r.URL,
		"method": r.Method,
	})

	// Get group ID from request path
	gid, err := uuid.Parse(mux.Vars(r)["groupid"])
	if err != nil {
		logger.WithError(err).Error("Can't parse group ID as UUID")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get member ID from request path
	mid, err := uuid.Parse(mux.Vars(r)["memberid"])
	if err != nil {
		logger.WithError(err).Error("Can't parse member ID as UUID")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	// Parse JSON
	var mb member.Member
	if err := json.NewDecoder(r.Body).Decode(&mb); err != nil {
		logger.WithError(err).Error("Can't parse request body as member")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if path and body IDs match
	if mb.ID != mid {
		logger.WithFields(log.Fields{
			"path": mid,
			"body": mb.ID,
		}).Error("Group IDs in path and body don't match")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	// Update member
	if err := m.UpdateMember(gid, &mb); err != nil {
		logger.WithError(err).Warn("Can't update member")

		if _, ok := err.(*NotFoundError); ok {
			rw.WriteHeader(http.StatusNotFound)
		} else {
			rw.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	rw.WriteHeader(http.StatusOK)
}

func deleteMemberHandler(m Manager, rw http.ResponseWriter, r *http.Request) {
	logger := log.WithFields(log.Fields{
		"uri":    r.URL,
		"method": r.Method,
	})

	// Get group ID from request path
	gid, err := uuid.Parse(mux.Vars(r)["groupid"])
	if err != nil {
		logger.WithError(err).Error("Can't parse group ID as UUID")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get member ID from request path
	mid, err := uuid.Parse(mux.Vars(r)["memberid"])
	if err != nil {
		logger.WithError(err).Error("Can't parse member ID as UUID")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	// Remove member
	if err := m.RemoveMember(gid, mid); err != nil {
		logger.WithError(err).Warn("Can't update member")

		if _, ok := err.(*NotFoundError); ok {
			rw.WriteHeader(http.StatusNotFound)
		} else if _, ok := err.(*BalanceError); ok {
			rw.WriteHeader(http.StatusConflict)
		} else {
			rw.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
