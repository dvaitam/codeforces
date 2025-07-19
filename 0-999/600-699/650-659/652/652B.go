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
	v := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &v[i])
	}
	sort.Ints(v)
	ans := make([]int, n)
	ptr := n - 1
	for i := 1; i < n; i += 2 {
		ans[i] = v[ptr]
		ptr--
	}
	for i := 0; i < n; i += 2 {
		ans[i] = v[ptr]
		ptr--
	}
	for i := 0; i < n; i++ {
		if i > 0 {
			fmt.Fprint(writer, " ")
		}
		fmt.Fprint(writer, ans[i])
	}
}
