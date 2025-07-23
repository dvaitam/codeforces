package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, s int
	if _, err := fmt.Fscan(in, &n, &s); err != nil {
		return
	}
	times := make([]int, n)
	for i := 0; i < n; i++ {
		var h, m int
		fmt.Fscan(in, &h, &m)
		times[i] = h*60 + m
	}

	if times[0] >= s+1 {
		fmt.Println("0 0")
		return
	}

	for i := 0; i < n-1; i++ {
		if times[i+1]-times[i] >= 2*s+2 {
			t := times[i] + s + 1
			fmt.Printf("%d %d\n", t/60, t%60)
			return
		}
	}

	t := times[n-1] + s + 1
	fmt.Printf("%d %d\n", t/60, t%60)
}
