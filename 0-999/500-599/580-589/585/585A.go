package main

import (
	"bufio"
	"fmt"
	"os"
)

// This program solves the problem described in problemA.txt (Gennady the Dentist).
// It simulates the queue of children applying the rules from the statement.
func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	v := make([]int, n)
	d := make([]int, n)
	p := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &v[i], &d[i], &p[i])
	}

	removed := make([]bool, n)
	res := []int{}
	for i := 0; i < n; i++ {
		if removed[i] || p[i] < 0 {
			continue
		}
		res = append(res, i+1)
		// Cry in the dentist's office with decreasing volume
		cur := v[i]
		for j := i + 1; j < n && cur > 0; j++ {
			if !removed[j] {
				p[j] -= cur
			}
			cur--
		}
		// Process chain reaction of children leaving the queue
		for j := i + 1; j < n; j++ {
			if removed[j] {
				continue
			}
			if p[j] < 0 {
				removed[j] = true
				for k := j + 1; k < n; k++ {
					if !removed[k] {
						p[k] -= d[j]
					}
				}
			}
		}
	}

	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	fmt.Fprintln(out, len(res))
	for i, x := range res {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, x)
	}
	if len(res) > 0 {
		fmt.Fprintln(out)
	}
}
