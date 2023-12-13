package main

import (
	"fmt"
	"log"
	"net/http"
)
func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Received request from %s\n", r.RemoteAddr)
		fmt.Printf("%s %s %s\n", r.Method, r.URL.Path, r.Proto)
		for name, headers := range r.Header {
			for _, h := range headers {
				fmt.Printf("%v: %v\n", name, h)
			}
		}
		
		// Respond with a hello message
		fmt.Fprintln(w, "Hello From Backend Server")
		fmt.Println("Replied with a hello message")
	})

	fmt.Println("Backend server is starting...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}