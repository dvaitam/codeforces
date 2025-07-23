package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}

	arr := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &arr[i])
	}

	// Precompute Fibonacci numbers modulo m up to n + 2
	fib := make([]int64, n+2)
	if n >= 1 {
		fib[1] = int64(1 % m)
	}
	if n >= 2 {
		fib[2] = int64(1 % m)
	}
	for i := 3; i <= n+1; i++ {
		fib[i] = (fib[i-1] + fib[i-2]) % int64(m)
	}

	var q int
	fmt.Fscan(reader, &q)
	for ; q > 0; q-- {
		var l, r int
		fmt.Fscan(reader, &l, &r)
		l--
		r--
		uniq := make(map[int64]struct{})
		for i := l; i <= r; i++ {
			if _, ok := uniq[arr[i]]; !ok {
				uniq[arr[i]] = struct{}{}
			}
		}
		values := make([]int64, 0, len(uniq))
		for k := range uniq {
			values = append(values, k)
		}
		sort.Slice(values, func(i, j int) bool { return values[i] < values[j] })
		res := int64(0)
		for i, v := range values {
			idx := i + 1
			for idx >= len(fib) {
				fib = append(fib, (fib[len(fib)-1]+fib[len(fib)-2])%int64(m))
			}
			res = (res + (v%int64(m))*fib[idx]%int64(m)) % int64(m)
		}
		fmt.Fprintln(writer, res%int64(m))
	}
}
