package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func isPerfectSquare(x int) bool {
	if x < 0 {
		return false
	}
	r := int(math.Sqrt(float64(x)))
	return r*r == x
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	ans := -1000001
	for i := 0; i < n; i++ {
		var v int
		fmt.Fscan(reader, &v)
		if !isPerfectSquare(v) && v > ans {
			ans = v
		}
	}
	fmt.Fprintln(writer, ans)
}
