package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int = 1000000007

// minusOne subtracts one from a non-negative decimal string.
func minusOne(s string) string {
	if s == "0" {
		return "0"
	}
	b := []byte(s)
	i := len(b) - 1
	for i >= 0 && b[i] == '0' {
		b[i] = '9'
		i--
	}
	if i >= 0 {
		b[i]--
	}
	j := 0
	for j < len(b)-1 && b[j] == '0' {
		j++
	}
	return string(b[j:])
}

// count returns the number of valid numbers of the same length as s
// that are <= s, divisible by m, and satisfy the d-magic property.
func count(s string, m, d int) int {
	n := len(s)
	dp := make([][2]int, m)
	dp[0][1] = 1
	for pos := 0; pos < n; pos++ {
		ndp := make([][2]int, m)
		limit := int(s[pos] - '0')
		for r := 0; r < m; r++ {
			for t := 0; t < 2; t++ {
				val := dp[r][t]
				if val == 0 {
					continue
				}
				maxD := 9
				if t == 1 {
					maxD = limit
				}
				if pos%2 == 1 { // even position (1-indexed)
					dig := d
					if dig <= maxD {
						nt := 0
						if t == 1 && dig == maxD {
							nt = 1
						}
						nr := (r*10 + dig) % m
						ndp[nr][nt] += val
						if ndp[nr][nt] >= MOD {
							ndp[nr][nt] -= MOD
						}
					}
				} else { // odd position
					start := 0
					if pos == 0 {
						start = 1
					}
					for dig := start; dig <= maxD; dig++ {
						if dig == d {
							continue
						}
						nt := 0
						if t == 1 && dig == maxD {
							nt = 1
						}
						nr := (r*10 + dig) % m
						ndp[nr][nt] += val
						if ndp[nr][nt] >= MOD {
							ndp[nr][nt] -= MOD
						}
					}
				}
			}
		}
		dp = ndp
	}
	res := dp[0][0] + dp[0][1]
	if res >= MOD {
		res -= MOD
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var m, d int
	fmt.Fscan(in, &m, &d)
	var a, b string
	fmt.Fscan(in, &a)
	fmt.Fscan(in, &b)

	L := len(b)
	resB := count(b, m, d)
	aMinus := minusOne(a)
	resA := 0
	if len(aMinus) == L {
		resA = count(aMinus, m, d)
	}
	ans := resB - resA
	if ans < 0 {
		ans += MOD
	}
	fmt.Fprintln(out, ans)
}
