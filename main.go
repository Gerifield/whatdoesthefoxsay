package main

import (
	"fmt"
	"os"
	"whatdoesthefoxsay/foxpost"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <tracking code>")

		return
	}

	trackingCode := os.Args[1]
	fmt.Println("Checking package:", trackingCode)

	fp := foxpost.New()
	status, err := fp.TrackPackage(trackingCode)
	if err != nil {
		fmt.Println("Error tracking package:", err)

		return
	}

	for _, update := range status {
		fmt.Printf("Status: %s\n", update.Title)
		fmt.Printf("Date: %s\n", update.Date)
		fmt.Printf("Description: %s\n", update.Desc)
		fmt.Println("-----------------------------")
	}
}
