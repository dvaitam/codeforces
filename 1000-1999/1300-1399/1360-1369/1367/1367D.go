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

	var q int
	fmt.Fscan(in, &q)
	for ; q > 0; q-- {
		var s string
		fmt.Fscan(in, &s)
		var m int
		fmt.Fscan(in, &m)
		b := make([]int, m)
		for i := 0; i < m; i++ {
			fmt.Fscan(in, &b[i])
		}

		freq := make([]int, 26)
		for _, ch := range s {
			freq[int(ch-'a')]++
		}

		res := make([]byte, m)
		used := make([]bool, m)
		remaining := m
		ch := 25
		for remaining > 0 {
			zeros := make([]int, 0)
			for i := 0; i < m; i++ {
				if !used[i] && b[i] == 0 {
					zeros = append(zeros, i)
				}
			}
			for ch >= 0 && freq[ch] < len(zeros) {
				ch--
			}
			for _, pos := range zeros {
				res[pos] = byte('a' + ch)
				used[pos] = true
			}
			freq[ch] -= len(zeros)
			ch--
			remaining -= len(zeros)
			for i := 0; i < m; i++ {
				if used[i] {
					continue
				}
				sum := 0
				for _, pos := range zeros {
					if i > pos {
						sum += i - pos
					} else {
						sum += pos - i
					}
				}
				b[i] -= sum
			}
		}
		fmt.Fprintln(out, string(res))
	}
}
