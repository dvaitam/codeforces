package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var s string
	if _, err := fmt.Fscan(reader, &s); err != nil {
		return
	}
	n := len(s)
	// nxt[i] is next index in the sequence starting from i
	nxt := make([]int, n+2)
	vst := make([]bool, n+2)
	// stacks of positions ending with '0' or '1'
	a0 := make([]int, 0, n)
	a1 := make([]int, 0, n)
	m := 0
	// 1-based positions
	for i := 1; i <= n; i++ {
		c := s[i-1]
		if c == '0' {
			if len(a1) > 0 {
				u := a1[len(a1)-1]
				a1 = a1[:len(a1)-1]
				nxt[u] = i
				a0 = append(a0, i)
			} else {
				a0 = append(a0, i)
				m++
			}
		} else if c == '1' {
			if len(a0) > 0 {
				u := a0[len(a0)-1]
				a0 = a0[:len(a0)-1]
				nxt[u] = i
				a1 = append(a1, i)
			} else {
				fmt.Fprintln(writer, -1)
				return
			}
		}
	}
	// any sequence ending with '1' is invalid
	if len(a1) > 0 {
		fmt.Fprintln(writer, -1)
		return
	}
	// output
	fmt.Fprintln(writer, m)
	for i := 1; i <= n; i++ {
		if vst[i] {
			continue
		}
		// build sequence starting at i
		seq := []int{}
		for u := i; u != 0; u = nxt[u] {
			seq = append(seq, u)
			vst[u] = true
		}
		// print sequence
		fmt.Fprint(writer, len(seq))
		for _, x := range seq {
			fmt.Fprint(writer, " ", x)
		}
		fmt.Fprintln(writer)
	}
}
