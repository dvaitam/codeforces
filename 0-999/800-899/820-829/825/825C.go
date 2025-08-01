package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	var k int64
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}

	sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })

	ans := 0
	for _, ai := range a {
		need := (ai + 1) / 2
		for k < need {
			k *= 2
			ans++
		}
		if ai > k {
			k = ai
		}
	}

	fmt.Fprintln(writer, ans)
}
