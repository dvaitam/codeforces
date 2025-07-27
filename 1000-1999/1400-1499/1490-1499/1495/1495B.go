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

	var n int
	fmt.Fscan(reader, &n)
	p := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &p[i])
	}

	l := make([]int, n)
	r := make([]int, n)
	for i := 0; i < n; i++ {
		l[i] = 1
		r[i] = 1
	}
	for i := 1; i < n; i++ {
		if p[i-1] < p[i] {
			l[i] = l[i-1] + 1
		}
	}
	for i := n - 2; i >= 0; i-- {
		if p[i+1] < p[i] {
			r[i] = r[i+1] + 1
		}
	}

	mxL := 0
	for _, v := range l {
		if v > mxL {
			mxL = v
		}
	}
	cntL := 0
	for _, v := range l {
		if v == mxL {
			cntL++
		}
	}

	mxR := 0
	for _, v := range r {
		if v > mxR {
			mxR = v
		}
	}
	cntR := 0
	for _, v := range r {
		if v == mxR {
			cntR++
		}
	}

	ans := 0
	if mxL == mxR && cntL == 1 && cntR == 1 {
		for i := 0; i < n; i++ {
			if l[i] == r[i] && l[i] == mxL && l[i]%2 == 1 {
				ans++
			}
		}
	}
	fmt.Fprintln(writer, ans)
}
