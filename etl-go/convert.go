package main

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

func convertProperty(document bson.M) map[string]interface{} {

	fmt.Println("Converting property:", document["uuid"])
	coverImage := ""

	if document["images"] != nil {
		// panic: interface conversion: interface {} is primitive.A, not []interface {}
		images := document["images"].(primitive.A)
		if len(images) > 0 {
			coverImage = images[0].(string)
		}
	}

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
		"coverImage":            coverImage,
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
