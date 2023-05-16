package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

// Глобальный Url бота
var botUrl string

// Структура респонса Telegram'а
type telegramResponse struct {
	Result []update `json:"result"`
}

// Структура апдейта
type update struct {
	UpdateId int     `json:"update_id"`
	Message  message `json:"message"`
}

// Структура сообщения
type message struct {
	MessageId int    `json:"message_id"`
	Chat      chat   `json:"chat"`
	Text      string `json:"text"`
}

// Структура чата
type chat struct {
	ChatId int `json:"id"`
}

// Структура для отправки сообщений
type sendMessage struct {
	ChatId              int    `json:"chat_id"`
	Text                string `json:"text"`
	ParseMode           string `json:"parse_mode"`
	HasProtectedContent bool   `json:"has_protected_content"`
}

// Структура для получения id сообщения после отправки
type sendMessageResponse struct {
	Result struct {
		MessageId int `json:"message_id"`
	} `json:"result"`
}

// Структура для удаления сообщений
type deleteMessage struct {
	ChatId    int `json:"chat_id"`
	MessageId int `json:"message_id"`
}

// Структура пароля
type password struct {
	Login    string
	Password string
	Service  string
}

// Функция инициализации конфига (всех токенов)
func initConfig() error {

	// Путь и имя файла
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")

	return viper.ReadInConfig()
}

// Функция отправки сообщения
func sendMsg(chatId int, text string) error {

	// Формирование сообщения
	buf, err := json.Marshal(sendMessage{
		ChatId:              chatId,
		Text:                text,
		ParseMode:           "HTML",
		HasProtectedContent: true,
	})
	if err != nil {
		log.Printf("in sendMsg: json.Marshal error: %s", err)
		return err
	}

	// Отправка сообщения
	resp, err := http.Post(botUrl+"/sendMessage", "application/json", bytes.NewBuffer(buf))
	if err != nil {
		log.Printf("in sendMsg: http.Post error: %s", err)
		return err
	}
	defer resp.Body.Close()

	// Получение id сообщения
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("in sendMsg: ioutil.ReadAll error: %s", err)
		return err
	}
	var msg sendMessageResponse
	err = json.Unmarshal(body, &msg)
	if err != nil {
		log.Printf("in sendMsg: json.Unmarshal error: %s", err)
		return err
	}

	// Запуск таймера для удаления сообщения
	go deleteMsg(chatId, msg.Result.MessageId)

	return nil
}

// Функция удаления сообщения (запускать только в горутине)
func deleteMsg(chatId, messageId int) error {

	// Ожидание 15ти секунд
	time.Sleep(time.Second * 15)

	// Формирование структуры удаления
	buf, err := json.Marshal(deleteMessage{
		ChatId:    chatId,
		MessageId: messageId,
	})
	if err != nil {
		log.Printf("in deleteMsg: json.Marshal error: %s", err)
		return err
	}

	// Удаление сообщения
	_, err = http.Post(botUrl+"/deleteMessage", "application/json", bytes.NewBuffer(buf))
	if err != nil {
		log.Printf("in deleteMsg: http.Post error: %s", err)
		return err
	}

	return nil
}

// Функция подключения к БД
func connectDB() (*sqlx.DB, error) {

	// Инициализация переменных окружения
	godotenv.Load()

	// Подключение к БД
	db, err := sqlx.Open("postgres",
		fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_USERNAME"),
			os.Getenv("DB_NAME"),
			os.Getenv("DB_PASSWORD")))
	if err != nil {
		return nil, fmt.Errorf("in sqlx.Open: %w", err)
	}

	// Проверка подключения
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("in db.Ping: %w", err)
	}

	return db, nil
}

