package main

import (
	"gomodel/cmd/seed/seed"
	"log/slog"
)

func Teardown() {

}

func Init(
	seed *seed.Seed,
	logger *slog.Logger,
) error {
	err := seed.Run()
	if err != nil {
		return err
	}

	logger.Info("Seed completed.")
	return nil
}
