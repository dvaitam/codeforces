package main

import (
	"fmt"
)

func main() {
	var n, m int64
	if _, err := fmt.Scan(&n, &m); err != nil {
		return
	}
	if n > m {
		n, m = m, n
	}
	if n == 1 {
		res := m
		r := m % 6
		if r <= 3 {
			res -= r
		} else if r == 4 {
			res -= 2
		} else if r == 5 {
			res -= 1
		}
		fmt.Println(res)
		return
	}
	if n == 2 {
		switch m {
		case 2:
			fmt.Println(0)
		case 3, 7:
			fmt.Println(2*m - 2)
		default:
			fmt.Println(2 * m)
		}
		return
	}
	total := n * m
	if (n&1) == 1 && (m&1) == 1 {
		total--
	}
	fmt.Println(total)
}
