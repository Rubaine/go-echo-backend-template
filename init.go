package main

import (
	"example.com/template/config"
	"example.com/template/models/postgresql"
	"time"

	"github.com/charmbracelet/log"
)

func init() {
	start := time.Now()

	postgresql.SQLCtx, postgresql.SQLConn = config.InitPgSQL()

	log.Debug("Initialization ended", "took", time.Since(start).Round(time.Millisecond).String())
}
