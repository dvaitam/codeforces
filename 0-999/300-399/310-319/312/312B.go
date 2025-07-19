package main

import (
	"fmt"
	"os"
)

func main() {
	var a, b, c, d float64
	for {
		_, err := fmt.Fscan(os.Stdin, &a, &b, &c, &d)
		if err != nil {
			break
		}
		p1 := a / b
		p2 := c / d
		ans := p1 * (1 / (1 - ((1 - p2) * (1 - p1))))
		fmt.Printf("%.10f\n", ans)
	}
}
