package main

import (
	"fmt"
)

func main() {

	// Получение входных данных
	var n int
	fmt.Scan(&n)
	balance := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Scan(&balance[i])
	}

	// Карта(множество) нормальных отрезков (пустая структура только в целях экономии памяти)
	normals := make(map[string]struct{})

	// Перебор всех возможных отрезков
	for start := 0; start < n-1; start++ {
		for end := start + 2; end <= n; end++ {

			// Сумма всех элементов в отрезке
			sum := 0
			for i := start; i < end; i++ {
				sum += balance[i]
			}

			// Проверка на разумность отрезка
			if sum == 0 {
				// Запись в карту(множество) нормальных отрезков
				for i := 0; i <= start; i++ {
					for j := end - 1; j < n; j++ {
						// Добавление записи формата "начало отрезка, конец"
						normals[fmt.Sprintf("%d,%d", i, j)] = struct{}{}
					}
				}
				break
			}
		}
	}

	// Вывод результата
	fmt.Println(len(normals))

}
