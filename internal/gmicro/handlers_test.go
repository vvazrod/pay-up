package gmicro_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/varrrro/pay-up/internal/gmicro/group"
	"github.com/varrrro/pay-up/internal/gmicro/member"
)

var router *mux.Router

func TestGetStatus(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fail()
	}

	router.ServeHTTP(rec, req)
	r := rec.Result()

	if r.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Wrong mime-type [Expected]: application/json [Actual]: %s", r.Header.Get("Content-Type"))
	}

	if r.StatusCode != http.StatusOK {
		t.Errorf("Wrong status code [Expected]: %d [Actual]: %d", http.StatusOK, r.StatusCode)
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Fail()
	}

	var data map[string]interface{}
	json.Unmarshal(body, &data)

	if _, ok := data["status"]; !ok {
		t.Error("Didn't return status")
	}
}

func TestPostGroup(t *testing.T) {
	rec := httptest.NewRecorder()
	reqBody := []byte(`{"name":"test"}`)
	req, err := http.NewRequest("POST", "/groups", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fail()
	}

	router.ServeHTTP(rec, req)
	r := rec.Result()

	if r.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Wrong mime-type [Expected]: application/json [Actual]: %s", r.Header.Get("Content-Type"))
	}

	if r.StatusCode != http.StatusCreated {
		t.Errorf("Wrong status code [Expected]: %d [Actual]: %d", http.StatusOK, r.StatusCode)
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Fail()
	}

	var g group.Group
	err = json.Unmarshal(body, &g)
	if err != nil {
		t.Error("Couldn't parse response as group.")
	}

	if g.Name != "test" {
		t.Error("Returned group name isn't correct.")
	}

	_, err = uuid.Parse(g.ID)
	if err != nil {
		t.Error("Returned ID isn't correct UUID.")
	}
}

func TestPostGroupBadRequest(t *testing.T) {
	rec := httptest.NewRecorder()
	reqBody := []byte(`{}`)
	req, err := http.NewRequest("POST", "/groups", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fail()
	}

	router.ServeHTTP(rec, req)
	r := rec.Result()

	if r.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Wrong mime-type [Expected]: application/json [Actual]: %s", r.Header.Get("Content-Type"))
	}

	if r.StatusCode != http.StatusBadRequest {
		t.Errorf("Wrong status code [Expected]: %d [Actual]: %d", http.StatusOK, r.StatusCode)
	}
}

func TestGetGroup(t *testing.T) {
	g, err := manager.CreateGroup("test")
	if err != nil {
		t.Fail()
	}

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/groups/"+g.ID, nil)
	if err != nil {
		t.Fail()
	}

	router.ServeHTTP(rec, req)
	r := rec.Result()

	if r.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Wrong mime-type [Expected]: application/json [Actual]: %s", r.Header.Get("Content-Type"))
	}

	if r.StatusCode != http.StatusOK {
		t.Errorf("Wrong status code [Expected]: %d [Actual]: %d", http.StatusOK, r.StatusCode)
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Fail()
	}

	var g2 group.Group
	err = json.Unmarshal(body, &g2)
	if err != nil {
		t.Error("Couldn't parse response as group.")
	}

	if g.ID != g2.ID || g.Name != g2.Name {
		t.Error("Returned group isn't correct.")
	}
}

func TestGetGroupNotFound(t *testing.T) {
	testID := uuid.New()

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/groups/"+testID.String(), nil)
	if err != nil {
		t.Fail()
	}

	router.ServeHTTP(rec, req)
	r := rec.Result()

	if r.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Wrong mime-type [Expected]: application/json [Actual]: %s", r.Header.Get("Content-Type"))
	}

	if r.StatusCode != http.StatusNotFound {
		t.Errorf("Wrong status code [Expected]: %d [Actual]: %d", http.StatusNotFound, r.StatusCode)
	}
}

func TestGetGroupBadRequest(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/groups/test", nil)
	if err != nil {
		t.Fail()
	}

	router.ServeHTTP(rec, req)
	r := rec.Result()

	if r.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Wrong mime-type [Expected]: application/json [Actual]: %s", r.Header.Get("Content-Type"))
	}

	if r.StatusCode != http.StatusBadRequest {
		t.Errorf("Wrong status code [Expected]: %d [Actual]: %d", http.StatusNotFound, r.StatusCode)
	}
}

