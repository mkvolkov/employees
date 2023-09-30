package aserver

import (
	"bytes"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"errors"

	"github.com/savsgio/atreugo/v11"
)

const (
	ResponseOK            = 200
	ResponseBadRequest    = 400
	ResponseAuthNeeded    = 401
	ResponseInternalError = 500
)

func (a *AServer) BeforePost(ctx *atreugo.RequestCtx) error {
	var err error

	auth := ctx.Request.Header.Peek("Authorization")
	if len(auth) == 0 {
		return ctx.TextResponse("Authorization needed\n", ResponseAuthNeeded)
	}

	crdts := extractUserPass(auth)
	if crdts == nil {
		return ctx.TextResponse("Error: couldn't extract user\n", ResponseInternalError)
	}

	user := crdts[0]
	pass := crdts[1]

	err = a.CsbEnforcer.LoadPolicy()
	if err != nil {
		return ctx.TextResponse(err.Error(), ResponseInternalError)
	}

	ok, err := a.CsbEnforcer.Enforce(string(user), "data", "write")
	if err != nil {
		return ctx.TextResponse(err.Error(), ResponseInternalError)
	}

	if !ok {
		return ctx.TextResponse("Casbin: unauthorized\n", ResponseAuthNeeded)
	}

	userStr := string(user)
	passStr, ok := a.MpAuth[userStr]
	if !ok {
		return ctx.TextResponse("Unknown user: unauthorized\n", ResponseAuthNeeded)
	}

	if checkAuth(user, pass, userStr, passStr) {
		return ctx.Next()
	} else {
		err = errors.New("HTTP error: wrong login and/or password\n")
	}

	return err
}

func (a *AServer) BeforeGet(ctx *atreugo.RequestCtx) error {
	var err error

	auth := ctx.Request.Header.Peek("Authorization")
	if len(auth) == 0 {
		return ctx.TextResponse("Authorization needed\n", ResponseAuthNeeded)
	}

	crdts := extractUserPass(auth)
	if crdts == nil {
		return ctx.TextResponse("Error: couldn't extract user\n", ResponseInternalError)
	}

	user := crdts[0]
	pass := crdts[1]

	err = a.CsbEnforcer.LoadPolicy()
	if err != nil {
		return ctx.TextResponse(err.Error(), ResponseInternalError)
	}

	ok, err := a.CsbEnforcer.Enforce(string(user), "data", "read")
	if err != nil {
		return ctx.TextResponse(err.Error(), ResponseInternalError)
	}

	if !ok {
		return ctx.TextResponse("Casbin: unauthorized\n", ResponseAuthNeeded)
	}

	userStr := string(user)
	passStr, ok := a.MpAuth[userStr]
	if !ok {
		return ctx.TextResponse("Unknown user: unauthorized\n", ResponseAuthNeeded)
	}

	if checkAuth(user, pass, userStr, passStr) {
		return ctx.Next()
	} else {
		err = errors.New("HTTP error: wrong login and/or password\n")
	}

	return err
}

func extractUserPass(authStr []byte) [][]byte {
	i := bytes.IndexByte(authStr, ' ')
	if i == -1 {
		return nil
	}

	if !bytes.EqualFold(authStr[:i], []byte("basic")) {
		return nil
	}

	decoded, err := base64.StdEncoding.DecodeString(string(authStr[i+1:]))
	if err != nil {
		return nil
	}

	crdts := bytes.Split(decoded, []byte(":"))
	if len(crdts) <= 1 {
		return nil
	}

	user := crdts[0]
	pass := crdts[1]

	res := make([][]byte, 0)
	res = append(res, user)
	res = append(res, pass)

	return res
}

func checkAuth(user, pass []byte, login, password string) bool {
	userHash := sha256.Sum256(user)
	passHash := sha256.Sum256(pass)

	expUserHash := sha256.Sum256([]byte(login))
	expPassHash := sha256.Sum256([]byte(password))

	userMatch := (subtle.ConstantTimeCompare(userHash[:], expUserHash[:]) == 1)
	passMatch := (subtle.ConstantTimeCompare(passHash[:], expPassHash[:]) == 1)

	if userMatch && passMatch {
		return true
	}

	return false
}

func (a *AServer) BeforeAll(ctx *atreugo.RequestCtx) error {
	ctype := ctx.Request.Header.Peek("Content-Type")
	if len(ctype) == 0 {
		return ctx.TextResponse("No content type\n", ResponseBadRequest)
	}

	ctypeStr := string(ctype)
	if ctypeStr == "application/json" {
		ctx.SetUserValue("ctype", "json")
	} else if ctypeStr == "application/xml" {
		ctx.SetUserValue("ctype", "xml")
	} else {
		return ctx.TextResponse("Unknown content type\n", ResponseBadRequest)
	}

	return ctx.Next()
}
