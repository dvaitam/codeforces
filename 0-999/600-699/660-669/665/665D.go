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
	a := make([]int, n)
	maxVal := 0
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
		if a[i] > maxVal {
			maxVal = a[i]
		}
	}
	// sieve up to twice the max value plus a margin
	limit := maxVal*2 + 2
	if limit < 3 {
		limit = 3
	}
	isPrime := make([]bool, limit)
	for i := 2; i < limit; i++ {
		isPrime[i] = true
	}
	for i := 2; i*i < limit; i++ {
		if isPrime[i] {
			for j := i * i; j < limit; j += i {
				isPrime[j] = false
			}
		}
	}
	// collect ones
	ans := make([]int, 0, n)
	for _, v := range a {
		if v == 1 {
			ans = append(ans, 1)
		}
	}
	// try adding one more element with v+1 prime
	for _, v := range a {
		if v > 1 && v+1 < limit && isPrime[v+1] {
			ans = append(ans, v)
			break
		}
	}
	if len(ans) > 1 {
		fmt.Fprintln(writer, len(ans))
		for i, v := range ans {
			if i > 0 {
				fmt.Fprint(writer, " ")
			}
			fmt.Fprint(writer, v)
		}
		fmt.Fprintln(writer)
		return
	}
	// clear ans and try any prime-sum pair
	ans = ans[:0]
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			s := a[i] + a[j]
			if s < limit && isPrime[s] {
				fmt.Fprintln(writer, 2)
				fmt.Fprintln(writer, a[i], a[j])
				return
			}
		}
	}
	// fallback: first element alone
	fmt.Fprintln(writer, 1)
	fmt.Fprintln(writer, a[0])
}
