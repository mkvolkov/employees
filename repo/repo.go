package repo

import (
	"context"
	"employees/model"

	"github.com/jmoiron/sqlx"
)

var queryHireEmployee string = `INSERT INTO employees
	(name, phone, gender, age, email, address)
	VALUES(?, ?, ?, ?, ?, ?)`

var queryFireEmployee string = `DELETE FROM employees
	WHERE id = ?`

var queryGetVdaysByID string = `SELECT vdays FROM employees
	WHERE id = ?`

var queryGetEmployeeByName string = `SELECT * FROM employees
	WHERE name LIKE ?`

func HireEmployee(db *sqlx.DB, ctx context.Context, emp *model.Employee) error {
	_, err := db.ExecContext(
		ctx,
		queryHireEmployee,
		emp.Name,
		emp.Phone,
		emp.Gender,
		emp.Age,
		emp.Email,
		emp.Address,
	)

	if err != nil {
		return err
	}

	return nil
}

func FireEmployee(db *sqlx.DB, ctx context.Context, id int) error {
	_, err := db.ExecContext(
		ctx,
		queryFireEmployee,
		id,
	)

	if err != nil {
		return err
	}

	return nil
}

func GetVacationDays(db *sqlx.DB, ctx context.Context, id int) (vd []model.Vdays, err error) {
	err = db.SelectContext(
		ctx,
		&vd,
		queryGetVdaysByID,
		id,
	)

	if err != nil {
		return nil, err
	}

	return vd, nil
}

func GetEmployeeByName(db *sqlx.DB, ctx context.Context, name string) (data []model.Employee, err error) {
	err = db.SelectContext(
		ctx,
		&data,
		queryGetEmployeeByName,
		name,
	)

	if err != nil {
		return nil, err
	}

	return data, nil
}
