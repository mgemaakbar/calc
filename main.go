package main

import (
	"calc/controller"
	"calc/handler"
	"calc/model"
	"calc/taxcodeinfo"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("Starting up....")
	dbConnStr := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable fallback_application_name=%s",
		"db", `5432`, `dev`, `postgres-dev`,
		`s3cr3tp4ssw0rd`, `wow`)
	db, err := sql.Open("postgres",
		dbConnStr)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS item (
		id serial primary key,
		name text,
		tax_code integer,
		price integer
	);`)

	if err != nil {
		fmt.Println(err)
		return
	}

	querier := model.NewQuerier(db)

	handler := handler.NewHandler(querier, taxcodeinfo.Map)
	controller := controller.NewController(handler)

	http.HandleFunc("/bill", controller.GetBill)
	http.HandleFunc("/create-item", controller.Create)

	fmt.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
