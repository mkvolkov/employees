package model

import (
	"encoding/xml"

	"github.com/guregu/null"
)

type Employee struct {
	ID      int      `db:"id"`
	Name    string   `db:"name" json:"name" xml:"name"`
	Phone   string   `db:"phone" json:"phone" xml:"phone"`
	Gender  string   `db:"gender" json:"gender" xml:"gender"`
	Age     int      `db:"age" json:"age" xml:"age"`
	Email   string   `db:"email" json:"email" xml:"email"`
	Address string   `db:"address" json:"address" xml:"address"`
	Vdays   null.Int `db:"vdays"`
}

type XEmployee struct {
	XMLName xml.Name   `xml:"data"`
	Empl    []Employee `xml:"employee"`
}

type ModifyID struct {
	ID int `json:"id"`
}

type XModifyID struct {
	XMLName xml.Name `xml:"data"`
	EmplID  struct {
		ID int `xml:"id"`
	} `xml:"empl_id"`
}

type ModifyName struct {
	Name string `json:"name"`
}

type XModifyName struct {
	XMLName  xml.Name `xml:"data"`
	EmplName struct {
		Name string `xml:"name"`
	} `xml:"empl_name"`
}

type Vdays struct {
	Vdays int `db:"vdays"`
}

type XVdays struct {
	XMLName xml.Name `xml:"data"`
	Vdays   struct {
		Days int `xml:"days"`
	} `xml:"vacation_days"`
}
