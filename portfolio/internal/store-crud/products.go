package crud

import (
	"database/sql"
	"errors"
	"fmt"
)

func InsertProducts(code int, name string, cat string, imageUrl string, price int) error {
	db, err := makeCN()
	if err != nil {
		return err
	}

	query := "INSERT INTO products (code, name, category, price, image_url) VALUES (?, ?, ?, ?, ?)"
	result, err := db.Exec(query, code, name, cat, price, imageUrl)
	if err != nil {
		if err.Error() == "Error 1062 (23000): Duplicate entry '1' for key 'products.PRIMARY'" {
			return errors.New("The product already exists")
		}

		return err
	}

	_, err = result.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}

func GetProductByCode(code int) (*Products, error) {
	db, err := makeCN()
	if err != nil {
		return nil, err
	}

	query := "SELECT * FROM products WHERE code = ?"
	row := db.QueryRow(query, code)

	p := &Products{}

	if err = row.Scan(&p.Code, &p.Name, &p.Cat, &p.Price, &p.ImageUrl); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		} else {
			return nil, err
		}
	}

	return p, err
}

func GetAllProducts() ([]*Products, error) {
	db, err := makeCN()
	if err != nil {
		return nil, err
	}

	products := []*Products{}

	query := "SELECT * FROM products ORDER BY code ASC"

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		p := &Products{}
		err := rows.Scan(&p.Code, &p.Name, &p.Cat, &p.ImageUrl, &p.Price)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, err
			} else {
				return nil, err
			}
		}

		products = append(products, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func UpdateProducts(columnEdit string, newValue any, code int) error {
	db, err := makeCN()
	if err != nil {
		return err
	}

	if _, err = GetProductByCode(code); err == sql.ErrNoRows {
		return errors.New("The product doesn't exists")
	}
	var query string
	switch columnEdit {
	case "code":
		query = "UPDATE products SET code = ? WHERE code = ?"
	case "name":
		query = "UPDATE products SET name = ? WHERE code = ?"
	case "category", "cat":
		query = "UPDATE products SET category = ? WHERE code = ?"
	case "price":
		query = "UPDATE products SET price = ? WHERE code = ?"
	case "image", "image-url", "url", "imageurl", "image url":
		query = "UPDATE products SET image_url = ? WHERE code = ?"
	default:
		err = fmt.Errorf("invalid column")
		return err
	}
	_, err = db.Exec(query, newValue, code)
	if err != nil {
		return err
	}

	return nil
}

func DeleteProducts(code int) error {
	db, err := makeCN()
	if err != nil {
		return err
	}

	if _, err = GetProductByCode(code); err != nil {
		return errors.New("The product doesn't exists")
	}

	query := "DELETE FROM products WHERE code = ?"
	_, err = db.Exec(query, code)
	if err != nil {
		return err
	}

	return nil
}
