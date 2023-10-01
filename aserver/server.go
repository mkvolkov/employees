package aserver

import (
	"context"
	"employees/cfg"
	"fmt"
	"os"

	"github.com/casbin/casbin/v2"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"github.com/savsgio/atreugo/v11"
)

type AServer struct {
	MpAuth      map[string]string
	AtrServer   *atreugo.Atreugo
	CsbEnforcer *casbin.Enforcer
	DB          *sqlx.DB
	Logger      zerolog.Logger
}

func NewServer(host, port string, cfg *cfg.Cfg, db *sqlx.DB) *AServer {
	addr := host + ":" + port
	aCfg := atreugo.Config{
		Addr: addr,

		GracefulShutdown: true,
	}

	aSrv := atreugo.New(aCfg)

	filePass, err := os.Open("login.txt")
	if err != nil {
		panic("Couldn't open file with users' passwords")
	}

	var user string
	var pass string

	var mpAuth = make(map[string]string)

	for {
		_, err := fmt.Fscanf(filePass, "%s %s\n", &user, &pass)
		if err != nil {
			break
		}

		mpAuth[user] = pass
	}

	cEnforcer, err := casbin.NewEnforcer("model.conf", "policy.csv")
	if err != nil {
		return nil
	}

	logger := zerolog.New(os.Stdout).Level(zerolog.InfoLevel)

	return &AServer{
		MpAuth:      mpAuth,
		AtrServer:   aSrv,
		CsbEnforcer: cEnforcer,
		DB:          db,
		Logger:      logger,
	}
}

func (s *AServer) Run(ctx context.Context) error {
	aHandlers := &RBase{db: s.DB}

	s.MapHandlers(aHandlers)

	if err := s.AtrServer.ListenAndServe(); err != nil {
		s.Logger.Fatal().Msgf("Failed Atreugo ListenAndServe(): %v\n", err)
	}

	return nil
}

func (s *AServer) MapHandlers(rs Routes) {
	s.AtrServer.POST("/hire", rs.HireEmployee()).UseBefore(s.BeforePost).UseBefore(s.BeforeAll)
	s.AtrServer.POST("/fire", rs.FireEmployee()).UseBefore(s.BeforePost).UseBefore(s.BeforeAll)
	s.AtrServer.GET("/getv", rs.GetVacationDays()).UseBefore(s.BeforeGet).UseBefore(s.BeforeAll)
	s.AtrServer.GET("/find", rs.GetEmployeeByName()).UseBefore(s.BeforeGet).UseBefore(s.BeforeAll)
}
