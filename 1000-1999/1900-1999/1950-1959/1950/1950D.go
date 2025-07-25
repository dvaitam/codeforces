package main

import (
	"bufio"
	"fmt"
	"os"
)

const maxN = 100000

var reachable [maxN + 1]bool

func isBinaryDecimal(x int) bool {
	for x > 0 {
		d := x % 10
		if d != 0 && d != 1 {
			return false
		}
		x /= 10
	}
	return true
}

func precompute() {
	// generate all binary decimals up to maxN
	nums := make([]int, 0)
	for i := 1; i <= maxN; i++ {
		if isBinaryDecimal(i) {
			nums = append(nums, i)
		}
	}
	reachable[1] = true
	queue := []int{1}
	for len(queue) > 0 {
		x := queue[0]
		queue = queue[1:]
		for _, b := range nums {
			if b == 1 {
				continue
			}
			y := x * b
			if y <= maxN && !reachable[y] {
				reachable[y] = true
				queue = append(queue, y)
			}
		}
	}
}

func main() {
	precompute()
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		if reachable[n] {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
