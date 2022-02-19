package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/alma-amirseitov/finance/internal/validator"
	"time"
)

type Payments struct {
	Id int64 `json:"id"`
	Name string `json:"name"`
	Price int32 `json:"price"`
	Date time.Time `json:"date"`
	PaymentType string `json:"payment_type"`
	Comment string `json:"comment"`
	CategoriesName  string `json:"categories_name"`
	Categories
}

func ValidatePayment(v *validator.Validator, p *Payments) {
	v.Check(p.Name != "", "name", "must be provided")
	v.Check(len(p.Name) <= 500, "name", "must not be more than 500 bytes long")
}

type PaymentsModel struct {
	DB *sql.DB
}

func (p PaymentsModel) Insert(payments *Payments) error {
	query := `
        INSERT INTO payments (name, payment_type, comment, category_id,price) 
        VALUES ($1, $2, $3, $4,$5)
        RETURNING id, name`

	args := []interface{}{payments.Name, payments.PaymentType, payments.Comment,payments.Categories.Id,payments.Price}


	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()


	return p.DB.QueryRowContext(ctx, query, args...).Scan(&payments.Id, &payments.Name)
}

func (p PaymentsModel) Get(id int64) (*Payments, error){
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
        SELECT payments.id, payments.name,price, date, payment_type, Comment,c.name
		FROM payments
		JOIN categories c ON c.id = payments.category_id
        where payments.id = $1
        order by payments`

	var payments Payments

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := p.DB.QueryRowContext(ctx, query, id).Scan(
		&payments.Id,
		&payments.Name,
		&payments.Price,
		&payments.Date,
		&payments.PaymentType,
		&payments.Comment,
		&payments.CategoriesName,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &payments, nil
}

func (p PaymentsModel) GetAll(name string ,filters Filters) ([]*Payments,error) {
	query := fmt.Sprintf(`
		SELECT  payments.id, payments.name, payment_type,date, comment, c.name, price
		FROM payments
		JOIN categories c ON c.id = payments.category_id
        WHERE (to_tsvector('simple', payments.name) @@ plainto_tsquery('simple', $1) OR $1 = '')
        ORDER BY id ASC`)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{name}

	rows, err := p.DB.QueryContext(ctx, query, args...)

	defer rows.Close()

	var payments []*Payments
	for rows.Next() {

		var payment Payments

		err := rows.Scan(
			&payment.Id,
			&payment.Name,
			&payment.PaymentType,
			&payment.Date,
			&payment.Comment,
			&payment.CategoriesName,
			&payment.Price,
		)
		if err != nil {
			return nil, err
		}

		payments = append(payments, &payment)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}


	return payments, nil
}

func (p PaymentsModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `DELETE FROM payments
        WHERE id = $1`


	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()


	result, err := p.DB.ExecContext(ctx,query, id)
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

func (p PaymentsModel) Update(payment *Payments) error {
	query := `
        UPDATE payments 
        SET name = $1, payment_type = $2, comment = $3, category_id = $4, price = $5
        WHERE id = $6
        RETURNING id`


	args := []interface{}{
		payment.Name,
		payment.PaymentType,
		payment.Comment,
		payment.Categories.Id,
		payment.Price,
		payment.Id,
	}


	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()



	err := p.DB.QueryRowContext(ctx,query, args...).Scan(&payment.Id)
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
