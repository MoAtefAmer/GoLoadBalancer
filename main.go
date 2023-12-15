package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

var (
	listOfServers = []string{"http://localhost:5201", "http://localhost:5200","http://localhost:5555"}
	currentIndex  = 0
	mutex         = &sync.Mutex{}
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Received request from %s\n", r.RemoteAddr)

		mutex.Lock()
		server := listOfServers[currentIndex]
		fmt.Println("Current Index", currentIndex)
		currentIndex = (currentIndex + 1) % len(listOfServers)
		mutex.Unlock()

		fmt.Printf("Redirecting request to %s\n", server)
		http.Redirect(w, r, server, 301)

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
	err := http.ListenAndServe(":5100", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
