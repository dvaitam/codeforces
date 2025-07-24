package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	l := make([]int, n)
	r := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &l[i], &r[i])
	}
	total := 0
	for a := 1; a <= m; a++ {
		for b := a; b <= m; b++ {
			valid := false
			ok := true
			for i := 0; i < n; i++ {
				start := a
				if l[i] > start {
					start = l[i]
				}
				end := b
				if r[i] < end {
					end = r[i]
				}
				length := end - start + 1
				if start > end {
					length = 0
				}
				if length > 0 {
					valid = true
				}
				if length%2 == 0 {
					if length != 0 {
						ok = false
						break
					}
				}
			}
			if ok && valid {
				total += b - a + 1
			}
		}
	}
	fmt.Println(total)
}
