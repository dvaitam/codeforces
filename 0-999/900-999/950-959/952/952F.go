package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var expr string
	fmt.Fscan(in, &expr)
	expr += "+" // sentinel to flush last number

	var result, cur int64
	sign := int64(1)
	for i := 0; i < len(expr); i++ {
		c := expr[i]
		if c == '+' || c == '-' {
			result += sign * cur
			cur = 0
		}
		if c == '-' {
			sign = -1
		}
		if c == '+' {
			sign = 1
		}
		cur = cur*10 + int64(int(c)-int('0'))
	}
	fmt.Print(result)
}
