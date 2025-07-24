package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int64
	var c1, c2 int64
	if _, err := fmt.Fscan(reader, &n, &c1, &c2); err != nil {
		return
	}
	var s string
	fmt.Fscan(reader, &s)

	// count adults represented by '1'
	adults := int64(strings.Count(s, "1"))
	maxGroups := adults
	if maxGroups > n {
		maxGroups = n
	}

	best := int64(1<<63 - 1)
	for k := int64(1); k <= maxGroups; k++ {
		q := n / k
		r := n % k
		// r groups of size q+1 and k-r groups of size q
		cost := r*(c1+c2*q*q) + (k-r)*(c1+c2*(q-1)*(q-1))
		if cost < best {
			best = cost
		}
	}
	fmt.Fprintln(writer, best)
}
