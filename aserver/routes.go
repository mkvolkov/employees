package aserver

import (
	"context"
	"employees/model"
	"employees/repo"

	"github.com/goccy/go-json"
	"github.com/jmoiron/sqlx"
	"github.com/savsgio/atreugo/v11"
)

type Routes interface {
	HireEmployee() atreugo.View
	FireEmployee() atreugo.View
	GetVacationDays() atreugo.View
	FindEmployeeByID() atreugo.View
}

type RBase struct {
	db *sqlx.DB
}

func (r *RBase) HireEmployee() atreugo.View {
	return func(ctx *atreugo.RequestCtx) error {
		var empl model.Employee

		ctxDb := context.Background()

		ctype := ctx.Value("ctype").(string)
		if ctype == "json" {
			err := json.Unmarshal(ctx.Request.Body(), &empl)
			if err != nil {
				return ctx.TextResponse(err.Error(), ResponseInternalError)
			}

			err = repo.HireEmployee(r.db, ctxDb, &empl)
			if err != nil {
				return ctx.ErrorResponse(err, ResponseInternalError)
			}
		}

		return ctx.TextResponse("Employee added successfully\n")
	}
}

func (r *RBase) FireEmployee() atreugo.View {
	return func(ctx *atreugo.RequestCtx) error {
		var delID model.DeleteID

		ctxDb := context.Background()

		ctype := ctx.Value("ctype").(string)
		if ctype == "json" {
			err := json.Unmarshal(ctx.Request.Body(), &delID)
			if err != nil {
				return ctx.TextResponse(err.Error(), ResponseInternalError)
			}

			err = repo.FireEmployee(r.db, ctxDb, delID.ID)
			if err != nil {
				return ctx.ErrorResponse(err, ResponseInternalError)
			}
		}

		return ctx.TextResponse("Employee deleted successfully\n")
	}
}

func (r *RBase) GetVacationDays() atreugo.View {
	return func(ctx *atreugo.RequestCtx) error {
		return ctx.TextResponse("GetVacationDays\n")
	}
}

func (r *RBase) FindEmployeeByID() atreugo.View {
	return func(ctx *atreugo.RequestCtx) error {
		return ctx.TextResponse("FindEmployeeByID\n")
	}
}
