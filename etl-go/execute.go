package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/meilisearch/meilisearch-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func execute() {

	log.Println("--------------------------------------------------------------")
	log.Println("Starting...")

	start := time.Now()
	count := 0

	log.Printf("The log level defined to %s", os.Getenv("LOG_LEVEL"))

	var client = connectToDatabase()
	defer disconnect()

	var collection *mongo.Collection
	var documents []bson.M

	db := client.Database(os.Getenv("MONGODB_DATABASE"))
	collection = db.Collection("properties")

	cursor, err := collection.Find(context.TODO(), bson.M{"tenantId": 1})
	if err != nil {
		log.Fatalf("Error finding documents: %v", err)
	}
	defer cursor.Close(context.TODO())

	if err = cursor.All(context.TODO(), &documents); err != nil {
		log.Fatalf("Error decoding documents: %v", err)
	}

	var properties []interface{}
	for _, document := range documents {
		count++
		properties = append(properties, convertProperty(document))
	}

	log.Printf("The system found %d properties", count)

	if os.Getenv("MEILISEARCH_HOST") == "" || os.Getenv("MEILISEARCH_MASTER_KEY") == "" {
		log.Fatal("Missing environment variables for MeiliSearch")
	}

	meiliClient := meilisearch.NewClient(meilisearch.ClientConfig{
		Host:   os.Getenv("MEILISEARCH_HOST"),
		APIKey: os.Getenv("MEILISEARCH_MASTER_KEY"),
	})

	task, err := meiliClient.DeleteIndex("properties")
	if err != nil {
		log.Fatalf("Error deleting documents from MeiliSearch: %v", err)
	}

	err = waitForTaskCompletion(meiliClient, task.TaskUID, 30*time.Second)
	if err != nil {
		log.Fatalf("Error waiting for deleting documents: %v", err)
	}

	_, err = meiliClient.GetIndex("properties")
	if err != nil {
		task, err = meiliClient.CreateIndex(&meilisearch.IndexConfig{
			Uid:        "properties",
			PrimaryKey: "uuid",
		})
		if err != nil {
			log.Fatalf("Error creating index in MeiliSearch: %v", err)
		}

		err = waitForTaskCompletion(meiliClient, task.TaskUID, 30*time.Second)
		if err != nil {
			log.Fatalf("Error waiting for creating index: %v", err)
		}
	}

	task, err = meiliClient.Index("properties").AddDocuments(properties, "uuid")
	if err != nil {
		log.Fatalf("Error adding documents to MeiliSearch: %v", err)
	}

	err = waitForTaskCompletion(meiliClient, task.TaskUID, 30*time.Second)
	if err != nil {
		log.Fatalf("Error waiting for adding documents: %v", err)
	}

	duration := time.Since(start)
	log.Printf("The system processed %d properties in %.2f seconds", count, duration.Seconds())

	log.Println("Finished")
}

func waitForTaskCompletion(meiliClient *meilisearch.Client, taskID int64, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	for {
		task, err := meiliClient.GetTask(taskID)
		if err != nil {
			return fmt.Errorf("error getting task status: %v", err)
		}

		if task.Status == "succeeded" {
			return nil
		}

		if task.Status == "failed" {
			return fmt.Errorf("task failed: %v", task)
		}

		select {
		case <-time.After(500 * time.Millisecond):
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for task completion")
		}
	}
}
