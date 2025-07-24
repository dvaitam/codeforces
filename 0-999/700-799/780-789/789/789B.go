package main

import (
	"bufio"
	"fmt"
	"os"
)

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var b1, q, l int64
	var m int
	if _, err := fmt.Fscan(in, &b1, &q, &l, &m); err != nil {
		return
	}
	bad := make(map[int64]bool, m)
	for i := 0; i < m; i++ {
		var x int64
		fmt.Fscan(in, &x)
		bad[x] = true
	}

	if abs(b1) > l {
		fmt.Println(0)
		return
	}

	// handle q == 0 separately
	if q == 0 {
		if !bad[b1] {
			if bad[0] {
				fmt.Println(1)
			} else {
				fmt.Println("inf")
			}
		} else {
			if bad[0] {
				fmt.Println(0)
			} else {
				fmt.Println("inf")
			}
		}
		return
	}

	// q == 1 => constant b1
	if q == 1 {
		if bad[b1] {
			fmt.Println(0)
		} else {
			fmt.Println("inf")
		}
		return
	}

	// q == -1 => alternate between b1 and -b1
	if q == -1 {
		if bad[b1] && bad[-b1] {
			fmt.Println(0)
		} else {
			fmt.Println("inf")
		}
		return
	}

	count := 0
	cur := b1
	for abs(cur) <= l {
		if !bad[cur] {
			count++
		}
		cur *= q
	}
	fmt.Println(count)
}
