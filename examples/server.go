package main

import (
	"fmt"
	"net/http"
	"os"
)

func handler(w http.ResponseWriter, r *http.Request) {
	hostName, err := os.Hostname()
	if err != nil {
		fmt.Fprintf(w, "Internal server error:%d", http.StatusInternalServerError)
	}
	fmt.Fprintf(w, "Hello from %s 👋\n", hostName)
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