func TestPutGroup(t *testing.T) {
	g, err := manager.CreateGroup("test")
	if err != nil {
		t.Fail()
	}

	rec := httptest.NewRecorder()
	reqBody := []byte(`{"name":"updated"}`)
	req, err := http.NewRequest("PUT", "/groups/"+g.ID, bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fail()
	}

	router.ServeHTTP(rec, req)
	r := rec.Result()

	if r.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Wrong mime-type [Expected]: application/json [Actual]: %s", r.Header.Get("Content-Type"))
	}

	if r.StatusCode != http.StatusOK {
		t.Errorf("Wrong status code [Expected]: %d [Actual]: %d", http.StatusOK, r.StatusCode)
	}

	if g, err = manager.FetchGroup(uuid.MustParse(g.ID)); err != nil {
		t.Fail()
	} else if g.Name != "updated" {
		t.Error("Group's name wasn't updated correctly.")
	}
}

func TestPutGroupNotFound(t *testing.T) {
	testID := uuid.New()

	rec := httptest.NewRecorder()
	reqBody := []byte(`{"name":"updated"}`)
	req, err := http.NewRequest("PUT", "/groups/"+testID.String(), bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fail()
	}

	router.ServeHTTP(rec, req)
	r := rec.Result()

	if r.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Wrong mime-type [Expected]: application/json [Actual]: %s", r.Header.Get("Content-Type"))
	}

	if r.StatusCode != http.StatusNotFound {
		t.Errorf("Wrong status code [Expected]: %d [Actual]: %d", http.StatusOK, r.StatusCode)
	}
}

func TestPutGroupBadRequestID(t *testing.T) {
	rec := httptest.NewRecorder()
	reqBody := []byte(`{"name":"updated"}`)
	req, err := http.NewRequest("PUT", "/groups/test", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fail()
	}

	router.ServeHTTP(rec, req)
	r := rec.Result()

	if r.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Wrong mime-type [Expected]: application/json [Actual]: %s", r.Header.Get("Content-Type"))
	}

	if r.StatusCode != http.StatusBadRequest {
		t.Errorf("Wrong status code [Expected]: %d [Actual]: %d", http.StatusOK, r.StatusCode)
	}
}

func TestPutGroupBadRequestBody(t *testing.T) {
	testID := uuid.New()

	rec := httptest.NewRecorder()
	reqBody := []byte(`{}`)
	req, err := http.NewRequest("PUT", "/groups/"+testID.String(), bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fail()
	}

	router.ServeHTTP(rec, req)
	r := rec.Result()

	if r.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Wrong mime-type [Expected]: application/json [Actual]: %s", r.Header.Get("Content-Type"))
	}

	if r.StatusCode != http.StatusBadRequest {
		t.Errorf("Wrong status code [Expected]: %d [Actual]: %d", http.StatusOK, r.StatusCode)
	}
}

func TestDeleteGroup(t *testing.T) {
	g, err := manager.CreateGroup("test")
	if err != nil {
		t.Fail()
	}

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("DELETE", "/groups/"+g.ID, nil)
	if err != nil {
		t.Fail()
	}

	router.ServeHTTP(rec, req)
	r := rec.Result()

	if r.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Wrong mime-type [Expected]: application/json [Actual]: %s", r.Header.Get("Content-Type"))
	}

	if r.StatusCode != http.StatusNoContent {
		t.Errorf("Wrong status code [Expected]: %d [Actual]: %d", http.StatusOK, r.StatusCode)
	}

	if _, err = manager.FetchGroup(uuid.MustParse(g.ID)); err == nil {
		t.Error("Group wasn't deleted correctly.")
	}
}

func TestDeleteGroupNotFound(t *testing.T) {
	testID := uuid.New()

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("DELETE", "/groups/"+testID.String(), nil)
	if err != nil {
		t.Fail()
	}

	router.ServeHTTP(rec, req)
	r := rec.Result()

	if r.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Wrong mime-type [Expected]: application/json [Actual]: %s", r.Header.Get("Content-Type"))
	}

	if r.StatusCode != http.StatusNotFound {
		t.Errorf("Wrong status code [Expected]: %d [Actual]: %d", http.StatusOK, r.StatusCode)
	}
}

func TestDeleteGroupBadRequest(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("DELETE", "/groups/test", nil)
	if err != nil {
		t.Fail()
	}

	router.ServeHTTP(rec, req)
	r := rec.Result()

	if r.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Wrong mime-type [Expected]: application/json [Actual]: %s", r.Header.Get("Content-Type"))
	}

	if r.StatusCode != http.StatusBadRequest {
		t.Errorf("Wrong status code [Expected]: %d [Actual]: %d", http.StatusOK, r.StatusCode)
	}
}