// Функция генерации и отправки ответов
func respond(update update) {

	// Всё кроме сообщений бот проигнорирует
	if update.Message.Text == "" {
		return
	}

	// Проверка на sql инъекции
	if strings.ContainsAny(update.Message.Text, "`'%") {
		sendMsg(update.Message.Chat.ChatId, "Пожалуйста не используйте символы <b>'</b> и <b>%</b>")
		return
	}

	// Разделение текста пользователя на слайс
	request := append(strings.Split(update.Message.Text, " "), "", "", "", "")

	// Обработчик команд
	switch request[0] {
	case "/help", "/start":
		sendMsg(update.Message.Chat.ChatId, "Привет! Вот список команд:\n /set <b>service login password</b> - установить пароль\n /get <b>service</b> - получить пароль\n /del <b>service</b> - удалить пароль\n\nУ каждого сервиса может быть только один пароль")
	case "/set":
		setPassword(update.Message.Chat.ChatId, request[1], request[2], request[3])
		go deleteMsg(update.Message.Chat.ChatId, update.Message.MessageId)
	case "/get":
		getPassword(update.Message.Chat.ChatId, request[1])
		go deleteMsg(update.Message.Chat.ChatId, update.Message.MessageId)
	case "/del":
		delPassword(update.Message.Chat.ChatId, request[1])
		go deleteMsg(update.Message.Chat.ChatId, update.Message.MessageId)
	default:
		sendMsg(update.Message.Chat.ChatId, "Я не понимаю, воспользуйтесь командой /help")
	}
}

// Функция добавления пароля
func setPassword(chatId int, serviceName, login, password string) {

	// Проверка на параметры
	if serviceName == "" || login == "" || password == "" {
		sendMsg(chatId, "Пожалуйста воспользуйся синтаксисом ниже:\n/set <b>service login password</b>")
		return
	}

	// Подключение к БД
	db, err := connectDB()
	if err != nil {
		log.Fatalf("in connectDB: %s", err)
	}

	// Добавление пароля в базу
	_, err = db.Exec("INSERT INTO passwords (chat_id, login, password, service) values ($1, $2, $3, $4)",
		chatId, login, password, serviceName)
	if err != nil {
		sendMsg(chatId, "Не получилось добавить")
		return
	}

	sendMsg(chatId, "Готово!")
}

// Функция получения пароля
func getPassword(chatId int, serviceName string) {

	// Проверка на параметр
	if serviceName == "" {
		sendMsg(chatId, "Пожалуйста воспользуйся синтаксисом ниже:\n/get <b>service</b>")
		return
	}

	// Подключение к БД
	db, err := connectDB()
	if err != nil {
		log.Fatalf("in connectDB: %s", err)
	}

	// Получение и проверка данных
	var result []password
	err = db.Select(&result, fmt.Sprintf("SELECT login, password, service FROM passwords WHERE chat_id=%d AND service='%s'", chatId, serviceName))
	if err != nil {
		sendMsg(chatId, "Не удалось найти пароль")
		return
	}

	if len(result) == 0 {
		sendMsg(chatId, "Пароль не найден")
		return
	}

	// Вывод
	sendMsg(chatId, fmt.Sprintf("Сервис - %s\nЛогин - %s\nПароль - %s", serviceName, result[0].Login, result[0].Password))
}

// Функция удаления пароля
func delPassword(chatId int, serviceName string) {

	// Проверка на параметр
	if serviceName == "" {
		sendMsg(chatId, "Пожалуйста воспользуйся синтаксисом ниже:\n/del <b>service</b>")
		return
	}

	// Подключение к БД
	db, err := connectDB()
	if err != nil {
		log.Fatalf("in connectDB: %s", err)
	}

	// Удаление пароля
	_, err = db.Exec("DELETE FROM passwords WHERE chat_id=$1 AND service=$2", chatId, serviceName)
	if err != nil {
		sendMsg(chatId, "Не получилось удалить")
		return
	}

	sendMsg(chatId, "Пароль удалён")
}

// Функция получения апдейтов
func getUpdates(offset int) ([]update, error) {

	// Запрос для получения апдейтов
	resp, err := http.Get(botUrl + "/getUpdates?offset=" + strconv.Itoa(offset))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Запись и обработка полученных данных
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var restResponse telegramResponse
	err = json.Unmarshal(body, &restResponse)
	if err != nil {
		return nil, err
	}

	return restResponse.Result, nil
}

func main() {

	// Инициализация конфига
	err := initConfig()
	if err != nil {
		log.Fatalf("initConfig error: %s", err)
		return
	}

	// Url бота для отправки и приёма сообщений
	botUrl = "https://api.telegram.org/bot" + viper.GetString("token")
	offSet := 0

	// Цикл работы бота
	for {

		// Получение апдейтов
		updates, err := getUpdates(offSet)
		if err != nil {
			log.Fatalf("getUpdates error: %s", err)
			return
		}

		// Обработка апдейтов
		for _, update := range updates {
			respond(update)
			offSet = update.UpdateId + 1
		}

	}

}
