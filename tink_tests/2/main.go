package main

import "fmt"

func main() {

	// Считывание переменных
	var n, m, k int
	fmt.Scan(&n, &m, &k)

	// Общее количество проверок нужное для решения
	checks := n * k

	// Общее количество проверок / количество проверок, которые может
	// провести один сеньор.
	fmt.Println((checks) / m)

}
