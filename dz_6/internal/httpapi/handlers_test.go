package httpapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"example.com/lesson8/internal/model"
	"example.com/lesson8/internal/repository"
	"example.com/lesson8/internal/service"
)

func newTestAPI() *API {
	return New(service.New(repository.New()), slog.New(slog.NewJSONHandler(io.Discard, nil)))
}

func request(h http.Handler, method, path, body, contentType string) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, path, strings.NewReader(body))
	if contentType != "" {
		r.Header.Set("Content-Type", contentType)
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	return w
}

func TestOrderContract(t *testing.T) {
	for _, router := range []string{"mux", "chi"} {
		t.Run(router, func(t *testing.T) {
			a := newTestAPI()
			var h http.Handler = a.Mux(false)
			if router == "chi" {
				h = a.Chi(false)
			}
			w := request(h, "POST", "/orders", `{"customer":" Анна "}`, "application/json; charset=utf-8")
			if w.Code != 201 || w.Header().Get("Location") != "/orders/1" {
				t.Fatalf("create: %d %s %v", w.Code, w.Body, w.Header())
			}
			var got model.Order
			if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
				t.Fatal(err)
			}
			if got != (model.Order{ID: 1, Customer: "Анна", Status: "new"}) {
				t.Fatalf("order=%+v", got)
			}
			w = request(h, "GET", "/orders/1", "", "")
			if w.Code != 200 || !strings.Contains(w.Body.String(), "Анна") {
				t.Fatalf("read: %d %s", w.Code, w.Body)
			}
			if w.Header().Get("X-Request-ID") == "" {
				t.Fatal("missing request id")
			}
			cases := []struct {
				name, method, path, body, ct string
				status                       int
			}{
				{"health", "GET", "/health", "", "", 200},
				{"missing", "GET", "/orders/999", "", "", 404},
				{"bad id", "GET", "/orders/abc", "", "", 400},
				{"zero id", "GET", "/orders/0", "", "", 400},
				{"negative id", "GET", "/orders/-1", "", "", 400},
				{"overflow id", "GET", "/orders/999999999999999999999999999", "", "", 400},
				{"empty body", "POST", "/orders", "", "application/json", 400},
				{"broken JSON", "POST", "/orders", "{", "application/json", 400},
				{"empty customer", "POST", "/orders", `{"customer":""}`, "application/json", 400},
				{"spaces", "POST", "/orders", `{"customer":" \t "}`, "application/json", 400},
				{"unicode spaces", "POST", "/orders", `{"customer":"\u00a0\u2003"}`, "application/json", 400},
				{"wrong type", "POST", "/orders", `{"customer":12}`, "application/json", 400},
				{"unknown field", "POST", "/orders", `{"customer":"А","status":"paid"}`, "application/json", 400},
				{"two values", "POST", "/orders", `{"customer":"А"} {}`, "application/json", 400},
				{"garbage suffix", "POST", "/orders", `{"customer":"А"} x`, "application/json", 400},
				{"null", "POST", "/orders", `null`, "application/json", 400},
				{"null field", "POST", "/orders", `{"customer":null}`, "application/json", 400},
				{"array", "POST", "/orders", `[]`, "application/json", 400},
				{"wrong media", "POST", "/orders", `{"customer":"А"}`, "text/plain", 415},
				{"missing media", "POST", "/orders", `{"customer":"А"}`, "", 415},
				{"too large", "POST", "/orders", `{"customer":"` + strings.Repeat("x", maxBody) + `"}`, "application/json", 413},
				{"large suffix", "POST", "/orders", `{"customer":"А"}` + strings.Repeat(" ", maxBody), "application/json", 413},
				{"unknown route", "GET", "/missing", "", "", 404},
				{"wrong method", "DELETE", "/orders/1", "", "", 405},
			}
			for _, tc := range cases {
				t.Run(tc.name, func(t *testing.T) {
					w := request(h, tc.method, tc.path, tc.body, tc.ct)
					if w.Code != tc.status {
						t.Fatalf("status=%d want=%d body=%s", w.Code, tc.status, w.Body)
					}
					if tc.status >= 400 && tc.name != "unknown route" && tc.name != "wrong method" {
						var body struct {
							Error string `json:"error"`
						}
						if w.Header().Get("Content-Type") != "application/json" || json.Unmarshal(w.Body.Bytes(), &body) != nil || body.Error == "" {
							t.Fatalf("bad error response: %v %s", w.Header(), w.Body)
						}
					}
					if tc.name == "wrong method" && w.Header().Get("Allow") == "" {
						t.Fatal("missing Allow")
					}
				})
			}
			w = request(h, "POST", "/orders", `{"customer":"Борис"}`+"\n\t", "application/json")
			if w.Code != 201 || w.Header().Get("Location") != "/orders/2" {
				t.Fatalf("invalid input changed store: %d %s", w.Code, w.Body)
			}
		})
	}
}

