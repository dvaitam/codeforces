package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	p = 998244353
	N = 200010
	K = 12
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)

	scanner.Scan()
	var n int64
	fmt.Sscan(scanner.Text(), &n)
	scanner.Scan()
	var k int
	fmt.Sscan(scanner.Text(), &k)

	f := make([][]int64, N)
	s := make([][]int64, N)
	for i := range f {
		f[i] = make([]int64, K)
		s[i] = make([]int64, K)
	}

	var ans int64
	for i := int64(1); i <= n; i++ {
		f[i][1] = 1
		s[i][1] = i
		if i > 1 {
			for j := 2; j <= k; j++ {
				term1 := (s[i-1][j-1] - f[i-2][j-1] + p) % p
				term2 := (s[i-1][j] - f[i-2][j] + p) % p
				if j == k {
					f[i][j] = (term1 + term2) % p
				} else {
					f[i][j] = term1 % p
				}
				s[i][j] = (s[i-1][j] + f[i][j]) % p
			}
		}
		if i < n {
			ans = (ans + f[i][k] + f[i][k-1]) % p
		}
	}

	fmt.Println(ans)
}
