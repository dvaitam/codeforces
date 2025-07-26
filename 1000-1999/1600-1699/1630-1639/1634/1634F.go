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

	var n, q int
	var mod int64
	fmt.Fscan(reader, &n, &q, &mod)

	a := make([]int64, n+1)
	b := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &b[i])
	}

	diff := make([]int64, n+3)
	arr := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		arr[i] = (a[i] - b[i]) % mod
		if arr[i] < 0 {
			arr[i] += mod
		}
	}
	notZero := 0
	for i := 1; i <= n; i++ {
		var val int64
		if i == 1 {
			val = arr[i]
		} else if i == 2 {
			val = arr[i] - arr[i-1]
		} else {
			val = arr[i] - arr[i-1] - arr[i-2]
		}
		val %= mod
		if val < 0 {
			val += mod
		}
		diff[i] = val
		if val != 0 {
			notZero++
		}
	}

	fib := make([]int64, n+3)
	fib[1], fib[2] = 1%mod, 1%mod
	for i := 3; i <= n+2; i++ {
		fib[i] = (fib[i-1] + fib[i-2]) % mod
	}

	modify := func(idx int, delta int64) {
		if idx >= len(diff) {
			return
		}
		old := diff[idx]
		diff[idx] = (diff[idx] + delta) % mod
		if diff[idx] < 0 {
			diff[idx] += mod
		}
		if idx <= n {
			if old == 0 && diff[idx] != 0 {
				notZero++
			} else if old != 0 && diff[idx] == 0 {
				notZero--
			}
		}
	}

	for ; q > 0; q-- {
		var c string
		var l, r int
		fmt.Fscan(reader, &c, &l, &r)
		sign := int64(1)
		if c == "B" {
			sign = -1
		}
		k := r - l + 1
		modify(l, sign)
		modify(r+1, -sign*fib[k+1])
		modify(r+2, -sign*fib[k])
		if notZero == 0 {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
