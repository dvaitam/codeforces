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
	L := 1
	R := n - 1
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		if u > v {
			u, v = v, u
		}
		if u > L {
			L = u
		}
		if v-1 < R {
			R = v - 1
		}
	}
	if L > R {
		fmt.Println(0)
	} else {
		fmt.Println(R - L + 1)
	}
}
