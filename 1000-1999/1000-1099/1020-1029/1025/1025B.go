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
	var a, b int
	fmt.Fscan(reader, &a, &b)
	x, y := a, b
	var primes []int
	for i := 2; i*i <= x || i*i <= y; i++ {
		if x%i == 0 || y%i == 0 {
			primes = append(primes, i)
			for x%i == 0 {
				x /= i
			}
			for y%i == 0 {
				y /= i
			}
		}
	}
	if x > 1 {
		primes = append(primes, x)
	}
	if y > 1 && y != x {
		primes = append(primes, y)
	}
	if len(primes) == 0 {
		fmt.Fprintln(writer, -1)
		return
	}
	ok := make([]bool, len(primes))
	for i := range ok {
		ok[i] = true
	}
	alive := len(primes)
	for i := 1; i < n; i++ {
		fmt.Fscan(reader, &a, &b)
		if alive == 0 {
			continue
		}
		for j, p := range primes {
			if ok[j] && a%p != 0 && b%p != 0 {
				ok[j] = false
				alive--
			}
		}
	}
	for j, p := range primes {
		if ok[j] {
			fmt.Fprintln(writer, p)
			return
		}
	}
	fmt.Fprintln(writer, -1)
}
