package main

import (
	"bufio"
	"fmt"
	"os"
)

// bruteForce computes the total sum of madness over all pairs (a,b)
// using exhaustive enumeration. It is correct but only feasible for
// very small n and m.
func bruteForce(n, m int, mod int) int {
	a := make([]int, m)
	b := make([]int, m)
	var dfs func(int)
	res := 0
	dfs = func(idx int) {
		if idx == m {
			seats := make([]bool, n)
			sum := 0
			for i := 0; i < m; i++ {
				pos := a[i]
				var seat int = -1
				if b[i] == 0 { // L -> search right
					for j := pos; j < n; j++ {
						if !seats[j] {
							seat = j
							break
						}
					}
				} else { // R -> search left
					for j := pos; j >= 0; j-- {
						if !seats[j] {
							seat = j
							break
						}
					}
				}
				if seat == -1 {
					return
				}
				seats[seat] = true
				if seat > pos {
					sum += seat - pos
				} else {
					sum += pos - seat
				}
			}
			res = (res + sum) % mod
			return
		}
		for i := 0; i < n; i++ {
			a[idx] = i
			for t := 0; t < 2; t++ {
				b[idx] = t
				dfs(idx + 1)
			}
		}
	}
	dfs(0)
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m, p int
	fmt.Fscan(in, &n, &m, &p)
	// The following brute force approach only works for very small
	// inputs. It is provided as a simple reference implementation
	// for the problem statement and is not intended to handle the
	// full constraints where n can be up to 500.
	ans := bruteForce(n, m, p)
	fmt.Println(ans)
}
