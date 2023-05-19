package main

import (
	"fmt"
)

// Функция вычисления количества вхождений числа в слайсе
func count(s []int, val int) int {
	result := 0
	for _, num := range s {
		if num == val {
			result++
		}
	}
	return result
}

// Функция записи словаря в слайс
func getMapValues(m map[int]int) []int {
	result := make([]int, len(m))
	i := 0
	for _, val := range m {
		result[i] = val
		i++
	}
	return result
}

// Функция проверки на скучность
func isBoring(s []int) bool {

	// Запись наличия чисел в карту
	m := make(map[int]bool)
	for _, num := range s {
		m[num] = true
	}

	// Проверка случая с картой в две записи (т.е все числа кроме одного
	// встречаются одинаковое количество раз)
	if len(m) == 2 {
		minVal, maxVal := 0, 0
		for num := range m {
			if minVal == 0 || num < minVal {
				minVal = num
			}
			if maxVal == 0 || num > maxVal {
				maxVal = num
			}
		}

		// Если одно из чисел встречается только один раз, то слайс является скучным
		if count(s, minVal) == 1 || count(s, maxVal) == 1 {
			return true
		}

	}

	// Проверка на слайсы, где все числа встречаются лиш раз
	if count(s, 1) == len(s) {
		return true
	}

	return false
}

func main() {

	// Получение входных данных
	var n int
	fmt.Scan(&n)
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Scan(&a[i])
	}

	// Результат и карта количества вхождений чисел
	result := 1
	m := make(map[int]int)

	// Перебор элементов слайса
	for i := 0; i < n; i++ {

		// Запись числа в  карту
		m[a[i]]++

		// Вычисление длинны скучного префикса
		if isBoring(getMapValues(m)) {
			result = i + 1
		}

	}

	// Вывод результата
	fmt.Println(result)

}
