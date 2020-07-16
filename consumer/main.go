package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/afex/hystrix-go/hystrix"
)

const commandName = "producer_api"

func main() {
	hystrix.ConfigureCommand(commandName, hystrix.CommandConfig{
		Timeout:                500,  // Time in milliseconds after which the caller will observe a timeout and walk away from the command execution, and performs fallback logic
		MaxConcurrentRequests:  100,  // maximum number of requests allowed into hystrix
		RequestVolumeThreshold: 3,    // Minimum number of requests in a rolling window that will trip the circuit.
		SleepWindow:            1000, // Amount of time, after tripping the circuit, to reject requests before allowing attempts again to determine if the circuit should again be closed
		ErrorPercentThreshold:  50,   // Error percentage at or above which the circuit should trip open and start short-circuiting requests to fallback logic
	})

	http.HandleFunc("/", logger(handle))
	log.Println("Consumer listening on :8080")
	http.ListenAndServe(":8080", nil)

}

func handle(w http.ResponseWriter, r *http.Request) {
	resultChan := make(chan string, 1)
	errChan := hystrix.Go(commandName, func() error {
		// talk to other services

		resp, err := http.Get("http://localhost:8081")
		if err != nil {
			return err
		}
		defer r.Body.Close()

		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		resultChan <- string(b)
		return nil
	}, nil)

	select {
	case result := <-resultChan:
		fmt.Println("success:", result)
		log.Printf("got response from service %v", result)
		w.WriteHeader(http.StatusOK)
	case err := <-errChan:
		fmt.Println("failure: ", err.Error())
		log.Printf("failure: %s", err)
		w.WriteHeader(http.StatusServiceUnavailable)
	}
}

// logger is Handler wrapper function for logging
func logger(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.Path, r.Method)
		fn(w, r)
	}
}
