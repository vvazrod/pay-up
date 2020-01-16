package gmicro_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/varrrro/pay-up/internal/gmicro/group"
	"github.com/varrrro/pay-up/internal/gmicro/member"
)

func TestStatusHandler(t *testing.T) {
	// Create request
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Errorf("Can't create request [Error]: %v", err)
	}

	// Serve test request
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	res := rec.Result() // get response
	defer res.Body.Close()

	// Check response Content-Type and status code
	if res.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Wrong content type [Expected]: %s [Actual]: %s", "application/json", res.Header.Get("Content-Type"))
	} else if res.StatusCode != http.StatusOK {
		t.Errorf("Wrong status code [Expected]: %d [Actual]: %d", http.StatusOK, res.StatusCode)
	}

	// Decode request body
	var data map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		t.Errorf("Can't decode JSON from response body [Error]: %v", err)
	}

	// Check if 'status' key is present
	if _, ok := data["status"]; !ok {
		t.Error("Key 'status' isn't present")
	}

	clearDB()
}

func TestGroupsHandler(t *testing.T) {
	g := group.Group{ID: uuid.New(), Name: "Test"}
	body, _ := json.Marshal(&g)

	cases := []struct {
		method     string
		reqBody    []byte
		statusCode int
	}{
		{"POST", body, http.StatusCreated},
		{"POST", []byte(""), http.StatusBadRequest},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%s %d", tc.method, tc.statusCode), func(t *testing.T) {
			// Create request
			req, err := http.NewRequest(tc.method, "/groups", bytes.NewBuffer(tc.reqBody))
			if err != nil {
				t.Errorf("Can't create request [Error]: %v", err)
			}

			// Serve test request
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)
			res := rec.Result() // get response
			defer res.Body.Close()

			// Check response Content-Type and status code
			if res.Header.Get("Content-Type") != "application/json" {
				t.Errorf("Wrong content type [Expected]: %s [Actual]: %s", "application/json", res.Header.Get("Content-Type"))
			} else if res.StatusCode != tc.statusCode {
				t.Errorf("Wrong status code [Expected]: %d [Actual]: %d", http.StatusOK, res.StatusCode)
			}
		})
	}
}

func TestGroupHandler(t *testing.T) {
	g := group.Group{ID: uuid.New(), Name: "Test"}
	gm.CreateGroup(&g)

	g2 := group.Group{ID: g.ID, Name: "Updated"}
	body1, _ := json.Marshal(&g2)

	g3 := group.Group{ID: uuid.New(), Name: "Fail"}
	g4 := group.Group{ID: g3.ID, Name: "UpdatedFail"}
	body2, _ := json.Marshal(&g4)

	cases := []struct {
		method     string
		gid        string
		reqBody    []byte
		statusCode int
		resBody    *group.Group
	}{
		{"GET", g.ID.String(), nil, http.StatusOK, &g},
		{"GET", "test", nil, http.StatusBadRequest, nil},
		{"GET", uuid.New().String(), nil, http.StatusNotFound, nil},

		{"PUT", g.ID.String(), body1, http.StatusOK, nil},
		{"PUT", "test", body1, http.StatusBadRequest, nil},
		{"PUT", g.ID.String(), []byte(""), http.StatusBadRequest, nil},
		{"PUT", uuid.New().String(), body1, http.StatusBadRequest, nil},
		{"PUT", g3.ID.String(), body2, http.StatusNotFound, nil},

		{"DELETE", g.ID.String(), nil, http.StatusNoContent, nil},
		{"DELETE", "test", nil, http.StatusBadRequest, nil},
		{"DELETE", g3.ID.String(), nil, http.StatusNotFound, nil},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%s %d", tc.method, tc.statusCode), func(t *testing.T) {
			// Create request
			var req *http.Request
			var err error
			if tc.method == "PUT" {
				req, err = http.NewRequest(tc.method, "/groups/"+tc.gid, bytes.NewBuffer(tc.reqBody))
			} else {
				req, err = http.NewRequest(tc.method, "/groups/"+tc.gid, nil)
			}
			if err != nil {
				t.Errorf("Can't create request [Error]: %v", err)
			}

			// Serve test request
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)
			res := rec.Result() // get response
			defer res.Body.Close()

			// Check response Content-Type and status code
			if res.Header.Get("Content-Type") != "application/json" {
				t.Errorf("Wrong content type [Expected]: %s [Actual]: %s", "application/json", res.Header.Get("Content-Type"))
			} else if res.StatusCode != tc.statusCode {
				t.Errorf("Wrong status code [Expected]: %d [Actual]: %d", http.StatusOK, res.StatusCode)
			}

			// Check response body
			if tc.method == "GET" &&
				tc.statusCode != http.StatusBadRequest &&
				tc.statusCode != http.StatusNotFound {
				// Decode response body
				var g group.Group
				if err := json.NewDecoder(res.Body).Decode(&g); err != nil {
					t.Errorf("Can't decode response body [Error]: %v", err)
				}

				// Check values
				if g.ID != tc.resBody.ID {
					t.Errorf("IDs don't match [Expected]: %v [Actual]: %v", tc.resBody.ID, g.ID)
				} else if g.Name != tc.resBody.Name {
					t.Errorf("Names don't match [Expected]: %s [Actual]: %s", tc.resBody.Name, g.Name)
				}
			}
		})
	}
}