func TestPostMember(t *testing.T) {
	g, err := manager.CreateGroup("test")
	if err != nil {
		t.Fail()
	}

	rec := httptest.NewRecorder()
	reqBody := []byte(`{"name":"test"}`)
	req, err := http.NewRequest("POST", "/groups/"+g.ID, bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fail()
	}

	router.ServeHTTP(rec, req)
	r := rec.Result()

	if r.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Wrong mime-type [Expected]: application/json [Actual]: %s", r.Header.Get("Content-Type"))
	}

	if r.StatusCode != http.StatusOK {
		t.Errorf("Wrong status code [Expected]: %d [Actual]: %d", http.StatusOK, r.StatusCode)
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Fail()
	}

	var m member.Member
	err = json.Unmarshal(body, &m)
	if err != nil {
		t.Error("Couldn't parse response as member.")
	}

	if m.Name != "test" || m.GroupID != g.ID {
		t.Error("Returned member isn't correct.")
	}

	_, err = uuid.Parse(m.ID)
	if err != nil {
		t.Error("Returned ID isn't correct UUID.")
	}
}

func TestPostMemberNotFound(t *testing.T) {
	testID := uuid.New()

	rec := httptest.NewRecorder()
	reqBody := []byte(`{"name":"test"}`)
	req, err := http.NewRequest("POST", "/groups/"+testID.String(), bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fail()
	}

	router.ServeHTTP(rec, req)
	r := rec.Result()

	if r.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Wrong mime-type [Expected]: application/json [Actual]: %s", r.Header.Get("Content-Type"))
	}

	if r.StatusCode != http.StatusNotFound {
		t.Errorf("Wrong status code [Expected]: %d [Actual]: %d", http.StatusOK, r.StatusCode)
	}
}

func TestPostMemberConflict(t *testing.T) {
	g, err := manager.CreateGroup("test")
	if err != nil {
		t.Fail()
	}

	_, err = manager.AddMember(uuid.MustParse(g.ID), "test")
	if err != nil {
		t.Fail()
	}

	rec := httptest.NewRecorder()
	reqBody := []byte(`{"name":"test"}`)
	req, err := http.NewRequest("POST", "/groups/"+g.ID, bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fail()
	}

	router.ServeHTTP(rec, req)
	r := rec.Result()

	if r.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Wrong mime-type [Expected]: application/json [Actual]: %s", r.Header.Get("Content-Type"))
	}

	if r.StatusCode != http.StatusConflict {
		t.Errorf("Wrong status code [Expected]: %d [Actual]: %d", http.StatusOK, r.StatusCode)
	}
}

func TestPostMemberBadRequestID(t *testing.T) {
	rec := httptest.NewRecorder()
	reqBody := []byte(`{"name":"test"}`)
	req, err := http.NewRequest("POST", "/groups/test", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fail()
	}

	router.ServeHTTP(rec, req)
	r := rec.Result()

	if r.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Wrong mime-type [Expected]: application/json [Actual]: %s", r.Header.Get("Content-Type"))
	}

	if r.StatusCode != http.StatusBadRequest {
		t.Errorf("Wrong status code [Expected]: %d [Actual]: %d", http.StatusOK, r.StatusCode)
	}
}

func TestPostMemberBadRequestBody(t *testing.T) {
	g, err := manager.CreateGroup("test")
	if err != nil {
		t.Fail()
	}

	rec := httptest.NewRecorder()
	reqBody := []byte(`{}`)
	req, err := http.NewRequest("POST", "/groups/"+g.ID, bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fail()
	}

	router.ServeHTTP(rec, req)
	r := rec.Result()

	if r.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Wrong mime-type [Expected]: application/json [Actual]: %s", r.Header.Get("Content-Type"))
	}

	if r.StatusCode != http.StatusBadRequest {
		t.Errorf("Wrong status code [Expected]: %d [Actual]: %d", http.StatusOK, r.StatusCode)
	}
}

func TestGetMember(t *testing.T) {
	g, err := manager.CreateGroup("test")
	if err != nil {
		t.Fail()
	}

	m, err := manager.AddMember(uuid.MustParse(g.ID), "test")
	if err != nil {
		t.Fail()
	}

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/groups/"+g.ID+"/members/"+m.ID, nil)
	if err != nil {
		t.Fail()
	}

	router.ServeHTTP(rec, req)
	r := rec.Result()

	if r.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Wrong mime-type [Expected]: application/json [Actual]: %s", r.Header.Get("Content-Type"))
	}

	if r.StatusCode != http.StatusOK {
		t.Errorf("Wrong status code [Expected]: %d [Actual]: %d", http.StatusOK, r.StatusCode)
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Fail()
	}

	var m2 member.Member
	err = json.Unmarshal(body, &m2)
	if err != nil {
		t.Error("Couldn't parse response as member.")
	}

	if m.Name != m2.Name {
		t.Error("Returned member isn't correct.")
	}
}

