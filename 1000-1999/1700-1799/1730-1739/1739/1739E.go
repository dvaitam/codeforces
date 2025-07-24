package main

import (
	"bufio"
	"fmt"
	"os"
)

func maxIndep(weights []int) int {
	n := len(weights)
	if n == 0 {
		return 0
	}
	dp0, dp1 := 0, weights[0]
	for i := 1; i < n; i++ {
		ndp0 := dp0
		if dp1 > ndp0 {
			ndp0 = dp1
		}
		ndp1 := dp0 + weights[i]
		dp0, dp1 = ndp0, ndp1
	}
	if dp0 > dp1 {
		return dp0
	}
	return dp1
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	var s1, s2 string
	fmt.Fscan(in, &s1)
	fmt.Fscan(in, &s2)

	t := make([]int, n)
	b := make([]int, n)
	for i := 0; i < n; i++ {
		if s1[i] == '1' {
			t[i] = 1
		}
		if s2[i] == '1' {
			b[i] = 1
		}
	}

	chain1 := []int{}
	chain2 := []int{}
	for i := 0; i < n; i++ {
		if (i+1)%2 == 1 {
			chain1 = append(chain1, t[i])
			chain2 = append(chain2, b[i])
		} else {
			chain1 = append(chain1, b[i])
			chain2 = append(chain2, t[i])
		}
	}

	ans := maxIndep(chain1) + maxIndep(chain2)
	fmt.Println(ans)
}
