package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

var a App

func TestMain(m *testing.M) {
	err := a.Initialise(DbUser, DbPassword, "test")
	if err != nil {
		log.Fatalln("error occured while initialising the database")
	}
	m.Run()

}

func createTable() {
	createTableQuery := `
        CREATE TABLE IF NOT EXISTS products (
            id int NOT NULL AUTO_INCREMENT,
            name varchar(255) NOT NULL,
            quantity int,
            price float(10,7),
            PRIMARY KEY (id)
        );`

	_, err := a.DB.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	_, err := a.DB.Exec("DELETE from products")
	_, err = a.DB.Exec("ALTER TABLE products AUTO_INCREMENT=1")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("clear table")
}

func addProduct(name string, quantity int, price float64) {
	query := fmt.Sprintf("INSERT into products(name,quantity,price) VALUES('%v',%v,%v)", name, quantity, price)
	_, err := a.DB.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

}

func TestGetProduct(t *testing.T) {
	createTable()
	clearTable()
	addProduct("Keyboard", 100, 500.00)
	request := httptest.NewRequest("GET", "/product/1", nil)
	response := sendRequest(request)
	checkStatusCode(t, http.StatusOK, response.Code)
}

func checkStatusCode(t *testing.T, expectedStatusCode int, actualStatusCode int) {
	if expectedStatusCode != actualStatusCode {
		t.Errorf("Expected status: %v, Recieved: %v", expectedStatusCode, actualStatusCode)

	}
}

func sendRequest(request *http.Request) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	a.Router.ServeHTTP(recorder, request)
	return recorder
}

func TestCreateProduct(t *testing.T) {
	clearTable()
	var product = []byte(`{"name":"chair", "quantity":1,"price":100}`)
	req, _ := http.NewRequest("POST", "/product", bytes.NewBuffer(product))
	req.Header.Set("Content-type", "application/json")

	response := sendRequest(req)
	checkStatusCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "chair" {
		t.Errorf("Excpected name : %v, Got : %v", "chair", m["name"])
	}
	if m["quantity"] != 1.0 {
		t.Errorf("Excpected quantity : %v, Got : %v", 1, m["Quantity"])
	}

}

func TestDeleteProduct(t *testing.T) {
	clearTable()

	addProduct("connector", 10, 10)

	req, _ := http.NewRequest("GET", "/product/1", nil)
	response := sendRequest(req)
	checkStatusCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/product/1", nil)
	response = sendRequest(req)
	checkStatusCode(t, http.StatusOK, response.Code)
}

func TestUpdateProduct(t *testing.T) {
	clearTable()
	addProduct("Keyboard", 100, 500.00)

	req, _ := http.NewRequest("GET", "/product/1", nil)
	response := sendRequest(req)
	checkStatusCode(t, http.StatusOK, response.Code)

	var oldValue map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &oldValue)

	var product = []byte(`{"name":"connector", "quantity":1,"price":10}`)
	req, _ = http.NewRequest("PUT", "/product/1", bytes.NewBuffer(product))
	req.Header.Set("Content-type", "application/json")

	response = sendRequest(req)
	var newValue map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &newValue)

	if oldValue["id"] != newValue["id"].(float64) {
		t.Errorf("Expected id : %v, Got : %v", oldValue["id"], newValue["id"])
	}

	if oldValue["name"] != newValue["name"] {
		t.Errorf("Expected name : %v, Got : %v", oldValue["name"], newValue["name"])
	}
}
