from time import sleep
import json

import requests


# Функция парсинга фильмов
def parse_kinopoisk():

    # Переменная для хранения текущей позиции
    left = 0

    # Список, в который будет вестись запись
    result = []

    # Переменная места в рейтинге
    place=1

    # Цикл обработки страниц
    for page in range(20):

        # Запрос на Кинопоиск
        resp = requests.get('https://www.kinopoisk.ru/lists/movies/top_1000/?sort=rating&page='+str(page+1))

        # Html страницы в формате str
        res = resp.text

        # Цикл обработки фильмов
        for _ in range(50):

            film = {}

            # Поиск и запись названия фильма
            left = res.find('"__typename":"Title","russian":"', left)+1
            film["title"] = res[left+31:res.find('"', left+31)]

            # Поиск и запись информации о доступности фильма в кинотеатре
            left = res.find(':213})":', left)+1
            film["ticket_available"] = ((res[left+7:res.find(',', left+7)])) == 'true'

            # Запись места в рейтинге
            film["place"] = place
            place+=1

            # Добавление фильма в список
            result.append(film)

        # Пауза на 3 секунды (Кинопоиск блокирует при частых запросах)
        sleep(3)

    # Запись списка в json файл
    file = open("films.json", "w")
    file.write(json.dumps(result))
    file.close()
