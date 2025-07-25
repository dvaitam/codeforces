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

	var T int
	if _, err := fmt.Fscan(reader, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(reader, &n)
		var s, t string
		fmt.Fscan(reader, &s)
		fmt.Fscan(reader, &t)

		// prefix count of ones in s
		prefS := make([]int, n+1)
		for i := 0; i < n; i++ {
			prefS[i+1] = prefS[i]
			if s[i] == '1' {
				prefS[i+1]++
			}
		}

		// compute B1 globally
		B1 := make([]bool, n)
		for i := 0; i < n; i++ {
			if t[i] == '1' {
				B1[i] = true
			}
			if i > 0 && i+1 < n && s[i-1] == '0' && s[i+1] == '0' {
				B1[i] = true
			}
		}

		// add[i] indicates that position i in s becomes 1 due to rule 2
		add := make([]int, n)
		for i := 1; i+1 < n; i++ {
			if s[i] == '0' && B1[i-1] && B1[i+1] {
				add[i] = 1
			}
		}

		prefAdd := make([]int, n+1)
		for i := 0; i < n; i++ {
			prefAdd[i+1] = prefAdd[i] + add[i]
		}

		var q int
		fmt.Fscan(reader, &q)
		for ; q > 0; q-- {
			var l, r int
			fmt.Fscan(reader, &l, &r)
			l--
			r--
			ones := prefS[r+1] - prefS[l]
			extra := 0
			if r-l >= 2 {
				extra = prefAdd[r] - prefAdd[l+1]
				j1 := l + 1
				old1 := add[j1]
				new1 := 0
				if s[j1] == '0' {
					left := t[l] == '1'
					right := B1[j1+1]
					if j1+1 == r {
						right = t[r] == '1'
					}
					if left && right {
						new1 = 1
					}
				}
				extra += new1 - old1
				j2 := r - 1
				if j2 != j1 {
					old2 := add[j2]
					new2 := 0
					if s[j2] == '0' {
						left := B1[j2-1]
						if j2-1 == l {
							left = t[l] == '1'
						}
						right := t[r] == '1'
						if left && right {
							new2 = 1
						}
					}
					extra += new2 - old2
				}
			}
			fmt.Fprintln(writer, ones+extra)
		}
	}
}
