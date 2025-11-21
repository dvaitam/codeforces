package main

import (
	"bufio"
	"fmt"
	"os"
)

func reachable(A, B string) bool {
	n := len(A)
	if n != len(B) {
		return false
	}
	start := 0
	target := 0
	for i := 0; i < n; i++ {
		if A[i] == '1' {
			start |= 1 << i
		}
		if B[i] == '1' {
			target |= 1 << i
		}
	}
	if start == target {
		return true
	}
	// exhaustive search over all states; intended for small n (used as reference)
	seen := make(map[int]struct{})
	queue := []int{start}
	seen[start] = struct{}{}
	for len(queue) > 0 {
		mask := queue[0]
		queue = queue[1:]
		for i := 0; i < n; i++ {
			if mask&(1<<i) == 0 {
				continue
			}
			next := mask
			if i > 0 {
				next ^= 1 << (i - 1)
			}
			if i+1 < n {
				next ^= 1 << (i + 1)
			}
			if _, ok := seen[next]; ok {
				continue
			}
			if next == target {
				return true
			}
			seen[next] = struct{}{}
			queue = append(queue, next)
		}
	}
	return false
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var N int
		var A, B string
		fmt.Fscan(in, &N)
		fmt.Fscan(in, &A)
		fmt.Fscan(in, &B)
		if reachable(A, B) {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
