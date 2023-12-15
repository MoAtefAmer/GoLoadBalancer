package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	allServers     = []string{"http://localhost:5201", "http://localhost:5200", "http://localhost:5555"} // Full list of servers
	listOfServers  = []string{} // List of servers that are up
	downedServers  = []string{} // List of servers that are down
	mutex          = &sync.Mutex{}
)

func healthCheck() {
	fmt.Println("Starting health check...")

	for range time.Tick(10 * time.Second) { // Run health check every 10 seconds
		mutex.Lock()
		var tempDownedServers []string

		// Use http.Client with a timeout
		client := http.Client{
			Timeout: 5 * time.Second,
		}

		// Check all servers every time
		for _, server := range allServers {
			_, err := client.Get(server)
			if err != nil {
				fmt.Println("-----------------------------------")
				fmt.Printf("Server %s is down\n", server)
				tempDownedServers = append(tempDownedServers, server)
			} else {
				fmt.Println("-----------------------------------")
				fmt.Printf("Server %s is up\n", server)
			}
		}

		// Update the lists with the current state of the servers
		listOfServers = []string{}
		downedServers = []string{}

		// Decide which servers are up and which are down
		for _, server := range allServers {
			foundDown := false
			for _, downedServer := range tempDownedServers {
				if server == downedServer {
					foundDown = true
					break
				}
			}
			if foundDown {
				downedServers = append(downedServers, server)
			} else {
				listOfServers = append(listOfServers, server)
			}
		}

		mutex.Unlock()
	}
}

func main() {
	// Initialize listOfServers with allServers
	listOfServers = allServers

	// Start the health check in the background
	go healthCheck()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		mutex.Lock() // Lock the mutex before accessing shared variables
		fmt.Fprintln(w, "Hello From Backend Server")
		fmt.Fprintln(w, "Available servers:", listOfServers)
		fmt.Fprintln(w, "Downed servers:", downedServers)
		mutex.Unlock() // Unlock the mutex when done accessing shared variables
	})

	fmt.Println("Backend server is starting...")
	err := http.ListenAndServe(":5100", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	// The "Backend server is shutting down..." log is unnecessary here,
	// as ListenAndServe will block indefinitely until an error occurs.
}