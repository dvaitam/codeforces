package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var w, l int
	if _, err := fmt.Fscan(in, &w, &l); err != nil {
		return
	}

	stones := make([]int64, w)
	for i := 1; i <= w-1; i++ {
		fmt.Fscan(in, &stones[i])
	}

	prefix := make([]int64, w)
	for i := 1; i <= w-1; i++ {
		prefix[i] = prefix[i-1] + stones[i]
	}

	ans := int64(1<<63 - 1)
	for start := 1; start <= w-l; start++ {
		sum := prefix[start+l-1] - prefix[start-1]
		if sum < ans {
			ans = sum
		}
	}

	fmt.Fprintln(out, ans)
}
