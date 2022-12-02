package main

import (
	"fmt"
	"strconv"
)

// Главная функция
func main() {

	// Три ведьмы
	var w1, w2, w3 int

	// Запись данных с консоли
	fmt.Scanln(&w1, &w2, &w3)

	// Вычисление и вывод результатов
	calculate(w1, w2, w3)
	calculate(w2, w1, w3)
	calculate(w3, w2, w1)

}

// Функция вычисления результата
func calculate(w1, w2, w3 int) {

	// Модификатор, который нужно прибавить к силе
	var magika int

	// Если значение силы находится посередине
	if (w1 > w2 && w1 < w3) || (w1 < w2 && w1 > w3) {

		// Ничего менять не надо, вывод результата и конец функции
		fmt.Println("0 2")
		return

	}

	// Если у всех ведьм равное значение силы
	if w1 == w2 && w1 == w3 {

		// Ничего менять не надо, вывод результата и конец функции
		fmt.Println("0 0")
		return

	}

	// Если значение силы совпадает с силой другой ведьмы
	if w1 == w2 || w1 == w3 {

		// Ничего менять не надо, вывод результата и конец функции
		fmt.Println("0 1")
		return

	}

	//Если у ведьмы больше силы, чем у остальных
	if w1 > w2 && w1 > w3 {

		// Поиск сильнейшей ведьмы из оставшихся и вычисление модификатора
		if w2 > w3 {
			magika = w2 - w1
		} else {
			magika = w3 - w1
		}

	} else {

		// Если ведьма самая слабая, поиск слабейшей из осташихся и вычисление модификатора
		if w2 > w3 {
			magika = w3 - w1
		} else {
			magika = w2 - w1
		}

	}

	// Изменение силы
	w1 += magika

	// Высчитывание и вывод последствий изменения силы
	if w1 == w3 && w1 == w2 {
		fmt.Println(strconv.Itoa(module(magika)) + " " + "0")
	} else if w1 != w2 && w1 != w3 {
		fmt.Println(strconv.Itoa(module(magika)) + " " + "2")
	} else {
		fmt.Println(strconv.Itoa(module(magika)) + " " + "1")
	}

}

// Функция высчитывания модуля числа
func module(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
