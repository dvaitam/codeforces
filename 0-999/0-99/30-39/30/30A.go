package main

import (
	"fmt"
	"math"
	"os"
)

func main() {
	var a, b, x int
	if _, err := fmt.Fscan(os.Stdin, &a, &b, &x); err != nil {
		return
	}
	for i := -1000; i <= 1000; i++ {
		v := float64(a) * math.Pow(float64(i), float64(x))
		if math.Abs(v-float64(b)) < 0.1 {
			fmt.Println(i)
			return
		}
	}
	fmt.Println("No solution")
}
