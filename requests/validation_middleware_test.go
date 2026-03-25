package requests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/LittleAksMax/bids-util/validation"
)

type testPayload struct {
	Name string `json:"name" validate:"required"`
	Age  int    `json:"age" validate:"nonnegative"`
}

func TestValidateRequestRejectsInvalidSingleObject(t *testing.T) {
	handler := ValidateRequest[testPayload]([]func(any) error{
		validation.ValidateRequiredFields,
		validation.ValidateNonNegativeFields,
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("next handler should not be called for invalid payload")
	}))

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name":"","age":5}`))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	var response APIResponse
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if response.Error != "name required" {
		t.Fatalf("expected validation error for missing name, got %q", response.Error)
	}
}

func TestValidateRequestRejectsInvalidSliceItem(t *testing.T) {
	handler := ValidateRequest[[]testPayload]([]func(any) error{
		validation.ValidateRequiredFields,
		validation.ValidateNonNegativeFields,
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("next handler should not be called for invalid payload")
	}))

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`[{"name":"ok","age":5},{"name":"","age":3}]`))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	var response APIResponse
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if response.Error != "name required" {
		t.Fatalf("expected validation error for missing name in slice item, got %q", response.Error)
	}
}

func TestValidateRequestStoresValidatedSliceInContext(t *testing.T) {
	expected := []testPayload{
		{Name: "alice", Age: 1},
		{Name: "bob", Age: 2},
	}

	handler := ValidateRequest[[]testPayload]([]func(any) error{
		validation.ValidateRequiredFields,
		validation.ValidateNonNegativeFields,
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body := GetRequestBody[[]testPayload](r)
		if body == nil {
			t.Fatal("expected request body in context")
		}
		if !reflect.DeepEqual(*body, expected) {
			t.Fatalf("expected request body %#v, got %#v", expected, *body)
		}
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`[{"name":"alice","age":1},{"name":"bob","age":2}]`))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, rec.Code)
	}
}
