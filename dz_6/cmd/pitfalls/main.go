// Эти случаи запускаются без открытого порта: go run ./cmd/pitfalls.
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
)

func main() {
	lateStatus()
	secondJSON()
	plainError()
	missingNext()
	muxPatterns()
}

func lateStatus() {
	w := httptest.NewRecorder()
	// Какой статус увидит клиент?
	fmt.Fprint(w, "order not found")
	w.WriteHeader(http.StatusNotFound)
	fmt.Printf("1. Запись тела перед статусом: %d\n", w.Code)
}

func secondJSON() {
	var in struct {
		Customer string `json:"customer"`
	}
	d := json.NewDecoder(strings.NewReader(`{"customer":"Анна"} {"customer":"Борис"}`))
	err := d.Decode(&in)
	// Проверил ли Decoder всё тело запроса?
	fmt.Printf("2. Один Decode: customer=%q, error=%v\n", in.Customer, err)
}

func plainError() {
	w := httptest.NewRecorder()
	http.Error(w, `{"error":"bad input"}`, http.StatusBadRequest)
	fmt.Printf("3. http.Error: Content-Type=%q\n", w.Header().Get("Content-Type"))
}

func missingNext() {
	w := httptest.NewRecorder()
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	h.ServeHTTP(w, httptest.NewRequest("GET", "/orders/1", nil))
	// Запрос зависнет или обработчик уже завершился?
	fmt.Printf("4. Handler без записи: status=%d, body=%q\n", w.Result().StatusCode, w.Body.String())
}

func muxPatterns() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /orders/{id}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, r.PathValue("id"))
	})
	mux.HandleFunc("GET /files/", func(w http.ResponseWriter, r *http.Request) {})
	for _, tc := range []struct{ method, path string }{{"HEAD", "/orders/1"}, {"POST", "/orders/1"}, {"GET", "/files"}, {"GET", "/files/a/b"}} {
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, httptest.NewRequest(tc.method, tc.path, nil))
		fmt.Printf("5. %s %s: status=%d, Allow=%q, Location=%q\n", tc.method, tc.path, w.Code, w.Header().Get("Allow"), w.Header().Get("Location"))
	}
	// Recorder вызывает обработчик напрямую. Отсутствие тела HEAD проверяем через настоящий сервер.
}
