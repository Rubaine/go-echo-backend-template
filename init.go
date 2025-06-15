package main

import (
	"time"

	"example.com/template/config"
	"example.com/template/models/postgresql"

	"github.com/charmbracelet/log"
)

func init() {
	start := time.Now()

	config.Init(Folder)
	postgresql.SQLCtx, postgresql.SQLConn = config.InitPgSQL()

	log.Debug("Initialization ended", "took", time.Since(start).Round(time.Millisecond).String())
}
