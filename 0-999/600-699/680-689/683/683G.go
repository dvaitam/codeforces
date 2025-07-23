package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func pow10(n int) int64 {
	p := int64(1)
	for i := 0; i < n; i++ {
		p *= 10
	}
	return p
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	s, _ := reader.ReadString('\n')
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return
	}
	dot := strings.IndexByte(s, '.')
	open := strings.IndexByte(s, '(')
	close := strings.IndexByte(s, ')')
	if dot == -1 || open == -1 || close == -1 {
		return
	}
	A := s[:dot]
	B := s[dot+1 : open]
	C := s[open+1 : close]

	var a int64
	fmt.Sscanf(A, "%d", &a)
	var b int64
	if len(B) > 0 {
		fmt.Sscanf(B, "%d", &b)
	}
	var bc int64
	if len(B)+len(C) > 0 {
		fmt.Sscanf(B+C, "%d", &bc)
	}
	k := len(B)
	m := len(C)
	powK := pow10(k)
	powM := pow10(m)
	denom := powK * (powM - 1)
	numer := bc - b
	totalNumer := a*denom + numer
	g := gcd(totalNumer, denom)
	fmt.Printf("%d/%d\n", totalNumer/g, denom/g)
}
