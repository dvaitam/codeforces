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
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	a := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &a[i])
	}

	// previous index with value > a[i]
	lg := make([]int, n+2)
	stack := make([]int, 0, n)
	for i := 1; i <= n; i++ {
		for len(stack) > 0 && a[stack[len(stack)-1]] <= a[i] {
			stack = stack[:len(stack)-1]
		}
		if len(stack) > 0 {
			lg[i] = stack[len(stack)-1]
		} else {
			lg[i] = 0
		}
		stack = append(stack, i)
	}

	// next index with value >= a[i]
	rg := make([]int, n+2)
	stack = stack[:0]
	for i := n; i >= 1; i-- {
		for len(stack) > 0 && a[stack[len(stack)-1]] < a[i] {
			stack = stack[:len(stack)-1]
		}
		if len(stack) > 0 {
			rg[i] = stack[len(stack)-1]
		} else {
			rg[i] = n + 1
		}
		stack = append(stack, i)
	}

	// boundaries where bits not present in a[i] appear on the left/right
	leftF := make([]int, n+2)
	last := make([]int, 31)
	for i := 1; i <= n; i++ {
		lf := 0
		x := a[i]
		for b := 0; b < 31; b++ {
			if (x>>b)&1 == 0 {
				if last[b] > lf {
					lf = last[b]
				}
			} else {
				last[b] = i
			}
		}
		leftF[i] = lf
	}

	rightF := make([]int, n+2)
	nxt := make([]int, 31)
	for b := 0; b < 31; b++ {
		nxt[b] = n + 1
	}
	for i := n; i >= 1; i-- {
		rf := n + 1
		x := a[i]
		for b := 0; b < 31; b++ {
			if (x>>b)&1 == 0 {
				if nxt[b] < rf {
					rf = nxt[b]
				}
			} else {
				nxt[b] = i
			}
		}
		rightF[i] = rf
	}

	var equal int64
	for i := 1; i <= n; i++ {
		L := lg[i]
		if leftF[i] > L {
			L = leftF[i]
		}
		L++
		if L > i {
			continue
		}
		R := rg[i]
		if rightF[i] < R {
			R = rightF[i]
		}
		R--
		if R < i {
			continue
		}
		cl := i - L + 1
		cr := R - i + 1
		equal += int64(cl) * int64(cr)
	}
	equal -= int64(n) // exclude length-1 segments

	totalPairs := int64(n) * int64(n-1) / 2
	result := totalPairs - equal
	fmt.Fprintln(writer, result)
}
