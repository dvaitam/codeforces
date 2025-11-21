package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	var c, d int64
	fmt.Fscan(in, &n, &c, &d)

	times := make([]int64, n)
	person := make([]byte, n)
	for i := 0; i < n; i++ {
		var ti int64
		var ch string
		fmt.Fscan(in, &ti, &ch)
		times[i] = ti
		person[i] = ch[0]
	}
	var tFinal int64
	fmt.Fscan(in, &tFinal)

	nextW := make([]int64, n)
	nextP := make([]int64, n)
	nextWTime := tFinal
	nextPTime := tFinal

	for i := n - 1; i >= 0; i-- {
		nextW[i] = nextWTime
		nextP[i] = nextPTime
		if person[i] == 'W' {
			nextWTime = times[i]
		} else {
			nextPTime = times[i]
		}
	}

	var ans int64
	for i := 0; i < n; i++ {
		var release int64
		if person[i] == 'W' {
			release = nextP[i]
		} else {
			release = nextW[i]
		}
		storeCost := c * (release - times[i])
		if storeCost < d {
			ans += storeCost
		} else {
			ans += d
		}
	}
	fmt.Fprintln(out, ans)
}
