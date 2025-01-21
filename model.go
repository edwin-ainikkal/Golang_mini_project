package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

type product struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"	`
}

var err error

func getProducts(db *sql.DB) ([]product, error) {
	query := "SELECT id,name,quantity,price from products"
	rows, err := db.Query(query)

	if err != nil {
		return nil, err
	}

	products := []product{}
	for rows.Next() {
		var p product
		err := rows.Scan(&p.ID, &p.Name, &p.Quantity, &p.Price)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil

}

func (p *product) getProduct(db *sql.DB) error {
	query := fmt.Sprintf("SELECT name,quantity,price from products where id = %v", p.ID)
	rows := db.QueryRow(query)
	err := rows.Scan(&p.Name, &p.Quantity, &p.Price)
	if err != nil {
		return err
	}
	return nil
}

func (p *product) createProduct(db *sql.DB) error {
	query := fmt.Sprintf("INSERT into products(name,quantity,price) values('%v',%v,%v)", p.Name, p.Quantity, p.Price)
	result, err := db.Exec(query)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	p.ID = int(id)
	return nil

}

func (p *product) updateProduct(db *sql.DB) error {
	query := fmt.Sprintf("UPDATE products set name='%v', quantity=%v, price=%v where id=%v", p.Name, p.Quantity, p.Price, p.ID)
	result, err := db.Exec(query)
	if err != nil {
		return err
	}
	log.Println(result.RowsAffected())
	rowsaff, err := result.RowsAffected()
	if rowsaff == 0 {
		return errors.New("no such row exists")
	}
	log.Println(result.RowsAffected())
	return nil

}

func (p *product) deleteProduct(db *sql.DB) error {
	query := fmt.Sprintf("DELETE FROM products WHERE id=%v", p.ID)
	_, err := db.Exec(query)
	if err != nil {
		log.Printf("Error deleting product with ID %v: %v", p.ID, err)
		return err
	}
	log.Printf("Product with ID %v deleted successfully", p.ID)
	return nil
}
