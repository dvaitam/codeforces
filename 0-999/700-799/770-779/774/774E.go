package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var s string
	if _, err := fmt.Fscan(reader, &s); err != nil {
		return
	}
	var m int64
	if _, err := fmt.Fscan(reader, &m); err != nil {
		return
	}
	n := len(s)
	// Precompute prefix values for s+s
	t := s + s
	pref := make([]int64, 2*n+1)
	for i := 0; i < 2*n; i++ {
		d := int64(t[i] - '0')
		pref[i+1] = (pref[i]*10 + d) % m
	}
	// pow10n = 10^n mod m
	pow10n := int64(1)
	for i := 0; i < n; i++ {
		pow10n = (pow10n * 10) % m
	}
	minRem := int64(-1)
	for i := 0; i < n; i++ {
		if s[i] == '0' {
			continue
		}
		val := (pref[i+n] - (pref[i]*pow10n)%m + m) % m
		if minRem == -1 || val < minRem {
			minRem = val
		}
	}
	fmt.Println(minRem)
}
