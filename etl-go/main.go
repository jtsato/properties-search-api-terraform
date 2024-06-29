package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/meilisearch/meilisearch-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var collection *mongo.Collection
var documents []bson.M

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	execute()
}

func execute() {
	log.Println("--------------------------------------------------------------")
	log.Println("Starting...")

	start := time.Now()
	count := 0

	log.Printf("The log level defined to %s", os.Getenv("LOG_LEVEL"))

	connectToDatabase()
	defer disconnect()

	database := client.Database(os.Getenv("MONGODB_DATABASE"))
	collection = database.Collection("properties")

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

	meiliClient := meilisearch.NewClient(meilisearch.ClientConfig{
		Host:   os.Getenv("MEILISEARCH_HOST"),
		APIKey: os.Getenv("MEILISEARCH_MANAGE_PROPERTIES_TOKEN"),
	})

	_, err = meiliClient.Index("properties").AddDocuments(properties, "uuid")
	if err != nil {
		log.Fatalf("Error adding documents to MeiliSearch: %v", err)
	}

	_, err = meiliClient.CreateIndex(&meilisearch.IndexConfig{
		Uid:        "properties",
		PrimaryKey: "uuid",
	})
	if err != nil {
		log.Fatalf("Error creating index in MeiliSearch: %v", err)
	}

	duration := time.Since(start)
	log.Printf("The system processed %d properties in %.2f seconds", count, duration.Seconds())

	log.Println("Finished")
}

func convertProperty(document bson.M) map[string]interface{} {
	return map[string]interface{}{
		"transactionText":       getTransactionText(document["transaction"].(string)),
		"typeText":              getTypeText(document["type"].(string)),
		"transactionTerms":      getTransactionTerms(document["transaction"].(string)),
		"district":              document["district"],
		"city":                  document["city"],
		"state":                 document["state"],
		"address":               document["address"],
		"numberOfBedroomsTerms": getNumberOfBedroomsTerms(int(document["numberOfBedrooms"].(int32))),
		"tenantName":            getTenantName(document["url"].(string)),
		"refId":                 document["refId"],
		"title":                 strings.ToLower(document["title"].(string)),
		"description":           strings.ToLower(document["description"].(string)),
		"typeTerms":             getTypeTerms(document["type"].(string)),
		"numberOfGaragesTerms":  getNumberOfGaragesTerms(int(document["numberOfGarages"].(int32))),
		"numberOfToiletsTerms":  getNumberOfToiletsTerms(int(document["numberOfToilets"].(int32))),
		"numberOfBedrooms":      document["numberOfBedrooms"],
		"rentalTotalPrice":      document["rentalTotalPrice"],
		"sellingPrice":          document["sellingPrice"],
		"numberOfGarages":       document["numberOfGarages"],
		"numberOfToilets":       document["numberOfToilets"],
		"builtArea":             document["builtArea"],
		"area":                  document["area"],
		"priceByM2":             document["priceByM2"],
		"ranking":               document["ranking"],
		"status":                document["status"],
		"transaction":           document["transaction"],
		"type":                  document["type"],
		"coverImage":            document["images"].([]interface{})[0],
		"uuid":                  document["uuid"],
	}
}

func getTenantName(url string) string {
	return url[8:strings.Index(url, ".com.br")]
}

func getTransactionText(transaction string) string {
	if transaction == "RENT" {
		return "Aluguel"
	}
	return "Venda"
}

func getTransactionTerms(transaction string) string {
	if transaction == "RENT" {
		return "📝, aluguel, alugar, locação, locar"
	}
	return "💲, venda, vender, compra, comprar"
}

func getTypeText(typeProp string) string {
	switch typeProp {
	case "APARTMENT":
		return "Apartamento"
	case "WAREHOUSE":
		return "Barracão"
	case "HOUSE":
		return "Casa"
	case "COUNTRY_HOUSE":
		return "Chácara"
	case "FARM":
		return "Fazenda"
	case "GARAGE":
		return "Garagem"
	case "LAND_DIVISION":
		return "Loteamento"
	case "BUSINESS_PREMISES":
		return "Ponto Comercial"
	case "OFFICE":
		return "Sala Comercial"
	case "TWO_STOREY_HOUSE":
		return "Sobrado"
	case "LAND":
		return "Terreno"
	default:
		return "Outro"
	}
}

