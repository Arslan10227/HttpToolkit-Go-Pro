package main

import (
	"fmt"
	"os"

	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/backup"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/config"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/server"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/settings"
)

func main() {
	cfg, err := config.LoadDefault()
	if err != nil {
		fmt.Fprintf(os.Stderr, "config: %v\n", err)
		os.Exit(1)
	}
	sm := settings.NewManager(cfg.ConfigDir)
	sm.SetBackup(backup.New(cfg))
	_ = sm.Load()
	app, err := server.Run(cfg, sm)
	if err != nil {
		fmt.Fprintf(os.Stderr, "server: %v\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "httptoolkit-go listening REST=%d admin=%d token=%s\n",
		cfg.ServerPort, cfg.AdminPort, app.AuthToken)
	select {}
}
