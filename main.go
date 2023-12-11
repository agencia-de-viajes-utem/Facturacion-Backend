package main

import (
	"fmt"
	"log"
	"webpaygo/api/handler"
	"webpaygo/api/models"

	"github.com/fenriz07/Golang-Transbank-WebPay-Rest/pkg/webpayplus"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
		return
	}

	// Accede a la variable DatoTransaction desde el paquete models
	transactions := models.DatoTransaction

	// Imprime el contenido de DatoTransaction en la consola
	fmt.Println("Contenido de DatoTransaction:")
	for _, transaction := range transactions {
		fmt.Printf("OrdenID: %s, SessionID: %s, Monto: %d, UrlRetorno: %s\n",
			transaction.OrdenID, transaction.SessionID, transaction.Monto, transaction.UrlRetorno)
	}

	webpayplus.SetEnvironmentIntegration()

	handler.Init()

}