func getTypeTerms(typeProp string) string {
	switch typeProp {
	case "TWO_STOREY_HOUSE":
		return "🏘️, sobrado, andares"
	case "APARTMENT":
		return "🏢, 🏬, apartamento, apartamentos, ap, ape, apt, apzinho, apezinho, apart, apto, flatinho, flat, kitnet, loft, quitinete, studio"
	case "HOUSE":
		return "🏠, 🏚️, casa, casinha, chalé, edícula, kaza, kza, mansão, vivenda"
	case "LAND":
		return "🏞️, 🌄, terreno, lote, terrenos, lotes"
	case "COUNTRY_HOUSE":
		return "🌳, 🏡, chácara, campo, chacarazinha, chacarazito, chacarinha, chacrinha, rural, sítio, sítiozinho, sítiozito, fazendinha"
	case "FARM":
		return "🚜, 🌾, 🐄, fazenda, sítio"
	case "GARAGE":
		return "🚗, 🚘, 🅿️, garagem, estacionamento, garage, vaga, carro"
	case "WAREHOUSE":
		return "🏭, 📦, barracão, armazém, armazem, galpão, galpao, depósito"
	case "OFFICE":
		return "🖥️, 🏛️, sala, sala comercial, sala_comercial, escritório, escritorio"
	case "BUSINESS_PREMISES":
		return "🏪, 🛍️, ponto, loja, comércio"
	case "LAND_DIVISION":
		return "🏞️, 🌄, loteamento, lote"
	default:
		return "❓, ❔, outro, outros"
	}
}

func getNumberOfBedroomsTerms(number int) string {
	numberAsString := convertNumberToPortugueseWords(number)
	return fmt.Sprintf("%d quartos, %d dormitórios, %s quartos, %s dormitórios", number, number, numberAsString, numberAsString)
}

func getNumberOfGaragesTerms(number int) string {
	numberAsString := convertNumberToPortugueseWords(number)
	return fmt.Sprintf("%d garagens, %d vagas, %d carros, %s garagens, %s vagas, %s carros", number, number, number, numberAsString, numberAsString, numberAsString)
}

func getNumberOfToiletsTerms(number int) string {
	numberAsString := convertNumberToPortugueseWords(number)
	return fmt.Sprintf("%d banheiros, %d toalete, %s banheiros, %s toalete", number, number, numberAsString, numberAsString)
}

func convertNumberToPortugueseWords(number int) string {
	words := []string{
		"sem", "um", "dois", "três", "quatro", "cinco", "seis", "sete", "oito", "nove", "dez",
		"onze", "doze", "treze", "quatorze", "quinze", "dezesseis", "dezessete", "dezoito", "dezenove", "vinte",
	}

	if number >= 0 && number <= 20 {
		return words[number]
	}

	return fmt.Sprintf("%d", number)
}

func connectToDatabase() {
	mongoUser := os.Getenv("MONGODB_USER")
	mongoPassword := os.Getenv("MONGODB_PASSWORD")
	mongoClusterUrl := os.Getenv("MONGODB_URL")
	mongoDatabase := os.Getenv("MONGODB_DATABASE")

	if mongoUser == "" || mongoPassword == "" || mongoClusterUrl == "" || mongoDatabase == "" {
		log.Fatalf("Missing environment variables for MongoDB")
	}

	uri := fmt.Sprintf("mongodb+srv://%s:%s@%s/%s?retryWrites=true&w=majority", mongoUser, mongoPassword, mongoClusterUrl, mongoDatabase)

	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalf("Error pinging MongoDB: %v", err)
	}

	log.Println("Connected to MongoDB")
}

func disconnect() {
	if client != nil {
		if err := client.Disconnect(context.TODO()); err != nil {
			log.Fatalf("Error disconnecting from MongoDB: %v", err)
		}
	}

	log.Println("Disconnected from MongoDB")
}

func info(message string) {
	if os.Getenv("LOG_LEVEL") == "DEBUG" || os.Getenv("LOG_LEVEL") == "INFO" {
		log.Printf("Hestia %s: %s", formatDate(time.Now()), message)
	}
}

func error(message string) {
	log.Printf("Hestia %s", message)
}

func padTo2Digits(num int) string {
	return fmt.Sprintf("%02d", num)
}

func formatDate(date time.Time) string {
	return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d.%03d",
		date.Year(), date.Month(), date.Day(),
		date.Hour(), date.Minute(), date.Second(), date.Nanosecond()/1e6)
}
