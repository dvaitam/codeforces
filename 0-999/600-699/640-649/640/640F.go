package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var a, b int
	if _, err := fmt.Fscan(reader, &a, &b); err != nil {
		return
	}
	if b < 2 {
		fmt.Println(0)
		return
	}
	if a < 2 {
		a = 2
	}
	// Sieve of Eratosthenes up to b
	n := b
	isPrime := make([]bool, n+1)
	for i := 2; i <= n; i++ {
		isPrime[i] = true
	}
	for p := 2; p*p <= n; p++ {
		if isPrime[p] {
			for m := p * p; m <= n; m += p {
				isPrime[m] = false
			}
		}
	}
	cnt := 0
	for i := a; i <= b; i++ {
		if isPrime[i] {
			cnt++
		}
	}
	fmt.Println(cnt)
}
