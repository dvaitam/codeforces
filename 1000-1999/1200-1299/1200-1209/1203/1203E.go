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
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}

	sort.Ints(a)

	const maxW = 150001
	used := make([]bool, maxW+2)
	ans := 0

	for _, w := range a {
		if w-1 > 0 && !used[w-1] {
			used[w-1] = true
			ans++
		} else if w <= maxW && !used[w] {
			used[w] = true
			ans++
		} else if w+1 <= maxW && !used[w+1] {
			used[w+1] = true
			ans++
		}
	}

	fmt.Fprintln(writer, ans)
}