func TestHEADOverNetwork(t *testing.T) {
	for _, chi := range []bool{false, true} {
		a := newTestAPI()
		h := a.Mux(false)
		if chi {
			h = a.Chi(false)
		}
		s := httptest.NewServer(h)
		resp, err := s.Client().Head(s.URL + "/health")
		if err != nil {
			s.Close()
			t.Fatal(err)
		}
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		s.Close()
		if err != nil || resp.StatusCode != 200 || len(body) != 0 {
			t.Fatalf("HEAD: %v %d %q", err, resp.StatusCode, body)
		}
	}
}

func TestConcurrentCreation(t *testing.T) {
	h := newTestAPI().Mux(false)
	var wg sync.WaitGroup
	ids := make(chan int, 40)
	for i := 0; i < 40; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			w := request(h, "POST", "/orders", `{"customer":"Анна"}`, "application/json")
			var out model.Order
			if w.Code != 201 || json.Unmarshal(w.Body.Bytes(), &out) != nil {
				t.Errorf("create: %d %s", w.Code, w.Body)
				return
			}
			ids <- out.ID
		}()
	}
	wg.Wait()
	close(ids)
	seen := map[int]bool{}
	for id := range ids {
		if seen[id] {
			t.Errorf("duplicate id=%d", id)
		}
		seen[id] = true
	}
	if len(seen) != 40 {
		t.Fatalf("orders=%d", len(seen))
	}
}

type failedOrders struct{}

func (failedOrders) Create(string) (model.Order, error) {
	return model.Order{}, errors.New("secret database details")
}
func (failedOrders) ByID(int) (model.Order, error) {
	return model.Order{}, errors.New("secret database details")
}

func TestInternalErrorIsPrivate(t *testing.T) {
	var logs bytes.Buffer
	a := New(failedOrders{}, slog.New(slog.NewJSONHandler(&logs, nil)))
	w := request(a.Mux(false), "POST", "/orders", `{"customer":"А"}`, "application/json")
	if w.Code != 500 || strings.Contains(w.Body.String(), "secret") || !strings.Contains(logs.String(), "secret database details") {
		t.Fatalf("response=%d %s log=%s", w.Code, w.Body, &logs)
	}
}

func TestPanicHasAccessLog(t *testing.T) {
	for _, chi := range []bool{false, true} {
		var logs bytes.Buffer
		a := New(service.New(repository.New()), slog.New(slog.NewJSONHandler(&logs, nil)))
		h := a.Mux(true)
		if chi {
			h = a.Chi(true)
		}
		w := request(h, "GET", "/demo/panic", "", "")
		if w.Code != 500 {
			t.Fatalf("status=%d", w.Code)
		}
		var access map[string]any
		for _, line := range strings.Split(strings.TrimSpace(logs.String()), "\n") {
			var row map[string]any
			if err := json.Unmarshal([]byte(line), &row); err != nil {
				t.Fatal(err)
			}
			if row["msg"] == "request" {
				access = row
			}
		}
		if access["status"] != float64(500) || access["request_id"] != w.Header().Get("X-Request-ID") {
			t.Fatalf("access log=%v", access)
		}
		if request(h, "GET", "/health", "", "").Code != 200 {
			t.Fatal("health failed after panic")
		}
	}
}

func TestLatePanicAbortsResponse(t *testing.T) {
	var logs bytes.Buffer
	a := New(nil, slog.New(slog.NewJSONHandler(&logs, nil)))
	h := a.logging(a.recovery(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("partial")); panic("late") })))
	w := httptest.NewRecorder()
	func() {
		defer func() {
			if got := recover(); got != http.ErrAbortHandler {
				t.Errorf("panic=%v", got)
			}
		}()
		h.ServeHTTP(w, httptest.NewRequest("GET", "/", nil))
		t.Error("expected abort")
	}()
	if w.Code != 200 || w.Body.String() != "partial" || !strings.Contains(logs.String(), `"aborted":true`) {
		t.Fatalf("%d %s logs=%s", w.Code, w.Body, &logs)
	}
}

func TestStatusWriterKeepsFirstFinalStatus(t *testing.T) {
	w := httptest.NewRecorder()
	rec := &statusWriter{ResponseWriter: w}
	rec.Write([]byte("ok"))
	rec.WriteHeader(500)
	if rec.status != 200 || w.Code != 200 {
		t.Fatalf("recorded=%d actual=%d", rec.status, w.Code)
	}
}