func TestGetMemberNotFound(t *testing.T) {
	testID := uuid.New()

	g, err := manager.CreateGroup("test")
	if err != nil {
		t.Fail()
	}

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/groups/"+g.ID+"/members/"+testID.String(), nil)
	if err != nil {
		t.Fail()
	}

	router.ServeHTTP(rec, req)
	r := rec.Result()

	if r.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Wrong mime-type [Expected]: application/json [Actual]: %s", r.Header.Get("Content-Type"))
	}

	if r.StatusCode != http.StatusNotFound {
		t.Errorf("Wrong status code [Expected]: %d [Actual]: %d", http.StatusOK, r.StatusCode)
	}
}

func TestGetMemberBadRequestGroup(t *testing.T) {
	testID := uuid.New()

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/groups/test/members/"+testID.String(), nil)
	if err != nil {
		t.Fail()
	}

	router.ServeHTTP(rec, req)
	r := rec.Result()

	if r.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Wrong mime-type [Expected]: application/json [Actual]: %s", r.Header.Get("Content-Type"))
	}

	if r.StatusCode != http.StatusBadRequest {
		t.Errorf("Wrong status code [Expected]: %d [Actual]: %d", http.StatusOK, r.StatusCode)
	}
}

func TestGetMemberBadRequestMember(t *testing.T) {
	testID := uuid.New()

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/groups/"+testID.String()+"/members/test", nil)
	if err != nil {
		t.Fail()
	}

	router.ServeHTTP(rec, req)
	r := rec.Result()

	if r.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Wrong mime-type [Expected]: application/json [Actual]: %s", r.Header.Get("Content-Type"))
	}

	if r.StatusCode != http.StatusBadRequest {
		t.Errorf("Wrong status code [Expected]: %d [Actual]: %d", http.StatusOK, r.StatusCode)
	}
}

func TestPutMember(t *testing.T) {
	g, err := manager.CreateGroup("test")
	if err != nil {
		t.Fail()
	}

	m, err := manager.AddMember(uuid.MustParse(g.ID), "test")
	if err != nil {
		t.Fail()
	}

	rec := httptest.NewRecorder()
	reqBody := []byte(`{"name":"updated"}`)
	req, err := http.NewRequest("PUT", "/groups/"+g.ID+"/members/"+m.ID, bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fail()
	}

	router.ServeHTTP(rec, req)
	r := rec.Result()

	if r.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Wrong mime-type [Expected]: application/json [Actual]: %s", r.Header.Get("Content-Type"))
	}

	if r.StatusCode != http.StatusOK {
		t.Errorf("Wrong status code [Expected]: %d [Actual]: %d", http.StatusOK, r.StatusCode)
	}

	if m, err = manager.FetchMember(uuid.MustParse(g.ID), uuid.MustParse(m.ID)); err != nil {
		t.Fail()
	} else if m.Name != "updated" {
		t.Error("Member's name wasn't updated correctly.")
	}
}

func TestPutMemberNotFound(t *testing.T) {
	testID := uuid.New()

	g, err := manager.CreateGroup("test")
	if err != nil {
		t.Fail()
	}

	rec := httptest.NewRecorder()
	reqBody := []byte(`{"name":"updated"}`)
	req, err := http.NewRequest("PUT", "/groups/"+g.ID+"/members/"+testID.String(), bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fail()
	}

	router.ServeHTTP(rec, req)
	r := rec.Result()

	if r.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Wrong mime-type [Expected]: application/json [Actual]: %s", r.Header.Get("Content-Type"))
	}

	if r.StatusCode != http.StatusNotFound {
		t.Errorf("Wrong status code [Expected]: %d [Actual]: %d", http.StatusOK, r.StatusCode)
	}
}

func TestPutMemberBadRequestGroupID(t *testing.T) {
	testID := uuid.New()

	rec := httptest.NewRecorder()
	reqBody := []byte(`{"name":"updated"}`)
	req, err := http.NewRequest("PUT", "/groups/test/members/"+testID.String(), bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fail()
	}

	router.ServeHTTP(rec, req)
	r := rec.Result()

	if r.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Wrong mime-type [Expected]: application/json [Actual]: %s", r.Header.Get("Content-Type"))
	}

	if r.StatusCode != http.StatusBadRequest {
		t.Errorf("Wrong status code [Expected]: %d [Actual]: %d", http.StatusOK, r.StatusCode)
	}
}

