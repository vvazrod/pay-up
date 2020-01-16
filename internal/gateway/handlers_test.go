package gateway_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/varrrro/pay-up/internal/tmicro/expense"
	"github.com/varrrro/pay-up/internal/tmicro/payment"
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
}

func TestExpensesHandler(t *testing.T) {
	e := expense.Expense{
		ID:         uuid.New(),
		GroupID:    uuid.New(),
		Amount:     25.4,
		Payer:      uuid.New(),
		Recipients: uuid.New().String() + ";" + uuid.New().String(),
	}
	body, _ := json.Marshal(&e)

	cases := []struct {
		method     string
		gid        string
		reqBody    []byte
		statusCode int
	}{
		{"POST", e.GroupID.String(), body, http.StatusAccepted},
		{"POST", e.GroupID.String(), []byte(`{"id":"test"}`), http.StatusBadRequest},
		{"POST", uuid.New().String(), body, http.StatusBadRequest},
		{"POST", e.GroupID.String(), []byte(`{"recipient":"test;"}`), http.StatusBadRequest},

		{"DELETE", e.GroupID.String(), nil, http.StatusAccepted},
		{"DELETE", "test", nil, http.StatusBadRequest},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%s %d", tc.method, tc.statusCode), func(t *testing.T) {
			// Create request
			var req *http.Request
			var err error
			if tc.method == "POST" {
				req, err = http.NewRequest(tc.method, "/groups/"+tc.gid+"/expenses", bytes.NewBuffer(tc.reqBody))
			} else {
				req, err = http.NewRequest(tc.method, "/groups/"+tc.gid+"/expenses", nil)
			}
			if err != nil {
				t.Errorf("Can't create request [Error]: %v", err)
			}

			// Serve test request
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)
			res := rec.Result() // get response
			defer res.Body.Close()

			// Check response status code
			if res.StatusCode != tc.statusCode {
				t.Errorf("Wrong status code [Expected]: %d [Actual]: %d", http.StatusOK, res.StatusCode)
			}
		})
	}
}

func TestPaymentsHandler(t *testing.T) {
	p := payment.Payment{
		ID:        uuid.New(),
		GroupID:   uuid.New(),
		Amount:    14.6,
		Payer:     uuid.New(),
		Recipient: uuid.New(),
	}
	body, _ := json.Marshal(&p)

	cases := []struct {
		method     string
		gid        string
		reqBody    []byte
		statusCode int
	}{
		{"POST", p.GroupID.String(), body, http.StatusAccepted},
		{"POST", p.GroupID.String(), []byte(`{"id":"test"}`), http.StatusBadRequest},
		{"POST", uuid.New().String(), body, http.StatusBadRequest},

		{"DELETE", p.GroupID.String(), nil, http.StatusAccepted},
		{"DELETE", "test", nil, http.StatusBadRequest},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%s %d", tc.method, tc.statusCode), func(t *testing.T) {
			// Create request
			var req *http.Request
			var err error
			if tc.method == "POST" {
				req, err = http.NewRequest(tc.method, "/groups/"+tc.gid+"/payments", bytes.NewBuffer(tc.reqBody))
			} else {
				req, err = http.NewRequest(tc.method, "/groups/"+tc.gid+"/payments", nil)
			}
			if err != nil {
				t.Errorf("Can't create request [Error]: %v", err)
			}

			// Serve test request
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)
			res := rec.Result() // get response
			defer res.Body.Close()

			// Check response status code
			if res.StatusCode != tc.statusCode {
				t.Errorf("Wrong status code [Expected]: %d [Actual]: %d", http.StatusOK, res.StatusCode)
			}
		})
	}
}
