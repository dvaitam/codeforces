package main

import (
	"bufio"
	"fmt"
	"os"
)

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	// generate primes up to 1000
	const maxP = 1000
	isPrime := make([]bool, maxP+1)
	for i := 2; i <= maxP; i++ {
		isPrime[i] = true
	}
	for i := 2; i*i <= maxP; i++ {
		if isPrime[i] {
			for j := i * i; j <= maxP; j += i {
				isPrime[j] = false
			}
		}
	}
	primes := make([]int, 0, 200)
	for i := 2; i <= maxP; i++ {
		if isPrime[i] {
			primes = append(primes, i)
		}
	}

	b := make([]int, n)
	ans := 0
	for i := 0; i+1 < n; i++ {
		if gcd(a[i], a[i+1]) != 1 {
			for _, p := range primes {
				if gcd(a[i], p) == 1 && gcd(p, a[i+1]) == 1 {
					b[i] = p
					break
				}
			}
			ans++
		}
	}
	// output
	fmt.Fprintln(writer, ans)
	for i := 0; i < n; i++ {
		fmt.Fprint(writer, a[i])
		fmt.Fprint(writer, " ")
		if i < n-1 && b[i] != 0 {
			fmt.Fprint(writer, b[i])
			fmt.Fprint(writer, " ")
		}
	}
}
