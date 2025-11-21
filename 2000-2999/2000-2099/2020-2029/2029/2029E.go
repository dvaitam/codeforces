package main

import (
	"bufio"
	"fmt"
	"os"
)

const maxVal = 400000

var spf [maxVal + 1]int
var isPrime [maxVal + 1]bool

func sieve() {
	for i := 2; i <= maxVal; i++ {
		if spf[i] == 0 {
			spf[i] = i
			isPrime[i] = true
			if i > maxVal/i {
				continue
			}
			for j := i * i; j <= maxVal; j += i {
				if spf[j] == 0 {
					spf[j] = i
				}
				isPrime[j] = false
			}
		}
	}
	isPrime[0], isPrime[1] = false, false
	for i := 2; i <= maxVal; i++ {
		if spf[i] == 0 {
			spf[i] = i
		}
	}
}

func main() {
	sieve()
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n)
		minVal := maxVal + 1
		primeCount := 0
		primeVal := -1
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
			if a[i] < minVal {
				minVal = a[i]
			}
			if isPrime[a[i]] {
				primeCount++
				primeVal = a[i]
			}
		}
		if primeCount == 0 {
			fmt.Fprintln(out, 2)
			continue
		}
		if primeCount > 1 {
			fmt.Fprintln(out, -1)
			continue
		}
		p := primeVal
		if p != minVal {
			fmt.Fprintln(out, -1)
			continue
		}
		ok := true
		twice := 2 * p
		for _, y := range a {
			if y == p {
				continue
			}
			if isPrime[y] {
				ok = false
				break
			}
			if y < twice {
				ok = false
				break
			}
			if y == twice {
				continue
			}
			if y == twice+1 {
				ok = false
				break
			}
			s := spf[y]
			if y-s < twice {
				ok = false
				break
			}
		}
		if ok {
			fmt.Fprintln(out, p)
		} else {
			fmt.Fprintln(out, -1)
		}
	}
}
