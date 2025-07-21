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
	return a
}

// prefixFunc computes prefix function (pi array) for KMP algorithm
func prefixFunc(s string) []int {
	n := len(s)
	pi := make([]int, n)
	for i := 1; i < n; i++ {
		j := pi[i-1]
		for j > 0 && s[i] != s[j] {
			j = pi[j-1]
		}
		if s[i] == s[j] {
			j++
		}
		pi[i] = j
	}
	return pi
}

// check that s is composed by repeating t
func isRepeat(s, t string) bool {
	lt := len(t)
	for i := range s {
		if s[i] != t[i%lt] {
			return false
		}
	}
	return true
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var s1, s2 string
	fmt.Fscan(reader, &s1, &s2)
	n1, n2 := len(s1), len(s2)
	g := gcd(n1, n2)
	t := s1[:g]
	if !isRepeat(s1, t) || !isRepeat(s2, t) {
		fmt.Fprintln(writer, 0)
		return
	}
	pi := prefixFunc(t)
	// minimal period of t
	p := g - pi[g-1]
	if g%p != 0 {
		p = g
	}
	// number of divisors of (g / p)
	d := g / p
	cnt := 0
	for i := 1; i*i <= d; i++ {
		if d%i == 0 {
			cnt++
			if i != d/i {
				cnt++
			}
		}
	}
	fmt.Fprintln(writer, cnt)
}
