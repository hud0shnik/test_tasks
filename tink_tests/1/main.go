package main

import "fmt"

func main() {

	// Запись роста каждого человека
	var h1, h2, h3, h4 int
	fmt.Scan(&h1, &h2, &h3, &h4)

	// Проверка на построение
	if h1 <= h2 && h2 <= h3 && h3 <= h4 {
		fmt.Println("YES")
	} else if h1 >= h2 && h2 >= h3 && h3 >= h4 {
		fmt.Println("YES")
	} else {
		fmt.Println("NO")
	}

}
