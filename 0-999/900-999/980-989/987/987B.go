package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var x, y float64
	if _, err := fmt.Fscan(in, &x, &y); err != nil {
		return
	}
	if x == y {
		fmt.Println("=")
		return
	}
	if x == 1 {
		if y == 1 {
			fmt.Println("=")
		} else {
			fmt.Println("<")
		}
		return
	}
	if y == 1 {
		fmt.Println(">")
		return
	}
	a := y * math.Log(x)
	b := x * math.Log(y)
	if math.Abs(a-b) <= 1e-10 {
		fmt.Println("=")
	} else if a < b {
		fmt.Println("<")
	} else {
		fmt.Println(">")
	}
}
