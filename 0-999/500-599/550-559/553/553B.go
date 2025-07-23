package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	var k uint64
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}
	// precompute Fibonacci numbers f[i] = number of valid permutations for i remaining numbers
	// f[0]=1, f[1]=1
	fib := make([]uint64, n+2)
	fib[0] = 1
	fib[1] = 1
	for i := 2; i <= n; i++ {
		fib[i] = fib[i-1] + fib[i-2]
	}
	if k > fib[n] {
		k = fib[n]
	}
	res := make([]int, 0, n)
	i := 1
	for i <= n {
		count := fib[n-i]
		if k <= count {
			res = append(res, i)
			i++
		} else {
			k -= count
			// pair i with i+1
			res = append(res, i+1, i)
			i += 2
		}
	}
	writer := bufio.NewWriter(os.Stdout)
	for idx, v := range res {
		if idx > 0 {
			writer.WriteByte(' ')
		}
		fmt.Fprint(writer, v)
	}
	writer.WriteByte('\n')
	writer.Flush()
}
