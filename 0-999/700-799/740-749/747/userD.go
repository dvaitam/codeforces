package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	t := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &t[i])
	}

	negIdx := make([]int, 0)
	for i, v := range t {
		if v < 0 {
			negIdx = append(negIdx, i)
		}
	}
	m := len(negIdx)
	if m == 0 {
		fmt.Println(0)
		return
	}
	if k < m {
		fmt.Println(-1)
		return
	}

	spare := k - m
	segments := 1
	gaps := make([]int, 0)

	for i := 0; i < m-1; i++ {
		gap := negIdx[i+1] - negIdx[i] - 1
		if gap > 0 {
			gaps = append(gaps, gap)
			segments++
		}
	}

	sort.Ints(gaps)
	for _, g := range gaps {
		if g <= spare {
			spare -= g
			segments--
		} else {
			break
		}
	}

	ans := 2 * segments
	lastNeg := negIdx[m-1]
	if lastNeg == n-1 {
		ans--
	} else {
		tail := n - 1 - lastNeg
		if spare >= tail {
			ans--
		}
	}

	fmt.Println(ans)
}