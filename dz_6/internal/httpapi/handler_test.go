package httpapi_test

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"lesson8/internal/httpapi"
	"lesson8/internal/model"
)

// До реализации health этот тест должен быть красным.
func TestHealth(t *testing.T) {
	h := httpapi.NewRouter(model.NewStore())
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d, want 200", rec.Code)
	}
	if rec.Header().Get("Content-Type") != "application/json" {
		t.Fatalf("Content-Type=%q", rec.Header().Get("Content-Type"))
	}
	var body struct {
		Status string `json:"status"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.Status != "ok" {
		t.Fatalf("status in body=%q", body.Status)
	}
}

// Добавьте проверки успешного создания и ошибочного входа.
func TestCreateOrder(t *testing.T) {
	cases := []struct {
		name       string
		body       string
		wantStatus int
		wantID     int
		wantName   string
	}{
		{"correct 1", `{"customer":"Test"}`, 201, 1, "Test"},
		{"empty name", `{"customer":" "}`, 400, 0, ""},
		{"correct 2", `{"customer":"Test2"}`, 201, 1, "Test2"},
		{"broken json", `{"customer":}`, 400, 0, ""},
		{"empty body", ``, 400, 0, ""},
		{"wrong type", `{"customer":42}`, 400, 0, ""},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			store := model.NewStore()
			h := httpapi.NewRouter(store)
			r := httptest.NewRequest(http.MethodPost, "/orders", strings.NewReader(tt.body))
			r.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, r)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status=%d, want %d", rec.Code, tt.wantStatus)
			}
			if tt := rec.Header().Get("Content-Type"); tt != "application/json" {
				t.Fatalf("Content-Type=%q", tt)
			}

			if tt.wantStatus == 201 {
				var o model.Order
				if err := json.Unmarshal(rec.Body.Bytes(), &o); err != nil {
					t.Fatal(err)
				}
				if o.ID != tt.wantID || o.Customer != tt.wantName || o.Status != "new" {
					t.Fatalf("body=%+v", o)
				}
				if loc := rec.Header().Get("Location"); loc != "/orders/"+strconv.Itoa(tt.wantID) {
					t.Fatalf("Location=%q", loc)
				}
			} else {
				var e struct {
					Error string `json:"error"`
				}
				if err := json.Unmarshal(rec.Body.Bytes(), &e); err != nil {
					t.Fatalf("decode error body: %v", err)
				}
				if e.Error == "" {
					t.Fatalf("empty error, body=%s", rec.Body.String())
				}
			}
		})
	}
}

func TestGetOrder(t *testing.T) {
	store := model.NewStore()
	order := store.Create("Bob")
	h := httpapi.NewRouter(store)

	cases := []struct {
		name string
		path string
		want int
	}{
		{"success", "/orders/" + strconv.Itoa(order.ID), 200},
		{"doesnot exist", "/orders/42", 404},
		{"not int", "/orders/qq", 400},
		{"zero", "/orders/0", 400},
		{"negative", "/orders/-1", 400},
		{"ovf", "/orders/9999999999999999999999", 400},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, tt.path, nil))

			if rec.Code != tt.want {
				t.Fatalf("status=%d, want %d", rec.Code, tt.want)
			}
			if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
				t.Fatalf("Content-Type=%q", ct)
			}

			if tt.want == 200 {
				var o model.Order
				_ = json.Unmarshal(rec.Body.Bytes(), &o)
				if o.ID != order.ID || o.Customer != "Bob" || o.Status != "new" {
					t.Fatalf("body=%+v", o)
				}
			}
		})
	}
}

func TestRecovery(t *testing.T) {
	var buf bytes.Buffer
	prev := slog.Default()
	slog.SetDefault(slog.New(slog.NewJSONHandler(&buf, nil)))
	defer slog.SetDefault(prev)

	panicmaster := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("Panic")
	})
	h := httpapi.RequestID(httpapi.Logger(httpapi.Recovery(panicmaster)))
	r := httptest.NewRecorder()
	h.ServeHTTP(r, httptest.NewRequest(http.MethodGet, "/any", nil))

	if r.Code != http.StatusInternalServerError {
		t.Fatalf("status=%d, want 500", r.Code)
	}
	if ct := r.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("Content-Type=%q", ct)
	}
}
