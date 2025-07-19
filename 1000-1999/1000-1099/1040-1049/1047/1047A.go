package main

import (
	"fmt"
)

func main() {
	var n int
	if _, err := fmt.Scan(&n); err != nil {
		return
	}
	d := n / 3
	r := n % 3
	a := d
	b := d
	c := d + r
	if a%3 == 0 && c%3 != 0 {
		a += 1
		b -= 1
	} else if a%3 == 0 && c%3 == 0 {
		a -= 1
		b -= 1
		c += 2
	} else if a%3 != 0 && c%3 == 0 {
		if (a+1)%3 == 0 {
			a += 2
			c -= 2
		} else {
			a += 1
			c -= 1
		}
	}
	fmt.Println(a, b, c)
}
