package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func lis(nums []int) []int {
	n := len(nums)
	if n == 0 {
		return nil
	}

	tailVals := make([]int, 0, n)
	tailIdx := make([]int, 0, n)
	prev := make([]int, n)
	for i := range prev {
		prev[i] = -1
	}

	for i, v := range nums {
		pos := sort.Search(len(tailVals), func(j int) bool { return tailVals[j] >= v })
		if pos > 0 {
			prev[i] = tailIdx[pos-1]
		}
		if pos == len(tailVals) {
			tailVals = append(tailVals, v)
			tailIdx = append(tailIdx, i)
		} else {
			tailVals[pos] = v
			tailIdx[pos] = i
		}
	}

	idx := tailIdx[len(tailVals)-1]
	seq := make([]int, 0, len(tailVals))
	for idx != -1 {
		seq = append(seq, idx)
		idx = prev[idx]
	}
	for i, j := 0, len(seq)-1; i < j; i, j = i+1, j-1 {
		seq[i], seq[j] = seq[j], seq[i]
	}
	return seq
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		m := n*n + 1
		nums := make([]int, m)
		for i := 0; i < m; i++ {
			fmt.Fscan(in, &nums[i])
		}

		target := n + 1
		seq := lis(nums)
		if len(seq) < target {
			neg := make([]int, m)
			for i, v := range nums {
				neg[i] = -v
			}
			seq = lis(neg)
		}

		if len(seq) > target {
			seq = seq[len(seq)-target:]
		}

		for i, idx := range seq {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, idx+1)
		}
		fmt.Fprint(out, "\n")
	}
}
