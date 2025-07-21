package main

import "fmt"

func main() {
	var s string
	if _, err := fmt.Scan(&s); err != nil {
		return
	}
	for _, c := range s {
		if c == 'H' || c == 'Q' || c == '9' {
			fmt.Println("YES")
			return
		}
	}
	fmt.Println("NO")
}