func TestPutMemberBadRequestMemberID(t *testing.T) {
	testID := uuid.New()

	rec := httptest.NewRecorder()
	reqBody := []byte(`{"name":"updated"}`)
	req, err := http.NewRequest("PUT", "/groups/"+testID.String()+"/members/test", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fail()
	}

	router.ServeHTTP(rec, req)
	r := rec.Result()

	if r.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Wrong mime-type [Expected]: application/json [Actual]: %s", r.Header.Get("Content-Type"))
	}

	if r.StatusCode != http.StatusBadRequest {
		t.Errorf("Wrong status code [Expected]: %d [Actual]: %d", http.StatusOK, r.StatusCode)
	}
}

func TestPutMemberBadRequestBody(t *testing.T) {
	g, err := manager.CreateGroup("test")
	if err != nil {
		t.Fail()
	}

	m, err := manager.AddMember(uuid.MustParse(g.ID), "test")
	if err != nil {
		t.Fail()
	}

	rec := httptest.NewRecorder()
	reqBody := []byte(`{}`)
	req, err := http.NewRequest("PUT", "/groups/"+g.ID+"/members/"+m.ID, bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fail()
	}

	router.ServeHTTP(rec, req)
	r := rec.Result()

	if r.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Wrong mime-type [Expected]: application/json [Actual]: %s", r.Header.Get("Content-Type"))
	}

	if r.StatusCode != http.StatusBadRequest {
		t.Errorf("Wrong status code [Expected]: %d [Actual]: %d", http.StatusOK, r.StatusCode)
	}
}

func TestDeleteMember(t *testing.T) {
	g, err := manager.CreateGroup("test")
	if err != nil {
		t.Fail()
	}

	m, err := manager.AddMember(uuid.MustParse(g.ID), "test")
	if err != nil {
		t.Fail()
	}

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("DELETE", "/groups/"+g.ID+"/members/"+m.ID, nil)
	if err != nil {
		t.Fail()
	}

	router.ServeHTTP(rec, req)
	r := rec.Result()

	if r.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Wrong mime-type [Expected]: application/json [Actual]: %s", r.Header.Get("Content-Type"))
	}

	if r.StatusCode != http.StatusNoContent {
		t.Errorf("Wrong status code [Expected]: %d [Actual]: %d", http.StatusOK, r.StatusCode)
	}

	if _, err = manager.FetchMember(uuid.MustParse(g.ID), uuid.MustParse(m.ID)); err == nil {
		t.Error("Member wasn't deleted correctly.")
	}
}

func TestDeleteMemberNotFound(t *testing.T) {
	testID := uuid.New()

	g, err := manager.CreateGroup("test")
	if err != nil {
		t.Fail()
	}

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("DELETE", "/groups/"+g.ID+"/members/"+testID.String(), nil)
	if err != nil {
		t.Fail()
	}

	router.ServeHTTP(rec, req)
	r := rec.Result()

	if r.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Wrong mime-type [Expected]: application/json [Actual]: %s", r.Header.Get("Content-Type"))
	}

	if r.StatusCode != http.StatusNotFound {
		t.Errorf("Wrong status code [Expected]: %d [Actual]: %d", http.StatusOK, r.StatusCode)
	}
}

func TestDeleteMemberBadRequestGroupID(t *testing.T) {
	testID := uuid.New()

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("DELETE", "/groups/test/members/"+testID.String(), nil)
	if err != nil {
		t.Fail()
	}

	router.ServeHTTP(rec, req)
	r := rec.Result()

	if r.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Wrong mime-type [Expected]: application/json [Actual]: %s", r.Header.Get("Content-Type"))
	}

	if r.StatusCode != http.StatusBadRequest {
		t.Errorf("Wrong status code [Expected]: %d [Actual]: %d", http.StatusOK, r.StatusCode)
	}
}

func TestDeleteMemberBadRequestMemberID(t *testing.T) {
	testID := uuid.New()

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("DELETE", "/groups/"+testID.String()+"/members/test", nil)
	if err != nil {
		t.Fail()
	}

	router.ServeHTTP(rec, req)
	r := rec.Result()

	if r.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Wrong mime-type [Expected]: application/json [Actual]: %s", r.Header.Get("Content-Type"))
	}

	if r.StatusCode != http.StatusBadRequest {
		t.Errorf("Wrong status code [Expected]: %d [Actual]: %d", http.StatusOK, r.StatusCode)
	}
}
