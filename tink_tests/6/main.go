package main

import (
	"fmt"
)

func main() {

	// Получение входных данных
	var n, s int
	fmt.Scan(&n, &s)
	l := make([]int, n)
	r := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Scan(&l[i], &r[i])
	}

	// Левая и правая границы поиска и результат
	left := 0
	right := s
	result := 0

	// Бинарный поиск
	for left <= right {

		// Вычисление медианы
		mid := (left + right) / 2

		// Количество учеников с баллами меньше медианы
		count := make([]int, s+1)
		for i := 0; i < n; i++ {
			if l[i] > mid {
				continue
			}
			count[min(r[i], mid)]++
		}

		sum := 0
		median := 0

		// Проход в обратном порядке
		for i := s; i >= 0; i-- {
			sum += count[i]
			if sum > n/2 {
				median = i
				break
			}
			if sum == n/2 && n%2 == 0 {
				for j := i - 1; j >= 0; j-- {
					if count[j] > 0 {
						median += j
						break
					}
				}
				median /= 2
				break
			}
		}

		if median > result {
			result = median
		}

		if sum > n/2 {
			right = mid - 1
		} else {
			left = mid + 1
		}
	}

	fmt.Println(result)
}

// Функция нахождения минимального числа
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
