package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int = 998244353

// checkGood returns true if binary string q satisfies the conditions for pattern p.
func checkGood(p, q string) bool {
	n := len(p)
	for i := 0; i < n; i++ {
		ch := p[i]
		found := false
		for l := 0; l <= i && !found; l++ {
			for r := i; r < n && !found; r++ {
				if l <= i && i <= r {
					sub := q[l : r+1]
					m := r - l + 1
					count := 0
					for j := 0; j < m; j++ {
						if sub[j] == ch {
							count++
						}
					}
					need := m / 2
					if m%2 != 0 {
						need = (m + 1) / 2
					}
					if count >= need {
						found = true
					}
				}
			}
		}
		if !found {
			return false
		}
	}
	return true
}

var total int

func gen(p []byte, cur []byte, idx int) {
	if idx == len(p) {
		if checkGood(string(p), string(cur)) {
			total++
			if total >= MOD {
				total -= MOD
			}
		}
		return
	}
	cur[idx] = '0'
	gen(p, cur, idx+1)
	cur[idx] = '1'
	gen(p, cur, idx+1)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	var p string
	fmt.Fscan(in, &p)

	total = 0
	cur := make([]byte, n)
	gen([]byte(p), cur, 0)
	fmt.Println(total % MOD)
}
