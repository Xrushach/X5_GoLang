
import "net/http"

func newRouter(store *Store) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", health)
	mux.HandleFunc("POST /orders", createOrder(store))
	return mux
}

func createOrder(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Прочитайте customer из JSON и проверьте имя до вызова store.Create.
		// Как завершить обработку сразу после ошибки?
		http.Error(w, "TODO: create order", http.StatusNotImplemented)
	}
}

