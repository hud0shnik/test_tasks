package config

import "github.com/joho/godotenv"

// Init загружает переменные окружения из .env файлов. Без параметров функция будет искать файл .env
func Init(filenames ...string) error {

	return godotenv.Load(filenames...)

}
