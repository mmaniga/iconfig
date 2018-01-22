package core

import (
	"fmt"
	"net/http"
	"strings"
)

//Start :
func Start() {
	configStore := newStore()
	http.Handle("/join/", configStore)

	go configStore.run()
	http.HandleFunc("/notify/", NotifyHandler(configStore))
	fmt.Println("Starting configuration server on 8091")
	if err := http.ListenAndServe(":8091", nil); err != nil {
		fmt.Println("ListenAndServe:", err)
	}
}

//NotifyHandler : "localhost://8091/notify/{key}/{value}"
func NotifyHandler(h *store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		args := strings.Split(r.URL.Path, "/")
		k := args[2]
		v := args[3]
		msg := &ConfigInfo{Key: k, Value: v}
		h.config <- msg
	}
}

//NotifyHandler : "localhost://8091/kv/{key}"
func GetKey(h *store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		args := strings.Split(r.URL.Path, "/")
		k := args[2]
		v := h.cmap[k]
		msg := &ConfigInfo{Key: k, Value: v}
		h.config <- msg
	}
}
