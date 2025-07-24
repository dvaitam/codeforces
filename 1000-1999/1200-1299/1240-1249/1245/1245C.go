package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 1000000007

// fib computes Fibonacci numbers modulo mod up to n.
func fib(n int) []int64 {
	f := make([]int64, n+2)
	f[0] = 1
	f[1] = 1
	for i := 2; i <= n; i++ {
		f[i] = (f[i-1] + f[i-2]) % mod
	}
	return f
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var s string
	fmt.Fscan(reader, &s)

	for _, c := range s {
		if c == 'w' || c == 'm' {
			fmt.Fprintln(writer, 0)
			return
		}
	}

	n := len(s)
	f := fib(n)
	res := int64(1)
	i := 0
	for i < n {
		if s[i] == 'u' || s[i] == 'n' {
			j := i
			for j < n && s[j] == s[i] {
				j++
			}
			length := j - i
			res = (res * f[length]) % mod
			i = j
		} else {
			i++
		}
	}
	fmt.Fprintln(writer, res)
}
