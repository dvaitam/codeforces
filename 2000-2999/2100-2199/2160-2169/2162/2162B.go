package main

import (
	"bufio"
	"fmt"
	"os"
)

func isNonDecreasing(seq []byte) bool {
	for i := 1; i < len(seq); i++ {
		if seq[i-1] > seq[i] {
			return false
		}
	}
	return true
}

func isPalindrome(seq []byte) bool {
	for l, r := 0, len(seq)-1; l < r; l, r = l+1, r-1 {
		if seq[l] != seq[r] {
			return false
		}
	}
	return true
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		var s string
		fmt.Fscan(in, &s)
		bytes := []byte(s)

		found := false
		var answer []int
		limit := 1 << n
		for mask := 0; mask < limit && !found; mask++ {
			var subseq []byte
			var indices []int
			for i := 0; i < n; i++ {
				if mask&(1<<i) != 0 {
					subseq = append(subseq, bytes[i])
					indices = append(indices, i+1)
				}
			}
			if !isNonDecreasing(subseq) {
				continue
			}

			x := make([]byte, 0, n-len(subseq))
			for i := 0; i < n; i++ {
				if mask&(1<<i) == 0 {
					x = append(x, bytes[i])
				}
			}
			if isPalindrome(x) {
				found = true
				answer = indices
			}
		}

		if !found {
			fmt.Fprintln(out, -1)
			continue
		}

		fmt.Fprintln(out, len(answer))
		if len(answer) == 0 {
			fmt.Fprintln(out)
		} else {
			for i, v := range answer {
				if i > 0 {
					fmt.Fprint(out, " ")
				}
				fmt.Fprint(out, v)
			}
			fmt.Fprintln(out)
		}
	}
}
