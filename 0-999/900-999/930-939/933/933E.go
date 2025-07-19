package main

import (
	"bufio"
	"fmt"
	"os"
)

func max(x, y int64) int64 {
	if x > y {
		return x
	}
	return y
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	a := make([]int64, n+2)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &a[i])
	}

	f := make([]int64, n+2)
	odd := make([]bool, n+2)
	for i := 1; i <= n; i++ {
		// option 1: operate at i
		fo := a[i]
		if i-2 >= 0 {
			fo += f[i-2]
		}
		// option 2: operate at i-1,i
		var fe int64
		if i-3 >= 0 {
			fe = f[i-3] + max(a[i], a[i-1])
		} else {
			fe = max(a[i], a[i-1])
		}
		if fo <= fe {
			f[i] = fo
			odd[i] = true
		} else {
			f[i] = fe
			odd[i] = false
		}
	}

	// Determine starting point
	p := 0
	if n > 0 && f[n-1] <= f[n] {
		p = n - 1
	} else {
		p = n
	}

	ans := make([]int, 0, n)
	// operate function
	operate := func(pos int) {
		if pos <= 0 || pos >= n {
			return
		}
		tmp := a[pos]
		if a[pos+1] < tmp {
			tmp = a[pos+1]
		}
		if tmp == 0 {
			return
		}
		ans = append(ans, pos)
		a[pos] -= tmp
		a[pos+1] -= tmp
	}

	for p > 0 {
		if odd[p] {
			operate(p)
			operate(p - 1)
			p -= 2
		} else {
			operate(p - 1)
			operate(p - 2)
			operate(p)
			p -= 3
		}
		if p < 0 {
			p = 0
		}
	}

	// Output
	fmt.Fprintln(writer, len(ans))
	for _, v := range ans {
		fmt.Fprintln(writer, v)
	}
}
