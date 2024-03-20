package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

type Product struct {
	ID     int     `json:"prd_id"`
	Name   string  `json:"prd_name"`
	Price  float64 `json:"prd_general_fix_price"`
	GameID int     `json:"prd_game_id"`
}

func InitDB(dsn string) error {
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	fmt.Println("MySQL bağlantısı başarılı!")

	createTableQuery := `
		CREATE TABLE IF NOT EXISTS products (
			prd_id INT AUTO_INCREMENT PRIMARY KEY,
			prd_name VARCHAR(255) NOT NULL,
			prd_general_fix_price FLOAT NOT NULL,
		    prd_game_id INT NOT NULL
		);
	`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		return err
	}

	return nil
}

func GetProducts() ([]Product, error) {
	rows, err := db.Query("SELECT prd_id, prd_name, prd_general_fix_price, prd_game_id FROM products WHERE prd_status > 0 AND prd_hide_stock = 0")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.GameID)
		if err != nil {
			return nil, err
		}

		products = append(products, p)
	}

	return products, nil
}

func GetProductByID(id int) (*Product, error) {
	var p Product
	err := db.QueryRow("SELECT prd_id, prd_name, prd_general_fix_price, prd_game_id FROM products WHERE prd_id = ? AND prd_status > 0 AND prd_hide_stock = 0", id).Scan(&p.ID, &p.Name, &p.Price, &p.GameID)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func GetProductsByGameID(gameID int) ([]Product, error) {
	rows, err := db.Query("SELECT prd_id, prd_name, prd_general_fix_price, prd_game_id FROM products WHERE prd_game_id = ? AND prd_status > 0 AND prd_hide_stock = 0", gameID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.GameID)
		if err != nil {
			return nil, err
		}

		products = append(products, p)
	}

	return products, nil
}
