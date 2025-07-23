package main

import (
	"bufio"
	"fmt"
	"os"
)

func isPalindrome(n int) bool {
	rev := 0
	tmp := n
	for tmp > 0 {
		rev = rev*10 + tmp%10
		tmp /= 10
	}
	return n == rev
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var p, q int
	if _, err := fmt.Fscan(in, &p, &q); err != nil {
		return
	}
	const N = 2000000
	primes := make([]bool, N+1)
	for i := 2; i <= N; i++ {
		primes[i] = true
	}
	for i := 2; i*i <= N; i++ {
		if primes[i] {
			for j := i * i; j <= N; j += i {
				primes[j] = false
			}
		}
	}
	pi := make([]int, N+1)
	rub := make([]int, N+1)
	for i := 1; i <= N; i++ {
		pi[i] = pi[i-1]
		rub[i] = rub[i-1]
		if primes[i] {
			pi[i]++
		}
		if isPalindrome(i) {
			rub[i]++
		}
	}
	ans := -1
	for n := 1; n <= N; n++ {
		if pi[n]*q <= rub[n]*p {
			ans = n
		}
	}
	if ans == -1 {
		fmt.Fprintln(out, "Palindromic tree is better than splay tree")
	} else {
		fmt.Fprintln(out, ans)
	}
}
