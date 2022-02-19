package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Categories struct {
	Id int64 `json:"id"`
	Name string `json:"name"`
}

type CategoriesModel struct {
	DB *sql.DB
}

func (c CategoriesModel) Insert(category *Categories) error {

	query := `INSERT INTO categories (name) 
        VALUES ($1)
        RETURNING id, name`

	args := []interface{}{category.Name}


	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()


	return c.DB.QueryRowContext(ctx, query, args...).Scan(&category.Id, &category.Name)
}

func (c CategoriesModel) Get(id int64) (*Categories, error){
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `SELECT categories.id, categories.name
		FROM categories
        where categories.id = $1`

	var category Categories

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := c.DB.QueryRowContext(ctx, query, id).Scan(
		&category.Id,
		&category.Name,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	fmt.Println(category)
	return &category, nil
}

func (c CategoriesModel) GetByName(name string) (*Categories, error){

	query := `SELECT categories.id, categories.name
		FROM categories
        where categories.name = $1`

	var category Categories

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := c.DB.QueryRowContext(ctx, query, name).Scan(
		&category.Id,
		&category.Name,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	fmt.Println(category)
	return &category, nil
}

func (c CategoriesModel) GetAll() ([]*Categories,error) {
	query := fmt.Sprintf(`SELECT  id, name FROM categories`)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()


	rows, err := c.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var categories []*Categories
	for rows.Next() {

		var category Categories

		err := rows.Scan(
			&category.Id,
			&category.Name,
		)
		if err != nil {
			return nil, err
		}

		categories = append(categories, &category)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func (c CategoriesModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `DELETE FROM categories
        WHERE id = $1`


	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()


	result, err := c.DB.ExecContext(ctx,query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (c CategoriesModel) Update(category *Categories) error {
	query := `
        UPDATE categories 
        SET name = $1
        WHERE id = $2
        RETURNING id`


	args := []interface{}{
		category.Name,
		category.Id,
	}


	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()



	err := c.DB.QueryRowContext(ctx,query, args...).Scan(&category.Id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

