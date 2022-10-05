package crud

import (
	"database/sql"
	"errors"
)

func InsertProducts(code int, name string, cat string, imageUrl string, price float64) error {
	db, err := makeCN()

	query := "INSERT INTO products (code,name,cat,price,image_url) VALUES ($1,$2,$3,$4,$5)"
	result, err := db.Exec(query, code, name, cat, price, imageUrl)
	if err != nil {
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

	query := "SELECT * FROM products WHERE code = $1"
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
		err := rows.Scan(&p.Code, &p.Name, &p.Cat, &p.Price, &p.ImageUrl)
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

	var query string
	switch columnEdit {
	case "code":
		query = "UPDATE products SET code = $1 WHERE code = $2"
	case "name":
		query = "UPDATE products SET name = $1 WHERE code = $2"
		break
	case "category", "cat":
		query = "UPDATE products SET category = $1 WHERE code = $2"
		break
	case "price":
		query = "UPDATE products SET price = $1 WHERE code = $2"
		break
	case "image", "image-url", "url", "imageurl", "image url":
		query = "UPDATE products SET image_url = $1 WHERE code = $2"
		break
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

	query := "DELETE FROM products WHERE code = $1"
	_, err = db.Exec(query, code)
	if err != nil {
		return err
	}

	return nil
}
