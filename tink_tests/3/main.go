package main

import (
	"fmt"
)

func main() {

	var n int
	var s string
	fmt.Scan(&n, &s)

	// Слайс наличия букв
	letters := make([]bool, 4)

	// Запись наличия букв
	for _, ch := range s {
		switch ch {
		case 'a':
			letters[0] = true
		case 'b':
			letters[1] = true
		case 'c':
			letters[2] = true
		case 'd':
			letters[3] = true
		}
	}

	// Проверка на наличие хотя бы одной хорошей подстроки
	for _, exist := range letters {
		if !exist {
			fmt.Println(-1)
			return
		}
	}

	// Минимальная длинна подстроки
	minLen := n

	// Поиск хорошей подстроки
	for i := 0; i < n-3; i++ {

		// Новый слайс наличия букв
		letters := make([]bool, 4)

		// Проход по всем возможным подстрокам
		for j := i; j < n; j++ {

			// Запись букв
			switch s[j] {
			case 'a':
				letters[0] = true
			case 'b':
				letters[1] = true
			case 'c':
				letters[2] = true
			case 'd':
				letters[3] = true
			}

			// Проверка на подстроки
			if letters[0] && letters[1] && letters[2] && letters[3] {

				// Вычисление длинны подстроки
				length := j - i + 1

				// Проверка на размер
				if length < minLen {
					minLen = length
				}

				// Выход из цикла подстрок
				break
			}
		}
	}

	// Вывод длинны минимальной подстроки
	fmt.Println(minLen)

}
