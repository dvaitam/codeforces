package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	m := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &m[i])
	}
	d := make([]int64, n)
	if n > 0 {
		d[n-1] = m[n-1]
		if d[n-1] < 0 {
			d[n-1] = 0
		}
		for i := n - 2; i >= 0; i-- {
			need := d[i+1] - 1
			if need < 0 {
				need = 0
			}
			if m[i] > need {
				d[i] = m[i]
			} else {
				d[i] = need
			}
		}
	}
	var sum int64
	for i := 0; i < n; i++ {
		sum += d[i]
	}
	fmt.Println(sum)
}
