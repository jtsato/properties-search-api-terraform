package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func SetupAttributes() {
	// Get the URL from the environment variable
	url := os.Getenv("URL")

	// Define the payload
	payload := map[string]interface{}{
		"filterableAttributes": []string{
			"transactionText",
			"typeText",
			"numberOfBedrooms",
			"rentalTotalPrice",
			"sellingPrice",
			"numberOfGarages",
			"numberOfToilets",
			"builtArea",
			"area",
			"priceByM2",
		},
		"sortableAttributes": []string{
			"rentalTotalPrice",
			"sellingPrice",
			"builtArea",
			"area",
			"priceByM2",
		},
	}

	// Marshal the payload into JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	// Create a new HTTP request
	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set the Content-Type header
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// Print the response status
	fmt.Println("Response Status:", resp.Status)
}
