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

	var n1, n2, n3 int
	if _, err := fmt.Fscan(in, &n1, &n2, &n3); err != nil {
		return
	}

	sums := make([]int64, 3)
	mins := []int64{1<<63 - 1, 1<<63 - 1, 1<<63 - 1}
	lens := []int{n1, n2, n3}
	for i := 0; i < 3; i++ {
		for j := 0; j < lens[i]; j++ {
			var v int64
			fmt.Fscan(in, &v)
			sums[i] += v
			if v < mins[i] {
				mins[i] = v
			}
		}
	}

	total := sums[0] + sums[1] + sums[2]
	cand := sums[0]
	if sums[1] < cand {
		cand = sums[1]
	}
	if sums[2] < cand {
		cand = sums[2]
	}
	pair := mins[0] + mins[1]
	if pair < cand {
		cand = pair
	}
	pair = mins[0] + mins[2]
	if pair < cand {
		cand = pair
	}
	pair = mins[1] + mins[2]
	if pair < cand {
		cand = pair
	}
	ans := total - 2*cand
	fmt.Fprintln(out, ans)
}
