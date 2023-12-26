package main

import (
	"fmt"
	"gomodel/internal/shared/util/password"
	"log/slog"

	"github.com/google/uuid"
)

func Teardown() {

}

func Init(
	logger *slog.Logger,
) error {
	fmt.Println(uuid.NewString())

	p := "admin123"
	hash, _ := password.HashAndSalt(p, p)
	fmt.Println(hash)

	logger.Info("Stub completed.")
	return nil
}
