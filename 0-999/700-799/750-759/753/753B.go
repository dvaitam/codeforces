package main

import "fmt"

func main() {
	var s string
	if _, err := fmt.Scan(&s); err != nil {
		return
	}
	fmt.Println(s)
}
