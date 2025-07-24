package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &arr[i])
	}

	prefix := make([]int, 0, n)
	for k := 0; k < n; k++ {
		// insert arr[k] into prefix keeping it sorted
		x := arr[k]
		idx := sort.SearchInts(prefix, x)
		prefix = append(prefix, 0)
		copy(prefix[idx+1:], prefix[idx:])
		prefix[idx] = x

		rating := 0
		for _, v := range prefix {
			if v > rating {
				rating++
			} else if v < rating {
				rating--
			}
		}

		if k > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, rating)
	}
	fmt.Fprintln(out)
}
