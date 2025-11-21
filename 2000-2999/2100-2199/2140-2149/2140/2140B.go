package main

import (
	"bufio"
	"fmt"
	"os"
)

func dac(len int) int64 {
	res := int64(1)
	for i := 0; i < len; i++ {
		res *= 10
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var x int64
		fmt.Fscan(in, &x)
		length := 0
		tmp := x
		for tmp > 0 {
			length++
			tmp /= 10
		}
		if length == 0 {
			length = 1
		}
		mod := dac(length)
		y := mod - (x % mod)
		if y == mod {
			y = 0
		}
		fmt.Fprintln(out, y)
	}
}

