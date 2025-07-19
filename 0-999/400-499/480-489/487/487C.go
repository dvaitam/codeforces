package main

import (
	"bufio"
	"fmt"
	"os"
)

// isprime checks if n is a prime number.
func isprime(n int) bool {
	if n < 2 {
		return false
	}
	for i := 2; i*i <= n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

// raizprimitiva finds a primitive root modulo p (p must be prime).
func raizprimitiva(p int) int {
	// brute-force search
	for g := 2; ; g++ {
		pot := g
		cnt := 1
		for pot != 1 {
			pot = (pot * g) % p
			cnt++
			if cnt > p {
				break
			}
		}
		if cnt == p-1 {
			return g
		}
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	switch {
	case n == 1:
		fmt.Fprintln(writer, "YES")
		fmt.Fprintln(writer, 1)
	case n == 4:
		fmt.Fprintln(writer, "YES")
		fmt.Fprintln(writer, 1)
		fmt.Fprintln(writer, 3)
		fmt.Fprintln(writer, 2)
		fmt.Fprintln(writer, 4)
	case !isprime(n):
		fmt.Fprintln(writer, "NO")
	case n == 2:
		fmt.Fprintln(writer, "YES")
		fmt.Fprintln(writer, 1)
		fmt.Fprintln(writer, 2)
	default:
		fmt.Fprintln(writer, "YES")
		g := raizprimitiva(n)
		gg := make([]int, n+1)
		gg[0] = 1
		for i := 1; i <= n; i++ {
			gg[i] = gg[i-1] * g % n
		}
		// print sequence
		fmt.Fprintln(writer, 1)
		a := 2
		for i := 0; i < n/2; i++ {
			pos := gg[n-a]
			qos := gg[a]
			if pos != 1 {
				fmt.Fprintln(writer, pos)
			}
			if qos != 1 {
				fmt.Fprintln(writer, qos)
			}
			a += 2
		}
		fmt.Fprintln(writer, n)
	}
}