func TestMembersHandler(t *testing.T) {
	g := group.Group{ID: uuid.New(), Name: "Test"}
	gm.CreateGroup(&g)

	m := member.Member{ID: uuid.New(), Name: "Test"}
	body, _ := json.Marshal(&m)

	cases := []struct {
		method     string
		gid        string
		reqBody    []byte
		statusCode int
	}{
		{"POST", g.ID.String(), body, http.StatusCreated},
		{"POST", "test", body, http.StatusBadRequest},
		{"POST", g.ID.String(), []byte(""), http.StatusBadRequest},
		{"POST", uuid.New().String(), body, http.StatusNotFound},
		{"POST", g.ID.String(), body, http.StatusConflict},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%s %d", tc.method, tc.statusCode), func(t *testing.T) {
			// Create request
			req, err := http.NewRequest(tc.method, "/groups/"+tc.gid+"/members", bytes.NewBuffer(tc.reqBody))
			if err != nil {
				t.Errorf("Can't create request [Error]: %v", err)
			}

			// Serve test request
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)
			res := rec.Result() // get response
			defer res.Body.Close()

			// Check response Content-Type and status code
			if res.Header.Get("Content-Type") != "application/json" {
				t.Errorf("Wrong content type [Expected]: %s [Actual]: %s", "application/json", res.Header.Get("Content-Type"))
			} else if res.StatusCode != tc.statusCode {
				t.Errorf("Wrong status code [Expected]: %d [Actual]: %d", http.StatusOK, res.StatusCode)
			}
		})
	}
}

func TestMemberHandler(t *testing.T) {
	g := group.Group{ID: uuid.New(), Name: "Test"}
	gm.CreateGroup(&g)

	m := member.Member{ID: uuid.New(), Name: "Test"}
	gm.AddMember(g.ID, &m)

	m2 := member.Member{ID: m.ID, Name: "Updated"}
	body, _ := json.Marshal(&m2)

	m3 := member.Member{ID: uuid.New(), Name: "Balance", Balance: 23.3}
	gm.AddMember(g.ID, &m3)

	cases := []struct {
		method     string
		gid        string
		mid        string
		reqBody    []byte
		statusCode int
		resBody    *member.Member
	}{
		{"GET", g.ID.String(), m.ID.String(), nil, http.StatusOK, &m},
		{"GET", "test", m.ID.String(), nil, http.StatusBadRequest, &m},
		{"GET", g.ID.String(), "test", nil, http.StatusBadRequest, &m},
		{"GET", uuid.New().String(), m.ID.String(), nil, http.StatusNotFound, &m},
		{"GET", g.ID.String(), uuid.New().String(), nil, http.StatusNotFound, &m},

		{"PUT", g.ID.String(), m.ID.String(), body, http.StatusOK, nil},
		{"PUT", "test", m.ID.String(), body, http.StatusBadRequest, nil},
		{"PUT", g.ID.String(), "test", body, http.StatusBadRequest, nil},
		{"PUT", g.ID.String(), uuid.New().String(), body, http.StatusBadRequest, nil},
		{"PUT", uuid.New().String(), m.ID.String(), body, http.StatusNotFound, nil},
		{"PUT", g.ID.String(), m.ID.String(), body, http.StatusConflict, nil},

		{"DELETE", g.ID.String(), m.ID.String(), nil, http.StatusNoContent, nil},
		{"DELETE", "test", m.ID.String(), nil, http.StatusBadRequest, nil},
		{"DELETE", g.ID.String(), "test", nil, http.StatusBadRequest, nil},
		{"DELETE", uuid.New().String(), m.ID.String(), nil, http.StatusNotFound, nil},
		{"DELETE", g.ID.String(), m3.ID.String(), nil, http.StatusConflict, nil},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%s %d", tc.method, tc.statusCode), func(t *testing.T) {
			// Create request
			var req *http.Request
			var err error
			if tc.method == "PUT" {
				req, err = http.NewRequest(tc.method, "/groups/"+tc.gid+"/members/"+tc.mid, bytes.NewBuffer(tc.reqBody))
			} else {
				req, err = http.NewRequest(tc.method, "/groups/"+tc.gid+"/members/"+tc.mid, nil)
			}
			if err != nil {
				t.Errorf("Can't create request [Error]: %v", err)
			}

			// Serve test request
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)
			res := rec.Result() // get response
			defer res.Body.Close()

			// Check response Content-Type and status code
			if res.Header.Get("Content-Type") != "application/json" {
				t.Errorf("Wrong content type [Expected]: %s [Actual]: %s", "application/json", res.Header.Get("Content-Type"))
			} else if res.StatusCode != tc.statusCode {
				t.Errorf("Wrong status code [Expected]: %d [Actual]: %d", http.StatusOK, res.StatusCode)
			}

			// Check response body
			if tc.method == "GET" &&
				tc.statusCode != http.StatusBadRequest &&
				tc.statusCode != http.StatusNotFound {
				// Decode response body
				var g group.Group
				if err := json.NewDecoder(res.Body).Decode(&g); err != nil {
					t.Errorf("Can't decode response body [Error]: %v", err)
				}

				// Check values
				if g.ID != tc.resBody.ID {
					t.Errorf("IDs don't match [Expected]: %v [Actual]: %v", tc.resBody.ID, g.ID)
				} else if g.Name != tc.resBody.Name {
					t.Errorf("Names don't match [Expected]: %s [Actual]: %s", tc.resBody.Name, g.Name)
				}
			}
		})
	}
}
