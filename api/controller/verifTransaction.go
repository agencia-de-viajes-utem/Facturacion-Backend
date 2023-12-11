package controller

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
	"webpaygo/api/config"
	"webpaygo/api/models"

	"github.com/fenriz07/Golang-Transbank-WebPay-Rest/pkg/transaction"
	"gorm.io/gorm"
)

// Definición de la estructura TransactionLog
type TransactionLog struct {
	NumberOrder string
	IdSession   string
	Response    string
	Error       string
}

func VerifTransaction(w http.ResponseWriter, r *http.Request) {
	/*
	 En caso del que pago sea anulado, comprobar si existe el parametro TBK_TOKEN.
	  Si existe el pago fue anulado por el usuario y debe comprobarse su estado con el Commit,
	  Si fue anulado adicionalmente tenemos los parametros TBK_ORDEN_COMPRA || TBK_ID_SESION
	*/

	log.Println("******************empieza*************")

	var token string = ""
	var numberOrder string = ""
	var idSession string = ""

	canceledToken := r.FormValue("TBK_TOKEN")

	if len(canceledToken) != 0 {
		token = canceledToken
		numberOrder = r.FormValue("TBK_ORDEN_COMPRA")
		idSession = r.FormValue("TBK_ID_SESION")

		log.Printf("Number Order: %s\n Id Session: %s\n", numberOrder, idSession)

	} else {
		token = r.FormValue("token_ws")
	}

	/*Commit de la transacción y resultado de la misma*/
	resp, err := transaction.Commit(token)

	if err != nil {
		fmt.Println(err)
	}

	// Crea un nuevo LogEntry
	newLogEntry := models.LogEntry{
		NumberOrder: resp.BuyOrder,
		IdSession:   resp.SessionID,
		Status:      resp.Status,
		Amount:      resp.Amount,
		//BuyOrder:        resp.BuyOrder,
		//SessionID:       resp.SessionID,
		AccountingDate:  resp.AccountingDate,
		TransactionDate: resp.TransactionDate,
		PaymentTypeCode: resp.PaymentTypeCode,
		CardDetail: models.CardDetail{
			CardNumber: resp.CardDetail.CardNumber,
		},
		AuthorizationCode: resp.AuthorizationCode,
	}

	// Verificar si el campo Status está vacío y establecerlo como "Anulado" si es necesario
	if newLogEntry.Status == "" {
		newLogEntry.Status = "Anulado"
	}
	// Asignación condicional para NumberOrder
	if numberOrder != "" {
		newLogEntry.NumberOrder = numberOrder
	} else {
		newLogEntry.NumberOrder = resp.BuyOrder
	}

	// Asignación condicional para IdSession
	if idSession != "" {
		newLogEntry.IdSession = idSession
	} else {
		newLogEntry.IdSession = resp.SessionID
	}

	// Asignación condicional para TransactionDate
	if resp.TransactionDate.IsZero() {
		newLogEntry.TransactionDate = time.Now()
	} else {
		newLogEntry.TransactionDate = resp.TransactionDate
	}

	// Verificar si el campo Status está vacío y establecerlo como "Anulado" si es necesario
	if newLogEntry.Status == "" {
		newLogEntry.Status = "Anulado"
	}

	db := config.InitDatabase()

	logTransactionData(db, newLogEntry)

	sqlDB, err := db.DB()
	if err != nil {
		log.Println("Error obteniendo la interfaz de la base de datos:", err)
	} else {
		// Cerrar la conexión a la base de datos después de usarla
		sqlDB.Close()
	}

	log.Println(resp)

	/*Obtención del status de la transacción*/
	resp2, err := transaction.GetStatus(token)

	log.Println(resp2)

	if err != nil {
		log.Println(err)
	}

	/*Anulación*/
	resp3, err := transaction.Refund(token, 1000)

	if err != nil {
		log.Println(err)
	}

	log.Println("Respuesta 3")
	log.Println(resp3)

	view := template.Must(template.ParseGlob("api/views/*"))

	err = view.ExecuteTemplate(w, "status.html", newLogEntry)

	if err != nil {
		http.Error(w, err.Error(), 500)
	}

}

func logTransactionData(db *gorm.DB, newLogEntry models.LogEntry) {
	// Almacenar el log en la base de datos
	err := storeLogEntry(db, newLogEntry)
	if err != nil {
		log.Println("Error storing log in the database:", err)
	}
}

func storeLogEntry(db *gorm.DB, newLogEntry models.LogEntry) error {
	//log.Printf("Inserting log entry into database: %+v\n", newLogEntry)

	err := db.Create(&newLogEntry).Error
	if err != nil {
		log.Println("Error database:", err)
		return err
	}

	return nil
}
