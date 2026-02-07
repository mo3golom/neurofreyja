package main

import (
	"context"

	"neurofreyja/internal/app"
	"neurofreyja/internal/processes/delete_messages"
)

func main() {
	application, err := app.Bootstrap()
	if err != nil {
		panic(err)
	}
	defer application.DB.Close()

	application.RegisterHandlers()

	runner := &delete_messages.Runner{
		History:   application.History,
		Messenger: application.Messenger,
		Logger:    application.Logger,
		Interval:  application.Config.DeleteSweepInterval,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go runner.Run(ctx)

	application.Logger.Info("bot started")
	application.Bot.Start()
}
