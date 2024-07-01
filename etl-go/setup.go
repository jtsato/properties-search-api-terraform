package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func setupAttributes() {
	host := os.Getenv("MEILISEARCH_HOST")
	url := fmt.Sprintf("%s/indexes/properties/settings", host)

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

	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("MEILISEARCH_MANAGE_PROPERTIES_TOKEN"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Response Status:", resp.Status)
}
