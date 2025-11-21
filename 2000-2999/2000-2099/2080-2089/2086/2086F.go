package main

import (
	"bufio"
	"fmt"
	"os"
)

type swap struct {
	l int
	r int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	if n <= 0 {
		return
	}

	var s string
	fmt.Fscan(in, &s)
	bytes := []byte(s)
	for len(bytes) < n {
		var chunk string
		if _, err := fmt.Fscan(in, &chunk); err != nil {
			break
		}
		bytes = append(bytes, chunk...)
	}
	if len(bytes) > n {
		bytes = bytes[:n]
	}

	cntA := 0
	for _, ch := range bytes {
		if ch == 'a' {
			cntA++
		}
	}
	cntB := n - cntA

	target := make([]byte, n)
	left, right := 0, n-1
	center := byte('a')
	if cntA%2 == 1 {
		center = 'a'
	} else {
		center = 'b'
	}
	for i := 0; i < cntA/2; i++ {
		target[left] = 'a'
		target[right] = 'a'
		left++
		right--
	}
	for i := 0; i < cntB/2; i++ {
		target[left] = 'b'
		target[right] = 'b'
		left++
		right--
	}
	target[n/2] = center

	current := append([]byte(nil), bytes...)
	var swaps []swap
	for i := 0; i < n; i++ {
		if current[i] == target[i] {
			continue
		}
		j := i + 1
		for ; j < n && current[j] != target[i]; j++ {
		}
		if j == n {
			return
		}
		current[i], current[j] = current[j], current[i]
		swaps = append(swaps, swap{l: i + 1, r: j + 1})
	}

	idx := 0
	for step := 1; step <= n; step++ {
		if idx < len(swaps) {
			ready := swaps[idx].r
			if swaps[idx].l > ready {
				ready = swaps[idx].l
			}
			if step >= ready {
				fmt.Fprintln(out, swaps[idx].l, swaps[idx].r)
				idx++
				continue
			}
		}
		fmt.Fprintln(out, 0, 0)
	}
}
