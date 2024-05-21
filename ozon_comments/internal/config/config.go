package config

import "github.com/joho/godotenv"

func Init(filenames ...string) error {

	return godotenv.Load(filenames...)

}
