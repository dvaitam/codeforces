package main

import (
	"bufio"
	"fmt"
	"os"
)

func sieve(n int) []int {
	isPrime := make([]bool, n+1)
	for i := 2; i <= n; i++ {
		isPrime[i] = true
	}
	for i := 2; i*i <= n; i++ {
		if isPrime[i] {
			for j := i * i; j <= n; j += i {
				isPrime[j] = false
			}
		}
	}
	primes := make([]int, 0)
	for i := 2; i <= n; i++ {
		if isPrime[i] {
			primes = append(primes, i)
		}
	}
	return primes
}

func squareFree(x int, primes []int) int {
	res := 1
	for _, p := range primes {
		if p*p > x {
			break
		}
		cnt := 0
		for x%p == 0 {
			x /= p
			cnt++
		}
		if cnt%2 == 1 {
			res *= p
		}
	}
	if x > 1 {
		res *= x
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	primes := sieve(3200)

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(reader, &n, &k)
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &arr[i])
			arr[i] = squareFree(arr[i], primes)
		}
		seg := 1
		seen := make(map[int]bool)
		for _, v := range arr {
			if seen[v] {
				seg++
				seen = make(map[int]bool)
			}
			seen[v] = true
		}
		fmt.Fprintln(writer, seg)
	}
}
