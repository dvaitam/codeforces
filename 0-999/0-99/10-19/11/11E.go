package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"
)

func calc(s string) float64 {
	var p, q, o, f int
	for i := 0; i < len(s); i++ {
		var w byte
		if q&1 != 0 {
			w = 'R'
		} else {
			w = 'L'
		}

		if s[i] == 'X' {
			q++
		} else if s[i] == w {
			p++
			q++
			f = 0
		} else {
			p++
			q += 2
			o += f
			if f == 0 {
				f = 1
			} else {
				f = 0
			}
		}
	}
	if q&1 != 0 {
		q++
		o += f
	}
	// Check if percentage > 50% to subtract overlap penalty
	// p/q > 0.5 => 2*p > q
	if 2*p > q {
		p -= o
		q -= 2 * o
	}
	return float64(p) / float64(q)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var t string
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}

	// Build string s by inserting X between duplicates
	var sb strings.Builder
	for i := 0; i < len(t); i++ {
		if i > 0 && t[i] == t[i-1] && t[i] != 'X' {
			sb.WriteByte('X')
		}
		sb.WriteByte(t[i])
	}
	s := sb.String()

	var ans float64
	// Check cyclic condition: if start and end match, we might need to pad X at start or end
	if len(t) > 0 && t[0] == t[len(t)-1] && t[0] != 'X' {
		v1 := calc("X" + s)
		v2 := calc(s + "X")
		if v1 > v2 {
			ans = v1
		} else {
			ans = v2
		}
	} else {
		ans = calc(s)
	}

	// Formatting logic as per C++ code: floor at 8 decimals, then print percentage
	ans = math.Floor(ans*1e8) / 1e8
	fmt.Printf("%.6f\n", ans*100.0)
}
