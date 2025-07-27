package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func isSquare(x int) bool {
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

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		ok := false
		if n%2 == 0 && isSquare(n/2) {
			ok = true
		}
		if !ok && n%4 == 0 && isSquare(n/4) {
			ok = true
		}
		if ok {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
