package main

import (
	"context"
	"employees/aserver"
	"employees/cfg"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/go-sql-driver/mysql"

	"github.com/jmoiron/sqlx"
)

func InitDB(cfg *cfg.Cfg) (*sqlx.DB, error) {
	connUrl := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		cfg.Mysql.User,
		cfg.Mysql.Password,
		cfg.Mysql.Host,
		cfg.Mysql.Port,
		cfg.Mysql.Dbname,
	)

	conn, err := sqlx.Connect(cfg.Mysql.Driver, connUrl)
	if err != nil {
		return nil, err
	}

	err = conn.Ping()
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func main() {
	cMainCfg := &cfg.Cfg{}
	err := cfg.LoadConfig(cMainCfg)
	if err != nil {
		log.Fatalln("Error in LoadConfig: ", err)
	}

	mainConn, err := InitDB(cMainCfg)
	if err != nil {
		log.Fatalln("Error in InitDB: ", err)
	}

	aCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	finishCh := make(chan struct{})
	go func() {
		s := <-signalCh
		log.Printf("got signal %v, graceful shutdown...", s)
		mainConn.Close()
		finishCh <- struct{}{}
	}()

	aServer := aserver.NewServer("localhost", "8080", cMainCfg, mainConn)
	err = aServer.Run(aCtx)
	if err != nil {
		log.Fatalln("couldn't run Atreugo server, exiting...")
	}

	<-finishCh
	log.Println("Finished shutdown")
}
