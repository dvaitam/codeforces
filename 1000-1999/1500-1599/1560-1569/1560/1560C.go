package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var k int64
		fmt.Fscan(reader, &k)
		r, c := position(k)
		fmt.Fprintf(writer, "%d %d\n", r, c)
	}
}

func position(k int64) (int64, int64) {
	root := int64(math.Sqrt(float64(k)))
	if root*root < k {
		root++
	}
	prev := (root - 1) * (root - 1)
	diff := k - prev
	if diff <= root {
		return diff, root
	}
	return root, root*2 - diff
}
