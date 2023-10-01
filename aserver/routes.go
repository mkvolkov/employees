package aserver

import (
	"context"
	"employees/model"
	"employees/repo"
	"encoding/xml"
	"fmt"

	"github.com/goccy/go-json"
	"github.com/jmoiron/sqlx"
	"github.com/savsgio/atreugo/v11"
)

type Routes interface {
	HireEmployee() atreugo.View
	FireEmployee() atreugo.View
	GetVacationDays() atreugo.View
	GetEmployeeByName() atreugo.View
}

type RBase struct {
	db *sqlx.DB
}

func (r *RBase) HireEmployee() atreugo.View {
	return func(ctx *atreugo.RequestCtx) error {
		var empl model.Employee
		var xempl model.XEmployee

		ctxDb := context.Background()

		ctype := ctx.Value("ctype").(string)
		if ctype == "json" {
			err := json.Unmarshal(ctx.Request.Body(), &empl)
			if err != nil {
				return ctx.TextResponse(err.Error(), ResponseBadRequest)
			}

			err = repo.HireEmployee(r.db, ctxDb, &empl)
			if err != nil {
				return ctx.ErrorResponse(err, ResponseInternalError)
			}
		}

		if ctype == "xml" {
			fmt.Println(string(ctx.Request.Body()))
			err := xml.Unmarshal(ctx.Request.Body(), &xempl)
			if err != nil {
				return ctx.TextResponse(err.Error(), ResponseBadRequest)
			}

			if len(xempl.Empl) > 0 {
				empl.Name = xempl.Empl[0].Name
				empl.Phone = xempl.Empl[0].Phone
				empl.Gender = xempl.Empl[0].Gender
				empl.Age = xempl.Empl[0].Age
				empl.Email = xempl.Empl[0].Email
				empl.Address = xempl.Empl[0].Address
			}

			err = repo.HireEmployee(r.db, ctxDb, &empl)
			if err != nil {
				return ctx.ErrorResponse(err, ResponseInternalError)
			}
		}

		return ctx.TextResponse("Employee added successfully\n", ResponseOK)
	}
}

func (r *RBase) FireEmployee() atreugo.View {
	return func(ctx *atreugo.RequestCtx) error {
		var delID model.ModifyID
		var xdelID model.XModifyID

		ctxDb := context.Background()

		ctype := ctx.Value("ctype").(string)
		if ctype == "json" {
			err := json.Unmarshal(ctx.Request.Body(), &delID)
			if err != nil {
				return ctx.TextResponse(err.Error(), ResponseBadRequest)
			}

			err = repo.FireEmployee(r.db, ctxDb, delID.ID)
			if err != nil {
				return ctx.ErrorResponse(err, ResponseInternalError)
			}
		}

		if ctype == "xml" {
			err := xml.Unmarshal(ctx.Request.Body(), &xdelID)
			if err != nil {
				return ctx.ErrorResponse(err, ResponseBadRequest)
			}

			err = repo.FireEmployee(r.db, ctxDb, xdelID.EmplID.ID)
			if err != nil {
				return ctx.ErrorResponse(err, ResponseInternalError)
			}
		}

		return ctx.TextResponse("Employee deleted successfully\n", ResponseOK)
	}
}

func (r *RBase) GetVacationDays() atreugo.View {
	return func(ctx *atreugo.RequestCtx) error {
		var getID model.ModifyID
		var xgetID model.XModifyID

		var data []model.Vdays
		var xvdata model.XVdays

		ctxDb := context.Background()

		ctype := ctx.Value("ctype").(string)
		if ctype == "json" {
			err := json.Unmarshal(ctx.Request.Body(), &getID)
			if err != nil {
				return ctx.TextResponse(err.Error(), ResponseBadRequest)
			}

			data, err = repo.GetVacationDays(r.db, ctxDb, getID.ID)
			if err != nil {
				return ctx.ErrorResponse(err, ResponseInternalError)
			}
		}

		if ctype == "xml" {
			err := xml.Unmarshal(ctx.Request.Body(), &xgetID)
			if err != nil {
				return ctx.ErrorResponse(err, ResponseBadRequest)
			}

			data, err = repo.GetVacationDays(r.db, ctxDb, xgetID.EmplID.ID)
			if err != nil {
				return ctx.ErrorResponse(err, ResponseInternalError)
			}

			if len(data) > 0 {
				xvdata.Vdays.Days = data[0].Vdays
			}

			respBytes, err := xml.Marshal(xvdata)
			if err != nil {
				return ctx.ErrorResponse(err, ResponseInternalError)
			}

			return ctx.TextResponse(string(respBytes), ResponseOK)
		}

		return ctx.JSONResponse(data)
	}
}

func (r *RBase) GetEmployeeByName() atreugo.View {
	return func(ctx *atreugo.RequestCtx) error {
		var getName model.ModifyName
		var xgetName model.XModifyName

		var data []model.Employee
		var xdata model.XEmployee

		ctxDb := context.Background()

		ctype := ctx.Value("ctype").(string)
		if ctype == "json" {
			err := json.Unmarshal(ctx.Request.Body(), &getName)
			if err != nil {
				return ctx.TextResponse(err.Error(), ResponseBadRequest)
			}

			modName := fmt.Sprintf("%%%s%%", getName.Name)

			data, err = repo.GetEmployeeByName(r.db, ctxDb, modName)
			if err != nil {
				return ctx.ErrorResponse(err, ResponseInternalError)
			}
		}

		if ctype == "xml" {
			err := xml.Unmarshal(ctx.Request.Body(), &xgetName)
			if err != nil {
				return ctx.ErrorResponse(err, ResponseBadRequest)
			}

			modName := fmt.Sprintf("%%%s%%", xgetName.EmplName.Name)

			data, err = repo.GetEmployeeByName(r.db, ctxDb, modName)
			if err != nil {
				return ctx.ErrorResponse(err, ResponseInternalError)
			}

			xdata.XMLName.Local = "data"
			xdata.XMLName.Space = "data"
			xdata.Empl = data

			respBytes, err := xml.Marshal(xdata)
			if err != nil {
				return ctx.ErrorResponse(err, ResponseInternalError)
			}

			return ctx.TextResponse(string(respBytes), ResponseOK)
		}

		return ctx.JSONResponse(data)
	}
}
