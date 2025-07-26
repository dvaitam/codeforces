package main

import (
	"bufio"
	"fmt"
	"os"
)

func precompute(p []int) []int {
	n := len(p)
	diff := make([]int, n+1)
	for j, v := range p {
		if v == n {
			continue
		}
		start := (j + 1) % n
		end := (j - v + n) % n
		if start <= end {
			diff[start]++
			diff[end+1]--
		} else {
			diff[start]++
			diff[n]--
			diff[0]++
			diff[end+1]--
		}
	}
	res := make([]int, n)
	cur := 0
	for i := 0; i < n; i++ {
		cur += diff[i]
		res[i] = cur
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	fmt.Fscan(reader, &n)
	p := make([]int, n)
	for i := range p {
		fmt.Fscan(reader, &p[i])
	}
	counts := precompute(p)
	rev := make([]int, n)
	for i, v := range p {
		rev[n-1-i] = v
	}
	countsRev := precompute(rev)

	var q int
	fmt.Fscan(reader, &q)

	offset := 0
	reversed := false
	if reversed {
		fmt.Fprintln(writer, countsRev[offset])
	} else {
		fmt.Fprintln(writer, counts[offset])
	}

	for ; q > 0; q-- {
		var t int
		fmt.Fscan(reader, &t)
		var k int
		if t == 1 || t == 2 {
			fmt.Fscan(reader, &k)
			k %= n
		}
		if t == 1 {
			offset = (offset + k) % n
		} else if t == 2 {
			offset = (offset - k) % n
			if offset < 0 {
				offset += n
			}
		} else {
			reversed = !reversed
			offset = (n - offset) % n
		}
		if reversed {
			fmt.Fprintln(writer, countsRev[offset])
		} else {
			fmt.Fprintln(writer, counts[offset])
		}
	}
}
