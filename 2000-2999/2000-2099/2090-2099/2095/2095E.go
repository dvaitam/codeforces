package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	var p, k int64
	if _, err := fmt.Fscan(in, &n, &p, &k); err != nil {
		return
	}

	a := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}

	sq := make([]int64, n)
	for i := 0; i < n; i++ {
		sq[i] = a[i] * a[i]
	}

	var ans int64
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			x := (a[i] ^ a[j]) % p
			y := (sq[i] ^ sq[j]) % p
			if (x*y)%p == k {
				ans++
			}
		}
	}

	fmt.Println(ans)
}
