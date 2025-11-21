package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func isPerfectOddSquare(x int) bool {
	if x <= 0 {
		return false
	}
	r := int(math.Round(math.Sqrt(float64(x))))
	return r*r == x && r%2 == 1
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		sum := 0
		happy := 0
		for i := 0; i < n; i++ {
			var a int
			fmt.Fscan(in, &a)
			sum += a
			if isPerfectOddSquare(sum) {
				happy++
			}
		}
		fmt.Fprintln(out, happy)
	}
}
