package main

import (
	"comments/internal/config"

	"github.com/sirupsen/logrus"
)

func main() {

	err := config.Init()
	if err != nil {
		logrus.Fatalf("InitConfig error: %v", err)
	}

}
