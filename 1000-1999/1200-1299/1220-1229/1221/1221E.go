package main

import (
	"bufio"
	"fmt"
	"os"
)

func canWin(a, b int, s string) bool {
	segs := []int{}
	n := len(s)
	i := 0
	for i < n {
		if s[i] == 'X' {
			i++
			continue
		}
		j := i
		for j < n && s[j] == '.' {
			j++
		}
		segs = append(segs, j-i)
		i = j
	}

	// check for segments Bob can play but Alice cannot
	for _, l := range segs {
		if l >= b && l < a {
			return false
		}
	}

	// helper to validate board after Alice move
	isValid := func(arr []int) bool {
		for _, l := range arr {
			if l >= b && l < a {
				return false
			}
			if l >= 2*b+a {
				return false
			}
		}
		return true
	}

	for idx, L := range segs {
		if L < a {
			continue
		}
		for pos := 0; pos <= L-a; pos++ {
			left := pos
			right := L - pos - a
			if (left >= b && left < a) || (right >= b && right < a) {
				continue
			}
			if left >= 2*b+a || right >= 2*b+a {
				continue
			}
			tmp := []int{}
			for j, l := range segs {
				if j == idx {
					continue
				}
				tmp = append(tmp, l)
			}
			if left > 0 {
				tmp = append(tmp, left)
			}
			if right > 0 {
				tmp = append(tmp, right)
			}
			if !isValid(tmp) {
				continue
			}
			cnt := 0
			for _, l := range tmp {
				if l >= a {
					cnt++
				}
			}
			if cnt%2 == 0 {
				return true
			}
		}
	}
	return false
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var q int
	fmt.Fscan(in, &q)
	for ; q > 0; q-- {
		var a, b int
		fmt.Fscan(in, &a, &b)
		var s string
		fmt.Fscan(in, &s)
		if canWin(a, b, s) {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
